package core

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"easyca/conf"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
var privatePerms os.FileMode = 0600
var publicPerms os.FileMode = 0644

// GetCertInfo 获取证书的信息
func GetCertInfo(path string) (x509.Certificate, error) {
	// rootPath := "certs/ca.pem"
	// rootData, err := ioutil.ReadFile(rootPath)
	// if err != nil {
	// 	log.Errorf("读取证书信息错误!", err)
	// }
	// roots := x509.NewCertPool()
	// roots.AppendCertsFromPEM(rootData)

	cert := &x509.Certificate{}

	certData, err := os.ReadFile(path)
	if err != nil {
		slog.Error("证书文件读取失败", "error", err)
		return *cert, err
	}
	block, _ := pem.Decode(certData)
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		slog.Error("证书文件解析失败", "error", err)
	}
	return *cert, err
}

// CreateCA path=ca
// dn=/CN=easycert-root-ca
// maxPathLenght=1 可以有一级中间证书
// validity=10*365 10年
// overwrite允许覆盖现有的证书
func CreateCA(path string, dn string, maxPathLength int, validity int, overwrite bool) error {

	slog.Info("Creating Certificate Authority with Subject: \n", "path", path, "dn", dn)

	if !overwrite {
		checkExisting(path)
	}

	//	ca := filepath.Dir(path)
	var caCert *x509.Certificate
	var err error
	var caKey *rsa.PrivateKey
	//不支持中间证书颁发机构
	//	isRoot := ca == filepath.Dir(conf.Default.CertPath)

	// if !isRoot {
	// 	caCert = parseCert(ca)
	// 	if !caCert.IsCA {
	// 		//log.Error("Certificate %s is not a certificate authority", ca)
	// 		return fmt.Errorf(fmt.Sprintf("上级证书 %s不是证书颁发机构", ca), err)
	// 	} else if !(caCert.MaxPathLen > 0) {
	// 		return fmt.Errorf("证书颁发机构 %s 不能再签发新的证书颁发机构(最大长度限制)", ca)
	// 		//log.Error("Certificate Authority %s can't sign other certificate authorities (maxPathLength exceeded)", ca)
	// 	}
	// 	maxPathLength = caCert.MaxPathLen - 1
	// 	caKey = parseKey(ca)
	// }

	key, derKey, err := generatePrivateKey()
	if err != nil {
		//log.Error("Error generating private key: %s", err)
		return fmt.Errorf("error generating private key: %s", err.Error())
	}

	notBefore := time.Now().UTC()
	notAfter := notBefore.AddDate(0, 0, validity)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		//log.Error("Failed to generate serial number: %s", err)
		return fmt.Errorf("failed to generate serial number: %s", err.Error())
	}

	subject, err := parseDn(caCert, dn)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            maxPathLength,
		MaxPathLenZero:        maxPathLength == 0,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	if caCert == nil {
		caCert = &template
		caKey = key
	}

	derCert, err := x509.CreateCertificate(rand.Reader, &template, caCert, &key.PublicKey, caKey)

	if err != nil {
		//log.Error("Failed to create CA Certificate: %s", err)
		return fmt.Errorf("failed to create CA Certificate: %s", err.Error())
	}
	saveCert(path, true, derCert)
	saveKey(path, true, derKey)

	copyFile(filepath.Join(path, "ca.crt"), filepath.Join(path, "ca.pem"), publicPerms)

	slog.Info("Finished Creating Certificate Authority with Subject: ", "path", path, "dn", dn)

	return nil
}

// CreateCRL 生成吊销证书信息
// 现在先生成个空的吊销证书列表，后续再实现逻辑
func CreateCRL() []byte {
	ca := filepath.Join(conf.Default.CertPath, "ca.crt")
	caKeyFile := filepath.Join(conf.Default.CertPath, "ca.key")
	caCert := parseCert(ca)
	caKey := parseKey(caKeyFile)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)
	temp := &x509.RevocationList{
		SignatureAlgorithm: caCert.SignatureAlgorithm,
		Number:             serialNumber,
		ThisUpdate:         time.Now(),
		NextUpdate:         time.Now().Add(2 * time.Hour),
	}
	result, err := x509.CreateRevocationList(rand.Reader, temp, caCert, caKey)
	if err != nil {
		slog.Error("生成吊销证书信息错误!", "error", err)
	}
	return result

}

