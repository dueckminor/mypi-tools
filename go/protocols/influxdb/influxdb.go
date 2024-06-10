package influxdb

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type Config struct {
	Uri          string `yaml:"uri"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	Bucket       string `yaml:"bucket"`
}

type Client interface {
	Close() error
	Flush()
	SendMetric(name string, value float64)
	SendMetricAtTs(name string, value float64, ts time.Time)
}

type client struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
}

func (cl *client) Close() error {
	cl.client.Close()
	return nil
}

func (cl *client) Flush() {
	cl.writeAPI.Flush(context.Background())
}

func (cl *client) SendMetric(name string, value float64) {
	cl.SendMetricAtTs(name, value, time.Now())
}

func (cl *client) SendMetricAtTs(name string, value float64, ts time.Time) {
	point := influxdb2.NewPointWithMeasurement("Wh")
	point.AddField("value", value)
	point.AddTag("device_class", "energy")
	point.AddTag("domain", "sensor")
	point.AddTag("device", "alphaess")
	point.AddTag("source", "mypi")
	point.AddTag("entity_id", name)
	point.SetTime(ts)
	err := cl.writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		fmt.Println("send metric to influxdb failed:", err)
	}
}

func NewClient(config Config) Client {
	result := &client{}
	result.client = influxdb2.NewClient(config.Uri, config.Token)

	result.writeAPI = result.client.WriteAPIBlocking(config.Organization, config.Bucket)

	return result
}
