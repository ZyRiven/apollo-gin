package setting

import (
	"apollo/consts"
	"apollo/utils"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemStateTemplate struct {
	MemUse  string `json:"memUse"`
	DiskUse string `json:"diskUse"`
	CpuUse  string `json:"cpuUse"`
	Ip      string `json:"ip"`
	SoftVer string `json:"softVer"`
	RunTime string `json:"runTime"` //累计时间
	// DeviceOnline     string `json:"deviceOnline"`     //设备在线率
	// DevicePacketLoss string `json:"devicePacketLoss"` //设备丢包率
}

var SystemState = SystemStateTemplate{
	Ip:      "127.0.0.1",
	MemUse:  "0%",
	DiskUse: "0%",
	CpuUse:  "0%",
	SoftVer: consts.BuildVersion,
	RunTime: "0",
	// DeviceOnline:     "0",
	// DevicePacketLoss: "0",
}
var (
	timeStart   time.Time
	CronProcess *cron.Cron
)

// StreamInit 获取内存、ip等系统信息
func StreamInit() {
	timeStart = time.Now()
	CollectSystemParam()
	CronProcess = cron.New(cron.WithSeconds())
	CronProcess.AddFunc("*/60 * * * * *", CollectSystemParam)
	CronProcess.Start()
}

func CollectSystemParam() {
	/************** 开机时间 ***********************/
	elapsed := time.Since(timeStart)
	sec := int64(elapsed.Seconds())
	day := sec / 86400
	hour := sec % 86400 / 3600
	min := sec % 3600 / 60
	sec = sec % 60
	SystemState.RunTime = fmt.Sprintf("%d天%d时%d分%d秒", day, hour, min, sec)
	/************** 获取IP    ***********************/
	SystemState.Ip = utils.GetIp().String()
	/************** 内存使用  ***********************/
	v, _ := mem.VirtualMemory()
	SystemState.MemUse = fmt.Sprintf("%3.1f", v.UsedPercent) + "%"
	/************** 磁盘使用  ***********************/
	exeCurDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	d, _ := disk.Usage(exeCurDir)
	SystemState.DiskUse = fmt.Sprintf("%3.1f", d.UsedPercent) + "%"
	/************** cpu使用  ***********************/
	percent, _ := cpu.Percent(time.Second, false)
	SystemState.CpuUse = fmt.Sprintf("%3.1f", percent[0]) + "%"
}
