package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"os-info/internal/sysinfo"
	"os-info/internal/ui"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&ui.CustomTheme{})

	w := a.NewWindow("System Information")

	sysInfo := sysinfo.New()

	content := ui.CreateInfoDisplay(sysInfo, w)

	tappable := ui.NewTappableContainer(content, func() {
		w.Close()
	})

	w.SetContent(tappable)

	w.SetOnClosed(func() {
		a.Quit()
	})

	w.SetFullScreen(true)

	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		w.Close()
	})

	w.ShowAndRun()
}
