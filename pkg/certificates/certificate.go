package certificates

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"net"
	"time"
)

func GetValidNotAfter(certPEM []byte) (*time.Time, error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &cert.NotAfter, nil
}

func GetCSR(serviceName *string, namespace *string, key *rsa.PrivateKey) ([]byte, error) {
	subject := pkix.Name{
		Organization: []string{"system:nodes"},
		CommonName:   fmt.Sprintf("system:node:%s.%s.svc", *serviceName, *namespace),
	}
	altNames := []string{
		*serviceName,
		fmt.Sprintf("%s.%s", *serviceName, *namespace),
		fmt.Sprintf("%s.%s.svc", *serviceName, *namespace),
		fmt.Sprintf("%s.%s.svc.cluster", *serviceName, *namespace),
		fmt.Sprintf("%s.%s.svc.cluster.local", *serviceName, *namespace),
	}
	ipAddresses := []net.IP{net.IPv4(127, 0, 0, 1)}

	csr, err := generateCSR(key, &subject, altNames, ipAddresses)
	if err != nil {
		return nil, err
	}

	return csr, nil
}

func generateCSR(privateKey *rsa.PrivateKey, subject *pkix.Name, dnsSANs []string, ipSANs []net.IP) (csr []byte, err error) {
	// Customize the signature for RSA keys, depending on the key size
	var sigType x509.SignatureAlgorithm

	keySize := privateKey.N.BitLen()
	switch {
	case keySize >= 4096:
		sigType = x509.SHA512WithRSA
	case keySize >= 3072:
		sigType = x509.SHA384WithRSA
	default:
		sigType = x509.SHA256WithRSA
	}

	keyUsage, err := marshalKeyUsage(x509.KeyUsage(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment))
	if err != nil {
		return nil, err
	}
	extKeyUsage, err := marshalExtKeyUsage([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, []asn1.ObjectIdentifier{})
	if err != nil {
		return nil, err
	}
	bc, err := marshalBasicConstraints(false, 0, false)
	if err != nil {
		return nil, err
	}

	template := &x509.CertificateRequest{
		Subject:            *subject,
		SignatureAlgorithm: sigType,
		DNSNames:           dnsSANs,
		IPAddresses:        ipSANs,
		ExtraExtensions: []pkix.Extension{
			bc,
			keyUsage,
			extKeyUsage,
		},
	}

	csr, err = x509.CreateCertificateRequest(cryptorand.Reader, template, privateKey)
	if err != nil {
		return nil, err
	}

	csrPemBlock := &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	}

	return pem.EncodeToMemory(csrPemBlock), nil
}