// CreateCertificate path
// createCertificate("ca/server", "/CN=server", "(服务器域名:*.idx.local),www.idx.local", 365+5,
//
//	x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment,
//	[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth})
func CreateCertificate(path string, dn string, san string, validity int, overwrite bool, keyUsage x509.KeyUsage, extKeyUsage []x509.ExtKeyUsage) (err error) {
	slog.Info("Creating Certificate with Subject:", "path", path, "dn", dn)

	if !overwrite {
		checkExisting(path)
	}

	caPath := filepath.Join(conf.Default.CertPath, "ca.crt")

	caCert := parseCert(caPath)
	if !caCert.IsCA {
		//errorLog.Error("Certificate %s is not a certificate authority", filepath.Dir(path))
		return fmt.Errorf("certificate %s is not a certificate authority", filepath.Dir(path))
	}
	caKeyFile := filepath.Join(conf.Default.CertPath, "ca.key")
	caKey := parseKey(caKeyFile)

	key, derKey, err := generatePrivateKey()
	if err != nil {
		//errorLog.Error("Error generating private key: %s", err)
		return fmt.Errorf("error generating private key: %s", err.Error())
	}

	notBefore := time.Now().UTC().Add(-10 * time.Minute) // -10 min to mitigate clock skew
	notAfter := notBefore.AddDate(0, 0, validity).Add(10 * time.Minute)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		//errorLog.Error("Failed to generate serial number: %s", err)
		return fmt.Errorf("faile to generate serial number: %s", err.Error())
	}

	subject, err := parseDn(caCert, dn)
	if err != nil {
		return err
	}
	crlPath := fmt.Sprintf("%s%s", conf.CurrentHost, conf.Default.CA.CRLURL)
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               *subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  false,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		EmailAddresses:        []string{},
		CRLDistributionPoints: []string{crlPath},
		//	OCSPServer:            []string{"http://crl.gdca.com.cn/crl/GDCA_TrustAUTH_R4_Primer_CA.crl"}, //吊销证书查阅路径
	}

	parseSubjectAlternativeNames(san, &template)

	derCert, err := x509.CreateCertificate(rand.Reader, &template, caCert, &key.PublicKey, caKey)
	if err != nil {
		//errorLog.Error("Failed to create Server Certificate %s: %s", path, err)
		return fmt.Errorf("failed to create Server Certificate %s: %s", path, err.Error())
	}

	saveCert(path, false, derCert)
	saveKey(path, false, derKey)
	savepfx(derCert, key, caCert, path, conf.Default.DefaultPWD)
	copyFile(filepath.Join(conf.Default.CertPath, "ca.pem"), filepath.Join(path, "ca.pem"), publicPerms)
	slog.Info("Finished Create Certificate with Subject:", "path", path, "dn", dn)

	return
}

func parseSubjectAlternativeNames(san string, template *x509.Certificate) {
	slog.Info("Parsing Subject Alternative Names:", "san", san)
	if san != "" {
		template.IPAddresses = []net.IP{}
		template.DNSNames = []string{}
		for _, h := range strings.Split(san, ",") {
			slog.Info("Parsing", "host", h)
			if ip := net.ParseIP(h); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else if email := parseEmailAddress(h); email != nil {
				template.EmailAddresses = append(template.EmailAddresses, email.Address)
			} else {
				template.DNSNames = append(template.DNSNames, h)

			}
		}
	}
}

// implemented as a seperate function because net.mail.ParseAddress
// panics on malformed addresses
func parseEmailAddress(address string) (email *mail.Address) {
	defer func() {
		if recover() != nil {
			email = nil
		}
	}()
	var err error
	email, err = mail.ParseAddress(address)
	if err == nil && email != nil {
		return email
	}
	return nil
}
func generatePrivateKey() (*rsa.PrivateKey, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	derKey, err := x509.MarshalPKCS8PrivateKey(key)

	if err != nil {
		return nil, nil, err
	}
	return key, derKey, nil
}

