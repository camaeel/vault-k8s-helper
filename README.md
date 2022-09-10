# vault-k8s-helper
Helper application for managing vault in kubernetes clusters.

Takes care of following aspects:
* preparing a certificate for vault pods signed by cluster CA
* initializing a new cluster
* automatic unsealing existing cluster

## Helm chart usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add vault-k8s-helper https://camaeel.github.io/vault-k8s-helper

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
vault-k8s-helper` to see the charts.

To install the vault-k8s-helper chart:

    helm install vault-k8s-helper vault-k8s-helper/vault-k8s-helper

To uninstall the chart:

    helm delete vault-k8s-helper
