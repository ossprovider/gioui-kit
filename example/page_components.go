package main

import (
	"fmt"

	"gioui.org/layout"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/modifier"
	"github.com/hongshengjie/gioui-kit/scaffold"
	"github.com/hongshengjie/gioui-kit/theme"
)

// ─── Page: Components ───────────────────────────────────────────────────────

func (a *App) pageComponents(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 32}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return pageHeader(th, gtx, "Components", "Dashboard / Components",
				"A complete showcase of all DaisyUI-inspired components.")
		}),

		// Sub-tabs
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.compTabs.Layout(gtx)
		}),

		// Tab content
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			switch a.compTabs.Selected() {
			case 1:
				return a.sectionBadgesChips(gtx)
			case 2:
				return a.sectionAlerts(gtx)
			case 3:
				return a.sectionAvatarsProgress(gtx)
			case 4:
				return a.sectionControls(gtx)
			case 5:
				return a.sectionDataDisplay(gtx)
			case 6:
				return a.sectionLayout(gtx)
			case 7:
				return a.sectionModifiers(gtx)
			default:
				return a.sectionButtons(gtx)
			}
		}),
	)
}

func (a *App) sectionButtons(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Buttons", "All variants, sizes, and states", func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 24}.Layout(gtx,
			// Variants
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return subSection(th, gtx, "Variants", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnPrimary, "Primary").WithVariant(component.BtnPrimary).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnSecondary, "Secondary").WithVariant(component.BtnSecondary).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnAccent, "Accent").WithVariant(component.BtnAccent).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnInfo, "Info").WithVariant(component.BtnInfo).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnSuccess, "Success").WithVariant(component.BtnSuccess).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnWarning, "Warning").WithVariant(component.BtnWarning).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnError, "Error").WithVariant(component.BtnError).Layout(gtx)
						}),
					)
				})
			}),
			// Ghost, Link, Outline
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return subSection(th, gtx, "Soft variants", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnOutline, "Outline").WithVariant(component.BtnOutline).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnGhost, "Ghost").WithVariant(component.BtnGhost).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnLink, "Link").WithVariant(component.BtnLink).Layout(gtx)
						}),
					)
				})
			}),
			// Sizes
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return subSection(th, gtx, "Sizes", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnXs, "XSmall").WithVariant(component.BtnPrimary).WithSize(component.BtnXs).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnSm, "Small").WithVariant(component.BtnPrimary).WithSize(component.BtnSm).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnPrimary, "Medium").WithVariant(component.BtnPrimary).WithSize(component.BtnMd).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnLg, "Large").WithVariant(component.BtnPrimary).WithSize(component.BtnLg).Layout(gtx)
						}),
					)
				})
			}),
			// States
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return subSection(th, gtx, "States", func(gtx layout.Context) layout.Dimensions {
					disBtn := component.NewButton(th, &a.btnCard1, "Disabled")
					disBtn.WithVariant(component.BtnPrimary)
					disBtn.Disabled = true
					loadBtn := component.NewButton(th, &a.btnCard2, "Loading...")
					loadBtn.WithVariant(component.BtnSecondary)
					loadBtn.Loading = true
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return disBtn.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return loadBtn.Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnModal, "Open Modal").WithVariant(component.BtnPrimary).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &a.btnToast, "Show Toast").WithVariant(component.BtnSuccess).Layout(gtx)
						}),
					)
				})
			}),
			// Icon Buttons
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return subSection(th, gtx, "Icon Buttons", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnPrimary, iconDashboard).WithVariant(component.BtnPrimary).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnSecondary, iconComponents).WithVariant(component.BtnSecondary).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnAccent, iconStar).WithVariant(component.BtnAccent).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnGhost, iconPerson).WithVariant(component.BtnGhost).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnOutline, iconList).WithVariant(component.BtnOutline).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnError, iconCheck).WithVariant(component.BtnError).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnPrimary, iconSettings).WithVariant(component.BtnPrimary).WithSize(component.BtnSm).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnPrimary, iconSettings).WithVariant(component.BtnPrimary).WithSize(component.BtnLg).Layout(gtx)
						}),
					)
				})
			}),

			// IconText Buttons
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return subSection(th, gtx, "IconText Buttons", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconTextButton(th, &a.btnIconText1, iconDashboard, "Dashboard").WithVariant(component.BtnPrimary).Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconTextButton(th, &a.btnIconText2, iconSettings, "Settings").WithVariant(component.BtnSecondary).Layout(gtx)
						}),
					)
				})
			}),
		)
	})(gtx)
}

