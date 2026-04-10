package main

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/scaffold"
	"github.com/hongshengjie/gioui-kit/theme"
)

func pageHeader(th *theme.Theme, gtx layout.Context, title, breadcrumb, subtitle string) layout.Dimensions {
	return kit.FlexCol{Gap: 8}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			parts := []string{"Home"}
			if breadcrumb != "" {
				parts = append(parts, title)
			}
			return scaffold.NewBreadcrumb(th, parts...).Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.NewText(th, title).H1().Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.NewText(th, subtitle).Sm().WithColor(theme.Gray500).Layout(gtx)
		}),
	)
}

// drawerNavItem renders a drawer navigation button with a real icon + label.
func drawerNavItem(th *theme.Theme, click *widget.Clickable, icon *widget.Icon, label string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			hovered := click.Hovered()
			col := th.BaseContent
			if hovered {
				col = th.Primary
			}
			return layout.Inset{Top: th.Space2, Bottom: th.Space2, Left: th.Space3, Right: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						iconSz := gtx.Dp(20)
						gtx.Constraints = layout.Exact(image.Pt(iconSz, iconSz))
						return icon.Layout(gtx, col)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Left: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, label).WithColor(col).Layout(gtx)
						})
					}),
				)
			})
		})
	}
}

func subSection(th *theme.Theme, gtx layout.Context, label string, body layout.Widget) layout.Dimensions {
	return kit.FlexCol{Gap: 8}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return component.NewText(th, label).H4().Layout(gtx)
		}),
		layout.Rigid(body),
	)
}

func sectionCard(th *theme.Theme, title, subtitle string, body layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 12}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 4}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, title).H2().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, subtitle).Sm().WithColor(theme.Gray400).Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewCard(th).WithBorder().Layout(gtx, body)
			}),
		)
	}
}

func colorSwatch(th *theme.Theme, name string, bg, fg color.NRGBA) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		card := component.NewCard(th)
		return card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return kit.Box{
				Background: bg,
				Radius:     th.RoundedMd,
				Padding:    layout.Inset{Top: th.Space3, Bottom: th.Space3, Left: th.Space3, Right: th.Space3},
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, name).Sm().WithColor(fg).Layout(gtx)
			})
		})
	}
}
