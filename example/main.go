// Command demo shows a complete example application using gioui-kit.
package main

import (
	"log" //nolint
	"net/http"
	_ "net/http/pprof"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/hongshengjie/gioui-kit/component"
	kit "github.com/hongshengjie/gioui-kit/layout"
	"github.com/hongshengjie/gioui-kit/scaffold"
	"github.com/hongshengjie/gioui-kit/theme"
)

// mustIcon panics if an iconvg icon fails to parse (should never happen with bundled data).
func mustIcon(data []byte) *widget.Icon {
	ic, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return ic
}

// Prebuilt Material Design icons used across the app.
var (
	iconMenu        = mustIcon(icons.NavigationMenu)
	iconClose       = mustIcon(icons.NavigationClose)
	iconDashboard   = mustIcon(icons.ActionDashboard)
	iconComponents  = mustIcon(icons.NavigationApps)
	iconLayout      = mustIcon(icons.ActionViewModule)
	iconForms       = mustIcon(icons.ContentCreate)
	iconSettings    = mustIcon(icons.ActionSettings)
	iconPerson      = mustIcon(icons.SocialPerson)
	iconCheck       = mustIcon(icons.ActionCheckCircle)
	iconStar        = mustIcon(icons.ActionStars)
	iconList        = mustIcon(icons.ActionList)
	iconInfo        = mustIcon(icons.ActionInfo)
	iconSuccess     = mustIcon(icons.ActionCheckCircle)
	iconWarning     = mustIcon(icons.AlertWarning)
	iconError       = mustIcon(icons.AlertError)
	iconChevronDown = mustIcon(icons.NavigationArrowDropDown)
	iconChevronUp   = mustIcon(icons.NavigationArrowDropUp)
	iconStarFilled  = mustIcon(icons.ToggleStar)
	iconStarBorder  = mustIcon(icons.ToggleStarBorder)
	iconHome        = mustIcon(icons.ActionHome)
	iconSearch      = mustIcon(icons.ActionSearch)
	iconPeople      = mustIcon(icons.SocialGroup)
	iconMoney       = mustIcon(icons.EditorAttachMoney)
	iconTrendingUp  = mustIcon(icons.ActionTrendingUp)
)

