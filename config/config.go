package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

var (
	TlsConfig *tls.Config
	C         map[string]string

	defaultValues = map[string]string{
		"API_URL":      "https://scalingo-api-production.appsdeck.eu",
		"API_PREFIX":   "/v1",
		"UNSECURE_SSL": "false",
	}
)

func init() {
	C = make(map[string]string)
	for varName, defaultValue := range defaultValues {
		C[varName] = os.Getenv(varName)
		if C[varName] == "" {
			C[varName] = defaultValue
		}
	}

	TlsConfig = &tls.Config{}
	if C["UNSECURE_SSL"] == "true" {
		TlsConfig.InsecureSkipVerify = true
	} else {
		certChain := decodePem(x509Chain)
		TlsConfig.RootCAs = x509.NewCertPool()
		for _, cert := range certChain.Certificate {
			x509Cert, err := x509.ParseCertificate(cert)
			if err != nil {
				panic(err)
			}
			TlsConfig.RootCAs.AddCert(x509Cert)
		}
		TlsConfig.BuildNameToCertificate()
	}
}

func GenTLSConfig(serverName string) *tls.Config {
	tlsConfig := &tls.Config{}
	*tlsConfig = *TlsConfig
	tlsConfig.ServerName = serverName
	return tlsConfig
}
