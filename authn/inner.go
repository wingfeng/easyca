package authn

import (
	"easyca/conf"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InnerBasic() gin.HandlerFunc {

	return func(c *gin.Context) {
		gin.BasicAuth(gin.Accounts{
			"caadmin": "pass@word1",
		})
		user := c.MustGet(gin.AuthUserKey).(string)
		userInfo := &UserClaims{
			Name:           "CA Admin",
			Subject:        "caadmin",
			Email:          "caadmin@easyca.local",
			PreferUserName: []string{user},
			Role:           []string{conf.Default.DefaultRoleName},
		}
		sessions.Default(c).Set("user", userInfo)
		sessions.Default(c).Save()
		c.Next()
	}
}
