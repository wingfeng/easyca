package controller

import (
	"easyca/model"
	"easyca/rbac"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllUser(context *gin.Context) {
	var results []model.ModifyRole

	users := rbac.GetAllSubjects()

	for _, user := range users {
		results = append(results, getRolesByUser(user, true))
	}

	context.JSON(http.StatusOK, results)
}

func GetRoleByUser(context *gin.Context) {
	var param model.User

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, getRolesByUser(param.User))
}

func ModifyRole(context *gin.Context) {
	var param model.ModifyRole

	if err := context.ShouldBindJSON(&param); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  err.Error(),
		})
		return
	}

	//默认当前角色
	if param.User == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"messages": http.StatusText(http.StatusBadRequest),
			"errInfo":  "user is not null",
		})
		return
	}

	if param.Del == "" {
		param.Del = "0"
	}

	err := rbac.ModifyPermissionsForUser(&param)
	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"messages": http.StatusText(http.StatusOK),
		})
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}
}

type roles struct {
	Url   string
	Label string
}

var RRoles = []roles{
	{
		Url:   "/api/casbin/get_all_user",
		Label: "获取所有角色",
	},
	{
		Url:   "/api/casbin/get_role_by_user",
		Label: "获取角色对应权限",
	},
	{
		Url:   "/api/casbin/modify_role",
		Label: "编辑角色对应权限",
	},
	{
		Url:   "/api/casbin/role_list",
		Label: "获取路由列表",
	},
	{
		Url:   "/api/ca/get_cert_info",
		Label: "获取证书信息",
	},
	{
		Url:   "/api/ca/create_ca",
		Label: "创建根证书",
	},
	{
		Url:   "/api/ca/get_base_cert",
		Label: "获取根证书",
	},
	{
		Url:   "/api/ca/download_base_cert",
		Label: "下载根证书",
	},
	{
		Url:   "/api/ca/create_certificate_personal",
		Label: "创建个人证书",
	},
	{
		Url:   "/api/ca/get_certificate_personal",
		Label: "获取个人证书",
	},
	{
		Url:   "/api/ca/download_certificate_personal",
		Label: "下载个人证书",
	},
	{
		Url:   "/api/ca/create_certificate_server",
		Label: "创建服务器证书",
	},
	{
		Url:   "/api/ca/get_certificate_server",
		Label: "获取服务器证书",
	},
	{
		Url:   "/api/ca/download_certificate_server",
		Label: "下载服务器证书",
	},
	{
		Url:   "/api/ca/create_certificate_sign",
		Label: "创建签名证书",
	},
	{
		Url:   "/api/ca/get_sign_server",
		Label: "获取签名证书",
	},
	{
		Url:   "/api/ca/download_certificate_sign",
		Label: "下载签名证书",
	},
	{
		Url:   "/api/ca/modify_role",
		Label: "编辑角色对应权限",
	},
	{
		Url:   "/api/ca/server_list",
		Label: "服务器证书列表",
	},
	{
		Url:   "/api/ca/sign_list",
		Label: "签名证书列表",
	},
}

var rolesMap = make(map[string]string)

func init() {
	for _, item := range RRoles {
		rolesMap[item.Url] = item.Label
	}
}

func RoleList(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"list": RRoles,
	})
}

func IsEnable(context *gin.Context) {
	user := context.Query("user")
	del := context.DefaultQuery("del", "0")

	roles := getRolesByUser(user)
	roles.Del = del

	err := rbac.ModifyPermissionsForUser(&roles)
	if err == nil {
		context.JSON(http.StatusOK, gin.H{
			"messages": http.StatusText(http.StatusOK),
		})
	} else {
		context.JSON(http.StatusInternalServerError, gin.H{
			"messages": http.StatusText(http.StatusInternalServerError),
			"errInfo":  err.Error(),
		})
	}
}

func getRolesByUser(user string, islabels ...bool) (result model.ModifyRole) {
	var islabel bool
	if len(islabels) != 0 {
		islabel = islabels[0]
	}

	result.User = user

	var userRole = rbac.GetPermissionsForUser(user)

	for _, item := range userRole {
		if !islabel {
			result.Roles = append(result.Roles, item[1])
		} else {
			result.Roles = append(result.Roles, rolesMap[item[1]])
		}
	}

	if len(userRole) != 0 {
		if len(userRole[0]) > 2 {
			result.Del = userRole[0][2]
		}
	}

	if result.Del == "" {
		result.Del = "0"
	}

	return
}
