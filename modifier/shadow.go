package modifier

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ShadowStyle defines a shadow preset (like Tailwind shadow-sm, shadow-lg, etc.).
type ShadowStyle struct {
	OffsetX unit.Dp
	OffsetY unit.Dp
	Blur    unit.Dp
	Spread  unit.Dp
	Color   color.NRGBA
}

// Preset shadows following Tailwind conventions.
var (
	ShadowSm = ShadowStyle{
		OffsetY: 1, Blur: 3,
		Color: color.NRGBA{A: 60},
	}
	ShadowMd = ShadowStyle{
		OffsetY: 4, Blur: 8, Spread: -1,
		Color: color.NRGBA{A: 80},
	}
	ShadowLg = ShadowStyle{
		OffsetY: 10, Blur: 15, Spread: -3,
		Color: color.NRGBA{A: 100},
	}
	ShadowXl = ShadowStyle{
		OffsetY: 20, Blur: 25, Spread: -5,
		Color: color.NRGBA{A: 120},
	}
	Shadow2xl = ShadowStyle{
		OffsetY: 25, Blur: 50, Spread: -12,
		Color: color.NRGBA{A: 150},
	}
)

// Shadow applies a drop shadow behind a widget.
type Shadow struct {
	Style  ShadowStyle
	Radius unit.Dp
}

func (s Shadow) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	// Draw shadow layers (approximation using multiple offset rects)
	blur := gtx.Dp(s.Style.Blur)
	offX := gtx.Dp(s.Style.OffsetX)
	offY := gtx.Dp(s.Style.OffsetY)
	spread := gtx.Dp(s.Style.Spread)
	rr := gtx.Dp(s.Radius)

	if blur > 0 {
		layers := blur / 2
		if layers < 2 {
			layers = 2
		}
		if layers > 8 {
			layers = 8
		}
		for i := 0; i < layers; i++ {
			t := float64(i+1) / float64(layers)
			expand := int(float64(blur)*t) + spread
			alpha := uint8(float64(s.Style.Color.A) * (1 - t))
			if alpha == 0 {
				continue
			}

			shadowRect := image.Rectangle{
				Min: image.Pt(offX-expand, offY-expand),
				Max: image.Pt(dims.Size.X+offX+expand, dims.Size.Y+offY+expand),
			}
			shadowColor := s.Style.Color
			shadowColor.A = alpha

			shadowRR := rr + expand
			if shadowRR < 0 {
				shadowRR = 0
			}

			stack := clip.UniformRRect(shadowRect, shadowRR).Push(gtx.Ops)
			paint.ColorOp{Color: shadowColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			stack.Pop()
		}
	}

	call.Add(gtx.Ops)
	return dims
}
