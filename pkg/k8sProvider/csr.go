package k8sProvider

import (
	"context"
	"fmt"
	"time"

	"github.com/mohae/deepcopy"
	log "github.com/sirupsen/logrus"
	apicertificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certificatesv1 "k8s.io/client-go/applyconfigurations/certificates/v1"
	applyMetaV1 "k8s.io/client-go/applyconfigurations/meta/v1"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	MAXCSRRETRIES int = 20
	CSRRETRYSLEEP int = 5
)

func CreateCSR(k8s *kubernetes.Clientset, ctx context.Context, name *string, csr []byte) (*apicertificatesv1.CertificateSigningRequest, error) {
	signerName := "kubernetes.io/kubelet-serving"
	kind := "CertificateSigningRequest"
	apiVersion := "certificates.k8s.io/v1"

	csrApply := certificatesv1.CertificateSigningRequestApplyConfiguration{
		ObjectMetaApplyConfiguration: &v1.ObjectMetaApplyConfiguration{
			Name: name,
		},
		Spec: &certificatesv1.CertificateSigningRequestSpecApplyConfiguration{
			Request:    csr,
			Groups:     []string{"system:authenticated"},
			SignerName: &signerName,
			Usages: []apicertificatesv1.KeyUsage{
				apicertificatesv1.UsageKeyEncipherment,
				apicertificatesv1.UsageDigitalSignature,
				apicertificatesv1.UsageServerAuth,
			},
		},
		TypeMetaApplyConfiguration: applyMetaV1.TypeMetaApplyConfiguration{
			Kind:       &kind,
			APIVersion: &apiVersion,
		},
	}
	crs, err := k8s.CertificatesV1().CertificateSigningRequests().Apply(ctx, &csrApply, metav1.ApplyOptions{FieldManager: FIELDMANAGER})
	return crs, err
}

func ApproveCSR(k8s *kubernetes.Clientset, ctx context.Context, name *string, csr *apicertificatesv1.CertificateSigningRequest) error {

	csrCopy := deepcopy.Copy(*csr).(apicertificatesv1.CertificateSigningRequest)
	csrCopy.Status.Conditions = append(csrCopy.Status.Conditions, apicertificatesv1.CertificateSigningRequestCondition{
		Type:           apicertificatesv1.CertificateApproved,
		Reason:         "Approved",
		Message:        "Approved by vault-k8s-helper",
		LastUpdateTime: metav1.Now(),
		Status:         corev1.ConditionTrue,
	})

	_, err := k8s.CertificatesV1().CertificateSigningRequests().UpdateApproval(ctx, *name, &csrCopy, metav1.UpdateOptions{})
	return err
}

func GetCSRCertificate(k8s *kubernetes.Clientset, ctx context.Context, name *string) (csrContent []byte, err error) {
	log.Infof("Waiting for CSR %s to have certificate generated...", *name)
	var csr *apicertificatesv1.CertificateSigningRequest
	var count int = 0

	for count = 0; count < MAXCSRRETRIES; count++ {
		csr, err = k8s.CertificatesV1().CertificateSigningRequests().Get(ctx, *name, metav1.GetOptions{})
		if err != nil {
			return
		}
		if len(csr.Status.Certificate) != 0 {
			csrContent = csr.Status.Certificate
			log.Infof("Received signed certificate after %d iterations", count+1)
			return
		}
		time.Sleep(time.Duration(CSRRETRYSLEEP) * time.Second)
	}
	err = fmt.Errorf("Certificate still not generated for CSR %s", *name)
	return

}
