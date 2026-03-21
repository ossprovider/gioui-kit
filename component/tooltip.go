package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/hongshengjie/gioui-kit/theme"
)

// TooltipPosition defines where the tooltip appears relative to the child.
type TooltipPosition int

const (
	TooltipTop    TooltipPosition = iota
	TooltipBottom
	TooltipLeft
	TooltipRight
)

// Tooltip wraps a widget and shows a text tooltip on hover.
type Tooltip struct {
	Label    string
	Position TooltipPosition
	bg       color.NRGBA
	hover    gesture.Hover
	th       *theme.Theme
}

// NewTooltip creates a new tooltip.
func NewTooltip(th *theme.Theme, label string) *Tooltip {
	return &Tooltip{
		Label:    label,
		Position: TooltipTop,
		bg:       th.Neutral,
		th:       th,
	}
}

// WithPosition sets the tooltip position.
func (t *Tooltip) WithPosition(p TooltipPosition) *Tooltip {
	t.Position = p
	return t
}

// WithBg sets the tooltip background color.
func (t *Tooltip) WithBg(c color.NRGBA) *Tooltip {
	t.bg = c
	return t
}

// Layout renders the child widget and its tooltip overlay when hovered.
func (t *Tooltip) Layout(gtx layout.Context, child layout.Widget) layout.Dimensions {
	// Render child first to get dimensions
	childDims := child(gtx)

	// Register hover input over the child area
	func() {
		defer clip.Rect{Max: childDims.Size}.Push(gtx.Ops).Pop()
		t.hover.Add(gtx.Ops)
	}()
	hovered := t.hover.Update(gtx.Source)

	// Render tooltip if hovered
	if hovered && t.Label != "" {
		// Record tooltip drawing
		macro := op.Record(gtx.Ops)
		tipDims := t.renderTip(gtx)
		tipCall := macro.Stop()

		// Determine offset
		var offset image.Point
		switch t.Position {
		case TooltipBottom:
			offset = image.Pt((childDims.Size.X-tipDims.Size.X)/2, childDims.Size.Y+gtx.Dp(4))
		case TooltipLeft:
			offset = image.Pt(-tipDims.Size.X-gtx.Dp(4), (childDims.Size.Y-tipDims.Size.Y)/2)
		case TooltipRight:
			offset = image.Pt(childDims.Size.X+gtx.Dp(4), (childDims.Size.Y-tipDims.Size.Y)/2)
		default: // TooltipTop
			offset = image.Pt((childDims.Size.X-tipDims.Size.X)/2, -tipDims.Size.Y-gtx.Dp(4))
		}

		// Draw tooltip on top via Defer — all push/pop must stay inside the macro
		macro2 := op.Record(gtx.Ops)
		stack := op.Offset(offset).Push(gtx.Ops)
		tipCall.Add(gtx.Ops)
		stack.Pop()
		op.Defer(gtx.Ops, macro2.Stop())
	}

	return childDims
}

func (t *Tooltip) renderTip(gtx layout.Context) layout.Dimensions {
	th := t.th
	padding := layout.Inset{Top: th.Space1, Bottom: th.Space1, Left: th.Space2, Right: th.Space2}
	radius := gtx.Dp(th.RoundedMd)
	bg := t.bg
	fg := th.NeutralContent

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			if sz.X > 0 && sz.Y > 0 {
				defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: bg}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
			}
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, t.Label, fg, th.XsSize, font.Normal)
			})
		}),
	)
}
