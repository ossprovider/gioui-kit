// Package scaffold provides app-level layout scaffolds for Gio UI.
//
// These are high-level layout patterns commonly used in applications:
//
//	scaffold.AppShell{} - Full app shell with navbar + sidebar + content
//	scaffold.Navbar{}   - Top navigation bar
//	scaffold.Sidebar{}  - Side navigation panel
//	scaffold.Drawer{}   - Slide-in drawer overlay
//	scaffold.Modal{}    - Modal dialog overlay
//	scaffold.BottomNav{} - Mobile bottom navigation
package scaffold

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/pkg/theme"
)

// ============================================================
// AppShell - Full application shell
// ============================================================

// AppShell provides a complete application layout with optional
// navbar, sidebar, and main content area.
//
// Layout structure:
//
//	┌─────────────── Navbar ───────────────┐
//	│                                      │
//	├──────────┬───────────────────────────┤
//	│          │                           │
//	│ Sidebar  │       Content Area        │
//	│          │                           │
//	│          │                           │
//	└──────────┴───────────────────────────┘
type AppShell struct {
	Navbar       layout.Widget
	Sidebar      layout.Widget
	SidebarWidth unit.Dp
	Content      layout.Widget
	th           *theme.Theme
}

func NewAppShell(th *theme.Theme) *AppShell {
	return &AppShell{
		SidebarWidth: 256,
		th:           th,
	}
}

func (a *AppShell) WithNavbar(w layout.Widget) *AppShell {
	a.Navbar = w
	return a
}

func (a *AppShell) WithSidebar(w layout.Widget, width unit.Dp) *AppShell {
	a.Sidebar = w
	a.SidebarWidth = width
	return a
}

func (a *AppShell) WithContent(w layout.Widget) *AppShell {
	a.Content = w
	return a
}

// Layout renders the full app shell.
func (a *AppShell) Layout(gtx layout.Context) layout.Dimensions {
	th := a.th

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Navbar
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if a.Navbar == nil {
				return layout.Dimensions{}
			}
			return a.Navbar(gtx)
		}),
		// Body: Sidebar + Content
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				// Sidebar
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if a.Sidebar == nil {
						return layout.Dimensions{}
					}
					sideW := gtx.Dp(a.SidebarWidth)
					gtx.Constraints.Min.X = sideW
					gtx.Constraints.Max.X = sideW
					gtx.Constraints.Min.Y = gtx.Constraints.Max.Y

					// Sidebar background
					sz := image.Pt(sideW, gtx.Constraints.Max.Y)
					paint.FillShape(gtx.Ops, th.Base200,
						clip.Rect{Max: sz}.Op(),
					)

					return a.Sidebar(gtx)
				}),
				// Divider
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if a.Sidebar == nil {
						return layout.Dimensions{}
					}
					sz := image.Pt(1, gtx.Constraints.Max.Y)
					paint.FillShape(gtx.Ops, th.Base300,
						clip.Rect{Max: sz}.Op(),
					)
					return layout.Dimensions{Size: sz}
				}),
				// Content
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if a.Content == nil {
						return layout.Dimensions{Size: gtx.Constraints.Max}
					}
					// Fill content bg
					sz := gtx.Constraints.Max
					paint.FillShape(gtx.Ops, th.Base100,
						clip.Rect{Max: sz}.Op(),
					)
					return a.Content(gtx)
				}),
			)
		}),
	)
}

// ============================================================
// Navbar - Top navigation bar
// ============================================================

// Navbar renders a top navigation bar (like DaisyUI navbar).
type Navbar struct {
	Height     unit.Dp
	Background color.NRGBA
	Bordered   bool
	th         *theme.Theme
}

func NewNavbar(th *theme.Theme) *Navbar {
	return &Navbar{
		Height:     56,
		Background: th.Base100,
		Bordered:   true,
		th:         th,
	}
}

func (n *Navbar) Layout(gtx layout.Context, start, center, end layout.Widget) layout.Dimensions {
	th := n.th
	h := gtx.Dp(n.Height)

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := image.Pt(gtx.Constraints.Max.X, h)

			// Background
			paint.FillShape(gtx.Ops, n.Background,
				clip.Rect{Max: sz}.Op(),
			)

			// Bottom border
			if n.Bordered {
				borderRect := image.Rect(0, h-1, sz.X, h)
				paint.FillShape(gtx.Ops, th.Base300,
					clip.Rect(borderRect).Op(),
				)
			}
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.Y = h
			gtx.Constraints.Max.Y = h
			return layout.Inset{Left: th.Space4, Right: th.Space4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Alignment: layout.Middle,
				}.Layout(gtx,
					// Start section
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if start == nil {
							return layout.Dimensions{}
						}
						return start(gtx)
					}),
					// Center spacer + center content
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if center == nil {
							return layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, 0)}
						}
						return layout.Center.Layout(gtx, center)
					}),
					// End section
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if end == nil {
							return layout.Dimensions{}
						}
						return end(gtx)
					}),
				)
			})
		}),
	)
}

