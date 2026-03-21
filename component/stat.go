package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Stat is a DaisyUI-style stat component for displaying metrics.
type Stat struct {
	Title  string
	Value  string
	Desc   string
	Figure string // optional emoji/icon shown on the right
	bg     color.NRGBA
	th     *theme.Theme
}

// NewStat creates a new stat component.
func NewStat(th *theme.Theme, title, value string) *Stat {
	return &Stat{Title: title, Value: value, th: th}
}

// WithDesc sets the description text below the value.
func (s *Stat) WithDesc(desc string) *Stat {
	s.Desc = desc
	return s
}

// WithFigure sets an icon/emoji shown on the right.
func (s *Stat) WithFigure(figure string) *Stat {
	s.Figure = figure
	return s
}

// WithBg sets a background color for the stat card.
func (s *Stat) WithBg(c color.NRGBA) *Stat {
	s.bg = c
	return s
}

// Layout renders the stat card.
func (s *Stat) Layout(gtx layout.Context) layout.Dimensions {
	th := s.th
	radius := gtx.Dp(th.RoundedXl)
	padding := layout.Inset{
		Top: th.Space4, Bottom: th.Space4,
		Left: th.Space6, Right: th.Space6,
	}
	bg := s.bg
	if bg == (color.NRGBA{}) {
		bg = th.Base100
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: bg}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return drawText(gtx, th, s.Title, theme.Opacity(th.BaseContent, 0.6), th.SmSize, font.Medium)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Top: th.Space1, Bottom: th.Space1}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return drawText(gtx, th, s.Value, th.BaseContent, th.H2Size, font.Bold)
								})
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if s.Desc == "" {
									return layout.Dimensions{}
								}
								return drawText(gtx, th, s.Desc, theme.Opacity(th.BaseContent, 0.5), th.XsSize, font.Normal)
							}),
						)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if s.Figure == "" {
							return layout.Dimensions{}
						}
						return layout.Inset{Left: th.Space4}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return drawText(gtx, th, s.Figure, theme.Opacity(th.BaseContent, 0.3), th.H1Size, font.Normal)
						})
					}),
				)
			})
		}),
	)
}

// StatGroup renders multiple stats horizontally.
type StatGroup struct {
	Stats    []*Stat
	Bordered bool
	th       *theme.Theme
}

// NewStatGroup creates a group of stats displayed side by side.
func NewStatGroup(th *theme.Theme, stats ...*Stat) *StatGroup {
	return &StatGroup{Stats: stats, th: th}
}

// WithBorder adds a border around the stat group.
func (sg *StatGroup) WithBorder() *StatGroup {
	sg.Bordered = true
	return sg
}

// Layout renders the stat group.
func (sg *StatGroup) Layout(gtx layout.Context) layout.Dimensions {
	th := sg.th
	radius := gtx.Dp(th.RoundedXl)

	inner := func(gtx layout.Context) layout.Dimensions {
		children := make([]layout.FlexChild, len(sg.Stats))
		for i, s := range sg.Stats {
			s := s
			children[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return s.Layout(gtx)
			})
		}
		return layout.Flex{}.Layout(gtx, children...)
	}

	if !sg.Bordered {
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
