package authn

import (
	"context"
	"crypto/sha1"
	"easyca/conf"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "log/slog"

	oidc "github.com/coreos/go-oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"golang.org/x/oauth2"
)

type UserClaims struct {
	Subject        string
	PreferUserName []string `json:"preferred_username"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	EmailVerified  bool     `json:"email_verified"`
	Role           []string `json:"role"`
}

// gin session key
const KEY = "easyca"

// 使用 Cookie 保存 session
func EnableCookieSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte(KEY))
	return sessions.Sessions("Default", store)
}

func Oidc() gin.HandlerFunc {

	var (
		clientID     = conf.Default.Oidc.ClientID
		clientSecret = conf.Default.Oidc.Secret
		ctx          = context.Background()
		oidcConfig   = &oidc.Config{
			ClientID: clientID,
		}
	)
	provider, err := oidc.NewProvider(ctx, conf.Default.Oidc.Issuer)
	if err != nil {
		log.Error("初始化OIDC Proiver错误,", "error", err)
	}
	verifier := provider.Verifier(oidcConfig)
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  conf.Default.Oidc.Callback,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "roles"},
	}

	state := fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v", time.Now()))))

	return func(context *gin.Context) {
		//设置当前服务器运行的主机信息，以便后面在给CA发证书是，以方便加上CRL的地址
		if strings.EqualFold("", conf.CurrentHost) {
			conf.CurrentHost = fmt.Sprintf("%s://%s", conf.Default.Oidc.CallbackProto, context.Request.Host)
		}

		//吊销证书列表不需要验证
		if context.Request.URL.Path == conf.Default.CA.CRLURL {
			return
		}
		session := sessions.Default(context)
		sessionValue := session.Get("userId")
		if sessionValue == nil && context.Request.URL.Path != "/signin-oidc" {

			config.RedirectURL = fmt.Sprintf("%s://%s%s", conf.Default.Oidc.CallbackProto, context.Request.Host, conf.Default.Oidc.Callback)
			context.Redirect(http.StatusFound, config.AuthCodeURL(state))
			return
		}
		//处理OpenID的回调
		if context.Request.URL.Path == "/signin-oidc" {

			context.Set("verifier", verifier)
			context.Set("state", state)
			context.Set("ctx", ctx)
			context.Set("OauthConfig", config)
			context.Set("oAuthProvider", provider)
		}

		context.Next()
	}
}

func Oidcsignin(c *gin.Context) {
	v, _ := c.Get("verifier")
	verifier := v.(*oidc.IDTokenVerifier)
	state := c.GetString("state")
	ct, _ := c.Get("ctx")
	ctx := ct.(context.Context)
	cf, _ := c.Get("OauthConfig")
	config := cf.(oauth2.Config)
	p, _ := c.Get("oAuthProvider")
	provider := p.(*oidc.Provider)
	if c.Request.URL.Query().Get("state") != state {
		log.Error("state did not match:", "state", state, " Query:", c.Request.URL.Query().Get("state"))
		return
	}

	oauth2Token, err := config.Exchange(ctx, c.Request.URL.Query().Get("code"))
	if err != nil {
		log.Error("Failed to exchange token: "+err.Error(), "error", http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Error("No id_token field in oauth2 token.")
		return
	}
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Error("Failed to verify ID Token: ", err.Error(), http.StatusInternalServerError)
		return
	}
	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))

	if err != nil {
		log.Error("Failed to get userinfo: " + err.Error())
		return
	}
	oauth2Token.AccessToken = "*REDACTED*" //隐藏Token信息
	claims := new(UserClaims)
	userInfo.Claims(claims)
	resp := struct {
		OAuth2Token   *oauth2.Token
		IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		UserInfo      *oidc.UserInfo
		UserClaims    *UserClaims
	}{oauth2Token, new(json.RawMessage), userInfo, claims}

	if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
		log.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		log.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	//	_ = data //use next timg

	log.Info(string(data))
	session := sessions.Default(c)
	claims.Subject = userInfo.Subject
	session.Set("user", claims)

	err = session.Save()
	if err != nil {
		log.Error("保存Session 信息错误", err)
	}
	token := fmt.Sprintf("Bearer %s", rawIDToken)
	//	c.SetCookie("Authorization", token, 600, "/", "", false, true)
	c.Writer.Header().Add("Authorization", token)
	c.Writer.WriteString("<html><script> window.location='/';</script></html>")
	c.Writer.Flush()
	//	c.AbortWithStatus(http.StatusFound)
}