// ============================================================
// Sidebar - Side navigation
// ============================================================

// SidebarItem represents a navigation item.
type SidebarItem struct {
	Label  string
	Icon   string // Unicode icon placeholder
	Active bool
	click  widget.Clickable
}

// Sidebar renders a vertical navigation sidebar.
type Sidebar struct {
	Items    []SidebarItem
	Header   layout.Widget
	Footer   layout.Widget
	Width    unit.Dp
	OnSelect func(index int)
	children []layout.FlexChild // reused across frames
	th       *theme.Theme
}

func NewSidebar(th *theme.Theme, items []SidebarItem) *Sidebar {
	return &Sidebar{
		Items: items,
		Width: 256,
		th:    th,
	}
}

func (s *Sidebar) Layout(gtx layout.Context) layout.Dimensions {
	th := s.th

	s.children = s.children[:0]

	// Header
	if s.Header != nil {
		s.children = append(s.children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Bottom: th.Space4}.Layout(gtx, s.Header)
		}))
	}

	// Menu items
	for i := range s.Items {
		i := i
		s.children = append(s.children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			item := &s.Items[i]
			if item.click.Clicked(gtx) && s.OnSelect != nil {
				s.OnSelect(i)
			}
			return s.layoutItem(gtx, item)
		}))
	}

	// Spacer + footer
	if s.Footer != nil {
		s.children = append(s.children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(0, gtx.Constraints.Max.Y)}
		}))
		s.children = append(s.children, layout.Rigid(s.Footer))
	}

	return layout.Inset{
		Top: th.Space4, Left: th.Space3, Right: th.Space3,
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, s.children...)
	})
}

func (s *Sidebar) layoutItem(gtx layout.Context, item *SidebarItem) layout.Dimensions {
	th := s.th
	radius := gtx.Dp(th.RoundedLg)

	return item.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		padding := layout.Inset{
			Top: th.Space2, Bottom: th.Space2,
			Left: th.Space3, Right: th.Space3,
		}
		return layout.Stack{Alignment: layout.W}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				sz := gtx.Constraints.Min
				bg := theme.Transparent
				if item.Active {
					bg = theme.WithAlpha(th.Primary, 25)
				} else if item.click.Hovered() {
					bg = th.Base300
				}
				if bg.A > 0 {
					defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: bg}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
				}
				pointer.CursorPointer.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					fg := th.BaseContent
					weight := font.Normal
					if item.Active {
						fg = th.Primary
						weight = font.SemiBold
					}
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if item.Icon == "" {
								return layout.Dimensions{}
							}
							return layout.Inset{Right: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return drawLabel(gtx, th, item.Icon, fg, th.FontSize, weight)
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return drawLabel(gtx, th, item.Label, fg, th.SmSize, weight)
						}),
					)
				})
			}),
		)
	})
}

// ============================================================
// Modal (DaisyUI modal)
// ============================================================

// Modal renders a centered dialog overlay.
type Modal struct {
	Visible  bool
	MaxWidth unit.Dp
	th       *theme.Theme
}

func NewModal(th *theme.Theme) *Modal {
	return &Modal{MaxWidth: 500, th: th}
}

func (m *Modal) Show() { m.Visible = true }
func (m *Modal) Hide() { m.Visible = false }

// Layout renders the modal overlay with content.
func (m *Modal) Layout(gtx layout.Context, content layout.Widget) layout.Dimensions {
	if !m.Visible {
		return layout.Dimensions{}
	}
	th := m.th

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		// Backdrop
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Max
			paint.FillShape(gtx.Ops, color.NRGBA{A: 128},
				clip.Rect{Max: sz}.Op(),
			)
			return layout.Dimensions{Size: sz}
		}),
		// Dialog
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			maxW := gtx.Dp(m.MaxWidth)
			if gtx.Constraints.Max.X > maxW {
				gtx.Constraints.Max.X = maxW
			}
			gtx.Constraints.Min.X = gtx.Constraints.Max.X

			radius := gtx.Dp(th.Rounded2xl)
			padding := layout.UniformInset(th.Space6)

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					sz := gtx.Constraints.Min
					defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					return layout.Dimensions{Size: sz}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return padding.Layout(gtx, content)
				}),
			)
		}),
	)
}

// ============================================================
// Drawer (slide-in panel)
// ============================================================

type DrawerSide int

const (
	DrawerLeft DrawerSide = iota
	DrawerRight
)

// Drawer renders a slide-in panel overlay.
type Drawer struct {
	Visible bool
	Side    DrawerSide
	Width   unit.Dp
	th      *theme.Theme
}

func NewDrawer(th *theme.Theme) *Drawer {
	return &Drawer{
		Width: 300,
		Side:  DrawerLeft,
		th:    th,
	}
}

func (d *Drawer) Toggle() { d.Visible = !d.Visible }
func (d *Drawer) Open()   { d.Visible = true }
func (d *Drawer) Close()  { d.Visible = false }

