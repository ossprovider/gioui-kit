// Package layout provides TailwindCSS-inspired layout primitives for Gio UI.
//
// Naming follows Tailwind conventions:
//   - Flex, FlexRow, FlexCol (flex layout)
//   - Grid (grid layout)
//   - Stack (absolute positioning / z-index)
//   - Container, Box (wrapper utilities)
//   - Gap, Padding, Margin modifiers
//
// Alignment uses Tailwind names:
//   - ItemsCenter, ItemsStart, ItemsEnd, ItemsStretch
//   - JustifyCenter, JustifyBetween, JustifyAround, JustifyEnd
package layout

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"image/color"
)

// ---------- Flex Layout (flex-row / flex-col) ----------

// FlexRow lays out children horizontally (like Tailwind `flex flex-row`).
type FlexRow struct {
	Gap       unit.Dp
	Alignment layout.Alignment // cross-axis: ItemsStart, ItemsCenter, ItemsEnd
	Spacing   layout.Spacing   // main-axis: JustifyStart, JustifyCenter, etc.
	Wrap      bool
}

// Layout renders children in a horizontal flex row.
func (f FlexRow) Layout(gtx layout.Context, children ...layout.FlexChild) layout.Dimensions {
	flex := layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: f.Alignment,
		Spacing:   f.Spacing,
	}
	if f.Gap > 0 {
		children = insertGaps(gtx, layout.Horizontal, f.Gap, children)
	}
	return flex.Layout(gtx, children...)
}

// FlexCol lays out children vertically (like Tailwind `flex flex-col`).
type FlexCol struct {
	Gap       unit.Dp
	Alignment layout.Alignment
	Spacing   layout.Spacing
}

// Layout renders children in a vertical flex column.
func (f FlexCol) Layout(gtx layout.Context, children ...layout.FlexChild) layout.Dimensions {
	flex := layout.Flex{
		Axis:      layout.Vertical,
		Alignment: f.Alignment,
		Spacing:   f.Spacing,
	}
	if f.Gap > 0 {
		children = insertGaps(gtx, layout.Vertical, f.Gap, children)
	}
	return flex.Layout(gtx, children...)
}

// insertGaps inserts spacer FlexChildren between items to simulate gap.
func insertGaps(gtx layout.Context, axis layout.Axis, gap unit.Dp, children []layout.FlexChild) []layout.FlexChild {
	if len(children) <= 1 {
		return children
	}
	result := make([]layout.FlexChild, 0, len(children)*2-1)
	gapPx := gtx.Dp(gap)
	spacer := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		if axis == layout.Horizontal {
			return layout.Dimensions{Size: image.Pt(gapPx, 0)}
		}
		return layout.Dimensions{Size: image.Pt(0, gapPx)}
	})
	for i, child := range children {
		if i > 0 {
			result = append(result, spacer)
		}
		result = append(result, child)
	}
	return result
}

// ---------- Flex child helpers (like Tailwind grow/shrink) ----------

// Rigid wraps layout.Rigid for convenience.
func Rigid(w layout.Widget) layout.FlexChild {
	return layout.Rigid(w)
}

// Flexed wraps layout.Flexed (like flex-grow).
func Flexed(weight float32, w layout.Widget) layout.FlexChild {
	return layout.Flexed(weight, w)
}

// Grow is shorthand for Flexed(1, w) (like Tailwind `flex-1`).
func Grow(w layout.Widget) layout.FlexChild {
	return layout.Flexed(1, w)
}

// ---------- Alignment constants (Tailwind naming) ----------

const (
	// Cross-axis alignment (items-*)
	ItemsStart   = layout.Start
	ItemsCenter  = layout.Middle
	ItemsEnd     = layout.End
	ItemsBaseline = layout.Baseline

	// Main-axis spacing (justify-*)
	JustifyStart   = layout.SpaceStart
	JustifyEnd     = layout.SpaceEnd
	JustifyCenter  = layout.SpaceSides
	JustifyBetween = layout.SpaceEvenly
	JustifyAround  = layout.SpaceAround
)

// ---------- Box / Container ----------

// Box is a simple rectangular container with optional padding, bg, and radius.
// Similar to a <div> with Tailwind classes like `p-4 bg-white rounded-lg`.
type Box struct {
	Padding    layout.Inset
	Background color.NRGBA
	Radius     unit.Dp
	MinWidth   unit.Dp
	MinHeight  unit.Dp
	MaxWidth   unit.Dp
	Border     Border
}

// Border represents a CSS-like border.
type Border struct {
	Color color.NRGBA
	Width unit.Dp
}

