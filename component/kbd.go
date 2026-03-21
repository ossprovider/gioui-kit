package component

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Kbd renders a DaisyUI-style keyboard key component.
type Kbd struct {
	Key string
	th  *theme.Theme
}

// NewKbd creates a new keyboard key display.
func NewKbd(th *theme.Theme, key string) *Kbd {
	return &Kbd{Key: key, th: th}
}

// Layout renders the keyboard key.
func (k *Kbd) Layout(gtx layout.Context) layout.Dimensions {
	th := k.th
	radius := gtx.Dp(th.RoundedMd)
	padding := layout.Inset{
		Top: th.Space1, Bottom: th.Space1,
		Left: th.Space2, Right: th.Space2,
	}
	shadowH := gtx.Dp(2)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			// Shadow / bottom border
			shadowRect := image.Rectangle{Max: image.Pt(sz.X, sz.Y+shadowH)}
			defer clip.UniformRRect(shadowRect, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			// Main surface
			mainRect := image.Rectangle{Max: sz}
			defer clip.UniformRRect(mainRect, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base200}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, k.Key, th.BaseContent, th.SmSize, font.SemiBold)
			})
		}),
	)
}

// KbdGroup renders a sequence of keyboard keys (e.g. Ctrl+C).
func KbdGroup(th *theme.Theme, keys ...string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		children := make([]layout.FlexChild, 0, len(keys)*2-1)
		for i, key := range keys {
			key := key
			if i > 0 {
				children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: 2, Right: 2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return drawText(gtx, th, "+", theme.Opacity(th.BaseContent, 0.5), th.SmSize, font.Normal)
					})
				}))
			}
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return NewKbd(th, key).Layout(gtx)
			}))
		}
		return layout.Flex{Alignment: layout.Middle}.Layout(gtx, children...)
	}
}
