# vault-k8s-helper
Helper application for managing vault with Raft storage & TLS in kubernetes clusters.

Takes care of following aspects:
* preparing a certificate for vault pods signed by cluster CA
* initializing a new cluster
* joining raft nodes to the first node
* automatic unsealing existing cluster

Application consists of 2 commands:
1. setup-tls - creates certificates for vault cluster signed by cluster's CA
1. vault-autounseal - initialzies, stores secrets and unseals vault.

## Helm chart usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add vault-k8s-helper https://camaeel.github.io/vault-k8s-helper

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
vault-k8s-helper` to see the charts.

Helm repository has 2 helm charts:
1. vault-cert-creator - which installs setup-tls tool and provides secrets for the vault. It will also renew certificates if they are near to be expired
1. vault-autounseal - this chart sets up vault-autounseal utility which is responsible for initializing and establishing a new cluster and unsealing sealed pods.


# Installation

Prefered way of instalation is using helm charts. Simplest setup can be achieved using following steps:
1. Install setup-tls:
    ```bash
    helm upgrade --install -n vault --create-namespace vault-cert-creator vault-cert-creator --repo https://camaeel.github.io/vault-k8s-helper/
    ```
1. Install vault
    ```bash
    helm upgrade --install -n vault --create-namespace vault vault --repo https://helm.releases.hashicorp.com/ --version 0.22.0 -f example/vault/vault-values.yaml
    ```
1. Install vault-autounseal
    ```bash
    helm upgrade --install -n vault-autounseal --create-namespace vault-autounseal vault-autounseal --repo https://camaeel.github.io/vault-k8s-helper/
    ```
