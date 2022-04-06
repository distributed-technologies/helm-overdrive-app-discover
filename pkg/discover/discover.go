package discover

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/distributed-technologies/helm-overdrive-app-discover/pkg/logging"
	"gopkg.in/yaml.v3"
)

// Structs for handling the information that comes from the `app.yaml`
type App struct {
	ArgocdAppNamespace            string
	ArgocdAppSourceRepoUrl        string
	ArgocdAppSourceTargetRevision string
	ArgocdAppSourcePath           string
	ApplicationFolder             string
	CreateNamespace               bool   `yaml:"createNamespace"`
	Name                          string `yaml:"name"`
	Namespace                     string `yaml:"namespace"`
	Project                       string `yaml:"project"`
	Source                        Source `yaml:"source"`
}

type Source struct {
	HelmRepo     string `yaml:"helm_repo"`
	ChartName    string `yaml:"chart_name"`
	ChartVersion string `yaml:"chart_version"`
}

// Unmarshals a .yaml file into the app struct
func (app *App) GetValuesFromYamlFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, &app)
}

// Gets a list of *.yaml files that contains `apiVersion: argocd-discover/v1alpha1` string and generates an ArgoCD application resource that is written to stdout
func Discover(folder string) error {
	logging.Debug("folder: %v\n", folder)

	yamlFiles, err := GetFiles(folder)
	if err != nil {
		return err
	}

	for _, path := range yamlFiles {
		logging.Debug("path: %s\n", path)

		var app App

		app.GetValuesFromYamlFile(path)
		app.getArgoCDEnvs()

		path, err = filepath.Rel(folder, filepath.Dir(path))
		if err != nil {
			return err
		}

		app.ApplicationFolder = path

		logging.Debug("App: %v\n", app)

		app.GenArgoCDApp()
	}

	return nil
}

// Using the Golang template engine to generate the ArgoCD app that is then written to stdout
func (app *App) GenArgoCDApp() error {
	template, err := template.New("template").Parse(getTempalte())
	if err != nil {
		return err
	}

	return template.Execute(os.Stdout, app)
}

// Trying to fetch ArgoCD environments
func (app *App) getArgoCDEnvs() {
	app.ArgocdAppNamespace = os.Getenv("ARGOCD_APP_NAMESPACE")
	app.ArgocdAppSourcePath = os.Getenv("ARGOCD_APP_SOURCE_PATH")
	app.ArgocdAppSourceRepoUrl = os.Getenv("ARGOCD_APP_SOURCE_REPO_URL")
	app.ArgocdAppSourceTargetRevision = os.Getenv("ARGOCD_APP_SOURCE_TARGET_REVISION")
}

// Looks up all files in the root, and checks if it contains the 'argocd-discover' apiVersion
func GetFiles(folder string) ([]string, error) {
	yamlFiles := []string{}
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".yaml") {

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "apiVersion: argocd-discover/v1alpha1") {
					yamlFiles = append(yamlFiles, path)
				}
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	logging.Debug("yamlFiles: %s", yamlFiles)
	return yamlFiles, nil
}

// Template for the ArgoCD app
func getTempalte() string {
	return `---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .Name }}
  namespace: {{ .ArgocdAppNamespace }}
  finalizers:
  - resources-finalizer.argocd.argoproj.io
spec:
  project: {{ .Project }}
  destination:
    server: https://kubernetes.default.svc
    namespace: {{ .Namespace }}
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
    syncOptions:
      - CreateNamespace={{ or .CreateNamespace false }}
  source:
    repoURL: {{ .ArgocdAppSourceRepoUrl }}
    targetRevision: {{ .ArgocdAppSourceTargetRevision }}
    path: {{ .ArgocdAppSourcePath }}
    plugin:
      name: helm-overdrive
      env:
      - name: HO_APPLICATION_FOLDER
        value: {{ .ApplicationFolder }}
      - name: HO_HELM_REPO
        value: {{ .Source.HelmRepo }}
      - name: HO_CHART_NAME
        value: {{ .Source.ChartName }}
      - name: HO_CHART_VERSION
        value: {{ .Source.ChartVersion }}
`
}
