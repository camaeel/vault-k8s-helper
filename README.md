# vault-k8s-helper
Helper application for managing vault in kubernetes clusters.

## Helm chart usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add vault-k8s-helper https://camaeel.github.io/vault-k8s-helper

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
vault-k8s-helper` to see the charts.

In the repo there are following charts:
* vault-cert-creator
* vault-autounseal


To install the <CHART> chart:

    helm install RELEASE-NAME vault-k8s-helper/<CHART>

To uninstall the chart:

    helm delete RELEASE-NAME
