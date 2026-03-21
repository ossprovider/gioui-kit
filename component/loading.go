package component

import (
	"image"
	"image/color"
	"math"
	"time"

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
	LoadingSpinner LoadingVariant = iota // rotating ring
	LoadingDots                          // three bouncing dots
	LoadingRing                          // thin ring
)

// Loading is a DaisyUI-style loading indicator.
type Loading struct {
	Variant  LoadingVariant
	Color    color.NRGBA
	Size     unit.Dp
	th       *theme.Theme
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
	// Invalidate every frame to animate
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

	// Animate rotation based on time
	t := float64(gtx.Now.UnixNano()) / float64(time.Second)
	angle := t * 2 * math.Pi // one full rotation per second

	// Draw arc segments to simulate spinner
	// We draw 8 dots around a circle with varying opacity
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
		defer clip.UniformRRect(rect, dotSize/2).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: col}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
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

	// Outer ring (faded track)
	track := image.Rectangle{Max: sz}
	paint.FillShape(gtx.Ops, theme.WithAlpha(l.Color, 30),
		clip.Stroke{
			Path:  clip.UniformRRect(track, size/2).Path(gtx.Ops),
			Width: float32(thick),
		}.Op(),
	)

	// Animate: rotating highlight segment (approximated by varying top arc)
	t := float64(gtx.Now.UnixNano()) / float64(time.Second)
	_ = t

	// Draw a quarter arc highlight at the top by masking with a rect
	// Simple approach: draw a colored arc using clipped rectangle
	innerSize := size - thick*2
	innerOff := thick
	innerRect := image.Rect(innerOff, innerOff, innerOff+innerSize, innerOff+innerSize)
	_ = innerRect

	// Full ring colored at primary (animated via rotation via ops.Transform not available here)
	// Fall back to just a full colored ring for simplicity
	paint.FillShape(gtx.Ops, l.Color,
		clip.Stroke{
			Path:  clip.UniformRRect(track, size/2).Path(gtx.Ops),
			Width: float32(thick),
		}.Op(),
	)

	return layout.Dimensions{Size: sz}
}

func (l *Loading) layoutDots(gtx layout.Context) layout.Dimensions {
	dotSize := gtx.Dp(l.Size / 3)
	gap := gtx.Dp(4)
	totalW := 3*dotSize + 2*gap
	totalH := dotSize

	t := float64(gtx.Now.UnixNano()) / float64(time.Second)

	for i := 0; i < 3; i++ {
		// Each dot bobs up/down with a phase offset
		phase := t*4 - float64(i)*0.4
		offset := int(float32(dotSize/3) * float32(math.Sin(phase)))
		x := i * (dotSize + gap)
		y := dotSize/2 - dotSize/2 + offset
		if y < 0 {
			y = 0
		}
		rect := image.Rect(x, y, x+dotSize, y+dotSize)
		defer clip.UniformRRect(rect, dotSize/2).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: image.Pt(totalW, totalH)}
}
