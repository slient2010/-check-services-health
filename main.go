package main

import (
	"check-services-health/common"
	"check-services-health/res"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	fileCfg  = "config.cfg"
	CfgData  Cfg
	TickTime int64
)

// 监控软件本身监控
type MonitorSrv struct {
	ChkSrvName string `json:"chksrv"`    // 服务健康检查名称
	ChkSrvUrl  string `json:"chksrvurl"` // 健康检查地址
	ReqstTimes int    `json:"int"`       // 健康检查次数
	ErrorTimes int    `json:"int"`       // 检查报错次数
}

// 健康检测结构体
type HealthChk struct {
	r     chan string      // 传递url地址
	w     chan *res.AppSrv // 读取url处理url地址
	url   res.GetUrl       // 获取url
	check CheckSrv         // 健康检查
}

// 处理获取到的health地址
func (hlthChk *HealthChk) Process() {
	hlturl := <-hlthChk.r
	// do request

	// 解析url得到app名称
	realurl, _ := url.Parse(string(hlturl))
	appname := strings.Split(realurl.Hostname(), ".")[0]

	// 健康检查
	result := common.HttpClientChkSrv(hlturl)

	appSrv := &res.AppSrv{}

	if err := json.Unmarshal([]byte(result), appSrv); err != nil {
		// log.Println(err)
		log.Fatalf("[Error]", err)
	}

	// 更新部分数据
	appSrv.AppName = appname
	// todo: 需要获取当前app的环境及版本
	appSrv.Env = "dev"
	appSrv.Version = "v0.0.1"

	// send result
	hlthChk.w <- appSrv
}

// 定义接口
type CheckSrv interface {
	SaveChkReslt(u chan *res.AppSrv)
}

type CheckData struct {
	ChkData string
}

// 健康检查并存入数据到influxdb
func (ChkData *CheckData) SaveChkReslt(u chan *res.AppSrv) {
	hltChkRes := <-u
	common.SaveToInfluxDb(hltChkRes)

}

type Cfg struct {
	CheckTickTime int64  `json:"CheckTickTime"`
	Notes         string `json:"Notes"`
	DataBase      struct {
		DbHost  string `json:"DbHost"`
		DbUser  string `json:"DbUser"`
		DbPass  string `json:"DbPass"`
		DbName  string `json:"DbName"`
		Charset string `json:"Charset"`
		DbPort  int    `json:"DbPort"`
		DbDebug bool   `json:"DbDebug"`
		Notes   string `json:"Notes"`
	} `json:"MySQL"`
	InfluxDb struct {
		DbUrl     string `json:"DbUrl"`
		DbName    string `json:"DbName"`
		Username  string `json:"Username"`
		Password  string `json:"Password"`
		Precision string `json:"Precision"`
		Notes     string `json:"Notes"`
	} `json:"InfluxDb"`
}

func initCfg() {
	data, err := ioutil.ReadFile(fileCfg)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	if err = json.Unmarshal(data, &CfgData); err != nil {
		log.Fatalf("解析配置文件 %s 失败，", fileCfg, err)
	}
	// mysql 配置信息
	res.DbHost = CfgData.DataBase.DbHost
	res.DbUser = CfgData.DataBase.DbUser
	res.DbPass = CfgData.DataBase.DbPass
	res.DbName = CfgData.DataBase.DbName
	res.Charset = CfgData.DataBase.Charset
	res.DbPort = CfgData.DataBase.DbPort
	res.DbDebug = CfgData.DataBase.DbDebug
	// influxdb配置信息
	common.DbUrl = CfgData.InfluxDb.DbUrl
	common.DbName = CfgData.InfluxDb.DbName
	common.Username = CfgData.InfluxDb.Username
	common.Password = CfgData.InfluxDb.Password
	common.Precision = CfgData.InfluxDb.Precision
	// 监控检测频度
	TickTime = CfgData.CheckTickTime
}

func exists(path string) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(r)
		}
	}()
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {
	fmt.Println("-----------------------------------")
	fileStatus, _ := exists(fileCfg)
	if !fileStatus {
		log.Fatalf("配置文件%s不存在", fileCfg)
	}
	log.Print("加载配置文件...")
	// 加载配置文件
	initCfg()

	r := &res.UrlData{
		Url: "",
	}

	w := &CheckData{
		ChkData: "",
	}

	h := &HealthChk{
		r:     make(chan string, 20),
		w:     make(chan *res.AppSrv),
		url:   r,
		check: w,
	}

	// 循环检查
	for {
		// 并发执行
		// 获取健康检查URL
		go h.url.GetUrlData(h.r)
		for i := 0; i <= 20; i++ {
			// 健康检查
			go h.Process()
		}
		for j := 0; j <= 30; j++ {
			// 保存到数据库
			go h.check.SaveChkReslt(h.w)
		}

		time.Sleep(time.Duration(TickTime) * time.Second)
	}
}
