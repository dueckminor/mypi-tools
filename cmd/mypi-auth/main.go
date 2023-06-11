package main

import (
	"encoding/base64"
	"flag"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/dueckminor/mypi-tools/go/auth"
	"github.com/dueckminor/mypi-tools/go/config"
	"github.com/dueckminor/mypi-tools/go/ginutil"
	"github.com/dueckminor/mypi-tools/go/pki"
	"github.com/dueckminor/mypi-tools/go/rand"
	"github.com/dueckminor/mypi-tools/go/restapi"
	"github.com/dueckminor/mypi-tools/go/users"

	"github.com/golang-jwt/jwt"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"

	// provide only the http rest api
	_ "github.com/dueckminor/mypi-tools/go/restapi/http"
)

var (
	mypiRoot = flag.String("mypi-root", "", "The root of the mypi filesystem")
	userCfg  *users.UserCfg
)

func init() {
	flag.Parse()
	if mypiRoot != nil && len(*mypiRoot) > 0 {
		config.InitApp(*mypiRoot)
	}

	var err error
	userCfg, err = users.ReadUserCfg()
	if err != nil {
		panic(err)
	}
	config.GetConfig()
}

func login(c *gin.Context) {
	var params struct {
		Username string
		Password string
	}
	err := c.BindJSON(&params)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !users.CheckPassword(params.Username, params.Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	session := sessions.Default(c)

	secret, _ := session.Get("secret").(string)
	//domain, _ := session.Get("domain").(string)

	if len(secret) == 0 {
		domain := ""
		host, _, _ := net.SplitHostPort(c.Request.Host)
		addr := net.ParseIP(host)
		if addr != nil {
			origin := c.Request.Header["Origin"]
			if len(origin) == 1 {
				uri, _ := url.Parse(origin[0])
				host = uri.Hostname()
				addr = net.ParseIP(host)
			}
		}

		if addr == nil {
			hostParts := strings.Split(host, ".")
			if len(hostParts) > 1 {
				domain = strings.Join(hostParts[1:], ".")
			}
		}

		secret, _ = rand.GetString(48)
		session.Set("secret", secret)
		session.Set("domain", domain)
		session.Set("username", params.Username)
		session.Save()
	}
	c.Data(http.StatusOK, "text/plain", []byte("OK"))
}

type ClaimsWithScope struct {
	Scopes []string `json:"scopes,omitempty"`
	jwt.StandardClaims
}

func handleOauthAuthorize(c *gin.Context) {
	session := sessions.Default(c)
	secret, _ := session.Get("secret").(string)

	if len(secret) > 0 {
		authRequest, err := auth.NewRequest()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		values := c.Request.URL.Query()
		authRequest.RedirectURI = values.Get("redirect_uri")
		redirectURI, err := url.Parse(authRequest.RedirectURI)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		values = redirectURI.Query()
		values.Add("code", authRequest.Id)
		redirectURI.RawQuery = values.Encode()
		c.Header("Location", redirectURI.String())
		c.AbortWithStatus(http.StatusFound)
	}

	c.Request.URL.Path = "/"
	c.Header("Location", c.Request.URL.String())
	c.AbortWithStatus(http.StatusFound)
}

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func basicAuth(c *gin.Context) string {
	s := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return ""
	}
	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return ""
	}
	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return ""
	}

	clientID := pair[0]
	if len(clientID) == 0 {
		return ""
	}
	if strings.Contains(clientID, "..") || strings.ContainsAny(clientID, "/\\") {
		return ""
	}

	clientConfig, err := config.ReadConfigFile(path.Join("etc/auth/clients", clientID+".yml"))
	if err != nil {
		return ""
	}
	clientSecret := clientConfig.GetString("client_secret")
	if clientSecret != pair[1] {
		return ""
	}

	return clientID
}

func handleOauthToken(c *gin.Context) {
	c.Request.ParseForm()
	code := c.Request.Form.Get("code")
	grantType := c.Request.Form.Get("grant_type")
	responseType := c.Request.Form.Get("response_type")
	redirectURI := c.Request.Form.Get("redirect_uri")
	clientID := c.Request.Form.Get("client_id")

	if grantType != "authorization_code" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if responseType != "token" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	username := basicAuth(c)
	if username != clientID {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	authRequest := auth.GetRequest(code)
	if authRequest.RedirectURI != redirectURI {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, ClaimsWithScope{})
	key, _ := config.ReadRSAPrivateKey("etc/auth/server/server_priv.pem")
	jwt, _ := token.SignedString(key)

	response := OauthTokenResponse{
		AccessToken: jwt,
	}
	c.AbortWithStatusJSON(http.StatusOK, response)
}

func handleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.AbortWithStatus(http.StatusAccepted)
}

type status struct {
	Username string `json:"username"`
}

func handleStatus(c *gin.Context) {
	session := sessions.Default(c)

	username, _ := session.Get("username").(string)

	c.AbortWithStatusJSON(http.StatusOK, status{
		Username: username,
	})
}

func GenerateRsaKeyPair(privFilename, pubFilename string) error {
	priv, pub := pki.GenerateRsaKeyPair()
	privPEM := pki.RsaPrivateKeyToPem(priv)
	pubPEM := pki.RsaPublicKeyToPem(pub)

	privFilename = config.GetFilename(privFilename)
	pubFilename = config.GetFilename(pubFilename)

	print(path.Dir(privFilename))
	os.MkdirAll(path.Dir(privFilename), 0755)
	os.MkdirAll(path.Dir(pubFilename), 0755)

	err := ioutil.WriteFile(privFilename, privPEM, 0600)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pubFilename, pubPEM, 0644)
}

func main() {
	if len(flag.Args()) > 0 {
		if flag.Args()[0] == "init" {
			err := GenerateRsaKeyPair("etc/auth/server/server_priv.pem", "etc/auth/server/server_pub.pem")
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
		if flag.Args()[0] == "create-client" {
			if len(flag.Args()) != 2 {
				panic("create-client needs exactly one arg")
			}
			os.MkdirAll(config.GetFilename("etc/auth/clients"), 0755)

			clientID := flag.Args()[1]
			pub, err := config.FileToString("etc/auth/server/server_pub.pem")
			if err != nil {
				panic(err)
			}
			clientConfig, err := config.GetOrCreateConfigFile(path.Join("etc/auth/clients", clientID+".yml"))
			clientConfig.SetString("server_key", pub)
			clientConfig.SetString("client_id", clientID)
			if len(clientConfig.GetString("client_secret")) == 0 {
				clientSecret, err := rand.GetString(32)
				if err != nil {
					panic(err)
				}
				clientConfig.SetString("client_secret", clientSecret)
			}

			err = clientConfig.Write()
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	}

	r := gin.Default()

	cfg, err := config.GetOrCreateConfigFile("mypi-auth.yml")
	if err != nil {
		panic(err)
	}

	err = ginutil.ConfigureSessionCookies(r, cfg, "MYPI_AUTH_SESSION")
	if err != nil {
		panic(err)
	}

	r.Use(cors.Default())

	r.POST("/login", login)
	r.POST("/logout", handleLogout)
	r.GET("/status", handleStatus)
	r.GET("/oauth/authorize", handleOauthAuthorize)
	r.POST("/oauth/token", handleOauthToken)

	restapi.Run(r)
}