// Layout renders the box.
func (b Box) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	// Apply max width constraint
	if b.MaxWidth > 0 {
		maxPx := gtx.Dp(b.MaxWidth)
		if gtx.Constraints.Max.X > maxPx {
			gtx.Constraints.Max.X = maxPx
		}
	}

	// Apply min dimensions
	if b.MinWidth > 0 {
		minPx := gtx.Dp(b.MinWidth)
		if gtx.Constraints.Min.X < minPx {
			gtx.Constraints.Min.X = minPx
		}
	}
	if b.MinHeight > 0 {
		minPx := gtx.Dp(b.MinHeight)
		if gtx.Constraints.Min.Y < minPx {
			gtx.Constraints.Min.Y = minPx
		}
	}

	return b.Padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			// Background
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				sz := gtx.Constraints.Min
				rr := gtx.Dp(b.Radius)
				rrect := clip.UniformRRect(image.Rectangle{Max: sz}, rr)
				defer rrect.Push(gtx.Ops).Pop()

				// Fill background
				if b.Background.A > 0 {
					paint.ColorOp{Color: b.Background}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
				}

				// Draw border
				if b.Border.Width > 0 {
					bw := gtx.Dp(b.Border.Width)
					drawBorder(gtx.Ops, sz, rr, bw, b.Border.Color)
				}

				return layout.Dimensions{Size: sz}
			}),
			// Content
			layout.Stacked(w),
		)
	})
}

// drawBorder draws a rectangular border.
func drawBorder(ops *op.Ops, sz image.Point, radius, width int, col color.NRGBA) {
	// Simple border: draw outer rect, then mask inner rect
	// For simplicity, we use stroke-based approach
	r := image.Rectangle{Max: sz}
	paint.FillShape(ops, col,
		clip.Stroke{
			Path:  clip.UniformRRect(r, radius).Path(ops),
			Width: float32(width),
		}.Op(),
	)
}

// Container is a centered max-width container (like Tailwind `container mx-auto`).
type Container struct {
	MaxWidth unit.Dp
	Padding  unit.Dp
}

// Layout centers the content with max-width.
func (c Container) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	maxW := gtx.Dp(c.MaxWidth)
	if maxW <= 0 {
		maxW = gtx.Dp(1280) // default xl breakpoint
	}
	pad := gtx.Dp(c.Padding)

	return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if gtx.Constraints.Max.X > maxW {
				gtx.Constraints.Max.X = maxW
			}
			return layout.Inset{
				Left:  unit.Dp(pad),
				Right: unit.Dp(pad),
			}.Layout(gtx, w)
		}),
	)
}

// ---------- Spacing Utilities (Tailwind p-*, m-*) ----------

// P returns a uniform Inset (like Tailwind `p-4`).
func P(dp unit.Dp) layout.Inset {
	return layout.UniformInset(dp)
}

// Px returns horizontal padding (like Tailwind `px-4`).
func Px(dp unit.Dp) layout.Inset {
	return layout.Inset{Left: dp, Right: dp}
}

// Py returns vertical padding (like Tailwind `py-4`).
func Py(dp unit.Dp) layout.Inset {
	return layout.Inset{Top: dp, Bottom: dp}
}

// Pt returns top padding.
func Pt(dp unit.Dp) layout.Inset {
	return layout.Inset{Top: dp}
}

// Pb returns bottom padding.
func Pb(dp unit.Dp) layout.Inset {
	return layout.Inset{Bottom: dp}
}

// Pl returns left padding.
func Pl(dp unit.Dp) layout.Inset {
	return layout.Inset{Left: dp}
}

// Pr returns right padding.
func Pr(dp unit.Dp) layout.Inset {
	return layout.Inset{Right: dp}
}

// Inset4 creates an inset with all 4 sides specified.
func Inset4(top, right, bottom, left unit.Dp) layout.Inset {
	return layout.Inset{Top: top, Right: right, Bottom: bottom, Left: left}
}

// ---------- Size Utilities ----------

// W sets a fixed width constraint.
func W(dp unit.Dp) func(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return func(gtx layout.Context, w layout.Widget) layout.Dimensions {
		px := gtx.Dp(dp)
		gtx.Constraints.Min.X = px
		gtx.Constraints.Max.X = px
		return w(gtx)
	}
}

// H sets a fixed height constraint.
func H(dp unit.Dp) func(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return func(gtx layout.Context, w layout.Widget) layout.Dimensions {
		px := gtx.Dp(dp)
		gtx.Constraints.Min.Y = px
		gtx.Constraints.Max.Y = px
		return w(gtx)
	}
}

