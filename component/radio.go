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

// RadioGroup renders a group of DaisyUI-style radio buttons.
type RadioGroup struct {
	Items    []string
	Selected int // index of selected item
	Variant  BtnVariant
	clicks   []widget.Clickable
	th       *theme.Theme
}

// NewRadioGroup creates a new radio group.
func NewRadioGroup(th *theme.Theme, items []string) *RadioGroup {
	return &RadioGroup{
		Items:   items,
		clicks:  make([]widget.Clickable, len(items)),
		Variant: BtnPrimary,
		th:      th,
	}
}

// WithVariant sets the radio accent color.
func (r *RadioGroup) WithVariant(v BtnVariant) *RadioGroup {
	r.Variant = v
	return r
}

func (r *RadioGroup) accentColor() color.NRGBA {
	th := r.th
	switch r.Variant {
	case BtnSecondary:
		return th.Secondary
	case BtnAccent:
		return th.Accent
	case BtnInfo:
		return th.Info
	case BtnSuccess:
		return th.Success
	case BtnWarning:
		return th.Warning
	case BtnError:
		return th.Error
	default:
		return th.Primary
	}
}

// Layout renders the radio group vertically.
func (r *RadioGroup) Layout(gtx layout.Context) layout.Dimensions {
	th := r.th
	accent := r.accentColor()
	outerSize := gtx.Dp(20)
	innerSize := gtx.Dp(10)

	// Handle clicks
	for i := range r.clicks {
		if r.clicks[i].Clicked(gtx) {
			r.Selected = i
		}
	}

	children := make([]layout.FlexChild, len(r.Items))
	for i, item := range r.Items {
		i, item := i, item
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			isSelected := i == r.Selected
			return layout.Inset{Bottom: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						sz := image.Pt(outerSize, outerSize)
						gtx.Constraints = layout.Exact(sz)
						return r.clicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							rect := image.Rectangle{Max: sz}
							// Outer circle border
							borderCol := th.Base300
							if isSelected {
								borderCol = accent
							}
							defer clip.UniformRRect(rect, outerSize/2).Push(gtx.Ops).Pop()
							paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
							paint.FillShape(gtx.Ops, borderCol,
								clip.Stroke{
									Path:  clip.UniformRRect(rect, outerSize/2).Path(gtx.Ops),
									Width: float32(gtx.Dp(2)),
								}.Op(),
							)
							// Inner dot when selected
							if isSelected {
								offset := (outerSize - innerSize) / 2
								innerRect := image.Rect(offset, offset, offset+innerSize, offset+innerSize)
								defer clip.UniformRRect(innerRect, innerSize/2).Push(gtx.Ops).Pop()
								paint.ColorOp{Color: accent}.Add(gtx.Ops)
								paint.PaintOp{}.Add(gtx.Ops)
							}
							pointer.CursorPointer.Add(gtx.Ops)
							return layout.Dimensions{Size: sz}
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							w := font.Normal
							col := th.BaseContent
							if isSelected {
								w = font.SemiBold
							}
							return drawText(gtx, th, item, col, th.FontSize, w)
						})
					}),
				)
			})
		})
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}
