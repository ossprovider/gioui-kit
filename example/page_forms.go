package main

import (
	"fmt"

	"gioui.org/layout"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/theme"
)

// ─── Page: Forms ────────────────────────────────────────────────────────────

func (a *App) pageForms(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 32}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return pageHeader(th, gtx, "Forms", "Dashboard / Forms",
				"Input fields, toggles, and interactive controls.")
		}),

		// Inputs
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Text Inputs", "All input variants", func(gtx layout.Context) layout.Dimensions {
				return kit.Grid{Cols: 1, MdCols: 2, Gap: 20}.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return kit.FlexCol{Gap: 16}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewInput(th, &a.editor1, "Default placeholder...").WithLabel("Default Input").Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewInput(th, &a.editor2, "Primary placeholder...").WithLabel("Primary Input").WithVariant(component.InputPrimary).Layout(gtx)
							}),
						)
					},
					func(gtx layout.Context) layout.Dimensions {
						return kit.FlexCol{Gap: 16}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewInput(th, &a.editor3, "Error state...").WithLabel("Error Input").WithVariant(component.InputError).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewInput(th, &a.editor4, "Ghost style (no border)...").WithLabel("Ghost Input").WithVariant(component.InputGhost).Layout(gtx)
							}),
						)
					},
				)
			})(gtx)
		}),

		// Toggles
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Toggle Switches", "Boolean controls", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 12}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewToggle(th, &a.toggle1, "Enable notifications (default on)").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewToggle(th, &a.toggle2, "Dark mode").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewToggle(th, &a.toggle3, "Auto-save drafts").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						state := "off"
						if a.toggle1.Value {
							state = "on"
						}
						return component.NewText(th, fmt.Sprintf("Notifications: %s", state)).Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Overlays demo
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Overlays", "Modal, Drawer, and Toast", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Modal").H4().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Click backdrop or the Close button to dismiss.").Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnModal, "Open Modal").WithVariant(component.BtnPrimary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.DividerH{Color: th.Base300}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Drawer").H4().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "A slide-in panel from the left edge.").Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnDrawer, "Open Drawer").WithVariant(component.BtnSecondary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.DividerH{Color: th.Base300}.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Toast").H4().Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Auto-dismisses after 3 seconds.").Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnToast, "Show Toast").WithVariant(component.BtnSuccess).Layout(gtx)
					}),
				)
			})(gtx)
		}),
	)
}
