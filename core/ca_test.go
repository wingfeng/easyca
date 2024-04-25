package core

import (
	"crypto/x509"
	"easyca/conf"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateCA(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	type args struct {
		path          string
		dn            string
		maxPathLength int
		validity      int
		overwrite     bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"root", args{
			path:          conf.Default.CertPath,
			dn:            fmt.Sprintf(conf.Default.CA.DefaultDN, "Easyca Root CA"),
			maxPathLength: 1,
			validity:      10 * 365,
			overwrite:     true,
		}},
		// {"Server ICA", args{
		// 	path:          conf.Default.CertPath + "/server",
		// 	dn:            "/CN=idx.local Server intermediate CA",
		// 	maxPathLength: 0,
		// 	validity:      5 * 365,
		// 	overwrite:     true,
		// }},
		// {"Personal ICA", args{
		// 	path:          "/mnt/d/certs/personal",
		// 	dn:            "/CN=idx.local Peronal intermediate CA",
		// 	maxPathLength: 0,
		// 	validity:      5 * 365,
		// 	overwrite:     true,
		// }},
		// {"Sign ICA", args{
		// 	path:          "/mnt/d/certs/sign",
		// 	dn:            "/CN=idx.local Sign Only intermediate CA",
		// 	maxPathLength: 0,
		// 	validity:      5 * 365,
		// 	overwrite:     true,
		// }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateCA(tt.args.path, tt.args.dn, tt.args.maxPathLength, tt.args.validity, tt.args.overwrite)
		})
	}
}

func TestCreateCertificate(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	type args struct {
		path        string
		dn          string
		san         string
		validity    int
		overwrite   bool
		keyUsage    x509.KeyUsage
		extKeyUsage []x509.ExtKeyUsage
	}
	tests := []struct {
		name string
		args args
	}{
		{"服务器签名", args{
			path:        conf.Default.CertPath + "/server/qqs",
			dn:          fmt.Sprintf(conf.Default.CA.DefaultDN, "easyca.local "),
			san:         "*.easyca.local,",
			validity:    365,
			overwrite:   true,
			keyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			extKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}},
		{"客户端签名", args{
			path:        conf.Default.CertPath + "/personal/wing",
			dn:          "/CN=wing@idx.local/CN=wing",
			san:         "WingFeng,wing@idx.local,192.168.0.1",
			validity:    365,
			overwrite:   true,
			keyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
			extKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageEmailProtection},
		}},
		{"文件签名", args{
			path:        conf.Default.CertPath + "/sign/pdf",
			dn:          fmt.Sprintf(conf.Default.CA.DefaultDN, "idx.local"),
			san:         "",
			keyUsage:    x509.KeyUsageDigitalSignature,
			validity:    365,
			extKeyUsage: nil,
			overwrite:   true,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateCertificate(tt.args.path, tt.args.dn, tt.args.san, tt.args.validity, tt.args.overwrite, tt.args.keyUsage, tt.args.extKeyUsage)
		})
	}
}

func TestGetCertInfo(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Root CA Info", args{
			path: "certs/certs.crt",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCertInfo(tt.args.path)
		})
	}
}

func TestCreateCRL(t *testing.T) {
	conf.InitConfig("../conf/config.yaml")
	result := CreateCRL()
	ioutil.WriteFile("/mnt/d/tmp.crl", result, os.ModePerm)
}
