package routers

import (
	"easyca/conf"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
)

func RegisterVue(r *gin.Engine)  {
	r.LoadHTMLFiles(fmt.Sprintf("%s/index.html",conf.Default.VuePath))
	files, _ := ioutil.ReadDir(conf.Default.VuePath)
	for _, file := range files {
		if path.Ext(file.Name()) == "html" {
			continue
		}
		if file.IsDir(){
			r.StaticFS(fmt.Sprintf("/%s",file.Name()),  http.Dir(fmt.Sprintf("%s/%s",conf.Default.VuePath, file.Name())))
		} else {
			r.StaticFile(fmt.Sprintf("/%s",file.Name()), fmt.Sprintf("%s/%s",conf.Default.VuePath,file.Name()))
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "企业证书管理系统",
		})
	})

	//前端路由注册
	v1 := r.Group("/v1")
	//前端路由层级
	const urllength = 3
	for i := 0; i <= urllength; i ++ {
		var url = "/"
		for j := 1; j <= i; j ++ {
			url = path.Join(url,":vue")
		}
		v1.GET(url, func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "企业证书管理系统",
			})
		})
	}
}
