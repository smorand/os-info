package sysinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

const requestTimeout = 30 * time.Second

// NetworkInfo represents network interface information
type NetworkInfo struct {
	ConnectionType string
	Country        string
	DNS            []string
	ESSID          string
	ExternalIP     string
	Gateway        string
	Interface      string
	IPAddress      string
	MACAddress     string
}

// GetNetworkInfoMultiLine returns network information as formatted lines
func (i *Info) GetNetworkInfoMultiLine() []string {
	if len(i.Networks) == 0 {
		return []string{"No network information available"}
	}

	n := i.Networks[0]
	lines := []string{
		fmt.Sprintf("%-15s %s (%s)", "Interface:", n.Interface, n.ConnectionType),
		fmt.Sprintf("%-15s %s", "IP Address:", n.IPAddress),
		fmt.Sprintf("%-15s %s", "MAC Address:", n.MACAddress),
	}

	if n.ConnectionType == "WiFi" && n.ESSID != "N/A" && n.ESSID != "" {
		lines = append(lines, fmt.Sprintf("%-15s %s", "ESSID:", n.ESSID))
	}

	lines = append(lines, fmt.Sprintf("%-15s %s", "Gateway:", n.Gateway))
	lines = append(lines, fmt.Sprintf("%-15s %s", "DNS Servers:", strings.Join(n.DNS, ", ")))
	lines = append(lines, fmt.Sprintf("%-15s %s", "External IP:", n.ExternalIP))
	lines = append(lines, fmt.Sprintf("%-15s %s", "Country:", n.Country))

	return lines
}

func (i *Info) collectNetworkInfo() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	defaultGateway := getDefaultGateway()
	activeInterface := getActiveInterface()

	for _, iface := range interfaces {
		if strings.HasPrefix(iface.Name, "lo") || len(iface.Addrs) == 0 {
			continue
		}

		if activeInterface != "" && iface.Name != activeInterface {
			continue
		}

		hasIPv4 := false
		for _, addr := range iface.Addrs {
			if strings.Contains(addr.Addr, ".") {
				hasIPv4 = true
				break
			}
		}
		if !hasIPv4 {
			continue
		}

		netInfo := NetworkInfo{
			Interface:  iface.Name,
			MACAddress: iface.HardwareAddr,
		}

		for _, addr := range iface.Addrs {
			if strings.Contains(addr.Addr, ".") {
				netInfo.IPAddress = addr.Addr
				break
			}
		}

		if strings.HasPrefix(iface.Name, "wl") || strings.HasPrefix(iface.Name, "wlan") {
			netInfo.ConnectionType = "WiFi"
			netInfo.ESSID = getWifiESSID(iface.Name)
		} else if strings.HasPrefix(iface.Name, "en") || strings.HasPrefix(iface.Name, "eth") {
			netInfo.ConnectionType = "Ethernet"
		} else {
			netInfo.ConnectionType = "Other"
		}

		netInfo.Gateway = defaultGateway
		netInfo.DNS = getDNSServers()
		netInfo.ExternalIP = "searching..."
		netInfo.Country = "searching..."

		i.Networks = append(i.Networks, netInfo)
		break
	}
}

func getActiveInterface() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "00000000" {
			return fields[0]
		}
	}

	return ""
}

func getDefaultGateway() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return "N/A"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[1] == "00000000" {
			gateway := fields[2]
			var a, b, c, d int
			_, _ = fmt.Sscanf(gateway, "%02x%02x%02x%02x", &d, &c, &b, &a)
			return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
		}
	}

	return "N/A"
}

func getDNSServers() []string {
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return []string{"N/A"}
	}

	var dnsServers []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				dnsServers = append(dnsServers, fields[1])
			}
		}
	}

	if len(dnsServers) == 0 {
		return []string{"N/A"}
	}

	return dnsServers
}

func getExternalIP() string {
	url := "http://api.ipify.org"

	client := &http.Client{
		Timeout: requestTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "N/A"
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "N/A"
		}
		ip := strings.TrimSpace(string(body))
		if ip != "" {
			return ip
		}
	}

	return "N/A"
}

func getCountry(externalIP string) string {
	if externalIP == "N/A" || externalIP == "" {
		return "N/A"
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country", externalIP)

	client := &http.Client{
		Timeout: requestTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "N/A"
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "N/A"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "N/A"
	}

	var result struct {
		Country string `json:"country"`
		Status  string `json:"status"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "N/A"
	}

	if result.Status == "success" && result.Country != "" {
		return result.Country
	}

	return "N/A"
}

func getWifiESSID(iface string) string {
	cmd := exec.Command("iwgetid", "-r", iface)
	output, err := cmd.Output()
	if err == nil {
		ssid := strings.TrimSpace(string(output))
		if ssid != "" {
			return ssid
		}
	}

	cmd = exec.Command("iw", "dev", iface, "link")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "SSID:") {
				parts := strings.Split(line, "SSID:")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	data, err := os.ReadFile("/proc/net/wireless")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, iface) {
				return "Connected"
			}
		}
	}

	return "N/A"
}