func checkExisting(path string) {
	fullPath := filepath.Join(path, filepath.Base(path))
	const errMsg = "Skipping creation of %s because file %s already exists.\nUse the \"-overwrite\" option to overwrite the existing file."
	if _, err := os.Stat(fullPath + ".crt"); err == nil {
		slog.Error(fmt.Sprintf(errMsg, path, "./"+fullPath+".crt"))
	}
	if _, err := os.Stat(fullPath + ".crt"); err == nil {
		slog.Error(fmt.Sprintf(errMsg, path, "./"+fullPath+".key"))
	}
	if _, err := os.Stat(filepath.Join(path, "ca.pem")); err == nil {
		slog.Error(fmt.Sprintf("Skipping creation of %s because file %s already exists.\nUse the \"-overwrite\" option to overwrite the existing file.", path, filepath.Join(path, "ca.pem")))
	}
}

func createDirectory(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		var publicPerms os.FileMode = 0755
		if err := os.MkdirAll(directory, publicPerms); err != nil {
			slog.Error("Error creating directory ", "directory", directory, "error", err)
		}
	}
}

func saveCert(directory string, isRoot bool, derCert []byte) {
	createDirectory(directory)
	fileName := filepath.Base(directory) + ".crt"
	if isRoot {
		fileName = "ca.crt"
	}
	fileName = filepath.Join(directory, fileName)

	slog.Info("Saving", "file", fileName)

	certFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, publicPerms)
	if err != nil {
		slog.Error("Failed to open %s for writing: %s", fileName, err)
	}
	defer func() {
		if err := certFile.Close(); err != nil {
			slog.Error("Failed to save %s: %s", fileName, err)
		}
	}()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derCert}); err != nil {
		slog.Error("Failed to marshall %s: %s", fileName, err)
	}

	if !isRoot {
		rootPath := conf.Default.CertPath
		caPath := filepath.Join(rootPath, "ca.crt")
		caFile, err := os.Open(caPath)
		if err != nil {
			slog.Error("Failed to open ca certificate:", "error", err)
		}
		defer func() {
			if err = caFile.Close(); err != nil {
				slog.Error("Failed to close:", "CA Path", caPath, "error", err)
			}
		}()
		_, err = io.Copy(certFile, caFile)
		if err != nil {
			slog.Error("Failed to concat ca certificates:", "error", err)
		}
		err = certFile.Sync()
		if err != nil {
			slog.Error("Failed to sync certificate file:", "error", err)
		}
	}
}

func saveKey(directory string, isRoot bool, derKey []byte) {

	fileName := filepath.Base(directory) + ".key"
	if isRoot {
		fileName = "ca.key"
	}
	fileName = filepath.Join(directory, fileName)

	slog.Info("Saving ", "filename", fileName)

	keyFile, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, privatePerms)
	if err != nil {
		slog.Error("Failed to open %s for writing: %s", fileName, err)
	}
	defer func() {
		if err := keyFile.Close(); err != nil {
			slog.Error("Failed to close %s: %s", fileName, err)
		}
	}()
	if err := pem.Encode(keyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: derKey}); err != nil {
		slog.Error("Failed to marshall %s: %s", fileName, err)
	}
}
func savepfx(crtBytes []byte, key interface{}, caCert *x509.Certificate, directory string, pwd string) {
	fileName := filepath.Join(directory, filepath.Base(directory)+".pfx")
	cert, err := x509.ParseCertificate(crtBytes)

	if err != nil {
		slog.Error("解析证书错误!")
	}
	cas := []*x509.Certificate{
		caCert,
	}
	pfx, err := pkcs12.Encode(rand.Reader, key, cert, cas, pwd)
	if err != nil {
		slog.Error("encode pfx error", "error", err)
	}
	err = os.WriteFile(fileName, pfx, os.ModeAppend)
	if err != nil {
		slog.Error("write pfx error", "error", err)
	}
}
func parseCert(fileName string) *x509.Certificate {

	// //逐级网上找，找到根目录位置
	// if _, err := os.Stat(fileName); err != nil && !strings.EqualFold("/", path) {
	// 	parentDir := filepath.Dir(path)
	// 	return parseCert(parentDir)
	// }
	der, err := os.ReadFile(fileName)
	if err != nil {
		errMsg := fmt.Sprintf("读取证书文件%s:%s失败。", fileName, err)
		slog.Error(errMsg)
		return nil
	}
	block, _ := pem.Decode(der)
	if block == nil || block.Type != "CERTIFICATE" {
		errMsg := fmt.Sprintf("解码证书%s:%s失败。", fileName, err)
		slog.Error(errMsg)
		return nil
	}
	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		errMsg := fmt.Sprintf("解析证书%s:%s失败。", fileName, err)
		slog.Error(errMsg)
		return nil

	}
	return crt
}

