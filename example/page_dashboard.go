package main

import (
	"image/color"

	"gioui.org/layout"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/scaffold"
	"github.com/hongshengjie/gioui-kit/theme"
)

// ─── Page: Dashboard ────────────────────────────────────────────────────────

func (a *App) pageDashboard(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 32}.Layout(gtx,
		// Hero
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return kit.FlexCol{Gap: 8}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return scaffold.NewBreadcrumb(th, "Home", "Dashboard").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, "Welcome to GioUI Kit").H1().Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th,
						"A comprehensive component library for Gio, inspired by TailwindCSS and DaisyUI.",
					).Sm().WithColor(theme.Gray500).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewButton(th, &a.btnOverview1, "Browse Components").WithVariant(component.BtnPrimary).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewButton(th, &a.btnOverview2, "View Source").WithVariant(component.BtnOutline).Layout(gtx)
							}),
						)
					})
				}),
			)
		}),

		// Stat cards — 1 col mobile → 2 sm → 4 md+
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return kit.Grid{Cols: 1, SmCols: 2, MdCols: 4, Gap: 16}.Layout(gtx,
				statCard(th, "12", "Components", theme.Blue500),
				statCard(th, "4", "Themes", theme.Purple500),
				statCard(th, "8", "Layout Types", theme.Emerald500),
				statCard(th, "6", "Form Controls", theme.Amber500),
			)
		}),

		// Recent activity
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Recent Activity", "Latest component updates", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 12}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAlert(th, "Button layout fixed — text no longer clipped by rounded corners.", component.AlertSuccess).WithIcon(iconSuccess).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAlert(th, "Modal backdrop click now closes the dialog.", component.AlertInfo).WithIcon(iconInfo).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAlert(th, "Toast auto-dismisses after 3 seconds.", component.AlertInfo).WithIcon(iconInfo).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAlert(th, "Memory explosion bug fixed — clip.Stroke replaced with paint-over.", component.AlertWarning).WithIcon(iconWarning).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Quick badges
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Tech Stack", "Built with", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewBadge(th, "Go 1.22+").WithVariant(component.BadgePrimary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewBadge(th, "Gio v0.9").WithVariant(component.BadgeAccent).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewBadge(th, "TailwindCSS").WithVariant(component.BadgeInfo).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewBadge(th, "DaisyUI").WithVariant(component.BadgeSecondary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewBadge(th, "MIT License").WithVariant(component.BadgeSuccess).Layout(gtx)
					}),
				)
			})(gtx)
		}),
	)
}

func statCard(th *theme.Theme, value, label string, accent color.NRGBA) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		card := component.NewCard(th).WithBorder()
		return card.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return kit.FlexCol{Gap: 4}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, value).H1().WithColor(accent).Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, label).Sm().WithColor(theme.Gray500).Layout(gtx)
				}),
			)
		})
	}
}
