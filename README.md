# Helm-overdrive-app-discover

Helm-overdrive-app-discover, shortened to App-discover, is a support tool build for Helm-overdrive.

App-discover currently serves two purposes:
1. Facilitate a way of generating ArgoCD applications without having to extend helm-overdrive
2. Simplify the bootstrap process

This is currently achived by having two commands:
* [Discover](#discover)
* [Install](#install)

The reason this is not build into Helm-overdrive, is that Helm-overdrive is an extension of Helm, and should be used for templating, not tool specific actions like generating ArgoCD applications.


## Discover
App-discover introduces a new file, usualy called `app.yaml`, the purpose of this file is to store the deployment information needed to build a [ArgoCD application resource](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/#applications).

```yaml
# app.yaml
apiVersion: argocd-discover/v1alpha1
name: argocd
namespace: argocd
createNamespace: true
project: default
source:
  HelmRepo: https://argoproj.github.io/argo-helm
  chartName: argo-cd
  chart_version: 4.4.0
```

The applications that App-discover generates are the same as described in Helm-overdrives [ArgoCD](
https://github.com/distributed-technologies/helm-overdrive/tree/feature/yaml-merge#argocd) section.

When App-discover is initiated, it "[walks](https://pkg.go.dev/path/filepath#Walk)" through a given filesystem starting in `./` and looks for files that contains `apiVersion: argocd-discover/v1alpha1`, for every file it finds, it uses the values in the file, and derives some from the path, to build a ArgoCD application resource, which is then written to stdout.

### flags:
| name | default | decription |
|------|---------|------------|
| --folder | ./ | Folder to find apps in (recursive) |


## Install
This command makes use of the `app.yaml` file and the `helm-overdrive.yaml` config file to get the needed values to run the helm-overdrive template engine on the chart specified `app.yaml` this then prints the templated application to stdout where it can be piped into `kubectl apply -f -`.

The reason for this commands existence is to make it easier to bootstrap ArgoCD for the first time, by running:
```bash
# Assuming the paths are correct
helm-overdrive-app-discover --app_file ./<base-folder>/<app_folder>/app.yaml -c ./helm-overdrive.yaml
```

### Flags
| name | default | decription |
|------|---------|------------|
| --app_file | ./ | The folder of the app you want to install |
| --config / -c | "" | Point to the helm-overdrive config (default is $HOME/.helm-overdrive.yaml) |