// MinW sets a minimum width.
func MinW(dp unit.Dp) func(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return func(gtx layout.Context, w layout.Widget) layout.Dimensions {
		px := gtx.Dp(dp)
		if gtx.Constraints.Min.X < px {
			gtx.Constraints.Min.X = px
		}
		return w(gtx)
	}
}

// MaxW sets a maximum width.
func MaxW(dp unit.Dp) func(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return func(gtx layout.Context, w layout.Widget) layout.Dimensions {
		px := gtx.Dp(dp)
		if gtx.Constraints.Max.X > px {
			gtx.Constraints.Max.X = px
		}
		return w(gtx)
	}
}

// WFull forces full width (like `w-full`).
func WFull(gtx layout.Context, w layout.Widget) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	return w(gtx)
}

// ---------- Grid Layout ----------

// Grid provides a simple grid layout (like Tailwind `grid grid-cols-N gap-N`).
type Grid struct {
	Cols int
	Gap  unit.Dp
}

// Layout renders widgets in a grid.
func (g Grid) Layout(gtx layout.Context, widgets ...layout.Widget) layout.Dimensions {
	if g.Cols <= 0 {
		g.Cols = 1
	}
	gapPx := gtx.Dp(g.Gap)
	totalGap := gapPx * (g.Cols - 1)
	colWidth := (gtx.Constraints.Max.X - totalGap) / g.Cols

	rows := (len(widgets) + g.Cols - 1) / g.Cols
	var totalHeight int

	for row := 0; row < rows; row++ {
		if row > 0 {
			totalHeight += gapPx
		}
		var rowHeight int
		for col := 0; col < g.Cols; col++ {
			idx := row*g.Cols + col
			if idx >= len(widgets) {
				break
			}
			// Position and size each cell
			offX := col * (colWidth + gapPx)
			macro := op.Record(gtx.Ops)
			cgtx := gtx
			cgtx.Constraints.Min.X = colWidth
			cgtx.Constraints.Max.X = colWidth
			dims := widgets[idx](cgtx)
			call := macro.Stop()

			offset := op.Offset(image.Pt(offX, totalHeight)).Push(gtx.Ops)
			call.Add(gtx.Ops)
			offset.Pop()

			if dims.Size.Y > rowHeight {
				rowHeight = dims.Size.Y
			}
		}
		totalHeight += rowHeight
	}

	return layout.Dimensions{
		Size: image.Pt(gtx.Constraints.Max.X, totalHeight),
	}
}

// ---------- Stack (z-index / absolute positioning) ----------

// Center centers a widget within its parent constraints.
func Center(gtx layout.Context, w layout.Widget) layout.Dimensions {
	return layout.Center.Layout(gtx, w)
}

// ---------- Divider ----------

// DividerH draws a horizontal divider line.
type DividerH struct {
	Color     color.NRGBA
	Thickness unit.Dp
	Inset     layout.Inset
}

func (d DividerH) Layout(gtx layout.Context) layout.Dimensions {
	thickness := gtx.Dp(d.Thickness)
	if thickness <= 0 {
		thickness = 1
	}
	return d.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		sz := image.Pt(gtx.Constraints.Max.X, thickness)
		paint.FillShape(gtx.Ops, d.Color,
			clip.Rect{Max: sz}.Op(),
		)
		return layout.Dimensions{Size: sz}
	})
}

// ---------- Spacer ----------

// SpaceH creates a horizontal spacer (width).
func SpaceH(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(gtx.Dp(dp), 0)}
	}
}

// SpaceV creates a vertical spacer (height).
func SpaceV(dp unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(0, gtx.Dp(dp))}
	}
}

// ---------- List (scrollable) ----------

// ScrollY is a convenience wrapper around layout.List for vertical scrolling.
type ScrollY struct {
	List layout.List
}

// NewScrollY creates a new vertical scrollable list.
func NewScrollY() *ScrollY {
	return &ScrollY{
		List: layout.List{
			Axis: layout.Vertical,
		},
	}
}

// Layout renders a scrollable list of widgets.
func (s *ScrollY) Layout(gtx layout.Context, count int, w layout.ListElement) layout.Dimensions {
	return s.List.Layout(gtx, count, w)
}

// ScrollX is a horizontal scrollable list.
type ScrollX struct {
	List layout.List
}

func NewScrollX() *ScrollX {
	return &ScrollX{
		List: layout.List{
			Axis: layout.Horizontal,
		},
	}
}

func (s *ScrollX) Layout(gtx layout.Context, count int, w layout.ListElement) layout.Dimensions {
	return s.List.Layout(gtx, count, w)
}
