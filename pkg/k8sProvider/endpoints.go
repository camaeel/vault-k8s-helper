package k8sProvider

import (
	"context"
	"fmt"

	"golang.org/x/exp/slices"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetServiceEndpoints returns alphabetically ordered list of pod names of given service
func GetServiceEndpoints(k8s *kubernetes.Clientset, ctx context.Context, serviceName *string, namespace *string) ([]string, error) {
	endpoints, err := k8s.CoreV1().Endpoints(*namespace).Get(ctx, *serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if endpoints == nil {
		return nil, fmt.Errorf("endpoints data not found for service %s in namespace %s", *serviceName, *namespace)
	}

	if len(endpoints.Subsets) < 1 {
		return nil, fmt.Errorf("endpoints subsests list is empty for service %s in namespace %s", *serviceName, *namespace)
	} else if len(endpoints.Subsets) > 1 {
		return nil, fmt.Errorf("endpoints list contains more than 1 Subset for service %s in namespace %s", *serviceName, *namespace)
	}

	addrs := endpoints.Subsets[0].Addresses
	if len(addrs) < 1 {
		return nil, fmt.Errorf("endpoint addresses empty for service %s in namespace %s", *serviceName, *namespace)
	}

	result := make([]string, 0)

	for a := range addrs {
		result = append(result, addrs[a].Hostname)
	}

	slices.Sort(result)

	return result, nil
}
