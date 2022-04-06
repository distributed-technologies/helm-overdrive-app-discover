# helm-overdrive

Helm-overdrive is a Cli tool that facilitates templating of a base and a environment config onto a application.

This app can template global value files on top of applications.
The app requires one `base_folder` that contains a `global_file`,
inside the `base_folder`, beside the `global_file`, there should also be an `application_folder`,
this should point to a folder that represents the helm chart that you want to deploy.
The `application_folder` then contains the `values_file`, this is the values that you want to change in the chart.
Also in the `application_folder` the can be an optional `additional_resources` folder,
here you can place fully templatable yaml files in case the chart is missing a resource you want beside it.

So the folder structure could look like this:
``` tree
📂base
┣ 📂applications
┃ ┗ 📂hello-world
┃   ┣ 📂additional_resources (optional)
┃   ┃ ┗ 📜cm.yaml
┃   ┗ 📜values.yaml
┗ 📜global.yaml
```
To see a full folder structure take a look at [config-example](/config-example/)

Here the `base_folder` is called `base`, with the `global_file` called `global.yaml` beside this is the `application_folder` is called `applications/hello-world`,
inside `applications/hello-world` is the `values_file` called `values.yaml` and the `additional_resources` folder containing an additional config map

## Using the tool
To template global on top of the `hello-world` application, using helm-overdrive, we have several options to define out folder structure when using the cli tool.

1. [flags](#using-flags)
2. [environment variables](#using-environment-variables)
3. [config file](#using-config)

In all three methods we will be setting the following values:
 - chart_name: The name of the helm chart you want to template
 - chart_version: The version of the chart you want to template
 - helm_repo: The repo url that the chart you want to template is located
 - base_folder: The path to the root folder of the base
 - global_file: The name of the global file
 - application_folder: The path from the `base_folder` to the application you want to tempalte
 - values_file: The name of the values file in the `application_folder`
 - additional_resources: The path from the `application_folder` to the folder containing the additional resources

The app is using the [viper](https://github.com/spf13/viper/) config package
so the priority for the input methods is:
  - flag
  - env
  - config

### Using flags
```
helm-overdrive template \
--chart_name hello-world \
--chart_version 0.1.0 \
--helm_repo https://helm.github.io/examples \
--base_folder base \
--global_file global.yaml \
--values_file values.yaml \
--application_folder applications/hello-world
--additional_resources additional_resources
```

Look at all the configs flag names [here](#config)

### Using environment variables
OBS! All envs for this app is prefixed with `HO`<br>

``` bash
export HO_CHART_NAME="hello-world" \
HO_CHART_VERSION="0.1.0" \
HO_HELM_REPO="https://helm.github.io/examples" \
HO_BASE_FOLDER="base" \
HO_GLOBAL_FILE="global.yaml" \
HO_VALUES_FILE="values.yaml" \
HO_APPLICATION_FOLDER="applications/hello-world" \
HO_ADDITIONAL_RESOURCES="additional_resources"

helm-overdrive template
```
Available environment vars can be found [here](#config)

### using-config
helm-overdrive.yaml would contain >
``` yaml
chart_name: hello-world
chart_version: "0.1.0"
helm_repo: https://helm.github.io/examples

base_folder: base
global_file: global.yaml
values_file: values.yaml
application_folder: applications/hello-world
additional_resources: additional_resources
```

After saving the config file
We can run the template command with the `--config / -c` flag pointing to the config file

``` bash
helm-overdrive --config helm-overdrive.yaml
```

Available keys can be found [here](#config)

## ArgoCD
Helm-overdrive is inteted to be used as a plugin source for [ArgoCD](https://argo-cd.readthedocs.io/en/stable/)

This can be done by making a HO repository, containing the `helm-overdrive.yaml` config and at least the base folder setup.

the `helm-overdrive.yaml` config would then contain the general folder structure.

then using [ArgoCD's guide](https://argo-cd.readthedocs.io/en/stable/operator-manual/custom_tools/) on how to add custom tools and add a [configManagementPlugin](https://argo-cd.readthedocs.io/en/stable/user-guide/config-management-plugins/)

We can add helm overdrive to ArgoCD to be used as a plugin.

Then for the app you want to install you create a application resource, pointing to the plugin, and fill out the last fields as environment variables
``` yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: hello-world
  namespace: argocd
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  destination:
    server: https://kubernetes.default.svc
    namespace: default
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
      allowEmpty: false
  source:
    repoURL: <pointing to the helm-overdrive repo you created>
    targetRevision: <the branch on that repo>
    path: <path to the root of the repo>
    plugin:
      name: helm-overdrive
      env:
      - name: HO_APPLICATION_FOLDER
        value: applications/hello-world
      - name: HO_HELM_REPO
        value: https://helm.github.io/examples
      - name: HO_CHART_NAME
        value: hello-world
      - name: HO_CHART_VERSION
        value: "0.1.0"
```

This should create a application called hello-world, that has been templated using helm-overdrive.

## code flow-chart
The diagram for the code flow can be seen here:<br>
[code-flow](docs/code-diagram.drawio.svg)

## config
| flag | env | config | description |
|------|-----|--------|-------------|
| --additional_resources | HO_ADDITIONAL_RESOURCES | additional_resources | Path to the folder that contains the additional resources, this has to be located within the <application_folder>, Same in base and env folders |
| --application_folder | HO_APPLICATION_FOLDER | application_folder | Path to the folder that contains the application, Same in base and env folders |
| --app_name | HO_APP_NAME | app_name | Name of the release |
| --base_folder | HO_BASE_FOLDER | base_folder | Path the folder containing the base config |
| --env_folder / -e | HO_ENV_FOLDER | env_folder | Name of the environment folder you with to deploy |
| --chart_version / -v | HO_CHART_VERSION | chart_version | Chart version |
| --chart_name / -n | HO_CHART_NAME | chart_name | Chart |
| --global_file | HO_GLOBAL_FILE | global_file | Name of the global files, same in base and env folders |
| --helm_repo | HO_HELM_REPO | helm_repo | Repo url |
| --values_file | HO_VALUE_FILES | values_file | Name of the value files in the application folder, Same in base and env folders |

## File structure
The file structure currently needed to use this app
``` tree
📦config-example
 ┣ 📂base
 ┃ ┣ 📂applications
 ┃ ┃ ┗ 📂hello-world
 ┃ ┃   ┣ 📂additional_resources (optional)
 ┃ ┃   ┃ ┗ 📜cm.yaml
 ┃ ┃   ┗ 📜values.yaml
 ┃ ┣ 📜global.yaml
 ┣ 📂env (optional)
 ┃ ┗ 📂test
 ┃   ┣ 📂applications
 ┃   ┃ ┗ 📂hello-world
 ┃   ┃   ┣ 📂additional_resources (optional)
 ┃   ┃   ┃ ┗ 📜sec.yaml
 ┃   ┃   ┗ 📜values.yaml
 ┃   ┗ 📜global.yaml
 ┗ 📜helm-overdrive.yaml
```
