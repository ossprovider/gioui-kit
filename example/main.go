// Command demo shows a complete example application using gioui-kit.
//
// It demonstrates all components in a real app layout with:
//   - AppShell with Navbar + Sidebar
//   - Buttons, Badges, Cards, Alerts, Inputs, Toggles
//   - Tabs, Progress bars, Avatars
//   - Grid and Flex layouts
//   - Modal and Toast overlays
package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/pkg/component"
	kit "github.com/hongshengjie/gioui-kit/pkg/layout"
	"github.com/hongshengjie/gioui-kit/pkg/scaffold"
	"github.com/hongshengjie/gioui-kit/pkg/theme"
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("GioUI Kit Demo"),
			app.Size(unit.Dp(1200), unit.Dp(800)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type App struct {
	th *theme.Theme

	// Shell
	shell   *scaffold.AppShell
	sidebar *scaffold.Sidebar
	navbar  *scaffold.Navbar

	// State
	modal      *scaffold.Modal
	drawer     *scaffold.Drawer
	toast      *scaffold.Toast
	scrollList kit.ScrollY

	// Buttons
	btnPrimary   widget.Clickable
	btnSecondary widget.Clickable
	btnAccent    widget.Clickable
	btnOutline   widget.Clickable
	btnGhost     widget.Clickable
	btnError     widget.Clickable
	btnSmall     widget.Clickable
	btnLarge     widget.Clickable
	btnModal      widget.Clickable
	btnToast      widget.Clickable
	btnModalClose widget.Clickable

	// Inputs
	editor1 widget.Editor
	editor2 widget.Editor
	toggle1 widget.Bool
	toggle2 widget.Bool

	// Tabs
	tabs *component.Tabs

	// Progress
	progress float32

	// Sidebar selection
	sidebarIndex int
}

func NewApp() *App {
	th := theme.Light()

	sideItems := []scaffold.SidebarItem{
		{Label: "Dashboard", Icon: "◉", Active: true},
		{Label: "Components", Icon: "◫"},
		{Label: "Layout", Icon: "⊞"},
		{Label: "Forms", Icon: "✎"},
		{Label: "Settings", Icon: "⚙"},
	}

	a := &App{
		th:       th,
		navbar:   scaffold.NewNavbar(th),
		sidebar:  scaffold.NewSidebar(th, sideItems),
		modal:    scaffold.NewModal(th),
		drawer:   scaffold.NewDrawer(th),
		toast:    scaffold.NewToast(th),
		tabs:     component.NewTabs(th, []string{"Overview", "Components", "Layout", "Source"}),
		progress: 0.65,
	}

	a.editor1.SingleLine = true
	a.editor2.SingleLine = true
	a.toggle1.Value = true

	a.sidebar.OnSelect = func(i int) {
		a.sidebarIndex = i
		for j := range a.sidebar.Items {
			a.sidebar.Items[j].Active = j == i
		}
	}

	a.shell = scaffold.NewAppShell(th)

	return a
}

func run(w *app.Window) error {
	a := NewApp()
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle interactions
			if a.btnModal.Clicked(gtx) {
				a.modal.Show()
			}
			if a.btnToast.Clicked(gtx) {
				a.toast.Show("Operation completed successfully!")
			}

			a.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	th := a.th

	// Configure shell
	a.shell.Navbar = func(gtx layout.Context) layout.Dimensions {
		return a.navbar.Layout(gtx,
			// Start: brand
			func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "GioUI Kit").H3().Bold().WithColor(th.Primary).Layout(gtx)
			},
			// Center: nil
			nil,
			// End: badges
			func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewBadge(th, "v0.1.0").WithVariant(component.BadgePrimary).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAvatar(th, "HS").Layout(gtx)
					}),
				)
			},
		)
	}

	a.shell.Sidebar = func(gtx layout.Context) layout.Dimensions {
		return a.sidebar.Layout(gtx)
	}
	a.shell.SidebarWidth = 240

	a.shell.Content = func(gtx layout.Context) layout.Dimensions {
		return a.layoutContent(gtx)
	}

	// Render shell
	dims := a.shell.Layout(gtx)

	// Overlay: Modal
	a.modal.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 16}.Layout(gtx,
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "Modal Dialog").H3().Bold().Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "This is a DaisyUI-style modal component rendered with Gio.").Sm().Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := component.NewButton(th, &a.btnModalClose, "Close")
				btn.WithVariant(component.BtnPrimary)
				return btn.Layout(gtx)
			}),
		)
	})

	// Overlay: Toast
	a.toast.Layout(gtx)

	return dims
}

