package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// MenuItem represents a single entry in a Menu.
type MenuItem struct {
	Label    string
	Icon     *widget.Icon
	Active   bool
	Disabled bool
	click    widget.Clickable
}

// NewMenuItem creates a menu item with a label.
func NewMenuItem(label string) *MenuItem {
	return &MenuItem{Label: label}
}

// WithIcon sets a real icon for the menu item.
func (m *MenuItem) WithIcon(icon *widget.Icon) *MenuItem {
	m.Icon = icon
	return m
}

// Clicked reports whether the menu item was clicked this frame.
func (m *MenuItem) Clicked(gtx layout.Context) bool {
	return m.click.Clicked(gtx)
}

// Menu is a DaisyUI-style vertical menu component.
type Menu struct {
	Items    []*MenuItem
	Compact  bool
	Bordered bool
	Rounded  bool
	th       *theme.Theme
}

// NewMenu creates a new menu.
func NewMenu(th *theme.Theme, items ...*MenuItem) *Menu {
	return &Menu{Items: items, Rounded: true, th: th}
}

// WithCompact reduces item padding.
func (m *Menu) WithCompact() *Menu {
	m.Compact = true
	return m
}

// WithBorder adds a border around the menu.
func (m *Menu) WithBorder() *Menu {
	m.Bordered = true
	return m
}

// Layout renders the menu.
func (m *Menu) Layout(gtx layout.Context) layout.Dimensions {
	th := m.th
	radius := gtx.Dp(th.RoundedLg)
	padding := layout.Inset{Top: th.Space2, Bottom: th.Space2, Left: th.Space3, Right: th.Space3}
	if m.Compact {
		padding = layout.Inset{Top: th.Space1, Bottom: th.Space1, Left: th.Space2, Right: th.Space2}
	}

	children := make([]layout.FlexChild, len(m.Items))
	for i, item := range m.Items {
		item := item
		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if item.Disabled {
				return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return menuItemContent(gtx, th, item, theme.Opacity(th.BaseContent, 0.35))
				})
			}

			return item.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						sz := gtx.Constraints.Min
						itemRadius := 0
						if m.Rounded {
							itemRadius = gtx.Dp(th.RoundedMd)
						}
						bg := th.Base100
						if item.Active {
							bg = theme.WithAlpha(th.Primary, 20)
						} else if item.click.Hovered() {
							bg = theme.WithAlpha(th.Primary, 10)
						}
						if sz.X > 0 && sz.Y > 0 {
							defer clip.UniformRRect(image.Rectangle{Max: sz}, itemRadius).Push(gtx.Ops).Pop()
							paint.ColorOp{Color: bg}.Add(gtx.Ops)
							paint.PaintOp{}.Add(gtx.Ops)
						}
						pointer.CursorPointer.Add(gtx.Ops)
						return layout.Dimensions{Size: sz}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Constraints.Max.X
							col := th.BaseContent
							if item.Active {
								col = th.Primary
							}
							return menuItemContent(gtx, th, item, col)
						})
					}),
				)
			})
		})
	}

	inner := func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: th.Space1, Bottom: th.Space1}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		})
	}

	if !m.Bordered {
		return inner(gtx)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			paint.FillShape(gtx.Ops, th.Base300,
				clip.Stroke{
					Path:  clip.UniformRRect(image.Rectangle{Max: sz}, radius).Path(gtx.Ops),
					Width: float32(gtx.Dp(1)),
				}.Op(),
			)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(inner),
	)
}

func menuItemContent(gtx layout.Context, th *theme.Theme, item *MenuItem, col color.NRGBA) layout.Dimensions {
	w := font.Normal
	if item.Active {
		w = font.SemiBold
	}
	if item.Icon == nil {
		return drawText(gtx, th, item.Label, col, th.FontSize, w)
	}
	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Right: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				iconPx := gtx.Sp(th.FontSize)
				gtx.Constraints = layout.Exact(image.Pt(iconPx, iconPx))
				return item.Icon.Layout(gtx, col)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawText(gtx, th, item.Label, col, th.FontSize, w)
		}),
	)
}
