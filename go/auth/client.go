package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthClient struct {
	AuthURI      string
	ClientID     string
	ClientSecret string
	ServerKey    *rsa.PublicKey
}

func (ac *AuthClient) RegisterHandler(e *gin.Engine) {
	e.GET("/login/callback", ac.handleLoginCallback)
	e.Use(ac.handleAuth)
}

func (ac *AuthClient) RegisterCallbackHandler(e *gin.Engine) {
	e.GET("/login/callback", ac.handleLoginCallback)
}

func (ac *AuthClient) GetHandler() gin.HandlerFunc {
	return ac.handleAuth
}
func (ac *AuthClient) GetHandlerForAPI() gin.HandlerFunc {
	return ac.handleAuthForApi
}

func (ac *AuthClient) GetHandlerIntegratedLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		prefixWithoutAuth := []string{
			"/js",
			"/fonts",
			"/login",
			"/logout",
			"/ws",
		}
		for _, prefix := range prefixWithoutAuth {
			if path == prefix ||
				strings.HasPrefix(path, prefix+"/") ||
				strings.HasPrefix(path, prefix+"?") {
				c.Next()
				return
			}
		}
		ac.handleAuth(c)
	}
}

func (ac *AuthClient) verifySession(c *gin.Context) bool {
	hostname := ginutil.GetHostname(c)

	session := sessions.Default(c)

	accessToken := session.Get("access_token")
	if nil != accessToken {
		if hostname != session.Get("hostname") {
			c.AbortWithStatus(http.StatusInternalServerError)
			return false
		}
		return true
	}
	return false
}

func (ac *AuthClient) handleAuthForApi(c *gin.Context) {
	ok := ac.verifySession(c)
	if c.IsAborted() {
		return
	}
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func (ac *AuthClient) handleAuth(c *gin.Context) {
	if ac.verifySession(c) || c.IsAborted() {
		return
	}

	scheme := ginutil.GetScheme(c)
	hostname := ginutil.GetHostname(c)

	authRequest, err := NewRequest()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	authRequest.Path = c.Request.URL.Path

	callbackURI := &url.URL{
		Scheme: scheme,
		Host:   hostname,
		Path:   "/login/callback",
	}
	values := url.Values{}
	values.Add("id", authRequest.Id)
	callbackURI.RawQuery = values.Encode()

	authRequest.RedirectURI = callbackURI.String()

	redirectToAuthURI, _ := url.Parse(ac.AuthURI)
	values = redirectToAuthURI.Query()
	values.Add("redirect_uri", authRequest.RedirectURI)
	values.Add("response_type", "code")
	values.Add("client_id", ac.ClientID)
	redirectToAuthURI.Path = "/oauth/authorize"
	redirectToAuthURI.RawQuery = values.Encode()

	c.Header("Location", redirectToAuthURI.String())
	c.AbortWithStatus(http.StatusFound)
}

func (ac *AuthClient) handleLoginCallback(c *gin.Context) {
	fmt.Println("login callback 2")
	session := sessions.Default(c)
	c.Request.ParseForm()
	code := c.Request.Form.Get("code")
	id := c.Request.Form.Get("id")

	authRequest := GetRequest(id)
	if nil == authRequest {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	path := authRequest.Path

	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	v.Set("response_type", "token")
	v.Set("redirect_uri", authRequest.RedirectURI)
	v.Set("client_id", ac.ClientID)

	//pass the values to the request's body

	authURIOauthToken, _ := url.Parse(ac.AuthURI)
	authURIOauthToken.Path = "oauth/token"

	req, err := http.NewRequest("POST", authURIOauthToken.String(), strings.NewReader(v.Encode()))
	req.SetBasicAuth(ac.ClientID, ac.ClientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	bodyText, err := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	json.Unmarshal(bodyText, &data)

	jwtToken, ok := data["access_token"].(string)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			return ac.ServerKey, nil
		}
		return nil, fmt.Errorf("Unexpected Signing Method")
	})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	fmt.Println(token)

	session.Set("access_token", data["access_token"])
	session.Set("hostname", ginutil.GetHostname(c))
	session.Save()

	c.Header("Location", path)
	c.AbortWithStatus(http.StatusFound)
}

///////////////////////////////////////////////////////////////////////////////

type AuthClientLocalSecret struct {
	LocalSecret string
}

func (ac *AuthClientLocalSecret) CreateLocalSecret() string {
	buf := make([]byte, 20)
	rand.Read(buf)
	ac.LocalSecret = base32.StdEncoding.EncodeToString(buf)
	return ac.LocalSecret
}

func (ac *AuthClientLocalSecret) GetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		q := c.Request.URL.Query()
		secret := q.Get("local_secret")
		if secret == ac.LocalSecret {
			session.Set("local_secret", true)
			session.Save()

			q.Del("local_secret")
			c.Request.URL.RawQuery = q.Encode()
			redirect := c.Request.URL.String()
			c.Header("Location", redirect)
			c.AbortWithStatus(http.StatusFound)
			return
		}

		if true != session.Get("local_secret") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}
