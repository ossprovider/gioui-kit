package component

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Select is a DaisyUI-style select/dropdown component.
// It expands inline to show options when opened.
type Select struct {
	Items    []string
	Selected int // index of selected item
	open     bool
	trigger  widget.Clickable
	clicks   []widget.Clickable
	th       *theme.Theme
}

// NewSelect creates a new select component.
func NewSelect(th *theme.Theme, items []string) *Select {
	return &Select{
		Items:    items,
		Selected: 0,
		clicks:   make([]widget.Clickable, len(items)),
		th:       th,
	}
}

// Value returns the currently selected item string.
func (s *Select) Value() string {
	if len(s.Items) == 0 {
		return ""
	}
	return s.Items[s.Selected]
}

// Layout renders the select component.
func (s *Select) Layout(gtx layout.Context) layout.Dimensions {
	th := s.th
	radius := gtx.Dp(th.RoundedLg)
	padding := layout.Inset{Top: th.Space2, Bottom: th.Space2, Left: th.Space3, Right: th.Space3}

	// Handle trigger click
	if s.trigger.Clicked(gtx) {
		s.open = !s.open
	}
	// Handle option clicks
	for i := range s.clicks {
		if s.clicks[i].Clicked(gtx) {
			s.Selected = i
			s.open = false
		}
	}

	label := "Select..."
	if len(s.Items) > 0 {
		label = s.Items[s.Selected]
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Trigger button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return s.trigger.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Stack{Alignment: layout.Center}.Layout(gtx,
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
						pointer.CursorPointer.Add(gtx.Ops)
						return layout.Dimensions{Size: sz}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
									return drawText(gtx, th, label, th.BaseContent, th.FontSize, font.Normal)
								}),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									arrow := "▾"
									if s.open {
										arrow = "▴"
									}
									return layout.Inset{Left: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
										return drawText(gtx, th, arrow, th.BaseContent, th.SmSize, font.Normal)
									})
								}),
							)
						})
					}),
				)
			})
		}),
		// Dropdown list
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !s.open {
				return layout.Dimensions{}
			}
			return layout.Inset{Top: th.Space1}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						children := make([]layout.FlexChild, len(s.Items))
						for i, item := range s.Items {
							i, item := i, item
							children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								isSelected := i == s.Selected
								return s.clicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									bg := th.Base100
									if isSelected {
										bg = theme.WithAlpha(th.Primary, 20)
									}
									if s.clicks[i].Hovered() {
										bg = theme.WithAlpha(th.Primary, 15)
									}
									return layout.Stack{}.Layout(gtx,
										layout.Expanded(func(gtx layout.Context) layout.Dimensions {
											sz := gtx.Constraints.Min
											defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
											paint.ColorOp{Color: bg}.Add(gtx.Ops)
											paint.PaintOp{}.Add(gtx.Ops)
											pointer.CursorPointer.Add(gtx.Ops)
											return layout.Dimensions{Size: sz}
										}),
										layout.Stacked(func(gtx layout.Context) layout.Dimensions {
											return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints.Min.X = gtx.Constraints.Max.X
												w := font.Normal
												col := th.BaseContent
												if isSelected {
													w = font.SemiBold
													col = th.Primary
												}
												return drawText(gtx, th, item, col, th.FontSize, w)
											})
										}),
									)
								})
							})
						}
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
					}),
				)
			})
		}),
	)
}
