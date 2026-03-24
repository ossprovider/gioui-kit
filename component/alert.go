package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

type AlertVariant int

const (
	AlertInfo AlertVariant = iota
	AlertSuccess
	AlertWarning
	AlertError
)

// Alert is a DaisyUI-style alert banner.
type Alert struct {
	Text    string
	Variant AlertVariant
	Icon    *widget.Icon
	th      *theme.Theme
}

func NewAlert(th *theme.Theme, text string, variant AlertVariant) *Alert {
	return &Alert{Text: text, Variant: variant, th: th}
}

// WithIcon sets a custom icon for the alert.
func (a *Alert) WithIcon(icon *widget.Icon) *Alert {
	a.Icon = icon
	return a
}

func (a *Alert) colors() (bg, fg, icon color.NRGBA) {
	th := a.th
	switch a.Variant {
	case AlertSuccess:
		return theme.WithAlpha(th.Success, 30), th.Success, th.Success
	case AlertWarning:
		return theme.WithAlpha(th.Warning, 30), th.Warning, th.Warning
	case AlertError:
		return theme.WithAlpha(th.Error, 30), th.Error, th.Error
	default:
		return theme.WithAlpha(th.Info, 30), th.Info, th.Info
	}
}

func (a *Alert) Layout(gtx layout.Context) layout.Dimensions {
	th := a.th
	bg, fg, _ := a.colors()
	radius := gtx.Dp(th.RoundedLg)
	padding := layout.Inset{
		Top: th.Space3, Bottom: th.Space3,
		Left: th.Space4, Right: th.Space4,
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			rrect := clip.UniformRRect(image.Rectangle{Max: sz}, radius)
			defer rrect.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: bg}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if a.Icon == nil {
							return layout.Dimensions{}
						}
						return layout.Inset{Right: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							iconPx := gtx.Sp(th.H3Size)
							gtx.Constraints = layout.Exact(image.Pt(iconPx, iconPx))
							return a.Icon.Layout(gtx, fg)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return drawText(gtx, th, a.Text, fg, th.SmSize, font.Medium)
					}),
				)
			})
		}),
	)
}
