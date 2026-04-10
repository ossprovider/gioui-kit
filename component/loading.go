package component

import (
	"image"
	"image/color"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"github.com/hongshengjie/gioui-kit/theme"
)

// LoadingVariant defines the visual style of the loading indicator.
type LoadingVariant int

const (
	LoadingSpinner LoadingVariant = iota // rotating ring of dots
	LoadingDots                          // three bouncing dots
	LoadingRing                          // thin spinning ring
)

// Loading is a DaisyUI-style loading indicator.
type Loading struct {
	Variant LoadingVariant
	Color   color.NRGBA
	Size    unit.Dp
	th      *theme.Theme
}

// NewLoading creates a new loading indicator.
func NewLoading(th *theme.Theme) *Loading {
	return &Loading{
		Variant: LoadingSpinner,
		Color:   th.Primary,
		Size:    32,
		th:      th,
	}
}

// WithVariant sets the loading style.
func (l *Loading) WithVariant(v LoadingVariant) *Loading {
	l.Variant = v
	return l
}

// WithColor sets the loading indicator color.
func (l *Loading) WithColor(c color.NRGBA) *Loading {
	l.Color = c
	return l
}

// WithSize sets the size in dp.
func (l *Loading) WithSize(s unit.Dp) *Loading {
	l.Size = s
	return l
}

// Layout renders the loading indicator.
func (l *Loading) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Execute(op.InvalidateCmd{At: gtx.Now.Add(16 * time.Millisecond)})

	switch l.Variant {
	case LoadingDots:
		return l.layoutDots(gtx)
	case LoadingRing:
		return l.layoutRing(gtx)
	default:
		return l.layoutSpinner(gtx)
	}
}

func (l *Loading) layoutSpinner(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(l.Size)
	sz := image.Pt(size, size)

	t := float64(gtx.Now.UnixNano()) / float64(time.Second)
	angle := t * 2 * math.Pi

	numDots := 8
	dotSize := size / 6
	radius := (size - dotSize) / 2
	cx, cy := size/2, size/2

	for i := 0; i < numDots; i++ {
		a := angle + float64(i)*2*math.Pi/float64(numDots)
		x := cx + int(float64(radius)*math.Cos(a)) - dotSize/2
		y := cy + int(float64(radius)*math.Sin(a)) - dotSize/2
		opacity := float32(i+1) / float32(numDots)
		col := theme.Opacity(l.Color, opacity)
		rect := image.Rect(x, y, x+dotSize, y+dotSize)
		s := clip.UniformRRect(rect, dotSize/2).Push(gtx.Ops)
		paint.ColorOp{Color: col}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		s.Pop()
	}

	return layout.Dimensions{Size: sz}
}

func (l *Loading) layoutRing(gtx layout.Context) layout.Dimensions {
	size := gtx.Dp(l.Size)
	sz := image.Pt(size, size)
	thick := size / 8
	if thick < 2 {
		thick = 2
	}

	cx := float32(size) / 2
	cy := float32(size) / 2
	outerR := cx
	innerR := cx - float32(thick)

	// Draw faded track (full ring)
	drawArcSegment(gtx, cx, cy, outerR, innerR, -math.Pi/2, 2*math.Pi, theme.WithAlpha(l.Color, 40))

	// Rotate and draw 270-degree colored arc
	t := float64(gtx.Now.UnixNano()) / float64(time.Second)
	angle := float32(math.Mod(t, 1.0) * 2 * math.Pi)
	aff := f32.Affine2D{}.Rotate(f32.Pt(cx, cy), angle)
	stack := op.Affine(aff).Push(gtx.Ops)
	drawArcSegment(gtx, cx, cy, outerR, innerR, -math.Pi/2, 1.5*math.Pi, l.Color)
	stack.Pop()

	return layout.Dimensions{Size: sz}
}

func (l *Loading) layoutDots(gtx layout.Context) layout.Dimensions {
	dotSize := gtx.Dp(l.Size / 3)
	gap := gtx.Dp(4)
	totalW := 3*dotSize + 2*gap
	bounce := dotSize / 2
	totalH := dotSize + bounce

	t := float64(gtx.Now.UnixNano()) / float64(time.Second)

	for i := 0; i < 3; i++ {
		phase := t*3.0 - float64(i)*0.3
		// offset: 0 at rest, negative = up
		offset := int(float32(bounce/2) * float32(math.Sin(phase)))
		x := i * (dotSize + gap)
		y := bounce/2 - offset
		rect := image.Rect(x, y, x+dotSize, y+dotSize)
		s := clip.UniformRRect(rect, dotSize/2).Push(gtx.Ops)
		paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		s.Pop()
	}

	return layout.Dimensions{Size: image.Pt(totalW, totalH)}
}
