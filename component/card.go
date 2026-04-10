package component

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Card is a DaisyUI-style card container.
type Card struct {
	Bordered bool
	Compact  bool
	th       *theme.Theme
}

func NewCard(th *theme.Theme) *Card {
	return &Card{th: th}
}

func (c *Card) WithBorder() *Card {
	c.Bordered = true
	return c
}

func (c *Card) WithCompact() *Card {
	c.Compact = true
	return c
}

func (c *Card) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	th := c.th
	padding := layout.UniformInset(th.Space6)
	if c.Compact {
		padding = layout.UniformInset(th.Space4)
	}

	inner := func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				sz := gtx.Constraints.Min
				defer clip.UniformRRect(image.Rectangle{Max: sz}, gtx.Dp(th.RoundedXl)).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return padding.Layout(gtx, body)
			}),
		)
	}

	if c.Bordered {
		return widget.Border{Color: th.Base300, CornerRadius: th.RoundedXl, Width: 1}.Layout(gtx, inner)
	}
	return inner(gtx)
}

// CardWithHeader renders a card with a separate title section.
func (c *Card) CardWithHeader(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	th := c.th
	return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawText(gtx, th, title, th.BaseContent, th.H3Size, font.Bold)
				})
			}),
			layout.Rigid(body),
		)
	})
}
