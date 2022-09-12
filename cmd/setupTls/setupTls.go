package main

import (
	"context"
	"flag"

	"github.com/camaeel/vault-k8s-helper/pkg/certificates"
	"github.com/camaeel/vault-k8s-helper/pkg/k8sProvider"
	log "github.com/sirupsen/logrus"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "", "Overwrite kubeconfig path")
	vaultNamespace := flag.String("vault-namespace", "vault", "Vault namespace")
	vaultSecretName := flag.String("vault-secret", "vault-server-tls", "Secret containing vault TLS certificates")
	vaultServiceName := flag.String("vault-service", "vault", "Vault service name")
	csrName := flag.String("vault-csr", "vault-server-tls", "Vault CSR name")
	flag.Parse()
	ctx := context.TODO()

	k8s, err := k8sProvider.GetClientSet(kubeconfig)
	if err != nil {
		panic(err)
	}

	keepCertificate, err := k8sProvider.CheckSecretValidity(k8s, ctx, vaultSecretName, vaultNamespace)
	if err != nil {
		panic(err)
	}

	if !keepCertificate {
		key := certificates.GenerateKey(4096)
		csr, err := certificates.GetCSR(vaultServiceName, vaultNamespace, key)
		if err != nil {
			panic(err)
		}

		err = k8sProvider.DeleteCSR(k8s, ctx, csrName)
		if err != nil {
			panic(err)
		}

		k8sCsr, err := k8sProvider.CreateCSR(k8s, ctx, csrName, csr)
		if err != nil {
			panic(err)
		}

		err = k8sProvider.ApproveCSR(k8s, ctx, csrName, k8sCsr)
		if err != nil {
			panic(err)
		}

		cert, err := k8sProvider.GetCSRCertificate(k8s, ctx, csrName)
		if err != nil {
			panic(err)
		}

		secretData := map[string][]byte{
			"vault.key": certificates.KeyToString(key),
			"vault.crt": cert,
		}
		err = k8sProvider.CreateNamespaceIfNotExists(k8s, ctx, vaultNamespace)
		if err != nil {
			panic(err)
		}

		err = k8sProvider.CreateOrReplaceSecret(k8s, ctx, vaultSecretName, vaultNamespace, secretData)
		if err != nil {
			panic(err)
		}

		log.Infof("Secret %s created in namespace %s", *vaultSecretName, *vaultNamespace)
	}

}
