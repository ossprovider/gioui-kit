package component

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/hongshengjie/gioui-kit/theme"
)

// RadialProgress renders a DaisyUI-style circular progress indicator.
type RadialProgress struct {
	Value   float32 // 0.0 to 1.0
	Size    unit.Dp
	Thick   unit.Dp
	Label   string // centered label (e.g. "75%")
	Variant ProgressVariant
	th      *theme.Theme
}

// NewRadialProgress creates a new radial (circular) progress component.
func NewRadialProgress(th *theme.Theme, value float32) *RadialProgress {
	return &RadialProgress{
		Value:   value,
		Size:    80,
		Thick:   8,
		Variant: ProgressPrimary,
		th:      th,
	}
}

// WithSize sets the diameter in dp.
func (r *RadialProgress) WithSize(s unit.Dp) *RadialProgress {
	r.Size = s
	return r
}

// WithThick sets the ring thickness in dp.
func (r *RadialProgress) WithThick(t unit.Dp) *RadialProgress {
	r.Thick = t
	return r
}

// WithLabel sets a centered label string.
func (r *RadialProgress) WithLabel(label string) *RadialProgress {
	r.Label = label
	return r
}

// WithVariant sets the progress color variant.
func (r *RadialProgress) WithVariant(v ProgressVariant) *RadialProgress {
	r.Variant = v
	return r
}

func (r *RadialProgress) fillColor() color.NRGBA {
	th := r.th
	switch r.Variant {
	case ProgressSecondary:
		return th.Secondary
	case ProgressAccent:
		return th.Accent
	case ProgressInfo:
		return th.Info
	case ProgressSuccess:
		return th.Success
	case ProgressWarning:
		return th.Warning
	case ProgressError:
		return th.Error
	default:
		return th.Primary
	}
}

// Layout renders the circular progress ring.
func (r *RadialProgress) Layout(gtx layout.Context) layout.Dimensions {
	th := r.th
	size := gtx.Dp(r.Size)
	thick := gtx.Dp(r.Thick)
	sz := image.Pt(size, size)
	accent := r.fillColor()

	cx := float32(size) / 2
	cy := float32(size) / 2
	outerR := cx
	innerR := cx - float32(thick)

	// Draw track (full circle)
	drawArcSegment(gtx, cx, cy, outerR, innerR, -math.Pi/2, 2*math.Pi, theme.WithAlpha(accent, 30))

	// Draw fill arc
	if r.Value > 0 {
		sweep := float64(r.Value) * 2 * math.Pi
		drawArcSegment(gtx, cx, cy, outerR, innerR, -math.Pi/2, sweep, accent)
	}

	// Center label
	if r.Label != "" {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, r.Label, th.BaseContent, th.SmSize, font.Bold)
			}),
		)
	}

	return layout.Dimensions{Size: sz}
}

// drawArcSegment draws a filled donut-slice (arc from startAngle, sweeping sweepAngle radians).
// Angles are in radians; 0 = right, -pi/2 = top (12 o'clock).
func drawArcSegment(gtx layout.Context, cx, cy, outerR, innerR float32, startAngle, sweepAngle float64, col color.NRGBA) {
	const segments = 64
	step := sweepAngle / segments

	var p clip.Path
	p.Begin(gtx.Ops)

	// Start point on outer arc
	a0 := startAngle
	x0 := cx + outerR*float32(math.Cos(a0))
	y0 := cy + outerR*float32(math.Sin(a0))
	p.MoveTo(f32.Pt(x0, y0))

	// Outer arc
	for i := 1; i <= segments; i++ {
		a := startAngle + step*float64(i)
		x := cx + outerR*float32(math.Cos(a))
		y := cy + outerR*float32(math.Sin(a))
		p.LineTo(f32.Pt(x, y))
	}

	// Inner arc (reverse)
	for i := segments; i >= 0; i-- {
		a := startAngle + step*float64(i)
		x := cx + innerR*float32(math.Cos(a))
		y := cy + innerR*float32(math.Sin(a))
		p.LineTo(f32.Pt(x, y))
	}

	p.Close()
	paint.FillShape(gtx.Ops, col, clip.Outline{Path: p.End()}.Op())
}
