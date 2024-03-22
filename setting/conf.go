package setting

import (
	"apollo/utils"
	"log"

	"gopkg.in/ini.v1"
)

var (
	AppMode        string
	HttpPort       string
	LogLevel       string
	LogToFile      bool
	LogFile        string
	LogFileMaxSize int
	LogFileBackup  int
	MqttOpen       bool
	MqttAddr       string
	MqttClientid   string
	TAddress       string
	TBaudRate      int
	TimeOut        int
	TReadInterval  int
	SFortime       int
)

// GetConf 获取配置文件
func GetConf() {
	path := "./config/config.ini"
	iniFile, err := ini.Load(path)
	if err != nil {
		cfg := ini.Empty()

		AppMode = "debug"
		HttpPort = ":8199"
		cfg.Section("server").Key("AppMode").SetValue(AppMode)
		cfg.Section("server").Key("HttpPort").SetValue(HttpPort)

		LogLevel = "debug"
		LogToFile = false
		LogFile = "./log/run.log"
		LogFileMaxSize = 5
		LogFileBackup = 3
		cfg.Section("logger").Key("LogLevel").SetValue(LogLevel)
		cfg.Section("logger").Key("LogToFile").MustBool(LogToFile)
		cfg.Section("logger").Key("LogFile").SetValue(LogFile)
		cfg.Section("logger").Key("LogFileMaxSize").MustInt(LogFileMaxSize)
		cfg.Section("logger").Key("LogFileBackup").MustInt(LogFileBackup)

		MqttOpen = false
		cfg.Section("mqtt").Key("Open").MustBool(MqttOpen)

		utils.DirIsExist("./config")
		err = cfg.SaveTo(path)
		if err != nil {
			log.Printf("写入config.ini失败 %v", err)
		}
		log.Printf("读取config.ini失败 %v,自动创建成功", err)

		return
	}

	LoadServer(iniFile)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8199")
	LogLevel = file.Section("logger").Key("LogLevel").MustString("debug")
	LogToFile = file.Section("logger").Key("LogToFile").MustBool(false)
	LogFile = file.Section("logger").Key("LogFile").MustString("./log/run.log")
	LogFileMaxSize = file.Section("logger").Key("LogFileMaxSize").MustInt(5)
	LogFileBackup = file.Section("logger").Key("LogFileBackup").MustInt(3)
	MqttOpen = file.Section("mqtt").Key("Open").MustBool(false)
	MqttAddr = file.Section("mqtt").Key("Addr").String()
	MqttClientid = file.Section("mqtt").Key("ClientId").String()
	TAddress = file.Section("tuatr").Key("Address").String()
	TBaudRate = file.Section("tuatr").Key("BaudRate").MustInt(9600)
	TimeOut = file.Section("tuatr").Key("TimeOut").MustInt(500)
	TReadInterval = file.Section("tuatr").Key("ReadInterval").MustInt(500)
	SFortime = file.Section("sensor").Key("ForTime").MustInt(3000)
}
