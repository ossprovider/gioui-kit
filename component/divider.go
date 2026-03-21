package component

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Divider is a DaisyUI-style divider line with optional label.
type Divider struct {
	Label    string
	Vertical bool
	th       *theme.Theme
}

// NewDivider creates a new divider.
func NewDivider(th *theme.Theme) *Divider {
	return &Divider{th: th}
}

// WithLabel sets a text label centered in the divider.
func (d *Divider) WithLabel(label string) *Divider {
	d.Label = label
	return d
}

// WithVertical makes the divider vertical.
func (d *Divider) WithVertical() *Divider {
	d.Vertical = true
	return d
}

// Layout renders the divider.
func (d *Divider) Layout(gtx layout.Context) layout.Dimensions {
	th := d.th
	lineThick := gtx.Dp(1)
	lineColor := th.Base300

	if d.Vertical {
		h := gtx.Constraints.Max.Y
		if h == 0 {
			h = gtx.Dp(40)
		}
		w := gtx.Dp(20)
		lineX := w / 2
		rect := image.Rect(lineX, 0, lineX+lineThick, h)
		defer clip.Rect(rect).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Pt(w, h)}
	}

	w := gtx.Constraints.Max.X
	totalH := gtx.Dp(32)
	midY := totalH / 2

	if d.Label == "" {
		rect := image.Rect(0, midY, w, midY+lineThick)
		defer clip.Rect(rect).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: image.Pt(w, totalH)}
	}

	// Label version: line — text — line
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			h := gtx.Constraints.Max.Y
			if h == 0 {
				h = totalH
			}
			lineY := h / 2
			rect := image.Rect(0, lineY, gtx.Constraints.Max.X, lineY+lineThick)
			defer clip.Rect(rect).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, h)}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: th.Space3, Right: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, d.Label, th.BaseContent, th.SmSize, font.Normal)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			h := gtx.Constraints.Max.Y
			if h == 0 {
				h = totalH
			}
			lineY := h / 2
			rect := image.Rect(0, lineY, gtx.Constraints.Max.X, lineY+lineThick)
			defer clip.Rect(rect).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: lineColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, h)}
		}),
	)
}
