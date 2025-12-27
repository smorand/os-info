package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TappableContainer is a container that executes a callback when tapped
type TappableContainer struct {
	widget.BaseWidget
	content  fyne.CanvasObject
	onTapped func()
}

// NewTappableContainer creates a new tappable container
func NewTappableContainer(content fyne.CanvasObject, onTapped func()) *TappableContainer {
	t := &TappableContainer{
		content:  content,
		onTapped: onTapped,
	}
	t.ExtendBaseWidget(t)
	return t
}

// CreateRenderer creates the widget renderer
func (t *TappableContainer) CreateRenderer() fyne.WidgetRenderer {
	return &tappableRenderer{
		container: t,
		objects:   []fyne.CanvasObject{t.content},
	}
}

// Tapped handles tap events
func (t *TappableContainer) Tapped(*fyne.PointEvent) {
	if t.onTapped != nil {
		t.onTapped()
	}
}

type tappableRenderer struct {
	container *TappableContainer
	objects   []fyne.CanvasObject
}

func (r *tappableRenderer) Destroy() {
}

func (r *tappableRenderer) Layout(size fyne.Size) {
	r.container.content.Resize(size)
	r.container.content.Move(fyne.NewPos(0, 0))
}

func (r *tappableRenderer) MinSize() fyne.Size {
	return r.container.content.MinSize()
}

func (r *tappableRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *tappableRenderer) Refresh() {
	r.container.content.Refresh()
}
