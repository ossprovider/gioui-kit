package component

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"github.com/hongshengjie/gioui-kit/theme"
)

// Table is a DaisyUI-style data table component.
type Table struct {
	Headers []string
	Rows    [][]string
	Zebra   bool // alternating row background
	Bordered bool
	th      *theme.Theme
}

// NewTable creates a new table component.
func NewTable(th *theme.Theme, headers []string, rows [][]string) *Table {
	return &Table{Headers: headers, Rows: rows, th: th}
}

// WithZebra enables alternating row colors.
func (t *Table) WithZebra() *Table {
	t.Zebra = true
	return t
}

// WithBorder adds a border to the table.
func (t *Table) WithBorder() *Table {
	t.Bordered = true
	return t
}

// Layout renders the table.
func (t *Table) Layout(gtx layout.Context) layout.Dimensions {
	th := t.th
	cellPad := layout.Inset{Top: th.Space3, Bottom: th.Space3, Left: th.Space4, Right: th.Space4}
	radius := gtx.Dp(th.RoundedXl)

	cols := len(t.Headers)
	if cols == 0 {
		return layout.Dimensions{}
	}

	renderRow := func(gtx layout.Context, cells []string, isHeader bool, rowIdx int) layout.Dimensions {
		rowChildren := make([]layout.FlexChild, cols)
		for j := 0; j < cols; j++ {
			j := j
			cell := ""
			if j < len(cells) {
				cell = cells[j]
			}
			rowChildren[j] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return cellPad.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					col := th.BaseContent
					w := font.Normal
					if isHeader {
						col = theme.Opacity(th.BaseContent, 0.7)
						w = font.SemiBold
					}
					return drawText(gtx, th, cell, col, th.SmSize, w)
				})
			})
		}

		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				sz := gtx.Constraints.Min
				bg := th.Base100
				if isHeader {
					bg = th.Base200
				} else if t.Zebra && rowIdx%2 == 1 {
					bg = th.Base200
				}
				defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
				paint.ColorOp{Color: bg}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{}.Layout(gtx, rowChildren...)
			}),
		)
	}

	allRows := make([]layout.FlexChild, 0, 1+len(t.Rows))
	// Header row
	allRows = append(allRows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return renderRow(gtx, t.Headers, true, -1)
	}))
	// Divider
	allRows = append(allRows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		lh := gtx.Dp(1)
		sz := image.Pt(gtx.Constraints.Max.X, lh)
		defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		return layout.Dimensions{Size: sz}
	}))
	// Data rows
	for i, row := range t.Rows {
		i, row := i, row
		allRows = append(allRows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return renderRow(gtx, row, false, i)
		}))
		if i < len(t.Rows)-1 {
			allRows = append(allRows, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lh := gtx.Dp(1)
				sz := image.Pt(gtx.Constraints.Max.X, lh)
				defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
				paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			}))
		}
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			if t.Bordered {
				paint.FillShape(gtx.Ops, th.Base300,
					clip.Stroke{
						Path:  clip.UniformRRect(image.Rectangle{Max: sz}, radius).Path(gtx.Ops),
						Width: float32(gtx.Dp(1)),
					}.Op(),
				)
			}
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx, allRows...)
		}),
	)
}
