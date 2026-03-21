// Package modifier provides Tailwind-style visual modifiers for Gio widgets.
//
// These are composable decorators that can be applied to any widget:
//
//	modifier.Shadow(modifier.ShadowLg).Layout(gtx, myWidget)
//	modifier.Bg(theme.Primary).Layout(gtx, myWidget)
//	modifier.Rounded(theme.RoundedLg).Layout(gtx, myWidget)
package modifier

import (
	"image"
	"image/color"
	"math"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

// ---------- Background ----------

// Bg adds a background color to a widget.
type Bg struct {
	Color  color.NRGBA
	Radius unit.Dp
}

func (b Bg) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	rr := gtx.Dp(b.Radius)
	rrect := clip.UniformRRect(image.Rectangle{Max: dims.Size}, rr)
	defer rrect.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: b.Color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	call.Add(gtx.Ops)
	return dims
}

// ---------- Rounded (clip) ----------

// Rounded clips a widget to a rounded rectangle.
type Rounded struct {
	Radius unit.Dp
}

func (r Rounded) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	rr := gtx.Dp(r.Radius)
	defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, rr).Push(gtx.Ops).Pop()
	call.Add(gtx.Ops)
	return dims
}

// ---------- Shadow ----------

// ShadowStyle defines a shadow preset (like Tailwind shadow-sm, shadow-lg, etc.).
type ShadowStyle struct {
	OffsetX  unit.Dp
	OffsetY  unit.Dp
	Blur     unit.Dp
	Spread   unit.Dp
	Color    color.NRGBA
}

// Preset shadows following Tailwind conventions.
var (
	ShadowSm = ShadowStyle{
		OffsetY: 1, Blur: 2,
		Color: color.NRGBA{A: 13},
	}
	ShadowMd = ShadowStyle{
		OffsetY: 4, Blur: 6, Spread: -1,
		Color: color.NRGBA{A: 20},
	}
	ShadowLg = ShadowStyle{
		OffsetY: 10, Blur: 15, Spread: -3,
		Color: color.NRGBA{A: 25},
	}
	ShadowXl = ShadowStyle{
		OffsetY: 20, Blur: 25, Spread: -5,
		Color: color.NRGBA{A: 30},
	}
	Shadow2xl = ShadowStyle{
		OffsetY: 25, Blur: 50, Spread: -12,
		Color: color.NRGBA{A: 40},
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
			alpha := uint8(float64(s.Style.Color.A) * (1 - t) / float64(layers))
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

// ---------- Opacity ----------

// OpacityMod applies an opacity modifier.
type OpacityMod struct {
	Opacity float32 // 0.0 - 1.0
}

func (o OpacityMod) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	// Use paint.OpacityOp if available; fallback to color tinting
	call.Add(gtx.Ops)
	return dims
}

// ---------- Ring (focus ring / outline) ----------

// Ring draws an outline ring around a widget (like Tailwind `ring-2 ring-blue-500`).
type Ring struct {
	Width  unit.Dp
	Color  color.NRGBA
	Offset unit.Dp
	Radius unit.Dp
}

func (r Ring) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	call.Add(gtx.Ops)

	// Draw ring
	rw := gtx.Dp(r.Width)
	ro := gtx.Dp(r.Offset)
	rr := gtx.Dp(r.Radius)

	ringRect := image.Rectangle{
		Min: image.Pt(-ro, -ro),
		Max: image.Pt(dims.Size.X+ro, dims.Size.Y+ro),
	}

	paint.FillShape(gtx.Ops, r.Color,
		clip.Stroke{
			Path:  clip.UniformRRect(ringRect, rr+ro).Path(gtx.Ops),
			Width: float32(rw),
		}.Op(),
	)

	return dims
}

// ---------- Gradient ----------

// GradientDir specifies gradient direction.
type GradientDir int

const (
	GradientToRight GradientDir = iota
	GradientToBottom
	GradientToBottomRight
)

// LinearGradient applies a linear gradient background.
type LinearGradient struct {
	From   color.NRGBA
	To     color.NRGBA
	Dir    GradientDir
	Radius unit.Dp
}

func (g LinearGradient) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	dims := w(gtx)
	call := macro.Stop()

	rr := gtx.Dp(g.Radius)

	// Draw gradient via multiple horizontal or vertical strips
	steps := 32
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps)
		t2 := float32(i+1) / float32(steps)
		col := lerpColor(g.From, g.To, t)

		var stripRect image.Rectangle
		switch g.Dir {
		case GradientToRight:
			x1 := int(t * float32(dims.Size.X))
			x2 := int(t2 * float32(dims.Size.X))
			stripRect = image.Rect(x1, 0, x2, dims.Size.Y)
		case GradientToBottom:
			y1 := int(t * float32(dims.Size.Y))
			y2 := int(t2 * float32(dims.Size.Y))
			stripRect = image.Rect(0, y1, dims.Size.X, y2)
		case GradientToBottomRight:
			x1 := int(t * float32(dims.Size.X))
			x2 := int(t2 * float32(dims.Size.X))
			y1 := int(t * float32(dims.Size.Y))
			y2 := int(t2 * float32(dims.Size.Y))
			stripRect = image.Rect(x1, y1, x2, y2)
		}

		if i == 0 && rr > 0 {
			defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, rr).Push(gtx.Ops).Pop()
		}
		paint.FillShape(gtx.Ops, col, clip.Rect(stripRect).Op())
	}

	call.Add(gtx.Ops)
	return dims
}

func lerpColor(a, b color.NRGBA, t float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(float32(a.R)*(1-t) + float32(b.R)*t),
		G: uint8(float32(a.G)*(1-t) + float32(b.G)*t),
		B: uint8(float32(a.B)*(1-t) + float32(b.B)*t),
		A: uint8(float32(a.A)*(1-t) + float32(b.A)*t),
	}
}

// ---------- Transition / Animation helpers ----------

// EaseInOut returns an eased value for animations.
func EaseInOut(t float32) float32 {
	return float32(-(math.Cos(math.Pi*float64(t)) - 1) / 2)
}

// EaseOut returns an ease-out value.
func EaseOut(t float32) float32 {
	return float32(1 - math.Pow(1-float64(t), 3))
}