func (d *Drawer) Layout(gtx layout.Context, content layout.Widget) layout.Dimensions {
	if !d.Visible {
		return layout.Dimensions{}
	}
	th := d.th
	w := gtx.Dp(d.Width)

	return layout.Stack{}.Layout(gtx,
		// Backdrop
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Max
			paint.FillShape(gtx.Ops, color.NRGBA{A: 100},
				clip.Rect{Max: sz}.Op(),
			)
			return layout.Dimensions{Size: sz}
		}),
		// Drawer panel
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			var offsetX int
			if d.Side == DrawerRight {
				offsetX = gtx.Constraints.Max.X - w
			}
			defer op.Offset(image.Pt(offsetX, 0)).Push(gtx.Ops).Pop()

			sz := image.Pt(w, gtx.Constraints.Max.Y)
			paint.FillShape(gtx.Ops, th.Base100,
				clip.Rect{Max: sz}.Op(),
			)

			// Draw right border for left drawer
			if d.Side == DrawerLeft {
				borderRect := image.Rect(w-1, 0, w, sz.Y)
				paint.FillShape(gtx.Ops, th.Base300,
					clip.Rect(borderRect).Op(),
				)
			}

			cgtx := gtx
			cgtx.Constraints.Min = sz
			cgtx.Constraints.Max = sz
			return content(cgtx)
		}),
	)
}

// ============================================================
// BottomNav - Mobile bottom navigation
// ============================================================

type BottomNavItem struct {
	Label  string
	Icon   string
	Active bool
	click  widget.Clickable
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
						return layout.Flex{
							Axis:      layout.Vertical,
							Alignment: layout.Middle,
							Spacing:   layout.SpaceSides,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								iconSize := unit.Sp(20)
								return drawLabel(gtx, th, item.Icon, fg, iconSize, font.Normal)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return drawLabel(gtx, th, item.Label, fg, th.XsSize, font.Medium)
							}),
						)
					})
				})
			}

			return layout.Flex{Alignment: layout.Middle}.Layout(gtx, bn.children...)
		}),
	)
}

// ============================================================
// Toast - Brief notification
// ============================================================

type ToastPosition int

const (
	ToastBottom ToastPosition = iota
	ToastTop
	ToastTopRight
	ToastBottomRight
)

type Toast struct {
	Text     string
	Variant  int // reuse AlertVariant
	Position ToastPosition
	Visible  bool
	th       *theme.Theme
}

func NewToast(th *theme.Theme) *Toast {
	return &Toast{Position: ToastBottom, th: th}
}

func (t *Toast) Show(text string) {
	t.Text = text
	t.Visible = true
}

func (t *Toast) Layout(gtx layout.Context) layout.Dimensions {
	if !t.Visible || t.Text == "" {
		return layout.Dimensions{}
	}
	th := t.th
	radius := gtx.Dp(th.RoundedLg)

	var alignment layout.Direction
	switch t.Position {
	case ToastTop:
		alignment = layout.N
	case ToastTopRight:
		alignment = layout.NE
	case ToastBottomRight:
		alignment = layout.SE
	default:
		alignment = layout.S
	}

	return alignment.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: th.Space4, Bottom: th.Space4,
			Left: th.Space4, Right: th.Space4,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					sz := gtx.Constraints.Min
					defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: th.Neutral}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					return layout.Dimensions{Size: sz}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top: th.Space3, Bottom: th.Space3,
						Left: th.Space4, Right: th.Space4,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return drawLabel(gtx, th, t.Text, th.NeutralContent, th.SmSize, font.Medium)
					})
				}),
			)
		})
	})
}

// ============================================================
// Breadcrumb
// ============================================================

type Breadcrumb struct {
	Items    []string
	children []layout.FlexChild // reused across frames
	th       *theme.Theme
}

func NewBreadcrumb(th *theme.Theme, items ...string) *Breadcrumb {
	return &Breadcrumb{Items: items, th: th}
}

func (b *Breadcrumb) Layout(gtx layout.Context) layout.Dimensions {
	th := b.th
	b.children = b.children[:0]

	for i, item := range b.Items {
		i, item := i, item
		if i > 0 {
			b.children = append(b.children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: th.Space2, Right: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawLabel(gtx, th, "/", theme.Opacity(th.BaseContent, 0.4), th.SmSize, font.Normal)
				})
			}))
		}
		b.children = append(b.children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			col := th.BaseContent
			weight := font.Normal
			if i == len(b.Items)-1 {
				weight = font.Medium
			} else {
				col = theme.Opacity(col, 0.6)
			}
			return drawLabel(gtx, th, item, col, th.SmSize, weight)
		}))
	}

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx, b.children...)
}

// ============================================================
// Helpers
// ============================================================

func drawLabel(gtx layout.Context, th *theme.Theme, txt string, col color.NRGBA, size unit.Sp, weight font.Weight) layout.Dimensions {
	lbl := widget.Label{MaxLines: 1}
	f := font.Font{Weight: weight}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	if th.Shaper == nil {
		th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Collection()))
	}
	return lbl.Layout(gtx, th.Shaper, f, size, txt, op.CallOp{})
}
