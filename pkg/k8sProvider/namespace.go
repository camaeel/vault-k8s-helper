package k8sProvider

import (
	"context"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespaceIfNotExists(k8s *kubernetes.Clientset, ctx context.Context, name *string) error {
	_, err := k8s.CoreV1().Namespaces().Get(ctx, *name, metav1.GetOptions{})

	if errors.IsNotFound(err) {
		log.Infof("Namespace %s not found, will be created.", *name)

		ns := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: *name,
			},
		}

		k8s.CoreV1().Namespaces().Create(ctx, &ns, metav1.CreateOptions{})
		return nil
	} else if err == nil {
		log.Infof("Namespace %s found", *name)
	}
	return err
}
