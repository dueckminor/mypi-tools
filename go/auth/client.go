package auth

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

func (ac *AuthClient) handleAuth(c *gin.Context) {
	scheme := ginutil.GetScheme(c)
	hostname := ginutil.GetHostname(c)

	session := sessions.Default(c)
	accessToken := session.Get("access_token")
	if nil != accessToken {
		if hostname != session.Get("hostname") {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

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

	fmt.Println(c.Request.URL.String())
	fmt.Println(redirectToAuthURI.String())
	fmt.Println("Header:", c.Request.Header)

	c.Header("Location", redirectToAuthURI.String())
	c.AbortWithStatus(http.StatusFound)
}

func (ac *AuthClient) handleLoginCallback(c *gin.Context) {
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
