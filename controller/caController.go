package controller

import (
	"bytes"
	"crypto/sha1"
	"crypto/x509"
	"easyca/authn"
	"easyca/conf"
	"easyca/core"
	"easyca/engine"
	"easyca/model"
	"easyca/utils"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GetCertInfo(context *gin.Context) {
	var param model.GetCertInfoParam

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	res, err := core.GetCertInfo(param.Path)
	if err == nil {
		context.JSON(http.StatusOK, res)
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}
}
func getUserInfo(c *gin.Context) *authn.UserClaims {
	userInfo := sessions.Default(c).Get("user").(*authn.UserClaims)
	return userInfo
}
func CreateCA(context *gin.Context) {
	var param model.CreateCAParam

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	//创建正式参数
	var createCAArgs = []model.CreateCAArgs{
		{
			Path:          fmt.Sprintf("%s", conf.Default.CertPath),
			Dn:            conf.Default.CA.DefaultDN,
			MaxPathLength: 1,
			Validity:      10,
			Overwrite:     true,
		},
		//不再创建中间证书因为创建中间签名证书的话，导入根证书的动作要做两次，一次导入到根证书，一次导入中间证书
		// {
		// 	Path:          fmt.Sprintf("%s/server", conf.Default.CertPath),
		// 	Dn:            "/CN=%s Server intermediate CA",
		// 	MaxPathLength: 0,
		// 	Validity:      5,
		// 	Overwrite:     true,
		// },
		// {
		// 	Path:          fmt.Sprintf("%s/personal", conf.Default.CertPath),
		// 	Dn:            "/CN=%s Peronal intermediate CA",
		// 	MaxPathLength: 0,
		// 	Validity:      5,
		// 	Overwrite:     true,
		// },
		// {
		// 	Path:          fmt.Sprintf("%s/sign", conf.Default.CertPath),
		// 	Dn:            "/CN=%s Sign Only intermediate CA",
		// 	MaxPathLength: 0,
		// 	Validity:      5,
		// 	Overwrite:     true,
		// },
	}

	for _, arg := range createCAArgs {
		err := core.CreateCA(arg.Path, fmt.Sprintf(arg.Dn, param.Dn), arg.MaxPathLength, arg.Validity*param.Validity, arg.Overwrite)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"messages": http.StatusText(http.StatusInternalServerError),
				"errInfo":  err.Error(),
			})
			return
		}
	}

	GetBaseCert(context)
}

// 获取根证书
func GetBaseCert(context *gin.Context) {
	var BaseCertPath = fmt.Sprintf("%s/ca.crt", conf.Default.CertPath)

	res, err := core.GetCertInfo(BaseCertPath)
	if err == nil {
		context.JSON(http.StatusOK, res)
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  "请创建根证书",
		})
	}
}

func DownloadBaseCert(context *gin.Context) {
	var srcPath = fmt.Sprintf("%s/%s", conf.Default.CertPath, "ca.crt")

	context.FileAttachment(srcPath, url.QueryEscape("根证书.crt"))
}

// 创建个人证书
func CreateCertificatePersonal(context *gin.Context) {
	userInfo := getUserInfo(context)
	userId := userInfo.Subject
	Email := userInfo.Email
	slog.Info("Create personal Cert", "user", userId)
	//创建个人证书
	var arg = model.CreateCertificateArgs{
		Path:        fmt.Sprintf("%s/personal/%s", conf.Default.CertPath, userId),
		Dn:          fmt.Sprintf("/CN=%s", Email),
		San:         "",
		Validity:    365,
		Overwrite:   true,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageEmailProtection},
	}

	if err := core.CreateCertificate(arg.Path, arg.Dn, arg.San, arg.Validity, arg.Overwrite, arg.KeyUsage, arg.ExtKeyUsage); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
		return
	}

	GetCertificatePersonal(context)
}

// 获取个人证书
func GetCertificatePersonal(context *gin.Context) {
	userId := getUserInfo(context).Subject

	var PersonalCertPath = fmt.Sprintf("%s/personal/%s/%s.crt", conf.Default.CertPath, userId, userId)

	res, err := core.GetCertInfo(PersonalCertPath)
	if err == nil {
		context.JSON(http.StatusOK, res)
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  "请创建个人证书",
		})
	}
}

func DownloadCertificatePersonal(c *gin.Context) {
	userId := getUserInfo(c).Subject

	var srcPath = path.Join(conf.Default.CertPath, "personal", userId) // fmt.Sprintf("%s/personal/%s", conf.Default.CertPath, userId)

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape("个人证书.zip")))
	c.Header("Content-Type", "application/octet-stream")

	utils.Zip(srcPath, c.Writer)

}

