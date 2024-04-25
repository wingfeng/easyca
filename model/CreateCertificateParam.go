package model

import "crypto/x509"

type CreateCertificateParam struct {
	Dn            string              `json:"dn"`
	Validity      int                 `json:"validity"`
}

type CreateCertificateArgs struct {
	Path          string              `json:"path"`
	Dn            string              `json:"dn"`
	San           string              `json:"san"`
	Overwrite     bool
	Validity      int                 `json:"validity"`
	KeyUsage      x509.KeyUsage       `json:"key_usage"`
	ExtKeyUsage   []x509.ExtKeyUsage  `json:"ext_key_usage"`
}
