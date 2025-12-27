package sysinfo

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func (i *Info) collectBatteryInfo() {
	i.BatteryPercent = 0
	i.BatteryStatus = "N/A"
	i.BatteryTemp = 0.0
	i.AdapterOnline = false

	if runtime.GOOS == "darwin" {
		readBatteryMacOS(i)
	} else {
		readBatteryLinux(i)
	}
}

func readBatteryLinux(i *Info) {
	batteries := []string{"BAT0", "BAT1"}

	for _, bat := range batteries {
		basePath := fmt.Sprintf("/sys/class/power_supply/%s", bat)

		if _, err := os.Stat(basePath); err == nil {
			readBatteryFromSysLinux(i, basePath)
			break
		}
	}

	adapters := []string{"AC", "AC0", "ADP0", "ADP1"}
	for _, adapter := range adapters {
		adapterPath := fmt.Sprintf("/sys/class/power_supply/%s/online", adapter)
		if data, err := os.ReadFile(adapterPath); err == nil {
			var online int
			_, _ = fmt.Sscanf(string(data), "%d", &online)
			i.AdapterOnline = online == 1
			break
		}
	}
}

func readBatteryMacOS(i *Info) {
	// macOS battery reading would require pmset command
	// Setting defaults for now
	i.BatteryPercent = 0
	i.BatteryStatus = "N/A"
	i.BatteryTemp = 0.0
	i.AdapterOnline = false
}

func readBatteryFromSysLinux(i *Info, basePath string) {
	if data, err := os.ReadFile(basePath + "/capacity"); err == nil {
		_, _ = fmt.Sscanf(string(data), "%d", &i.BatteryPercent)
	}

	if data, err := os.ReadFile(basePath + "/status"); err == nil {
		i.BatteryStatus = strings.TrimSpace(string(data))
	}

	if data, err := os.ReadFile(basePath + "/temp"); err == nil {
		var temp int
		_, _ = fmt.Sscanf(string(data), "%d", &temp)
		i.BatteryTemp = float64(temp) / 10.0
	}
}
