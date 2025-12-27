package sysinfo

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskInfo represents information about a disk partition
type DiskInfo struct {
	Free        uint64
	MountPoint  string
	Total       uint64
	Used        uint64
	UsedPercent float64
}

// GetDiskInfoTable returns disk information as a formatted table
func (i *Info) GetDiskInfoTable() []string {
	if len(i.Disks) == 0 {
		return []string{"No disk information available"}
	}

	var lines []string

	lines = append(lines, fmt.Sprintf("%-15s %12s %12s %12s %10s",
		"Mount Point", "Total", "Used", "Free", "Usage"))

	lines = append(lines, strings.Repeat("-", 70))

	for _, d := range i.Disks {
		line := fmt.Sprintf("%-15s %10.1f GB %10.1f GB %10.1f GB %9.1f%%",
			d.MountPoint,
			float64(d.Total)/1024/1024/1024,
			float64(d.Used)/1024/1024/1024,
			float64(d.Free)/1024/1024/1024,
			d.UsedPercent)
		lines = append(lines, line)
	}

	return lines
}

func (i *Info) collectDiskInfo() {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return
	}

	for _, partition := range partitions {
		if isVirtualFilesystem(partition.Fstype) {
			continue
		}

		if isVirtualDevice(partition.Device) {
			continue
		}

		if strings.HasPrefix(partition.Mountpoint, "/boot") {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		i.Disks = append(i.Disks, DiskInfo{
			Free:        usage.Free,
			MountPoint:  partition.Mountpoint,
			Total:       usage.Total,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
		})
	}
}

func isVirtualFilesystem(fstype string) bool {
	virtualFs := []string{
		"tmpfs", "devtmpfs", "sysfs", "proc", "devpts",
		"cgroup", "cgroup2", "pstore", "bpf", "tracefs",
		"debugfs", "securityfs", "sockfs", "pipefs",
		"configfs", "selinuxfs", "autofs", "mqueue",
		"hugetlbfs", "fusectl", "fuse.gvfsd-fuse",
		"fuse.portal", "nsfs", "overlay", "squashfs",
	}

	for _, vfs := range virtualFs {
		if strings.Contains(fstype, vfs) {
			return true
		}
	}
	return false
}

func isVirtualDevice(device string) bool {
	if strings.HasPrefix(device, "/dev/loop") {
		return true
	}

	if strings.Contains(device, "/snap/") {
		return true
	}

	virtualPrefixes := []string{
		"none", "udev", "tmpfs", "cgmfs", "overlay",
	}

	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(device, prefix) {
			return true
		}
	}

	return false
}
