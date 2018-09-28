package common

import (
	"log"
	"strings"
	"time"

	"check-services-health/res"
	"github.com/influxdata/influxdb/client/v2"
	"reflect"
)

var (
	DbUrl     = "http://127.0.0.1:8086"
	DbName    = "influxdb"
	Username  = "influxdb"
	Password  = "influxdb"
	Precision = "s"
)

var (
	appSrv res.AppSrv
)

// func SaveToInfluxDb(healthData *res.LogSrv) {
func SaveToInfluxDb(healthData interface{}) {

	dataType := reflect.TypeOf(healthData)
	switch dataType {
	case reflect.TypeOf(&appSrv):
		data := healthData.(*res.AppSrv)
		// Create a new HTTPClient
		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     DbUrl,
			Username: Username,
			Password: Password,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		// Create a new point batch
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  DbName,
			Precision: Precision,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a point and add to batch
		// todo: 获取当前健康检查的环境.dev
		appName := "monitor"

		appdata := strings.Split(data.AppName, "-")
		for _, v := range appdata {
			appName = appName + "_" + v
		}

		// appName = data.Env + appName
		tags := map[string]string{"appname": appName}
		fields := map[string]interface{}{
			"app.status":      data.Message,
			"app.environment": data.Env,
			"app.code":        data.Code,
			"app.version":     data.Version,
			"disk.status":     data.Data.Details.DiskSpace.Status,
			"disk.total":      data.Data.Details.DiskSpace.Details.Total,
			"disk.free":       data.Data.Details.DiskSpace.Details.Free,
			"disk.threshold":  data.Data.Details.DiskSpace.Details.Threshold,
			"mysql.status":    data.Data.Details.Db.Mysql.Status,
			"mysql.version":   data.Data.Details.Db.Mysql.Version,
			"redis.status":    data.Data.Details.Db.Redis.Status,
			"redis.version":   data.Data.Details.Db.Redis.Version,
		}

		// pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
		pt, err := client.NewPoint(appName, tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// Write the batch
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// Close client resources
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Println("[Info] 存入influxdb数据解析错误！")
	}

}
