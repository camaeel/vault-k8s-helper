package k8sProvider

import (
	"context"
	"fmt"
	"time"

	"github.com/camaeel/vault-k8s-helper/pkg/certificates"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	applyMetaV1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var MINVALIDDAYS int = 30

func CheckSecretValidity(k8s *kubernetes.Clientset, ctx context.Context, name *string, namespace *string) (bool, error) {
	minCertificateValid, err := time.ParseDuration(fmt.Sprintf("%dh", MINVALIDDAYS*24))
	if err != nil {
		panic(err)
	}

	secret, err := k8s.CoreV1().Secrets(*namespace).Get(ctx, *name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Infof("Secret %s not found in namespace %s, so it will be created", *name, *namespace)
			return false, nil
		}
		return false, err
	}

	log.Infof("Secret %s found in namespace %s", *name, *namespace)
	if crt, ok := secret.Data["vault.crt"]; ok {
		validBefore, err := certificates.GetValidNotAfter(crt)
		if err != nil {
			log.Warnf("Couldn't parse certificate from secret, will recreate it")
			return false, nil
		} else if time.Until(*validBefore) < minCertificateValid {
			log.Infof("Certificate is valid until %s, will be recreated", *validBefore)
			return false, nil
		} else {
			log.Infof("Certificate is valid until %s, won't be recreated", *validBefore)
			return true, nil
		}
	}

	log.Warnf("Secret found, but without certificate data. Secret will be recreated")
	return false, nil

}

func CreateOrReplaceSecret(k8s *kubernetes.Clientset, ctx context.Context, name *string, namespace *string, data map[string][]byte) error {
	kind := "Secret"
	apiVersion := "v1"
	secretData := v1.SecretApplyConfiguration{
		Data: data,
		// Immutable: &immutable,
		// Type:      &secretType,
		ObjectMetaApplyConfiguration: &applyMetaV1.ObjectMetaApplyConfiguration{
			Namespace: namespace,
			Name:      name,
		},
		TypeMetaApplyConfiguration: applyMetaV1.TypeMetaApplyConfiguration{
			Kind:       &kind,
			APIVersion: &apiVersion,
		},
	}
	// options := metav1.PatchOptions{
	// 	// Force: &o.ForceConflicts,
	// }

	_, err := k8s.CoreV1().Secrets(*namespace).Apply(ctx, &secretData, metav1.ApplyOptions{FieldManager: FIELDMANAGER})
	if err != nil {
		panic(err)
	}
	return nil
}
