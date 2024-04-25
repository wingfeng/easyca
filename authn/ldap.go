package authn

import (
	"easyca/conf"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
)

func LdapBasic() gin.HandlerFunc {

	return func(c *gin.Context) {
		slog.Info("Begin LDAP Basic Authentication")
		authHeader := c.Request.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			slog.Warn("Authorization Header is empty!")
			returnUnAuth(c)
			return
		}

		authHeader = strings.TrimPrefix(authHeader, "Basic ")
		decodedHeader, err := base64.StdEncoding.DecodeString(authHeader)

		if err != nil {
			slog.Error("Decode header error", "error", err)
			returnUnAuth(c)
			return

		}
		user := strings.Split(string(decodedHeader), ":")[0]
		password := strings.Split(string(decodedHeader), ":")[1]
		entry := validateUser(user, password)
		if entry == nil {
			returnUnAuth(c)
			return
		}
		c.Set(gin.AuthUserKey, user)
		slog.Info("user logined", "user", user)
		session := sessions.Default(c)
		groups := entry.GetAttributeValue("memberOf")
		slog.Info("MemberOf", "Group", groups)
		userInfo := &UserClaims{
			Name:           entry.GetAttributeValue("cn"),
			Subject:        entry.GetAttributeValue("uid"),
			PreferUserName: entry.GetAttributeValues("uid"),
			Email:          entry.GetAttributeValue("mail"),
			Role:           []string{conf.Default.DefaultRoleName},
		}
		session.Set("user", userInfo)
		session.Save()
		c.Next()
	}
}
func returnUnAuth(c *gin.Context) {
	realm := "Authorization Required"
	realm = "Basic realm=" + strconv.Quote(realm)
	c.Header("WWW-Authenticate", realm)
	c.AbortWithStatus(http.StatusUnauthorized)
}

func validateUser(username string, password string) *ldap.Entry {
	opt := conf.Default
	l, err := ldap.DialURL(opt.LDAP.URL)
	l.Bind(opt.LDAP.BindDN, opt.LDAP.BindPassword)

	if err != nil {

		slog.Error(err.Error())
		return nil
	}
	defer l.Close()
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		opt.LDAP.BaseDN,
		ldap.ScopeChildren, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", ldap.EscapeFilter(username)),
		[]string{"dn", "cn", "uid", "mail", "displayName", "givenName", "sn", "memberOf"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	if len(sr.Entries) != 1 {
		slog.Error("User does not exist or too many entries returned")
		return nil
	}
	entry := sr.Entries[0]
	userdn := entry.DN

	// Bind as the user to verify their password
	err = l.Bind(userdn, password)
	if err != nil {
		slog.Error("verify user password error", "error", err.Error())
		return nil
	}
	// entry.PrettyPrint(3)
	slog.Info("User login:", "user", entry)
	return entry
}
