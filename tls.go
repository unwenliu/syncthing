package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base32"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"time"
)

const (
	tlsRSABits = 2048
	tlsName    = "syncthing"
)

func loadCert(dir string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(path.Join(dir, "cert.pem"), path.Join(dir, "key.pem"))
}

func certId(bs []byte) string {
	hf := sha1.New()
	hf.Write(bs)
	id := hf.Sum(nil)
	return base32.StdEncoding.EncodeToString(id)
}

func newCertificate(dir string) {
	priv, err := rsa.GenerateKey(rand.Reader, tlsRSABits)
	fatalErr(err)

	notBefore := time.Now()
	notAfter := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName: tlsName,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	fatalErr(err)

	certOut, err := os.Create(path.Join(dir, "cert.pem"))
	fatalErr(err)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	okln("Created TLS certificate file")

	keyOut, err := os.OpenFile(path.Join(dir, "key.pem"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	fatalErr(err)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	okln("Created TLS key file")
}
