package main

import (
	"os"
	"strings"

	"github.com/dueckminor/mypi-tools/go/protocols/influxdb"
	"github.com/dueckminor/mypi-tools/go/util"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InfluxDB influxdb.Config `yaml:"influxdb"`
}

func main() {
	var cfg Config

	if (len(os.Args) == 3) && !strings.HasPrefix(os.Args[1], "-") && util.FileExists(os.Args[1]) {
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			panic(err)
		}
	}

	// var influx influxdb.Client
	// influx = influxdb.NewClient(cfg.InfluxDB)
	// influx.SendMetric()
}
