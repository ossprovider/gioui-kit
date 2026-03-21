package component

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Range is a DaisyUI-style range/slider component backed by widget.Float.
type Range struct {
	Float   *widget.Float // Value in [0, 1]
	Variant BtnVariant
	th      *theme.Theme
}

// NewRange creates a new range slider.
func NewRange(th *theme.Theme, f *widget.Float) *Range {
	return &Range{Float: f, Variant: BtnPrimary, th: th}
}

// WithVariant sets the slider accent color.
func (r *Range) WithVariant(v BtnVariant) *Range {
	r.Variant = v
	return r
}

func (r *Range) accentColor() color.NRGBA {
	th := r.th
	switch r.Variant {
	case BtnSecondary:
		return th.Secondary
	case BtnAccent:
		return th.Accent
	case BtnInfo:
		return th.Info
	case BtnSuccess:
		return th.Success
	case BtnWarning:
		return th.Warning
	case BtnError:
		return th.Error
	default:
		return th.Primary
	}
}

// Layout renders the range slider.
func (r *Range) Layout(gtx layout.Context) layout.Dimensions {
	th := r.th
	trackH := gtx.Dp(6)
	thumbSize := gtx.Dp(20)
	totalH := thumbSize
	w := gtx.Constraints.Max.X
	accent := r.accentColor()

	size := image.Pt(w, totalH)
	gtx.Constraints = layout.Exact(size)

	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		// Register drag input spanning the full track
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return r.Float.Layout(gtx, layout.Horizontal, th.Space2)
		}),
		// Visual: track + fill + thumb
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			val := r.Float.Value

			// Track background
			trackY := (totalH - trackH) / 2
			trackRect := image.Rect(0, trackY, w, trackY+trackH)
			defer clip.UniformRRect(trackRect, trackH/2).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Fill
			fillW := int(float32(w) * val)
			if fillW > 0 {
				fillRect := image.Rect(0, trackY, fillW, trackY+trackH)
				defer clip.UniformRRect(fillRect, trackH/2).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: accent}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
			}

			// Thumb
			thumbX := int(float32(w-thumbSize) * val)
			thumbRect := image.Rect(thumbX, 0, thumbX+thumbSize, thumbSize)
			defer clip.UniformRRect(thumbRect, thumbSize/2).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: accent}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// White center dot
			dotSize := gtx.Dp(8)
			dotOff := (thumbSize - dotSize) / 2
			dotRect := image.Rect(thumbX+dotOff, dotOff, thumbX+dotOff+dotSize, dotOff+dotSize)
			defer clip.UniformRRect(dotRect, dotSize/2).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: theme.White}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			return layout.Dimensions{Size: size}
		}),
	)
}