func (a *App) sectionBadgesChips(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 24}.Layout(gtx,
		// Badges
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Badges", "Status indicators and tags", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, "Variants", func(gtx layout.Context) layout.Dimensions {
							return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Default").Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Primary").WithVariant(component.BadgePrimary).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Secondary").WithVariant(component.BadgeSecondary).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Accent").WithVariant(component.BadgeAccent).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Info").WithVariant(component.BadgeInfo).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Success").WithVariant(component.BadgeSuccess).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Warning").WithVariant(component.BadgeWarning).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Error").WithVariant(component.BadgeError).Layout(gtx)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Ghost").WithVariant(component.BadgeGhost).Layout(gtx)
								}),
							)
						})
					}),
				)
			})(gtx)
		}),
		// Chips
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Chips", "Removable tag components", func(gtx layout.Context) layout.Dimensions {
				return kit.WrapRow{Gap: 8, RowGap: 8}.Layout(gtx,
					func(gtx layout.Context) layout.Dimensions {
						return component.NewChip(th, "Design").Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewChip(th, "Engineering").Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewChip(th, "Product").Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewChip(th, "Go").Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewChip(th, "Gio UI").Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewChip(th, "TailwindCSS").Layout(gtx)
					},
				)
			})(gtx)
		}),
		// Skeleton
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Skeleton", "Loading placeholder", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 8}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewSkeleton(th).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						s := component.NewSkeleton(th)
						s.Width = 300
						s.Height = 14
						return s.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						s := component.NewSkeleton(th)
						s.Width = 160
						s.Height = 14
						return s.Layout(gtx)
					}),
				)
			})(gtx)
		}),
	)
}

func (a *App) sectionAlerts(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Alerts", "Notification banners for user feedback", func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 12}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "This is an informational message.", component.AlertInfo).WithIcon(iconInfo).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "Operation completed successfully! Your changes have been saved.", component.AlertSuccess).WithIcon(iconSuccess).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "Please review your input before submitting.", component.AlertWarning).WithIcon(iconWarning).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "An error occurred. Please try again or contact support.", component.AlertError).WithIcon(iconError).Layout(gtx)
			}),
		)
	})(gtx)
}

func (a *App) sectionAvatarsProgress(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 24}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Avatars", "User profile circles", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 16, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						av := component.NewAvatar(th, "XS")
						av.Size = component.AvatarXs
						return av.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						av := component.NewAvatar(th, "SM")
						av.Size = component.AvatarSm
						return av.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAvatar(th, "MD").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						av := component.NewAvatar(th, "LG")
						av.Size = component.AvatarLg
						return av.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "XS  SM  MD  LG").Sm().WithColor(theme.Gray400).Layout(gtx)
					}),
				)
			})(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Progress", "Progress indicators", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 12}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, fmt.Sprintf("%.0f%%", a.progress*100)).Sm().Bold().Layout(gtx)
							}),
							kit.Grow(func(gtx layout.Context) layout.Dimensions {
								return component.NewProgress(th, a.progress).Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						p := component.NewProgress(th, 0.4)
						p.Variant = component.ProgressSuccess
						return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "40%").Sm().WithColor(th.Success).Layout(gtx)
							}),
							kit.Grow(func(gtx layout.Context) layout.Dimensions { return p.Layout(gtx) }),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						p := component.NewProgress(th, 0.75)
						p.Variant = component.ProgressWarning
						return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "75%").Sm().WithColor(th.Warning).Layout(gtx)
							}),
							kit.Grow(func(gtx layout.Context) layout.Dimensions { return p.Layout(gtx) }),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						p := component.NewProgress(th, 0.25)
						p.Variant = component.ProgressError
						return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "25%").Sm().WithColor(th.Error).Layout(gtx)
							}),
							kit.Grow(func(gtx layout.Context) layout.Dimensions { return p.Layout(gtx) }),
						)
					}),
				)
			})(gtx)
		}),
	)
}