// 创建服务器证书
func CreateCertificateServer(context *gin.Context) {
	var param model.CreateCAParam

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	var newServerCa model.ServerCa
	hash := sha1.New()
	//hash a Normalized(lowercase and trime space) server name
	hash.Write([]byte(strings.TrimSpace(strings.ToLower(param.Dn))))
	hashcode := fmt.Sprintf("%x", hash.Sum(nil))
	slog.Info("Create server cert", "server", param.Dn, "hash", hashcode)

	newServerCa.Id = hashcode
	newServerCa.UserId = getUserInfo(context).Subject
	newServerCa.Dn = param.Dn

	//创建服务器证书
	var arg = model.CreateCertificateArgs{
		Path:        fmt.Sprintf("%s/server/%s", conf.Default.CertPath, newServerCa.Id),
		Dn:          fmt.Sprintf("/CN=%s", param.Dn),
		San:         param.Dn, //证书的域名和DN一致
		Validity:    365,
		Overwrite:   true,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	if err := core.CreateCertificate(arg.Path, arg.Dn, arg.San, arg.Validity, arg.Overwrite, arg.KeyUsage, arg.ExtKeyUsage); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
		return
	}
	engine.Db.Find(&newServerCa)
	//插入数据
	if _, err := engine.Db.InsertOne(&newServerCa); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
		return
	}

	//返回请求请求文件
	var reqInfo = model.Id{
		Id: newServerCa.Id,
	}
	reqByte, _ := json.Marshal(reqInfo)
	context.Request.Body = io.NopCloser(bytes.NewReader(reqByte))

	GetCertificateServer(context)
}

// 获取服务器证书
func GetCertificateServer(context *gin.Context) {
	var param model.Id

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	var ServerCertPath = fmt.Sprintf("%s/server/%s/%s.crt", conf.Default.CertPath, param.Id, param.Id)

	res, err := core.GetCertInfo(ServerCertPath)
	if err == nil {
		context.JSON(http.StatusOK, res)
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}
}

func DownloadCertificateServer(context *gin.Context) {
	var param model.Id

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	var srcPath = path.Join(conf.Default.CertPath, "server", param.Id) // fmt.Sprintf("%s/server/%s", conf.Default.CertPath, param.Id)

	context.Header("Content-Disposition", "attachment;filename=server_cert.zip")
	context.Header("Content-Type", "application/octet-stream")

	utils.Zip(srcPath, context.Writer)
	context.Writer.Flush()
}

// 创建签名证书
func CreateCertificateSign(context *gin.Context) {
	var param model.CreateCAParam

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	var newSignCa model.SignCa

	newSignCa.Id = uuid.NewV4().String()
	newSignCa.UserId = getUserInfo(context).Subject
	newSignCa.Dn = param.Dn

	//创建签名证书
	var arg = model.CreateCertificateArgs{
		Path:      fmt.Sprintf("%s/sign/%s", conf.Default.CertPath, newSignCa.Id),
		Dn:        fmt.Sprintf("/CN=%s", param.Dn),
		San:       "",
		Validity:  365,
		Overwrite: true,
		KeyUsage:  x509.KeyUsageDigitalSignature,
	}

	if err := core.CreateCertificate(arg.Path, arg.Dn, arg.San, arg.Validity, arg.Overwrite, arg.KeyUsage, arg.ExtKeyUsage); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
		return
	}

	//插入数据
	if _, err := engine.Db.InsertOne(&newSignCa); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
		return
	}

	//返回请求请求文件
	var reqInfo = model.Id{
		Id: newSignCa.Id,
	}
	reqByte, _ := json.Marshal(reqInfo)
	context.Request.Body = ioutil.NopCloser(bytes.NewReader(reqByte))

	GetCertificateSign(context)
}

// 获取签名证书
func GetCertificateSign(context *gin.Context) {
	var param model.Id

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	var SignCertPath = fmt.Sprintf("%s/sign/%s/%s.crt", conf.Default.CertPath, param.Id, param.Id)

	res, err := core.GetCertInfo(SignCertPath)
	if err == nil {
		context.JSON(http.StatusOK, res)
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}
}

func DownloadCertificateSign(context *gin.Context) {
	var param model.Id

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	var srcPath = path.Join(conf.Default.CertPath, "sign", param.Id) // fmt.Sprintf("%s/sign/%s", conf.Default.CertPath, param.Id)

	context.Header("Content-Disposition", fmt.Sprintf("attachment;filename=\"%s\"", url.QueryEscape("签名证书.zip")))
	context.Header("Content-Type", "application/octet-stream")

	utils.Zip(srcPath, context.Writer)
}

func ServerList(context *gin.Context) {
	userId := getUserInfo(context).Subject

	page, rows := engine.PageNum(context)

	var err error
	var result []model.ServerCa
	var total int64

	e := engine.Db.Where("user_id = ?", userId)

	err = e.Limit(rows, (page-1)*rows).Find(&result)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}

	total, err = e.Count(new(model.ServerCa))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}

	context.JSON(http.StatusOK, gin.H{
		"list":  result,
		"total": total,
	})
}

func SignList(context *gin.Context) {
	userId := getUserInfo(context).Subject

	page, rows := engine.PageNum(context)

	var err error
	var result []model.SignCa
	var total int64
	e := engine.Db.Where("user_id = ?", userId)

	err = e.Limit(rows, (page-1)*rows).Find(&result)
	if err != nil {
		slog.Error("Sign List error", "error", err)
	}
	total, err = e.Count(new(model.SignCa))

	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"list":  result,
			"total": total,
		})
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}
}

func GetCRL(context *gin.Context) {
	crl := core.CreateCRL()
	context.Header("content-type", "application/pkix-crl")
	context.Writer.Write(crl)
}
