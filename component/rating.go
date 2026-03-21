package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Rating is a DaisyUI-style star rating component.
type Rating struct {
	Value   int // 1..Max (0 = none)
	Max     int
	Variant BtnVariant
	clicks  []widget.Clickable
	th      *theme.Theme
}

// NewRating creates a new rating component with max stars.
func NewRating(th *theme.Theme, max int) *Rating {
	if max <= 0 {
		max = 5
	}
	return &Rating{
		Max:     max,
		clicks:  make([]widget.Clickable, max),
		Variant: BtnWarning,
		th:      th,
	}
}

// WithVariant sets the star color variant.
func (r *Rating) WithVariant(v BtnVariant) *Rating {
	r.Variant = v
	return r
}

func (r *Rating) starColor() color.NRGBA {
	th := r.th
	switch r.Variant {
	case BtnPrimary:
		return th.Primary
	case BtnSecondary:
		return th.Secondary
	case BtnAccent:
		return th.Accent
	case BtnInfo:
		return th.Info
	case BtnSuccess:
		return th.Success
	case BtnError:
		return th.Error
	default: // BtnWarning
		return th.Warning
	}
}

// Layout renders the star rating.
func (r *Rating) Layout(gtx layout.Context) layout.Dimensions {
	th := r.th
	starSize := gtx.Dp(28)
	gap := gtx.Dp(2)
	accent := r.starColor()

	// Handle clicks
	for i := range r.clicks {
		if r.clicks[i].Clicked(gtx) {
			if r.Value == i+1 {
				r.Value = 0 // clicking same star deselects
			} else {
				r.Value = i + 1
			}
		}
	}

	totalW := r.Max*starSize + (r.Max-1)*gap
	totalH := starSize

	children := make([]layout.FlexChild, r.Max)
	for i := 0; i < r.Max; i++ {
		i := i
		if i > 0 {
			// handled by spacing in layout
		}
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			var inset layout.Inset
			if i > 0 {
				inset.Left = 2 // gap in dp
			}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				sz := image.Pt(starSize, starSize)
				gtx.Constraints = layout.Exact(sz)
				return r.clicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					filled := i < r.Value
					col := th.Base300
					if filled {
						col = accent
					}
					// Draw star "★"
					pointer.CursorPointer.Add(gtx.Ops)
					return layout.Stack{Alignment: layout.Center}.Layout(gtx,
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{Size: sz}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return drawStar(gtx, th, col)
						}),
					)
				})
			})
		})
	}

	_ = totalW
	_ = totalH

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx, children...)
}

// drawStar renders a star glyph centered in the given size.
func drawStar(gtx layout.Context, th *theme.Theme, col color.NRGBA) layout.Dimensions {
	return drawText(gtx, th, "★", col, 20, font.Normal)
}