// ─── Section: Controls ───────────────────────────────────────────────────────

func (a *App) sectionControls(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 24}.Layout(gtx,
		// Checkboxes
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Checkbox", "Multi-select boolean controls", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 12}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewCheckbox(th, &a.check1, "Accept terms and conditions").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewCheckbox(th, &a.check2, "Subscribe to newsletter").WithVariant(component.BtnSecondary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewCheckbox(th, &a.check3, "Enable dark mode").WithVariant(component.BtnAccent).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Radio
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Radio", "Single-select option group", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 8}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.radioGroup.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, fmt.Sprintf("Selected: %s", a.radioGroup.Items[a.radioGroup.Selected()])).Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Select
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Select", "Inline expanding dropdown", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 8}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.W(300)(gtx, func(gtx layout.Context) layout.Dimensions {
							return a.selectComp.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, fmt.Sprintf("Value: %s", a.selectComp.Value())).Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Range
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Range", "Draggable slider", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, fmt.Sprintf("Primary — %.0f%%", a.rangeFloat1.Value*100), func(gtx layout.Context) layout.Dimensions {
							return component.NewRange(th, &a.rangeFloat1).Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return subSection(th, gtx, fmt.Sprintf("Secondary — %.0f%%", a.rangeFloat2.Value*100), func(gtx layout.Context) layout.Dimensions {
							return component.NewRange(th, &a.rangeFloat2).WithVariant(component.BtnSecondary).Layout(gtx)
						})
					}),
				)
			})(gtx)
		}),

		// Rating
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Rating", "Star rating selector", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 8}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.rating.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, fmt.Sprintf("Rating: %d / %d stars", a.rating.Value, a.rating.Max)).Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Accordion
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Accordion", "Collapsible content sections", func(gtx layout.Context) layout.Dimensions {
				bodies := []layout.Widget{
					func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "GioUI Kit is a TailwindCSS and DaisyUI inspired component library for the Gio immediate-mode UI framework.").Sm().WithColor(theme.Gray500).Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Add github.com/hongshengjie/gioui-kit to your go.mod and import the packages you need.").Sm().WithColor(theme.Gray500).Layout(gtx)
					},
					func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Yes! Pass any *theme.Theme to component constructors. Use theme.Light(), theme.Dark(), theme.Cupcake(), or theme.Nord().").Sm().WithColor(theme.Gray500).Layout(gtx)
					},
				}
				return a.accordion.Layout(gtx, bodies)
			})(gtx)
		}),

		// Menu
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Menu", "Vertical navigation menu", func(gtx layout.Context) layout.Dimensions {
				return kit.W(200)(gtx, func(gtx layout.Context) layout.Dimensions {
					return component.NewMenu(th, a.menuItems...).WithBorder().Layout(gtx)
				})
			})(gtx)
		}),

		// Tooltip
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Tooltip", "Hover to reveal tooltip", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 16, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.tooltip.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return component.NewBadge(th, "Hover me").WithVariant(component.BadgePrimary).Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						tip2 := component.NewTooltip(th, "Bottom tooltip").WithPosition(component.TooltipBottom)
						return tip2.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return component.NewBadge(th, "Bottom tip").WithVariant(component.BadgeSecondary).Layout(gtx)
						})
					}),
				)
			})(gtx)
		}),

		// Divider
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Divider", "Horizontal and labeled dividers", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewDivider(th).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewDivider(th).WithLabel("OR").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewDivider(th).WithLabel("Section Break").Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Kbd
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Keyboard Keys", "Keyboard shortcut display", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 16, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewKbd(th, "Ctrl").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.KbdGroup(th, "Ctrl", "C")(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.KbdGroup(th, "Ctrl", "Shift", "Z")(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewKbd(th, "Enter").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewKbd(th, "Esc").Layout(gtx)
					}),
				)
			})(gtx)
		}),
	)
}

