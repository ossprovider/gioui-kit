package main

import (
	"image/color"

	"gioui.org/layout"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/theme"
)

// ─── Page: Layout ───────────────────────────────────────────────────────────

func (a *App) pageLayout(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 32}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return pageHeader(th, gtx, "Layout", "Dashboard / Layout",
				"Grid and Flex layout primitives inspired by TailwindCSS.")
		}),

		// Cards section
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Cards", "Content containers", func(gtx layout.Context) layout.Dimensions {
				return kit.Grid{Cols: 1, MdCols: 2, LgCols: 3, Gap: 16}.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return component.NewCard(th).WithBorder().CardWithHeader(gtx, "Bordered Card",
							func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "This card has a visible border. Suitable for lighter backgrounds.").Sm().WithColor(theme.Gray500).Layout(gtx)
							})
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewCard(th).CardWithHeader(gtx, "Default Card",
							func(gtx layout.Context) layout.Dimensions {
								return kit.FlexCol{Gap: 12}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return component.NewText(th, "Cards can contain any content including nested components.").Sm().WithColor(theme.Gray500).Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return component.NewButton(th, &a.btnCard1, "Action").WithVariant(component.BtnPrimary).WithSize(component.BtnSm).Layout(gtx)
									}),
								)
							})
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewCard(th).WithCompact().CardWithHeader(gtx, "Compact Card",
							func(gtx layout.Context) layout.Dimensions {
								return kit.FlexCol{Gap: 8}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return component.NewText(th, "Compact variant with less padding.").Sm().WithColor(theme.Gray500).Layout(gtx)
									}),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										return kit.FlexRow{Gap: 4}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return component.NewBadge(th, "Go").WithVariant(component.BadgePrimary).Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return component.NewBadge(th, "Gio").WithVariant(component.BadgeAccent).Layout(gtx)
											}),
										)
									}),
								)
							})
					},
				)
			})(gtx)
		}),

		// Grid examples
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Grid Layout", "grid-cols-1 / grid-cols-2 / grid-cols-3", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 20}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "1 Column", func(gtx layout.Context) layout.Dimensions {
							return kit.Grid{Cols: 1, Gap: 8}.Layout(gtx,
								gridBox(th, "Col 1 / 1", theme.Blue100),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "2 Columns", func(gtx layout.Context) layout.Dimensions {
							return kit.Grid{Cols: 2, Gap: 8}.Layout(gtx,
								gridBox(th, "Col 1 / 2", theme.Blue100),
								gridBox(th, "Col 2 / 2", theme.Blue200),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "3 Columns", func(gtx layout.Context) layout.Dimensions {
							return kit.Grid{Cols: 3, Gap: 8}.Layout(gtx,
								gridBox(th, "Col 1 / 3", theme.Indigo100),
								gridBox(th, "Col 2 / 3", theme.Indigo200),
								gridBox(th, "Col 3 / 3", theme.Indigo300),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "4 Columns", func(gtx layout.Context) layout.Dimensions {
							return kit.Grid{Cols: 4, Gap: 8}.Layout(gtx,
								gridBox(th, "1/4", theme.Purple100),
								gridBox(th, "2/4", theme.Purple200),
								gridBox(th, "3/4", theme.Purple300),
								gridBox(th, "4/4", theme.Purple400),
							)
						})
					}),
				)
			})(gtx)
		}),

		// Flex layout examples
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Flex Layout", "FlexRow and FlexCol with gap and alignment", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 20}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "FlexRow — gap-8, items-center", func(gtx layout.Context) layout.Dimensions {
							return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Item A").WithVariant(component.BadgePrimary).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Item B").WithVariant(component.BadgeSecondary).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Item C").WithVariant(component.BadgeAccent).Layout(gtx)
								}),
								kit.Grow(func(gtx layout.Context) layout.Dimensions {
									return component.NewText(th, "← flex-1 spacer →").Sm().WithColor(theme.Gray400).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewButton(th, &a.btnCard3, "End").WithVariant(component.BtnOutline).WithSize(component.BtnSm).Layout(gtx)
								}),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "FlexCol — gap-8", func(gtx layout.Context) layout.Dimensions {
							return kit.FlexCol{Gap: 8}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewAlert(th, "Row 1 — stretched full width", component.AlertInfo).WithIcon(iconInfo).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewAlert(th, "Row 2 — stretched full width", component.AlertSuccess).WithIcon(iconSuccess).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewAlert(th, "Row 3 — stretched full width", component.AlertWarning).WithIcon(iconWarning).Layout(gtx)
								}),
							)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "Dividers & Spacers", func(gtx layout.Context) layout.Dimensions {
							return kit.FlexCol{Gap: 12}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewText(th, "Section A").Sm().WithColor(theme.Gray500).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return kit.DividerH{Color: th.Base300}.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewText(th, "Section B").Sm().WithColor(theme.Gray500).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return kit.DividerH{Color: th.Base300}.Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewText(th, "Section C").Sm().WithColor(theme.Gray500).Layout(gtx)
								}),
							)
						})
					}),
				)
			})(gtx)
		}),
	)
}

func gridBox(th *theme.Theme, label string, bg color.NRGBA) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return kit.Box{
			Background: bg,
			Radius:     th.RoundedMd,
			Padding:    layout.UniformInset(th.Space3),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return component.NewText(th, label).Sm().Layout(gtx)
		})
	}
}
