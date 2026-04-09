package component

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
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
					// Draw vector checkmark
					drawCheckmark(gtx.Ops, sz, checkCol)
					return layout.Dimensions{Size: sz}
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

// drawCheckmark draws a vector checkmark centered in sz using the given color.
func drawCheckmark(ops *op.Ops, sz image.Point, col color.NRGBA) {
	w := float32(sz.X)
	h := float32(sz.Y)
	// Checkmark points: short leg down-left, long leg up-right
	x1, y1 := w*0.20, h*0.50
	x2, y2 := w*0.42, h*0.72
	x3, y3 := w*0.80, h*0.28
	strokeWidth := w * 0.14

	var path clip.Path
	path.Begin(ops)
	path.MoveTo(f32.Pt(x1, y1))
	path.LineTo(f32.Pt(x2, y2))
	path.LineTo(f32.Pt(x3, y3))
	paint.FillShape(ops, col,
		clip.Stroke{Path: path.End(), Width: strokeWidth}.Op(),
	)
}
