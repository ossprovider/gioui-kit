package component

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Steps is a DaisyUI-style steps progress indicator.
type Steps struct {
	Items   []string
	Current int // index of current step (0-based), -1 = none completed
	Variant BtnVariant
	th      *theme.Theme
}

// NewSteps creates a new steps component.
func NewSteps(th *theme.Theme, items []string) *Steps {
	return &Steps{Items: items, Current: 0, Variant: BtnPrimary, th: th}
}

// WithCurrent sets the current active step index.
func (s *Steps) WithCurrent(i int) *Steps {
	s.Current = i
	return s
}

// WithVariant sets the accent color for completed steps.
func (s *Steps) WithVariant(v BtnVariant) *Steps {
	s.Variant = v
	return s
}

func (s *Steps) accentColor() color.NRGBA {
	th := s.th
	switch s.Variant {
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

// Layout renders the steps indicator horizontally.
func (s *Steps) Layout(gtx layout.Context) layout.Dimensions {
	th := s.th
	circleSize := gtx.Dp(32)
	lineH := gtx.Dp(2)
	accent := s.accentColor()
	accentContent := th.PrimaryContent
	switch s.Variant {
	case BtnSecondary:
		accentContent = th.SecondaryContent
	case BtnAccent:
		accentContent = th.AccentContent
	case BtnInfo:
		accentContent = th.InfoContent
	case BtnSuccess:
		accentContent = th.SuccessContent
	case BtnWarning:
		accentContent = th.WarningContent
	case BtnError:
		accentContent = th.ErrorContent
	}

	children := make([]layout.FlexChild, len(s.Items))
	for i, item := range s.Items {
		i, item := i, item
		isDone := i <= s.Current
		isActive := i == s.Current

		children[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				// Circle + connector line row
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						// Left line (not for first item)
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if i == 0 {
								return layout.Dimensions{Size: image.Pt(0, circleSize)}
							}
							w := gtx.Constraints.Max.X
							if w <= 0 {
								return layout.Dimensions{}
							}
							col := th.Base300
							if i <= s.Current {
								col = accent
							}
							rect := image.Rect(0, (circleSize-lineH)/2, w, (circleSize-lineH)/2+lineH)
							defer clip.Rect(rect).Push(gtx.Ops).Pop()
							paint.ColorOp{Color: col}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
							return layout.Dimensions{Size: image.Pt(w, circleSize)}
						}),
						// Circle
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							sz := image.Pt(circleSize, circleSize)
							gtx.Constraints = layout.Exact(sz)
							rect := image.Rectangle{Max: sz}
							circleBg := th.Base300
							circleText := th.BaseContent
							if isDone {
								circleBg = accent
								circleText = accentContent
							}
							defer clip.UniformRRect(rect, circleSize/2).Push(gtx.Ops).Pop()
							paint.ColorOp{Color: circleBg}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
							_ = isActive
							return layout.Stack{Alignment: layout.Center}.Layout(gtx,
								layout.Expanded(func(gtx layout.Context) layout.Dimensions {
									return layout.Dimensions{Size: sz}
								}),
								layout.Stacked(func(gtx layout.Context) layout.Dimensions {
									label := fmt.Sprintf("%d", i+1)
									if isDone && i < s.Current {
										label = "✓"
									}
									return drawText(gtx, th, label, circleText, th.SmSize, font.Bold)
								}),
							)
						}),
						// Right line (not for last item)
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							if i == len(s.Items)-1 {
								return layout.Dimensions{Size: image.Pt(0, circleSize)}
							}
							w := gtx.Constraints.Max.X
							if w <= 0 {
								return layout.Dimensions{}
							}
							col := th.Base300
							if i < s.Current {
								col = accent
							}
							rect := image.Rect(0, (circleSize-lineH)/2, w, (circleSize-lineH)/2+lineH)
							defer clip.Rect(rect).Push(gtx.Ops).Pop()
							paint.ColorOp{Color: col}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
							return layout.Dimensions{Size: image.Pt(w, circleSize)}
						}),
					)
				}),
				// Label
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: th.Space1}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						col := th.BaseContent
						if !isDone {
							col = theme.Opacity(th.BaseContent, 0.5)
						}
						w := font.Normal
						if isActive {
							w = font.SemiBold
						}
						return drawText(gtx, th, item, col, th.XsSize, w)
					})
				}),
			)
		})
	}

	return layout.Flex{}.Layout(gtx, children...)
}
