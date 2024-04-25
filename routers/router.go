package routers

import (
	"easyca/authn"
	"easyca/conf"
	"easyca/controller"

	api2 "easyca/routers/api"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	r.Use(authn.EnableCookieSession())
	if conf.Default.UseOIDC {
		r.Use(authn.Oidc())
		r.GET(conf.Default.Oidc.Callback, authn.Oidcsignin)
	}
	//
	r.Use(authn.LdapBasic())
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
