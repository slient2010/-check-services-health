package res

import (
	"log"
)

// 定义接口
type GetUrl interface {
	GetUrlData(u chan string)
}

type UrlData struct {
	Url string
}

// 实现接口
func (url *UrlData) GetUrlData(u chan string) {
	// 从数据库中获取应用信息
	data := GetData()

	for _, v := range data {
		// 获取app url，构造health检查地址
		// todo: dev环境，获取jenkins状态，要是在打包则跳过检查，等等
		if v.HealthUrl == "" {
			log.Printf("[warn] %s has not health check url", v.Name)
			continue
		}
		healthCheckUrl := string(v.HealthUrl)
		u <- healthCheckUrl
		// fmt.Println(v.AppUrl)
	}

}

// app service struct
/*
healthData = `
{
	"appname": "tmc-service",   // 存入influxdb, 则后面处理加上环境.如开发环境（dev)，为:dev-tmc-service
	"code": 0,		 // 状态码，需要明确说明，是整个服务ok，还是其他关联的服务，如mysql，redis等服务ok？ 0 -- 服务都正常，1 -- 服务异常？
	"message": "OK",		 // 和code区别？
	"data": {                  // 关联的内部服务？mysql，redis，是否还有其他的关联？
		"db": {
			"mysql": {
				"status": 0,    // 以状态码方式说明，0 -- 服务都正常，1 -- 服务异常？
				"version": "5.7.0",
			},
			"redis": {
				"status": 0,    // 以状态码方式说明，0 -- 服务都正常，1 -- 服务异常？
				"version": "4.0.10"
			}
		}
	}
}
`

*/

type AppSrv struct {
	AppName string `json:"appname"` // 应用名称
	Message string `json:"message"` // 应用状态, "OK"?
	Env     string `json:"env"`
	Version string `json:"version"`
	Code    int    `json:"code"` // 应用状态码
	Data    struct {
		Status  string `json:"status"`
		Details struct {
			DiskSpace struct {
				Status  string `json:"status"`
				Details struct {
					Total     float64 `json:"total"`
					Free      float64 `json:"free"`
					Threshold float64 `json:"threshold"`
				}
			}
			Db struct {
				Mysql struct {
					Status  int    `json:"Status"`
					Version string `json:"version"`
				} `json:"mysql"`
				Redis struct {
					Status  int    `json:"Status"`
					Version string `json:"version"`
				} `json:"redis"`
			} `json:"db"`
		} `json:"details"`
	} `json:"data"`
}
