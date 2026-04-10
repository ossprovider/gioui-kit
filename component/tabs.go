package component

import (
	"image"
	"strconv"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

type TabVariant int

const (
	TabBoxed TabVariant = iota
	TabBordered
	TabLifted
)

// Tabs manages a tabbed interface backed by widget.Enum for selection state.
type Tabs struct {
	Items    []string
	Enum     widget.Enum
	Variant  TabVariant
	children []layout.FlexChild // reused across frames to avoid per-frame alloc
	th       *theme.Theme
}

func NewTabs(th *theme.Theme, items []string) *Tabs {
	t := &Tabs{
		Items: items,
		th:    th,
	}
	// Default to first tab selected.
	if len(items) > 0 {
		t.Enum.Value = "0"
	}
	return t
}

// Selected returns the index of the currently active tab.
func (t *Tabs) Selected() int {
	i, _ := strconv.Atoi(t.Enum.Value)
	return i
}

func (t *Tabs) Layout(gtx layout.Context) layout.Dimensions {
	th := t.th

	if cap(t.children) < len(t.Items) {
		t.children = make([]layout.FlexChild, len(t.Items))
	}
	t.children = t.children[:len(t.Items)]
	for i, item := range t.Items {
		i, item := i, item
		key := strconv.Itoa(i)
		t.children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			isActive := t.Enum.Value == key
			return t.Enum.Layout(gtx, key, func(gtx layout.Context) layout.Dimensions {
				padding := layout.Inset{
					Top: th.Space2, Bottom: th.Space2,
					Left: th.Space4, Right: th.Space4,
				}
				return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{Alignment: layout.S}.Layout(gtx,
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							fg := th.BaseContent
							weight := font.Normal
							if isActive {
								fg = th.Primary
								weight = font.SemiBold
							}
							return drawText(gtx, th, item, fg, th.SmSize, weight)
						}),
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							if !isActive {
								return layout.Dimensions{}
							}
							sz := gtx.Constraints.Min
							// Active indicator underline
							indicatorH := gtx.Dp(2)
							indicatorRect := image.Rect(0, sz.Y-indicatorH, sz.X, sz.Y)
							paint.FillShape(gtx.Ops, th.Primary,
								clip.Rect(indicatorRect).Op(),
							)
							return layout.Dimensions{Size: sz}
						}),
					)
				})
			})
		})
	}

	return layout.Flex{}.Layout(gtx, t.children...)
}
