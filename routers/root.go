package routers

import (
	"easyca/core"
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

var certPath string

// 查看证书
func ViewCert(c *gin.Context) {
	p := c.Query("path")
	cert, err := core.GetCertInfo(p)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"type":  "cmd",
		"tmpl":  "server",
		"cert":  cert,
		"path":  p,
		"Error": err,
	})
}

func CreateRootCert(c *gin.Context) {
	t := c.DefaultQuery("type", "root")
	c.HTML(http.StatusOK, "newca.html", gin.H{
		"type": t,
	})
}

func CreateServerCert(c *gin.Context) {

}
func CreateSignCert(c *gin.Context) {

}

func DownloadCert(c *gin.Context) {
	p := c.Query("path")
	n := path.Base(p)
	s := fmt.Sprint("attachment; filename=\"\"", n)
	c.Writer.Header().Add("Content-Disposition", s)
	c.File(p)

}
