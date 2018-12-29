package main

import (
	"fmt"

	"github.com/dueckminor/mypi-api/go/config"
	"github.com/dueckminor/mypi-api/go/debug"
	"github.com/dueckminor/mypi-api/go/webhandler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/tiaguinho/gosoap"
)

type GetGeoIPResponse struct {
	GetGeoIPResult GetGeoIPResult
}

type GetGeoIPResult struct {
	ReturnCode        string
	IP                string
	ReturnCodeDetails string
	CountryName       string
	CountryCode       string
}

var (
	r GetGeoIPResponse
)

func main() {
	{
		soap, err := gosoap.SoapClient("http://fritz.box:49000/tr64desc.xml")
		if err != nil {
			panic(fmt.Errorf("error not expected: %s", err))
		}

		params := gosoap.Params{
			"IPAddress": "8.8.8.8",
		}

		err = soap.Call("GetInfo", params)
		if err != nil {
			panic(fmt.Errorf("error in soap call: %s", err))
		}

		soap.Unmarshal(&r)
		if r.GetGeoIPResult.CountryCode != "USA" {
			panic(fmt.Errorf("error: %+v", r))
		}
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.GetString("config", "root"))

	r := gin.Default()
	webhandler.SetupEndpoints(r)

	r.Use(static.Serve("/", static.LocalFile("dist", false)))

	debug.SetupProxy(r, "http://localhost:8081")

	r.Run()
}
