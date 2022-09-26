package k8sProvider

import (
	"context"
	"fmt"
	"time"

	"github.com/mohae/deepcopy"
	log "github.com/sirupsen/logrus"
	apicertificatesv1 "k8s.io/api/certificates/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	MAXCSRRETRIES int = 20
	CSRRETRYSLEEP int = 5
)

// deletes CSR if it exists.
func DeleteCSR(k8s *kubernetes.Clientset, ctx context.Context, name *string) error {
	err := k8s.CertificatesV1().CertificateSigningRequests().Delete(ctx, *name, metav1.DeleteOptions{})
	if errors.IsNotFound(err) {
		log.Debugf("CSR delete didn't succeed because it was no longer there.")
		return nil
	}
	return err
}

// creates new CSR
func CreateCSR(k8s *kubernetes.Clientset, ctx context.Context, name *string, csr []byte) (*apicertificatesv1.CertificateSigningRequest, error) {
	signerName := "kubernetes.io/kubelet-serving"

	csrApply := apicertificatesv1.CertificateSigningRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: *name,
		},
		Spec: apicertificatesv1.CertificateSigningRequestSpec{
			Request:    csr,
			Groups:     []string{"system:authenticated"},
			SignerName: signerName,
			Usages: []apicertificatesv1.KeyUsage{
				apicertificatesv1.UsageKeyEncipherment,
				apicertificatesv1.UsageDigitalSignature,
				apicertificatesv1.UsageServerAuth,
			},
		},
	}
	crs, err := k8s.CertificatesV1().CertificateSigningRequests().Create(ctx, &csrApply, metav1.CreateOptions{FieldManager: FIELDMANAGER})
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

	for count := 0; count < MAXCSRRETRIES; count++ {
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
