package scaffold

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// SidebarItem represents a navigation item.
type SidebarItem struct {
	Label    string
	Icon     string        // Unicode/emoji icon (used when IconData is nil)
	IconData *widget.Icon  // Material Design iconvg icon (takes priority over Icon)
	Active   bool
	click    widget.Clickable
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
							if item.IconData != nil {
								return layout.Inset{Right: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									iconSz := gtx.Dp(20)
									gtx.Constraints = layout.Exact(image.Pt(iconSz, iconSz))
									return item.IconData.Layout(gtx, fg)
								})
							}
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
