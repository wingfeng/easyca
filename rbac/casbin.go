package rbac

import (
	"easyca/authn"
	"easyca/conf"
	"fmt"
	"net/http"

	log "log/slog"

	"github.com/casbin/casbin/v2"
	xormadapter "github.com/casbin/xorm-adapter/v2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var (
	globalEnforcer *casbin.Enforcer
	//rLock   *sync.RWMutex
)

func InitEnforcer() {
	adapter, err := xormadapter.NewAdapter(conf.Default.DB.DBDriver, conf.Default.DB.DBConnection, true)
	if err != nil {
		panic(fmt.Sprintf("adapter err is %v", err.Error()))
	}
	//初始化权限配置
	globalEnforcer, err = casbin.NewEnforcer(conf.Default.RbacModelPath, adapter)
	if err != nil {
		panic(fmt.Sprintf("NewEnforcer err is %v", err.Error()))
	}
}

func DefaultEnforcer() *casbin.Enforcer {
	return globalEnforcer
}

// 权限认证中间件
func RbacHandle() gin.HandlerFunc {
	return func(context *gin.Context) {
		//rLock.RLock()
		//defer rLock.RUnlock()

		//加载权限
		_ = globalEnforcer.LoadPolicy()

		//获取请求接口
		reqUrl := context.Request.URL.Path
		//获取请求方法
		//reqAct := context.Request.Method
		session := sessions.Default(context)
		user := session.Get("user").(*authn.UserClaims)
		//将rbac现有用户的角色清除

		globalEnforcer.DeleteRolesForUser(user.Subject)
		for _, role := range user.Role {
			globalEnforcer.AddRoleForUser(user.Subject, role)

		}
		//判断用户是否有权限
		ok, _ := globalEnforcer.Enforce(user.Subject, reqUrl, "0")

		//判断策略中是否存在
		if ok {
			context.Next()
			return

		}

		log.Debug("权限验证不通过")
		context.JSON(http.StatusForbidden, gin.H{
			"messages": http.StatusText(http.StatusForbidden),
			"errInfo":  "抱歉, 您无权访问",
		})
		context.Abort()
	}
}