func main() {
	go func() {
		log.Println("pprof listening on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
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

// App holds all application state.
type App struct {
	th *theme.Theme

	// Shell
	shell   *scaffold.AppShell
	sidebar *scaffold.Sidebar
	navbar  *scaffold.Navbar

	// Overlays
	modal  *scaffold.Modal
	drawer *scaffold.Drawer
	toast  *scaffold.Toast

	// Scrolling
	scroll     kit.ScrollY
	demoScroll kit.ScrollY

	// Navigation
	pageIndex int // 0=Dashboard 1=Components 2=Layout 3=Forms 4=Settings

	// Components page: sub-tabs
	compTabs *component.Tabs

	// --- Button demo clickables ---
	btnPrimary   widget.Clickable
	btnSecondary widget.Clickable
	btnAccent    widget.Clickable
	btnOutline   widget.Clickable
	btnGhost     widget.Clickable
	btnLink      widget.Clickable
	btnError     widget.Clickable
	btnInfo      widget.Clickable
	btnSuccess   widget.Clickable
	btnWarning   widget.Clickable
	btnXs        widget.Clickable
	btnSm        widget.Clickable
	btnLg        widget.Clickable

	// --- Action buttons ---
	btnModal       widget.Clickable
	btnModalClose  widget.Clickable
	btnToast       widget.Clickable
	btnDrawer      widget.Clickable
	btnDrawerClose widget.Clickable

	// --- Mobile nav (hamburger + drawer items) ---
	btnHamburger      widget.Clickable
	btnDrawerDash     widget.Clickable
	btnDrawerComp     widget.Clickable
	btnDrawerLayout   widget.Clickable
	btnDrawerForms    widget.Clickable
	btnDrawerSettings widget.Clickable

	// --- Icon buttons ---
	ibnHamburger widget.Clickable
	ibnClose     widget.Clickable
	ibnPrimary   widget.Clickable
	ibnSecondary widget.Clickable
	ibnAccent    widget.Clickable
	ibnGhost     widget.Clickable
	ibnOutline   widget.Clickable
	ibnError     widget.Clickable

	// --- IconText buttons ---
	btnIconText1 widget.Clickable
	btnIconText2 widget.Clickable

	// --- Overview quick-start ---
	btnOverview1 widget.Clickable
	btnOverview2 widget.Clickable

	// --- Card embedded buttons ---
	btnCard1 widget.Clickable
	btnCard2 widget.Clickable
	btnCard3 widget.Clickable

	// --- Settings: theme picker ---
	btnLight   widget.Clickable
	btnDark    widget.Clickable
	btnCupcake widget.Clickable
	btnNord    widget.Clickable

	// --- Form editors ---
	editor1 widget.Editor
	editor2 widget.Editor
	editor3 widget.Editor
	editor4 widget.Editor

	// --- Toggles ---
	toggle1 widget.Bool
	toggle2 widget.Bool
	toggle3 widget.Bool

	// --- Checkboxes ---
	check1 widget.Bool
	check2 widget.Bool
	check3 widget.Bool

	// --- Radio group ---
	radioGroup *component.RadioGroup

	// --- Select ---
	selectComp *component.Select

	// --- Range sliders ---
	rangeFloat1 widget.Float
	rangeFloat2 widget.Float

	// --- Rating ---
	rating *component.Rating

	// --- Accordion ---
	accordion *component.Accordion

	// --- Menu ---
	menuItems []*component.MenuItem

	// --- Tooltip ---
	tooltip *component.Tooltip

	// --- Steps ---
	steps     *component.Steps
	stepIndex int

	// --- Progress ---
	progress float32
}

func NewApp() *App {
	th := theme.Light()

	sideItems := []scaffold.SidebarItem{
		{Label: "Dashboard", IconData: iconDashboard, Active: true},
		{Label: "Components", IconData: iconComponents},
		{Label: "Layout", IconData: iconLayout},
		{Label: "Forms", IconData: iconForms},
		{Label: "Settings", IconData: iconSettings},
	}

	a := &App{
		th:         th,
		navbar:     scaffold.NewNavbar(th),
		sidebar:    scaffold.NewSidebar(th, sideItems),
		modal:      scaffold.NewModal(th),
		drawer:     scaffold.NewDrawer(th),
		toast:      scaffold.NewToast(th),
		compTabs:   component.NewTabs(th, []string{"Buttons", "Badges & Chips", "Alerts", "Avatars & Progress", "Controls", "Data Display", "Layout", "Modifiers"}),
		progress:   0.65,
		radioGroup: component.NewRadioGroup(th, []string{"Option A", "Option B", "Option C"}),
		selectComp: component.NewSelect(th, []string{"Apple", "Banana", "Cherry", "Durian", "Elderberry"}).WithChevrons(iconChevronDown, iconChevronUp),
		rating:     component.NewRating(th, 5).WithStarIcons(iconStarFilled, iconStarBorder),
		accordion: component.NewAccordion(th,
			component.NewAccordionItem("What is GioUI Kit?"),
			component.NewAccordionItem("How do I install it?"),
			component.NewAccordionItem("Can I use custom themes?"),
		),
		tooltip:   component.NewTooltip(th, "This is a tooltip!"),
		steps:     component.NewSteps(th, []string{"Account", "Profile", "Review", "Done"}),
		stepIndex: 1,
		menuItems: []*component.MenuItem{
			component.NewMenuItem("Dashboard").WithIcon(iconDashboard),
			component.NewMenuItem("Components").WithIcon(iconComponents),
			component.NewMenuItem("Settings").WithIcon(iconSettings),
		},
	}
	a.menuItems[0].Active = true
	a.rangeFloat1.Value = 0.4
	a.rangeFloat2.Value = 0.7
	a.rating.Value = 3

	a.editor1.SingleLine = true
	a.editor2.SingleLine = true
	a.editor3.SingleLine = true
	a.editor4.SingleLine = true
	a.toggle1.Value = true

	a.sidebar.OnSelect = func(i int) { a.selectPage(i) }

	a.shell = scaffold.NewAppShell(th)
	a.scroll.List.Axis = layout.Vertical // zero-value is Horizontal; must set explicitly
	return a
}

// selectPage switches the active page and syncs the sidebar.
func (a *App) selectPage(i int) {
	a.pageIndex = i
	for j := range a.sidebar.Items {
		a.sidebar.Items[j].Active = j == i
	}
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

			// Hamburger → open drawer
			if a.btnHamburger.Clicked(gtx) || a.ibnHamburger.Clicked(gtx) {
				a.drawer.Open()
			}
			// Drawer nav items
			if a.btnDrawerDash.Clicked(gtx) {
				a.selectPage(0)
				a.drawer.Close()
			}
			if a.btnDrawerComp.Clicked(gtx) {
				a.selectPage(1)
				a.drawer.Close()
			}
			if a.btnDrawerLayout.Clicked(gtx) {
				a.selectPage(2)
				a.drawer.Close()
			}
			if a.btnDrawerForms.Clicked(gtx) {
				a.selectPage(3)
				a.drawer.Close()
			}
			if a.btnDrawerSettings.Clicked(gtx) {
				a.selectPage(4)
				a.drawer.Close()
			}

			// Modal
			if a.btnModal.Clicked(gtx) {
				a.modal.Show()
			}
			if a.btnModalClose.Clicked(gtx) {
				a.modal.Hide()
			}
			// Toast
			if a.btnToast.Clicked(gtx) {
				a.toast.Show("Operation completed successfully!")
			}
			// Drawer
			if a.btnDrawer.Clicked(gtx) {
				a.drawer.Open()
			}
			if a.btnDrawerClose.Clicked(gtx) || a.ibnClose.Clicked(gtx) {
				a.drawer.Close()
			}
			// Theme switcher
			if a.btnLight.Clicked(gtx) {
				*a.th = *theme.Light()
			}
			if a.btnDark.Clicked(gtx) {
				*a.th = *theme.Dark()
			}
			if a.btnCupcake.Clicked(gtx) {
				*a.th = *theme.Cupcake()
			}
			if a.btnNord.Clicked(gtx) {
				*a.th = *theme.Nord()
			}

			a.layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (a *App) layout(gtx layout.Context) layout.Dimensions {
	th := a.th

	a.shell.Navbar = func(gtx layout.Context) layout.Dimensions {
		mobile := kit.ScreenBreakpoint(gtx) < kit.BreakpointLg
		return a.navbar.Layout(gtx,
			// start: hamburger on mobile, brand on desktop
			func(gtx layout.Context) layout.Dimensions {
				if mobile {
					return component.NewIconButton(th, &a.ibnHamburger, iconMenu).WithVariant(component.BtnGhost).Layout(gtx)
				}
				return component.NewText(th, "GioUI Kit").H3().Bold().WithColor(th.Primary).Layout(gtx)
			},
			// center: brand on mobile, empty on desktop
			func(gtx layout.Context) layout.Dimensions {
				if !mobile {
					return layout.Dimensions{}
				}
				return component.NewText(th, "GioUI Kit").H3().Bold().WithColor(th.Primary).Layout(gtx)
			},
			// end: badge+avatar on desktop, avatar only on mobile
			func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if mobile {
							return layout.Dimensions{}
						}
						return component.NewBadge(th, "v0.1.0").WithVariant(component.BadgePrimary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewAvatar(th, "GK").Layout(gtx)
					}),
				)
			},
		)
	}

	a.shell.Sidebar = func(gtx layout.Context) layout.Dimensions {
		return a.sidebar.Layout(gtx)
	}
	a.shell.SidebarWidth = 220

	a.shell.Content = func(gtx layout.Context) layout.Dimensions {
		return a.layoutContent(gtx)
	}

	dims := a.shell.Layout(gtx)

	// Overlay: Drawer
	a.drawer.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		th := a.th
		return layout.Inset{
			Top: th.Space4, Left: th.Space4, Right: th.Space4,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return kit.FlexCol{Gap: 16}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return kit.FlexRow{Gap: 8, Alignment: kit.ItemsCenter}.Layout(gtx,
						kit.Grow(func(gtx layout.Context) layout.Dimensions {
							return component.NewText(th, "Navigation").H3().Bold().Layout(gtx)
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return component.NewIconButton(th, &a.ibnClose, iconClose).WithVariant(component.BtnGhost).WithSize(component.BtnSm).Layout(gtx)
						}),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return kit.DividerH{Color: th.Base300}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawerNavItem(th, &a.btnDrawerDash, iconDashboard, "Dashboard")(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawerNavItem(th, &a.btnDrawerComp, iconComponents, "Components")(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawerNavItem(th, &a.btnDrawerLayout, iconLayout, "Layout")(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawerNavItem(th, &a.btnDrawerForms, iconForms, "Forms")(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return drawerNavItem(th, &a.btnDrawerSettings, iconSettings, "Settings")(gtx)
				}),
			)
		})
	})

	// Overlay: Modal
	a.modal.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		th := a.th
		return kit.FlexCol{Gap: 16}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "Modal Dialog").H3().Bold().Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return component.NewText(th, "This is a DaisyUI-style modal rendered with Gio. Click the backdrop or Close button to dismiss.").Sm().WithColor(theme.Gray500).Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return kit.FlexRow{Gap: 8}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnModalClose, "Close").WithVariant(component.BtnPrimary).Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return component.NewButton(th, &a.btnModalClose, "Cancel").WithVariant(component.BtnGhost).Layout(gtx)
					}),
				)
			}),
		)
	})

	// Overlay: Toast
	a.toast.Layout(gtx)

	return dims
}

// layoutContent routes to the selected page.
func (a *App) layoutContent(gtx layout.Context) layout.Dimensions {
	return a.scroll.List.Layout(gtx, 1, func(gtx layout.Context, _ int) layout.Dimensions {
		th := a.th
		hPad := th.Space8
		if kit.ScreenBreakpoint(gtx) < kit.BreakpointMd {
			hPad = th.Space4
		}
		return layout.Inset{
			Top: th.Space6, Bottom: th.Space8,
			Left: hPad, Right: hPad,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			switch a.pageIndex {
			case 1:
				return a.pageComponents(gtx)
			case 2:
				return a.pageLayout(gtx)
			case 3:
				return a.pageForms(gtx)
			case 4:
				return a.pageSettings(gtx)
			default:
				return a.pageDashboard(gtx)
			}
		})
	})
}