func (a *App) layoutContent(gtx layout.Context) layout.Dimensions {
	th := a.th
	list := &a.scrollList

	return list.List.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		return layout.Inset{
			Top: th.Space6, Bottom: th.Space6,
			Left: th.Space8, Right: th.Space8,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return kit.FlexCol{Gap: 32}.Layout(gtx,
				// Title + breadcrumb
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return kit.FlexCol{Gap: 8}.Layout(gtx,
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return scaffold.NewBreadcrumb(th, "Home", "Components", "Demo").Layout(gtx)
						}),
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Component Gallery").H1().Layout(gtx)
						}),
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "A comprehensive showcase of all gioui-kit components inspired by TailwindCSS and DaisyUI.").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
					)
				}),

				// Tabs
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.tabs.Layout(gtx)
				}),

				// Section: Buttons
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.sectionButtons(gtx)
				}),

				// Section: Badges
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.sectionBadges(gtx)
				}),

				// Section: Alerts
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.sectionAlerts(gtx)
				}),

				// Section: Cards
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.sectionCards(gtx)
				}),

				// Section: Forms
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.sectionForms(gtx)
				}),

				// Section: Progress
				kit.Rigid(func(gtx layout.Context) layout.Dimensions {
					return a.sectionProgress(gtx)
				}),
			)
		})
	})
}

func (a *App) sectionButtons(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Buttons", "DaisyUI-style buttons with variants and sizes", func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 16}.Layout(gtx,
			// Variant row
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "Variants").H4().Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnPrimary, "Primary").WithVariant(component.BtnPrimary).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnSecondary, "Secondary").WithVariant(component.BtnSecondary).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnAccent, "Accent").WithVariant(component.BtnAccent).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnOutline, "Outline").WithVariant(component.BtnOutline).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnGhost, "Ghost").WithVariant(component.BtnGhost).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnError, "Error").WithVariant(component.BtnError).Layout(gtx)
					}),
				)
			}),
			// Size row
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "Sizes").H4().Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnSmall, "Small").WithVariant(component.BtnPrimary).WithSize(component.BtnSm).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &widget.Clickable{}, "Medium").WithVariant(component.BtnPrimary).WithSize(component.BtnMd).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnLarge, "Large").WithVariant(component.BtnPrimary).WithSize(component.BtnLg).Layout(gtx)
					}),
				)
			}),
			// Action row
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "Actions").H4().Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnModal, "Open Modal").WithVariant(component.BtnPrimary).Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnToast, "Show Toast").WithVariant(component.BtnSuccess).Layout(gtx)
					}),
				)
			}),
		)
	})(gtx)
}

func (a *App) sectionBadges(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Badges", "Status indicators and tags", func(gtx layout.Context) layout.Dimensions {
		return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Default").Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Primary").WithVariant(component.BadgePrimary).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Secondary").WithVariant(component.BadgeSecondary).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Accent").WithVariant(component.BadgeAccent).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Success").WithVariant(component.BadgeSuccess).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Warning").WithVariant(component.BadgeWarning).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewBadge(th, "Error").WithVariant(component.BadgeError).Layout(gtx)
			}),
		)
	})(gtx)
}

func (a *App) sectionAlerts(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Alerts", "Notification banners", func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 12}.Layout(gtx,
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "This is an informational message.", component.AlertInfo).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "Operation completed successfully!", component.AlertSuccess).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "Please check your input values.", component.AlertWarning).Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewAlert(th, "An error occurred while processing.", component.AlertError).Layout(gtx)
			}),
		)
	})(gtx)
}

