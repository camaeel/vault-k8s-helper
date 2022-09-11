package k8sProvider

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var FIELDMANAGER string = "application/vault-k8s-helper"

func GetClientSet(overwriteKubeconfig *string) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()

	if err != nil {

		kubeconfig :=
			clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
		if *overwriteKubeconfig != "" {
			kubeconfig = *overwriteKubeconfig
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}
