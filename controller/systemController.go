package controller

import (
	"easyca/authn"
	"easyca/conf"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Current(context *gin.Context) {
	session := sessions.Default(context)
	user := session.Get("user").(*authn.UserClaims)

	context.JSON(http.StatusOK, user)
}

func Version(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"AppName":      conf.AppName,
		"AppVersion":   conf.AppVersion,
		"BuildVersion": conf.BuildVersion,
		"BuildTime":    conf.BuildTime,
		"GitRevision":  conf.GitRevision,
		"GitBranch":    conf.GitBranch,
		"GoVersion":    conf.GoVersion,
	})
}

// Menu 获取系统菜单
func Menu(c *gin.Context) {

	//ToDO: 加上权限检查
	menu := conf.GetMenu("conf/menu.yaml")

	c.JSON(http.StatusOK, menu)

}