func (a *App) sectionCards(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Cards", "Content containers with different styles", func(gtx layout.Context) layout.Dimensions {
		return kit.Grid{Cols: 3, Gap: 16}.Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				card := component.NewCard(th).WithBorder()
				return card.CardWithHeader(gtx, "Bordered Card", func(gtx layout.Context) layout.Dimensions {
					return component.NewText(th, "This card has a visible border around it, suitable for lighter backgrounds.").Sm().WithColor(theme.Gray500).Layout(gtx)
				})
			},
			func(gtx layout.Context) layout.Dimensions {
				card := component.NewCard(th)
				return card.CardWithHeader(gtx, "Default Card", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexCol{Gap: 12}.Layout(gtx,
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Cards can contain any content including other components.").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewButton(th, &widget.Clickable{}, "Action").WithVariant(component.BtnPrimary).WithSize(component.BtnSm).Layout(gtx)
						}),
					)
				})
			},
			func(gtx layout.Context) layout.Dimensions {
				card := component.NewCard(th).WithCompact()
				return card.CardWithHeader(gtx, "Compact Card", func(gtx layout.Context) layout.Dimensions {
					return kit.FlexCol{Gap: 8}.Layout(gtx,
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Compact variant has reduced padding.").Sm().WithColor(theme.Gray500).Layout(gtx)
						}),
						kit.Rigid(func(gtx layout.Context) layout.Dimensions {
							return kit.FlexRow{Gap: 4}.Layout(gtx,
								kit.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Go").WithVariant(component.BadgePrimary).Layout(gtx)
								}),
								kit.Rigid(func(gtx layout.Context) layout.Dimensions {
									return component.NewBadge(th, "Gio").WithVariant(component.BadgeAccent).Layout(gtx)
								}),
							)
						}),
					)
				})
			},
		)
	})(gtx)
}

func (a *App) sectionForms(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Form Elements", "Input fields, toggles, and controls", func(gtx layout.Context) layout.Dimensions {
		return kit.Grid{Cols: 2, Gap: 24}.Layout(gtx,
			// Left column: Inputs
			func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewInput(th, &a.editor1, "Enter your name...").WithLabel("Name").Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewInput(th, &a.editor2, "Enter email...").WithLabel("Email").WithVariant(component.InputPrimary).Layout(gtx)
					}),
				)
			},
			// Right column: Toggles
			func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 16}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, "Toggle Switches").H4().Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewToggle(th, &a.toggle1, "Enable notifications").Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewToggle(th, &a.toggle2, "Dark mode").Layout(gtx)
					}),
				)
			},
		)
	})(gtx)
}

func (a *App) sectionProgress(gtx layout.Context) layout.Dimensions {
	th := a.th
	return sectionCard(th, "Progress", "Progress indicators", func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 12}.Layout(gtx,
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, fmt.Sprintf("%.0f%%", a.progress*100)).Sm().Bold().Layout(gtx)
					}),
					kit.Grow(func(gtx layout.Context) layout.Dimensions {
						return component.NewProgress(th, a.progress).Layout(gtx)
					}),
				)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				p := component.NewProgress(th, 0.4)
				p.Variant = component.ProgressSuccess
				return p.Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				p := component.NewProgress(th, 0.8)
				p.Variant = component.ProgressWarning
				return p.Layout(gtx)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				p := component.NewProgress(th, 0.25)
				p.Variant = component.ProgressError
				return p.Layout(gtx)
			}),
		)
	})(gtx)
}

// sectionCard wraps a demo section in a card with title.
func sectionCard(th *theme.Theme, title, subtitle string, body layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return kit.FlexCol{Gap: 12}.Layout(gtx,
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexCol{Gap: 4}.Layout(gtx,
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, title).H2().Layout(gtx)
					}),
					kit.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewText(th, subtitle).Sm().WithColor(theme.Gray400).Layout(gtx)
					}),
				)
			}),
			kit.Rigid(func(gtx layout.Context) layout.Dimensions {
				card := component.NewCard(th).WithBorder()
				return card.Layout(gtx, body)
			}),
		)
	}
}
