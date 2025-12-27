package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"os-info/internal/sysinfo"
)

// CreateInfoDisplay creates the main information display
func CreateInfoDisplay(info *sysinfo.Info, w fyne.Window) *fyne.Container {
	title := widget.NewLabelWithStyle("System Information", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	dateTimeSection := createDateTimeSection(
		info.DateTime,
		info.Uptime,
		color.RGBA{R: 100, G: 149, B: 237, A: 255},
	)

	systemSection := createSystemSection(
		info.OSType,
		info.Distribution,
		info.OSVersion,
		color.RGBA{R: 60, G: 179, B: 113, A: 255},
	)

	diskSection := createColoredSectionMultiLineMonospaceWithIcon(
		theme.StorageIcon(),
		"Disk",
		info.GetDiskInfoTable(),
		color.RGBA{R: 255, G: 140, B: 0, A: 255},
	)

	adapterStatus := "offline"
	if info.AdapterOnline {
		adapterStatus = "online"
	}
	batterySection := createBatterySection(
		info.BatteryPercent,
		info.BatteryStatus,
		adapterStatus,
		info.BatteryTemp,
		color.RGBA{R: 220, G: 20, B: 60, A: 255},
	)

	networkTextBinding := binding.NewString()
	_ = networkTextBinding.Set(strings.Join(info.GetNetworkInfoMultiLine(), "\n"))

	networkSection := createDynamicColoredSectionMultiLineMonospaceWithIcon(
		theme.MailSendIcon(),
		"Network",
		networkTextBinding,
		color.RGBA{R: 147, G: 112, B: 219, A: 255},
	)

	info.UpdateExternalNetworkInfo(func() {
		_ = networkTextBinding.Set(strings.Join(info.GetNetworkInfoMultiLine(), "\n"))
	})

	content := container.NewVBox(
		title,
		widget.NewSeparator(),
		dateTimeSection,
		systemSection,
		diskSection,
		batterySection,
		networkSection,
	)

	return content
}

func createBatterySection(percent int, status string, adapterStatus string, temp float64, bgColor color.Color) fyne.CanvasObject {
	icon := widget.NewIcon(theme.WarningIcon())

	batteryBold := canvas.NewText(fmt.Sprintf("Battery: %d%%", percent), color.White)
	batteryBold.TextStyle = fyne.TextStyle{Bold: true}

	statusText := canvas.NewText(fmt.Sprintf("(%s - Adapter %s)", status, adapterStatus), color.White)

	line1 := container.NewHBox(icon, batteryBold, statusText)

	tempLabel := widget.NewLabel(fmt.Sprintf("Temperature: %.1f°C", temp))

	vbox := container.NewVBox(line1, tempLabel)

	rect := canvas.NewRectangle(bgColor)
	rect.SetMinSize(fyne.NewSize(680, 10))

	paddedContent := container.NewPadded(vbox)

	section := container.NewStack(rect, paddedContent)

	return section
}

func createColoredSectionMultiLineMonospaceWithIcon(icon fyne.Resource, title string, lines []string, bgColor color.Color) fyne.CanvasObject {
	var contentObjects []fyne.CanvasObject

	iconWidget := widget.NewIcon(icon)

	if title != "" {
		titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		header := container.NewHBox(iconWidget, titleLabel)
		contentObjects = append(contentObjects, header)
	} else {
		contentObjects = append(contentObjects, iconWidget)
	}

	var textContent string
	for i, line := range lines {
		if line != "" {
			if i > 0 {
				textContent += "\n"
			}
			textContent += line
		}
	}

	if textContent != "" {
		label := widget.NewLabelWithStyle(textContent, fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
		contentObjects = append(contentObjects, label)
	}

	vbox := container.NewVBox(contentObjects...)

	rect := canvas.NewRectangle(bgColor)
	rect.SetMinSize(fyne.NewSize(680, 10))

	paddedContent := container.NewPadded(vbox)

	section := container.NewStack(rect, paddedContent)

	return section
}

func createDateTimeSection(dateTime string, uptime string, bgColor color.Color) fyne.CanvasObject {
	icon := widget.NewIcon(theme.InfoIcon())

	dateText := canvas.NewText(dateTime, color.White)
	dateText.TextStyle = fyne.TextStyle{Bold: true}

	uptimeText := canvas.NewText(fmt.Sprintf("Uptime: %s", uptime), color.White)
	uptimeText.TextSize = 12

	header := container.NewHBox(icon, dateText)

	vbox := container.NewVBox(header, uptimeText)

	rect := canvas.NewRectangle(bgColor)
	rect.SetMinSize(fyne.NewSize(680, 10))

	paddedContent := container.NewPadded(vbox)

	section := container.NewStack(rect, paddedContent)

	return section
}

func createDynamicColoredSectionMultiLineMonospaceWithIcon(icon fyne.Resource, title string, textBinding binding.String, bgColor color.Color) fyne.CanvasObject {
	var contentObjects []fyne.CanvasObject

	iconWidget := widget.NewIcon(icon)

	if title != "" {
		titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		header := container.NewHBox(iconWidget, titleLabel)
		contentObjects = append(contentObjects, header)
	} else {
		contentObjects = append(contentObjects, iconWidget)
	}

	label := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Monospace: true})
	label.Bind(textBinding)
	contentObjects = append(contentObjects, label)

	vbox := container.NewVBox(contentObjects...)

	rect := canvas.NewRectangle(bgColor)
	rect.SetMinSize(fyne.NewSize(680, 10))

	paddedContent := container.NewPadded(vbox)

	section := container.NewStack(rect, paddedContent)

	return section
}

func createSystemSection(osType string, distribution string, osVersion string, bgColor color.Color) fyne.CanvasObject {
	icon := widget.NewIcon(theme.ComputerIcon())

	var systemBold *canvas.Text
	var detailsText *canvas.Text

	if osType == "macOS" {
		systemBold = canvas.NewText("System: macOS", color.White)
		detailsText = canvas.NewText(fmt.Sprintf(" (%s, kernel %s)", distribution, osVersion), color.White)
	} else {
		systemBold = canvas.NewText("System: Linux", color.White)
		detailsText = canvas.NewText(fmt.Sprintf(" (%s, kernel %s)", distribution, osVersion), color.White)
	}
	systemBold.TextStyle = fyne.TextStyle{Bold: true}

	hbox := container.NewHBox(icon, systemBold, detailsText)

	rect := canvas.NewRectangle(bgColor)
	rect.SetMinSize(fyne.NewSize(680, 10))

	paddedContent := container.NewPadded(hbox)

	section := container.NewStack(rect, paddedContent)

	return section
}
