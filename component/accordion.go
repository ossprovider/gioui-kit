package component

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// AccordionItem is a single collapsible section in an Accordion.
type AccordionItem struct {
	Title  string
	open   bool
	click  widget.Clickable
}

// NewAccordionItem creates an accordion item with a title.
func NewAccordionItem(title string) *AccordionItem {
	return &AccordionItem{Title: title}
}

// Open returns whether this item is expanded.
func (a *AccordionItem) Open() bool { return a.open }

// Accordion is a DaisyUI-style collapsible accordion component.
type Accordion struct {
	Items []*AccordionItem
	th    *theme.Theme
}

// NewAccordion creates a new accordion.
func NewAccordion(th *theme.Theme, items ...*AccordionItem) *Accordion {
	return &Accordion{Items: items, th: th}
}

// Layout renders the accordion.
func (a *Accordion) Layout(gtx layout.Context, bodies []layout.Widget) layout.Dimensions {
	th := a.th
	radius := gtx.Dp(th.RoundedLg)

	// Handle clicks
	for i := range a.Items {
		if a.Items[i].click.Clicked(gtx) {
			a.Items[i].open = !a.Items[i].open
		}
	}

	n := len(a.Items)
	if len(bodies) < n {
		n = len(bodies)
	}

	children := make([]layout.FlexChild, n)
	for i := 0; i < n; i++ {
		i := i
		item := a.Items[i]
		body := bodies[i]
		isLast := i == n-1

		children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			borderBottom := !isLast

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					sz := gtx.Constraints.Min
					// Top-rounded only if first, bottom-rounded only if last and closed
					rr := 0
					if i == 0 {
						rr = radius
					}
					_ = rr
					defer clip.UniformRRect(image.Rectangle{Max: sz}, 0).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					// Bottom border
					if borderBottom {
						lh := gtx.Dp(1)
						rect := image.Rect(0, sz.Y-lh, sz.X, sz.Y)
						defer clip.Rect(rect).Push(gtx.Ops).Pop()
						paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
					}
					return layout.Dimensions{Size: sz}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						// Header
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return item.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								padding := layout.Inset{
									Top: th.Space3, Bottom: th.Space3,
									Left: th.Space4, Right: th.Space4,
								}
								return layout.Stack{Alignment: layout.Center}.Layout(gtx,
									layout.Expanded(func(gtx layout.Context) layout.Dimensions {
										sz := gtx.Constraints.Min
										if item.click.Hovered() {
											defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
											paint.ColorOp{Color: theme.WithAlpha(th.Primary, 10)}.Add(gtx.Ops)
											paint.PaintOp{}.Add(gtx.Ops)
										}
										pointer.CursorPointer.Add(gtx.Ops)
										return layout.Dimensions{Size: sz}
									}),
									layout.Stacked(func(gtx layout.Context) layout.Dimensions {
										return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
											return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
												layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
													return drawText(gtx, th, item.Title, th.BaseContent, th.FontSize, font.SemiBold)
												}),
												layout.Rigid(func(gtx layout.Context) layout.Dimensions {
													sz := gtx.Dp(10)
													return drawTriangle(gtx.Ops, item.open, theme.Opacity(th.BaseContent, 0.5), sz)
												}),
											)
										})
									}),
								)
							})
						}),
						// Body (when open)
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if !item.open {
								return layout.Dimensions{}
							}
							return layout.Inset{
								Left: th.Space4, Right: th.Space4,
								Bottom: th.Space4,
							}.Layout(gtx, body)
						}),
					)
				}),
			)
		})
	}

	// Wrap in rounded border
	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			bw := gtx.Dp(1)
			inner := image.Rect(bw, bw, sz.X-bw, sz.Y-bw)
			if inner.Dx() > 0 && inner.Dy() > 0 {
				ir := radius - bw
				if ir < 0 {
					ir = 0
				}
				defer clip.UniformRRect(inner, ir).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
			}
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
		}),
	)
}

// drawTriangle draws a filled up or down triangle of the given pixel size.
func drawTriangle(ops *op.Ops, up bool, col color.NRGBA, size int) layout.Dimensions {
	s := float32(size)
	half := s / 2
	var path clip.Path
	path.Begin(ops)
	if up {
		path.MoveTo(f32.Pt(0, s))
		path.LineTo(f32.Pt(s, s))
		path.LineTo(f32.Pt(half, 0))
	} else {
		path.MoveTo(f32.Pt(0, 0))
		path.LineTo(f32.Pt(s, 0))
		path.LineTo(f32.Pt(half, s))
	}
	path.Close()
	paint.FillShape(ops, col, clip.Outline{Path: path.End()}.Op())
	return layout.Dimensions{Size: image.Pt(size, size)}
}