func parseKey(filename string) *rsa.PrivateKey {

	der, err := os.ReadFile(filename)
	if err != nil {
		slog.Error("Failed to read private key file", "file", filename, "error", err)
	}
	block, _ := pem.Decode(der)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		slog.Error("Failed to decode private key for", "file", filename, "error", err)
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("Failed to parse private key for ", "file", filename, "error", err)
	}
	return key.(*rsa.PrivateKey)
}

func parseDn(ca *x509.Certificate, dn string) (*pkix.Name, error) {
	slog.Info("Parsing distinguished name:", "dn", dn)
	var caName pkix.Name
	if ca != nil {
		caName = ca.Subject
	} else {
		caName = pkix.Name{}
	}
	newName := &pkix.Name{}
	for _, element := range strings.Split(strings.Trim(dn, "/"), "/") {
		value := strings.Split(element, "=")
		if len(value) != 2 {
			return nil, fmt.Errorf("failed to parse distinguised name: malformed element %s in dn", element)
			//log.Error("Failed to parse distinguised name: malformed element %s in dn", element)
		}
		switch strings.ToUpper(value[0]) {
		case "CN": // commonName
			newName.CommonName = value[1]
		case "C": // countryName
			if value[1] == "" {
				caName.Country = []string{}
			} else {
				newName.Country = append(newName.Country, value[1])
			}
		case "L": // localityName
			if value[1] == "" {
				caName.Locality = []string{}
			} else {
				newName.Locality = append(newName.Locality, value[1])
			}
		case "ST": // stateOrProvinceName
			if value[1] == "" {
				caName.Province = []string{}
			} else {
				newName.Province = append(newName.Province, value[1])
			}
		case "O": // organizationName
			if value[1] == "" {
				caName.Organization = []string{}
			} else {
				newName.Organization = append(newName.Organization, value[1])
			}
		case "OU": // organizationalUnitName
			if value[1] == "" {
				caName.OrganizationalUnit = []string{}
			} else {
				newName.OrganizationalUnit = append(newName.OrganizationalUnit, value[1])
			}
		default:
			return nil, fmt.Errorf("Failed to parse distinguised name: unknown element %s", element)
			//log.Error("Failed to parse distinguised name: unknown element %s", element)
		}
	}
	if ca != nil {
		newName.Country = append(caName.Country, newName.Country...)
		newName.Locality = append(caName.Locality, newName.Locality...)
		newName.Province = append(caName.Province, newName.Province...)
		newName.Organization = append(caName.Organization, newName.Organization...)
		newName.OrganizationalUnit = append(caName.OrganizationalUnit, newName.OrganizationalUnit...)
	}
	return newName, nil
}

func copyFile(source string, dest string, perms os.FileMode) {
	sourceFile, err := os.Open(source)
	if err != nil {
		slog.Error("Failed to open %s for reading: %s", source, err)
	}
	defer func() {
		if err = sourceFile.Close(); err != nil {
			slog.Error("Failed to close %s: %s", source, err)
		}
	}()
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, perms)
	if err != nil {
		slog.Error("Failed to open %s for writing: %s", dest, err)
	}
	defer func() {
		if err = destFile.Close(); err != nil {
			slog.Error("Failed to close %s: %s", dest, err)
		}
	}()
	if _, err = io.Copy(destFile, sourceFile); err != nil {
		slog.Error("Failed to copy %s: %s", source, err)
	}
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read file %s: %s", path, err)
	}
	return string(data)
}
