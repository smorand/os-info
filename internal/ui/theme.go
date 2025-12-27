package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme increases the text size by 1.5x
type CustomTheme struct{}

// Color returns theme colors
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

// Font returns theme fonts
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns theme icons
func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns theme sizes with 1.5x multiplier for text
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	baseSize := theme.DefaultTheme().Size(name)
	return baseSize * 1.5
}
