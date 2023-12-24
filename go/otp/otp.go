package otp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pquerna/otp"

	"image/png"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

// nolint: all
func main() {

	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/qr", func(c *gin.Context) {
		key, _ := totp.Generate(totp.GenerateOpts{
			AccountName: "bar",
			Issuer:      "foo.dueckminor.de",
			SecretSize:  64,
			Algorithm:   otp.AlgorithmSHA512,
		})
		fmt.Println(key.Secret())
		image, _ := key.Image(256, 256)
		buf := new(bytes.Buffer)
		png.Encode(buf, image)
		c.Data(200, "image/png", buf.Bytes())
	})
	r.Run() // listen and serve on 0.0.0.0:8080

	dat, err := ioutil.ReadFile("config/mypi.json")
	fmt.Print(string(dat))

	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

type Resp struct {
	Cockie string
	Error  string
}

// nolint: all
func mainCookie() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		value, err := c.Cookie("mypi")
		errText := ""
		if err != nil {
			errText = err.Error()
		}
		if len(value) == 0 {
			c.SetCookie("mypi", "bar", 3600, "/", "localhost", false, false)
		} else {
			c.SetCookie("mypi", "", 0, "/", "localhost", false, false)
		}
		c.JSON(200, Resp{
			Cockie: value,
			Error:  errText,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
