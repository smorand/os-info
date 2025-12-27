package sysinfo

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/host"
)

// Info contains all system information
type Info struct {
	AdapterOnline  bool
	BatteryPercent int
	BatteryStatus  string
	BatteryTemp    float64
	DateTime       string
	Disks          []DiskInfo
	Distribution   string
	Networks       []NetworkInfo
	OSType         string
	OSVersion      string
	Uptime         string
}

// New creates and populates a new Info instance
func New() *Info {
	info := &Info{}

	info.collectDateTimeInfo()
	info.collectOSInfo()
	info.collectDiskInfo()
	info.collectBatteryInfo()
	info.collectNetworkInfo()

	return info
}

// UpdateExternalNetworkInfo updates the external IP and country asynchronously
func (i *Info) UpdateExternalNetworkInfo(callback func()) {
	if len(i.Networks) == 0 {
		return
	}

	go func() {
		externalIP := getExternalIP()
		country := getCountry(externalIP)

		i.Networks[0].ExternalIP = externalIP
		i.Networks[0].Country = country

		if callback != nil {
			callback()
		}
	}()
}

func (i *Info) collectDateTimeInfo() {
	now := time.Now()

	day := now.Day()
	suffix := getDaySuffix(day)
	i.DateTime = fmt.Sprintf("%s %d%s %s %d - %s",
		now.Format("Monday"),
		day,
		suffix,
		now.Format("January"),
		now.Year(),
		now.Format("15:04:05"))

	hostInfo, err := host.Info()
	if err == nil {
		uptime := time.Duration(hostInfo.Uptime) * time.Second
		days := int(uptime.Hours() / 24)
		hours := int(uptime.Hours()) % 24
		minutes := int(uptime.Minutes()) % 60

		i.Uptime = fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else {
		i.Uptime = "Unknown"
	}
}

func (i *Info) collectOSInfo() {
	osType := runtime.GOOS

	hostInfo, err := host.Info()
	if err == nil {
		if osType == "darwin" {
			i.OSType = "macOS"
			i.Distribution = hostInfo.PlatformVersion
			i.OSVersion = hostInfo.KernelVersion
		} else {
			i.OSType = "Linux"
			i.Distribution = fmt.Sprintf("%s %s", hostInfo.Platform, hostInfo.PlatformVersion)
			i.OSVersion = hostInfo.KernelVersion
		}
	} else {
		i.OSType = osType
		i.OSVersion = "Unknown"
		i.Distribution = "Unknown"
	}
}

func getDaySuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}
