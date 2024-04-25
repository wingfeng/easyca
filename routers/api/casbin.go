package api

import (
	"easyca/controller"
	"easyca/rbac"

	"github.com/gin-gonic/gin"
)

//后端api
func RegisterCasBin(r *gin.RouterGroup) {
	api := r.Group("/casbin")
	api.Use(rbac.RbacHandle())
	{
		api.Any("/get_all_user", controller.GetAllUser)

		api.Any("/get_role_by_user", controller.GetRoleByUser)

		api.Any("/modify_role", controller.ModifyRole)

		api.Any("/role_list", controller.RoleList)

		api.Any("/is_enable", controller.IsEnable)
	}
}
