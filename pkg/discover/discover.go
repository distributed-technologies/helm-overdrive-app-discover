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

type App struct {
	Argocd_app_namespace              string
	Argocd_app_source_repo_url        string
	Argocd_app_source_target_revision string
	Argocd_app_source_path            string
	Application_folder                string
	Name                              string `yaml:"name"`
	Namespace                         string `yaml:"namespace"`
	Project                           string `yaml:"project"`
	Source                            Source `yaml:"source"`
}

type Source struct {
	Helm_repo     string `yaml:"helm_repo"`
	Chart_name    string `yaml:"chart_name"`
	Chart_version string `yaml:"chart_version"`
}

func (app *App) GetValuesFromYamlFile(path string) error {

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &app)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) GenArgoCDApp() error {

	template, err := template.ParseFiles("resources/argocd_application.yaml")
	if err != nil {
		return err
	}

	err = template.Execute(os.Stdout, app)
	if err != nil {
		return err
	}

	return nil
}

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

		// Remove the base from the path
		tmpString := strings.ReplaceAll(path, folder+"/", "")
		tmpStringArray := strings.Split(tmpString, "/")

		// We can then assume that the last two entris in the array represents the app.yaml file and the application folder
		// this gives us the app path without having to know that the name in the app.yaml file and the name of the folder it's in is equal
		tmpString = strings.Join(tmpStringArray[:len(tmpStringArray)-1], "/")

		app.Application_folder = tmpString

		logging.Debug("App: %v\n", app)

		app.GenArgoCDApp()
	}

	return nil
}

func (app *App) getArgoCDEnvs() {

	app.Argocd_app_namespace = os.Getenv("ARGOCD_APP_NAMESPACE")
	app.Argocd_app_source_path = os.Getenv("ARGOCD_APP_SOURCE_PATH")
	app.Argocd_app_source_repo_url = os.Getenv("ARGOCD_APP_SOURCE_REPO_URL")
	app.Argocd_app_source_target_revision = os.Getenv("ARGOCD_APP_SOURCE_TARGET_REVISION")

}

func GetFiles(folder string) ([]string, error) {
	yamlFiles := []string{}
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, ".yaml") {

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