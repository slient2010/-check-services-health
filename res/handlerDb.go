package res

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

type AppInfo struct {
	gorm.Model
	AppID     int
	Name      string `gorm:"size:100:unique; not null;"` // 业务名称
	Port      string `gorm:"size:50;"`                   // 项目生命周期
	AppUrl    string `gorm:"size:100;"`                  // 业务对外地址，接口
	Env       string `gorm:"size:50;"`                   // 业务当前环境
	AppType   string `gorm:"size:50:not null;"`          // 业务生命周期
	Version   string `gorm:"size:50;"`                   // 业务对外地址，接口
	AppDocs   string `gorm:"size:50;"`                   // 接口地址，文档
	HealthUrl string `gorm:"size:100;"`                  // 健康检查地址
	Status    string `gorm:"size:45;"`                   // 项目状态
	Notes     string `gorm:"size:100;"`                  // 备注
	Operator  string `gorm:"size:50;"`                   // 记录操作人, 方便定位

}

var (
	DbHost  = "127.0.0.1"
	DbUser  = "devops"
	DbPass  = "password"
	DbName  = "devops"
	Charset = "utf8"
	DbPort  = 3306
	DbDebug = false
)

func dbInfo() *gorm.DB {
	db, err := gorm.Open("mysql", DbUser+":"+DbPass+"@tcp("+DbHost+")/"+DbName+"?charset="+Charset+"&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	db.LogMode(DbDebug)
	err = db.DB().Ping()
	if err != nil {
		db.DB().Close()
		log.Fatal("Error", err.Error())
	}
	return db
}

func GetData() []AppInfo {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Error", r)
		}
	}()
	db := dbInfo()
	defer db.Close()

	if !db.HasTable(&AppInfo{}) {
		return nil
	}
	var appinfo []AppInfo
	all := db.Find(&appinfo)
	jsonData, _ := json.Marshal(all.Value)
	// fmt.Println(string(jsonData))
	var appData []AppInfo
	json.Unmarshal([]byte(jsonData), &appData)

	return appData
}