// ─── Section: Data Display ───────────────────────────────────────────────────

func (a *App) sectionDataDisplay(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 24}.Layout(gtx,
		// Stat
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Stat", "Metric display cards", func(gtx layout.Context) layout.Dimensions {
				return kit.Grid{Cols: 1, SmCols: 2, MdCols: 4, Gap: 12}.Layout(gtx,
					component.NewStat(th, "Total Users", "89,400").WithDesc("↑ 12% from last month").WithFigureIcon(iconPeople).Layout,
					component.NewStat(th, "Revenue", "$45,231").WithDesc("↑ 8% from last month").WithFigureIcon(iconMoney).Layout,
					component.NewStat(th, "Active Sessions", "1,429").WithDesc("→ stable").WithFigureIcon(iconTrendingUp).Layout,
					component.NewStat(th, "Issues", "12").WithDesc("↓ 3 resolved today").WithFigureIcon(iconWarning).Layout,
				)
			})(gtx)
		}),

		// Steps
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Steps", "Progress wizard indicator", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.steps.WithCurrent(a.stepIndex).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, fmt.Sprintf("Step %d of %d", a.stepIndex+1, len(a.steps.Items))).Sm().WithColor(theme.Gray500).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Radial Progress
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Radial Progress", "Circular progress indicators", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 24, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewRadialProgress(th, 0.7).WithLabel("70%").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewRadialProgress(th, 0.45).WithLabel("45%").WithVariant(component.ProgressSecondary).WithSize(96).WithThick(10).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewRadialProgress(th, 1.0).WithLabel("100%").WithVariant(component.ProgressSuccess).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewRadialProgress(th, 0.2).WithLabel("20%").WithVariant(component.ProgressError).Layout(gtx)
					}),
				)
			})(gtx)
		}),

		// Loading
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Loading", "Animated loading indicators", func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 32, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.FlexCol{Gap: 8}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "Spinner").Xs().WithColor(theme.Gray400).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewLoading(th).WithVariant(component.LoadingSpinner).Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.FlexCol{Gap: 8}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "Dots").Xs().WithColor(theme.Gray400).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewLoading(th).WithVariant(component.LoadingDots).WithSize(40).Layout(gtx)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return kit.FlexCol{Gap: 8}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewText(th, "Ring").Xs().WithColor(theme.Gray400).Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return component.NewLoading(th).WithVariant(component.LoadingRing).WithColor(th.Secondary).Layout(gtx)
							}),
						)
					}),
				)
			})(gtx)
		}),

		// Table
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return sectionCard(th, "Table", "Data table with headers and rows", func(gtx layout.Context) layout.Dimensions {
				headers := []string{"Name", "Role", "Status", "Joined"}
				rows := [][]string{
					{"Alice Chen", "Admin", "Active", "2023-01"},
					{"Bob Smith", "Developer", "Active", "2023-03"},
					{"Carol Wu", "Designer", "Away", "2023-06"},
					{"Dan Park", "DevOps", "Active", "2024-01"},
					{"Eve Johnson", "QA", "Offline", "2024-02"},
				}
				return component.NewTable(th, headers, rows).WithZebra().WithBorder().Layout(gtx)
			})(gtx)
		}),
	)
}

