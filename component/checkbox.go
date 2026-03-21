package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Checkbox is a DaisyUI-style checkbox component.
type Checkbox struct {
	Bool     *widget.Bool
	Label    string
	Variant  BtnVariant
	Disabled bool
	th       *theme.Theme
}

// NewCheckbox creates a new checkbox with theme, state, and label.
func NewCheckbox(th *theme.Theme, b *widget.Bool, label string) *Checkbox {
	return &Checkbox{Bool: b, Label: label, Variant: BtnPrimary, th: th}
}

// WithVariant sets the checkbox accent color variant.
func (c *Checkbox) WithVariant(v BtnVariant) *Checkbox {
	c.Variant = v
	return c
}

func (c *Checkbox) colors() (fill, check color.NRGBA) {
	th := c.th
	switch c.Variant {
	case BtnSecondary:
		return th.Secondary, th.SecondaryContent
	case BtnAccent:
		return th.Accent, th.AccentContent
	case BtnInfo:
		return th.Info, th.InfoContent
	case BtnSuccess:
		return th.Success, th.SuccessContent
	case BtnWarning:
		return th.Warning, th.WarningContent
	case BtnError:
		return th.Error, th.ErrorContent
	default:
		return th.Primary, th.PrimaryContent
	}
}

// Layout renders the checkbox.
func (c *Checkbox) Layout(gtx layout.Context) layout.Dimensions {
	th := c.th
	boxSize := gtx.Dp(20)
	radius := gtx.Dp(th.RoundedMd)
	fillCol, checkCol := c.colors()

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			sz := image.Pt(boxSize, boxSize)
			gtx.Constraints = layout.Exact(sz)
			return c.Bool.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				rect := image.Rectangle{Max: sz}
				if c.Bool.Value {
					defer clip.UniformRRect(rect, radius).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: fillCol}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					// Centered checkmark
					return layout.Stack{Alignment: layout.Center}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: sz}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return drawText(gtx, c.th, "✓", checkCol, c.th.XsSize, font.Bold)
						}),
					)
				}
				// Unchecked: border box
				defer clip.UniformRRect(rect, radius).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				paint.FillShape(gtx.Ops, th.Base300,
					clip.Stroke{
						Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
						Width: float32(gtx.Dp(2)),
					}.Op(),
				)
				pointer.CursorPointer.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if c.Label == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Left: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				col := th.BaseContent
				if c.Disabled {
					col = theme.Opacity(col, 0.5)
				}
				return drawText(gtx, th, c.Label, col, th.FontSize, font.Normal)
			})
		}),
	)
}
