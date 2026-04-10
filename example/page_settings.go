package main

import (
	"gioui.org/layout"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/theme"
)

// ─── Page: Settings ─────────────────────────────────────────────────────────

func (a *App) pageSettings(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 32}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return pageHeader(th, gtx, "Settings", "Dashboard / Settings",
				"Customize the application theme and appearance.")
		}),

		// Theme picker
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Theme", "Switch between built-in themes", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Select a theme to apply it globally.").Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewButton(th, &a.btnLight, "Light").WithVariant(component.BtnOutline).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewButton(th, &a.btnDark, "Dark").WithVariant(component.BtnOutline).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewButton(th, &a.btnCupcake, "Cupcake").WithVariant(component.BtnOutline).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewButton(th, &a.btnNord, "Nord").WithVariant(component.BtnOutline).Layout(gtx)
							}),
						)
					}),
				)
			})(gtx)
		}),

		// Typography preview
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Typography", "Font size scale", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 8}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Heading 1 — H1").H1().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Heading 2 — H2").H2().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Heading 3 — H3").H3().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Heading 4 — H4").H4().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Body — default font size").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Small — sm size, muted").Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "XSmall — xs size, muted").Xs().WithColor(theme.Gray400).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Color palette preview
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Color Palette", "Current theme semantic colors", func(gtx layout.Context) layout.Dimensions {
				return kit.Grid{Cols: 2, MdCols: 4, Gap: 12}.Layout(gtx,
					colorSwatch(th, "Primary", th.Primary, th.PrimaryContent),
					colorSwatch(th, "Secondary", th.Secondary, th.SecondaryContent),
					colorSwatch(th, "Accent", th.Accent, th.AccentContent),
					colorSwatch(th, "Neutral", th.Neutral, th.NeutralContent),
					colorSwatch(th, "Info", th.Info, th.InfoContent),
					colorSwatch(th, "Success", th.Success, th.SuccessContent),
					colorSwatch(th, "Warning", th.Warning, th.WarningContent),
					colorSwatch(th, "Error", th.Error, th.ErrorContent),
				)
			})(gtx)
		}),
	)
}
