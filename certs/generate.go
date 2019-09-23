package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	rsaBits int = 2048
)

type Cert struct {
	Cert []byte
	Key  []byte
}

type GenArg struct {
	CA          Cert
	ValidFor    time.Duration
	ExtKeyUsage []x509.ExtKeyUsage
	IPAddresses []net.IP
}

func Genereate(a GenArg) (Cert, error) {
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)

	if err != nil {
		return Cert{}, fmt.Errorf("failed to generate private key: %s", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(a.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return Cert{}, fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Stark & Wayne"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		ExtKeyUsage:           a.ExtKeyUsage,
		BasicConstraintsValid: true,
		IPAddresses:           a.IPAddresses,
	}

	var signCert *x509.Certificate
	var signKey *rsa.PrivateKey

	template.IsCA = len(a.CA.Cert) == 0
	if template.IsCA {
		template.KeyUsage = x509.KeyUsageCertSign
		signCert = &template
		signKey = priv
	} else {
		template.KeyUsage = x509.KeyUsageKeyEncipherment
		signCertPem, _ := pem.Decode(a.CA.Cert)
		signCert, err = x509.ParseCertificate(signCertPem.Bytes)
		if err != nil {
			return Cert{}, fmt.Errorf("failed to parse CA cert: %s", err)
		}
		signKeyPem, _ := pem.Decode(a.CA.Key)
		signKey, err = x509.ParsePKCS1PrivateKey(signKeyPem.Bytes)
		if err != nil {
			return Cert{}, fmt.Errorf("failed to parse CA private key: %s", err)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, signCert, &priv.PublicKey, signKey)
	if err != nil {
		return Cert{}, fmt.Errorf("failed to create Certificate: %s", err)
	}

	return Cert{
		Cert: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE",
			Bytes: derBytes}),
		Key: pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv)}),
	}, nil

}
