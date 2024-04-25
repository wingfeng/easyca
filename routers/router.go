package routers

import (
	"easyca/authn"
	"easyca/conf"
	"easyca/controller"
	"log/slog"
	"strings"

	api2 "easyca/routers/api"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	r.Use(authn.EnableCookieSession())
	//根据配置切换验证方式
	auth := strings.TrimSpace(strings.ToLower(conf.Default.Authn))
	switch auth {
	case "inner":
		r.Use(gin.BasicAuth(gin.Accounts{
			"caadmin": "pass@word1",
		}))
		r.Use(authn.InnerBasic())
		slog.Info("系统采用inner方式验证用户")
	case "ldap":
		r.Use(authn.LdapBasic())
		slog.Info("系统采用LDAP验证用户")
	case "oidc":
		r.Use(authn.Oidc())
		r.GET(conf.Default.Oidc.Callback, authn.Oidcsignin)
		slog.Info("系统采用OIDC方式验证用户")
	}

	r.GET(conf.Default.CA.CRLURL, controller.GetCRL)
	//vue前端
	RegisterVue(r)

	api := r.Group("/api")
	//后端接口
	api2.RegisterCa(api)
	//权限管理
	api2.RegisterCasBin(api)
	//系统
	api2.RegisterSystem(api)

	return r
}

// func defaultPage(c *gin.Context) {
// 	p := "cert/certs/certs.crt"
// 	cert, err := core.GetCertInfo(p)
// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"type":  "cmd",
// 		"tmpl":  "server",
// 		"cert":  cert,
// 		"path":  p,
// 		"Error": err,
// 	})
// }
