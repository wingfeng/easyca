package api

import (
	"easyca/controller"
	"easyca/rbac"

	"github.com/gin-gonic/gin"
)

//后端api
func RegisterCa(r *gin.RouterGroup) {
	api := r.Group("/ca")
	api.Use(rbac.RbacHandle())
	{
		api.GET("/get_cert_info", controller.GetCertInfo)

		api.POST("/create_ca", controller.CreateCA)

		api.Any("/create_certificate_personal", controller.CreateCertificatePersonal)

		api.Any("/get_certificate_personal", controller.GetCertificatePersonal)

		api.GET("/download_certificate_personal", controller.DownloadCertificatePersonal)

		api.POST("/create_certificate_server", controller.CreateCertificateServer)

		api.POST("/get_certificate_server", controller.GetCertificateServer)

		api.POST("/download_certificate_server", controller.DownloadCertificateServer)

		api.POST("/create_certificate_sign", controller.CreateCertificateSign)

		api.POST("/get_certificate_sign", controller.GetCertificateSign)

		api.POST("/download_certificate_sign", controller.DownloadCertificateSign)

		api.GET("/get_base_cert", controller.GetBaseCert)

		api.GET("/download_base_cert", controller.DownloadBaseCert)

		api.Any("/server_list", controller.ServerList)

		api.Any("/sign_list", controller.SignList)
	}
}
