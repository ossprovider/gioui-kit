package scaffold

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

type BottomNavItem struct {
	Label    string
	Icon     string       // Unicode icon (used when IconData is nil)
	IconData *widget.Icon // Material Design iconvg icon (takes priority over Icon)
	Active   bool
	click    widget.Clickable
}

// BottomNav renders a mobile-style bottom navigation bar.
type BottomNav struct {
	Items    []BottomNavItem
	OnSelect func(index int)
	children []layout.FlexChild // reused across frames
	th       *theme.Theme
}

func NewBottomNav(th *theme.Theme, items []BottomNavItem) *BottomNav {
	return &BottomNav{Items: items, th: th}
}

func (bn *BottomNav) Layout(gtx layout.Context) layout.Dimensions {
	th := bn.th
	h := gtx.Dp(64)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := image.Pt(gtx.Constraints.Max.X, h)
			// Background
			paint.FillShape(gtx.Ops, th.Base100,
				clip.Rect{Max: sz}.Op(),
			)
			// Top border
			borderRect := image.Rect(0, 0, sz.X, 1)
			paint.FillShape(gtx.Ops, th.Base300,
				clip.Rect(borderRect).Op(),
			)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = h
			gtx.Constraints.Max.Y = h

			if cap(bn.children) < len(bn.Items) {
				bn.children = make([]layout.FlexChild, len(bn.Items))
			}
			bn.children = bn.children[:len(bn.Items)]
			for i := range bn.Items {
				i := i
				bn.children[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					item := &bn.Items[i]
					if item.click.Clicked(gtx) && bn.OnSelect != nil {
						bn.OnSelect(i)
					}

					return item.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						fg := th.BaseContent
						if item.Active {
							fg = th.Primary
						}
						pointer.CursorPointer.Add(gtx.Ops)
						return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{
								Axis:      layout.Vertical,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									if item.IconData != nil {
										iconSz := gtx.Dp(24)
										gtx.Constraints = layout.Exact(image.Pt(iconSz, iconSz))
										return item.IconData.Layout(gtx, fg)
									}
									return drawLabel(gtx, th, item.Icon, fg, unit.Sp(20), font.Normal)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return layout.Inset{Top: 2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return drawLabel(gtx, th, item.Label, fg, th.XsSize, font.Medium)
									})
								}),
							)
						})
					})
				})
			}

			return layout.Flex{Alignment: layout.Middle}.Layout(gtx, bn.children...)
		}),
	)
}