func (a *App) sectionLayout(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 24}.Layout(gtx,

		// Grid
		layout.Rigid(sectionCard(th, "Grid", "Responsive grid layout", func(gtx layout.Context) layout.Dimensions {
			grid := kit.Grid{Cols: 1, MdCols: 2, LgCols: 3, Gap: 16}
			return grid.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return component.NewCard(th).WithBorder().CardWithHeader(gtx, "Card 1",
						layout.Widget(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Responsive grid item").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
					)
				},
				func(gtx layout.Context) layout.Dimensions {
					return component.NewCard(th).WithBorder().CardWithHeader(gtx, "Card 2",
						layout.Widget(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Another grid item").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
					)
				},
				func(gtx layout.Context) layout.Dimensions {
					return component.NewCard(th).WithBorder().CardWithHeader(gtx, "Card 3",
						layout.Widget(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Third item").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
					)
				},
				func(gtx layout.Context) layout.Dimensions {
					return component.NewCard(th).WithBorder().CardWithHeader(gtx, "Card 4",
						layout.Widget(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Fourth item").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
					)
				},
			)
		})),

		// Box
		layout.Rigid(sectionCard(th, "Box", "Simple container with padding and background", func(gtx layout.Context) layout.Dimensions {
			return kit.Box{Padding: layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}, Background: th.Base200, Radius: th.RoundedMd}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, "This is inside a Box with padding and background.").Layout(gtx)
				},
			)
		})),

		// Container
		layout.Rigid(sectionCard(th, "Container", "Centered max-width container", func(gtx layout.Context) layout.Dimensions {
			return kit.Container{MaxWidth: 600, Padding: 16}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, "This content is centered with a max width of 600dp.").Layout(gtx)
				},
			)
		})),

		// ScrollY
		layout.Rigid(sectionCard(th, "ScrollY", "Vertical scrollable list", func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Max.Y = 200 // limit height for demo
			return a.demoScroll.Layout(gtx, 10, func(gtx layout.Context, i int) layout.Dimensions {
				return component.NewText(th, fmt.Sprintf("Scrollable item %d", i+1)).Layout(gtx)
			})
		})),

		// BottomNav
		layout.Rigid(sectionCard(th, "BottomNav", "Mobile-style bottom navigation", func(gtx layout.Context) layout.Dimensions {
			items := []scaffold.BottomNavItem{
				{Label: "Home", IconData: iconHome, Active: true},
				{Label: "Search", IconData: iconSearch},
				{Label: "Profile", IconData: iconPerson},
			}
			return scaffold.NewBottomNav(th, items).Layout(gtx)
		})),
	)
}

func (a *App) sectionModifiers(gtx layout.Context) layout.Dimensions {
	th := a.th
	return kit.FlexCol{Gap: 24}.Layout(gtx,

		// Rounded
		layout.Rigid(sectionCard(th, "Rounded", "Rounded corners modifier", func(gtx layout.Context) layout.Dimensions {
			return modifier.Rounded{Radius: th.RoundedLg}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return kit.Box{Padding: layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}, Background: th.Primary}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Rounded box").WithColor(th.PrimaryContent).Layout(gtx)
						},
					)
				},
			)
		})),

		// Shadow
		layout.Rigid(sectionCard(th, "Shadow", "Drop shadow modifier", func(gtx layout.Context) layout.Dimensions {
			return modifier.Shadow{Style: modifier.ShadowMd, Radius: th.RoundedMd}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return kit.Box{Padding: layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}, Background: th.Base100, Radius: th.RoundedMd}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Card with shadow").Layout(gtx)
						},
					)
				},
			)
		})),

		// Opacity
		layout.Rigid(sectionCard(th, "Opacity", "Opacity modifier", func(gtx layout.Context) layout.Dimensions {
			return kit.FlexRow{Gap: 16}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, "Normal").Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return modifier.OpacityMod{Opacity: 0.5}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "50% opacity").Layout(gtx)
						},
					)
				}),
			)
		})),

		// Ring
		layout.Rigid(sectionCard(th, "Ring", "Ring border modifier", func(gtx layout.Context) layout.Dimensions {
			return modifier.Ring{Color: th.Primary, Width: 2}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return kit.Box{Padding: layout.Inset{Top: 16, Bottom: 16, Left: 16, Right: 16}, Background: th.Base100, Radius: th.RoundedMd}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Box with ring").Layout(gtx)
						},
					)
				},
			)
		})),

		// Gradient
		layout.Rigid(sectionCard(th, "Gradient", "Gradient background modifier", func(gtx layout.Context) layout.Dimensions {
			grad := modifier.LinearGradient{
				From: th.Primary,
				To:   th.Secondary,
				Dir:  modifier.GradientToBottom,
			}
			return grad.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					return kit.Box{Padding: layout.Inset{Top: 32, Bottom: 32, Left: 16, Right: 16}, Radius: th.RoundedMd}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Gradient background").WithColor(theme.White).Bold().Layout(gtx)
						},
					)
				},
			)
		})),
	)
}
