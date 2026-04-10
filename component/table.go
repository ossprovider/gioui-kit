package component

import (
	"fmt"
	"image"
	"sort"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// TableCol configures a single column.
type TableCol struct {
	Title    string
	Flex     float32 // relative width weight; 0 defaults to 1
	Sortable bool
}

// TableAction is a button shown in the selection action bar when rows are selected.
type TableAction struct {
	Label   string
	Variant BtnVariant
	// OnClick is called with the original (pre-sort) indices of selected rows.
	// After the call the selection is NOT cleared automatically — clear it by
	// modifying t.selected if needed.
	OnClick func(rows []int)
	click   widget.Clickable
}

// NewTableAction creates a TableAction with a primary variant.
func NewTableAction(label string, fn func([]int)) *TableAction {
	return &TableAction{Label: label, Variant: BtnPrimary, OnClick: fn}
}

// WithVariant sets the button variant for this action.
func (a *TableAction) WithVariant(v BtnVariant) *TableAction { a.Variant = v; return a }

// Table is a DaisyUI-style advanced data table with sorting, row selection, and virtual scroll.
type Table struct {
	Cols     []TableCol
	Rows     [][]string

	Zebra      bool
	Bordered   bool
	Compact    bool
	Selectable bool    // prepend a checkbox column for row selection
	MaxHeight  unit.Dp // 0 = natural height; > 0 = fixed height with internal scroll

	// OnRowClick is called with the original (pre-sort) row index on row click.
	OnRowClick func(row int)

	// actions shown in the selection toolbar (only when Selectable == true).
	actions []*TableAction

	// internal state
	sortCol      int
	sortAsc      bool
	sortOrder    []int
	colBtns      []*widget.Clickable
	rowBtns      []*widget.Clickable
	selected     []*widget.Bool
	selectAll    widget.Bool
	prevSelectAll bool   // previous Value of selectAll — used to detect toggles
	list         widget.List
	th           *theme.Theme
}

// NewTable creates a Table from plain string headers (backward-compatible).
func NewTable(th *theme.Theme, headers []string, rows [][]string) *Table {
	cols := make([]TableCol, len(headers))
	for i, h := range headers {
		cols[i] = TableCol{Title: h, Flex: 1}
	}
	return NewDataTable(th, cols, rows)
}

// NewDataTable creates a Table with full column configuration.
func NewDataTable(th *theme.Theme, cols []TableCol, rows [][]string) *Table {
	for i := range cols {
		if cols[i].Flex <= 0 {
			cols[i].Flex = 1
		}
	}
	t := &Table{
		Cols:    cols,
		Rows:    rows,
		sortCol: -1,
		th:      th,
	}
	t.list.Axis = layout.Vertical
	return t
}

func (t *Table) WithZebra() *Table                  { t.Zebra = true; return t }
func (t *Table) WithBorder() *Table                 { t.Bordered = true; return t }
func (t *Table) WithCompact() *Table                { t.Compact = true; return t }
func (t *Table) WithSelectable() *Table             { t.Selectable = true; return t }
func (t *Table) WithMaxHeight(h unit.Dp) *Table     { t.MaxHeight = h; return t }
func (t *Table) WithOnRowClick(fn func(int)) *Table { t.OnRowClick = fn; return t }

// WithActions sets the bulk-action buttons shown in the selection toolbar.
// Requires WithSelectable().
func (t *Table) WithActions(actions ...*TableAction) *Table {
	t.actions = actions
	return t
}

// SelectedRows returns the original row indices of all currently selected rows.
func (t *Table) SelectedRows() []int {
	var out []int
	for i, s := range t.selected {
		if s.Value {
			out = append(out, i)
		}
	}
	return out
}

// ClearSelection deselects all rows and resets the select-all checkbox.
func (t *Table) ClearSelection() {
	for _, s := range t.selected {
		s.Value = false
	}
	t.selectAll.Value = false
	t.prevSelectAll = false
}

// ─── internal state management ──────────────────────────────────────────────

func (t *Table) ensureState() {
	for len(t.colBtns) < len(t.Cols) {
		t.colBtns = append(t.colBtns, new(widget.Clickable))
	}
	n := len(t.Rows)
	for len(t.rowBtns) < n {
		t.rowBtns = append(t.rowBtns, new(widget.Clickable))
	}
	for len(t.selected) < n {
		t.selected = append(t.selected, new(widget.Bool))
	}
	if len(t.sortOrder) != n {
		t.sortOrder = make([]int, n)
		for i := range t.sortOrder {
			t.sortOrder[i] = i
		}
		if t.sortCol >= 0 {
			t.applySort()
		}
	}
}

func (t *Table) applySort() {
	if t.sortCol < 0 || t.sortCol >= len(t.Cols) {
		return
	}
	col, asc := t.sortCol, t.sortAsc
	sort.SliceStable(t.sortOrder, func(i, j int) bool {
		a, b := "", ""
		if col < len(t.Rows[t.sortOrder[i]]) {
			a = t.Rows[t.sortOrder[i]][col]
		}
		if col < len(t.Rows[t.sortOrder[j]]) {
			b = t.Rows[t.sortOrder[j]][col]
		}
		if asc {
			return a < b
		}
		return a > b
	})
}

// processEvents handles column-header sort clicks, row clicks, and action buttons.
func (t *Table) processEvents(gtx layout.Context) {
	// Column sort
	for i, btn := range t.colBtns {
		if i >= len(t.Cols) {
			break
		}
		if t.Cols[i].Sortable && btn.Clicked(gtx) {
			if t.sortCol == i {
				t.sortAsc = !t.sortAsc
			} else {
				t.sortCol = i
				t.sortAsc = true
			}
			t.applySort()
		}
	}
	// Row clicks
	if t.OnRowClick != nil {
		for i, btn := range t.rowBtns {
			if i >= len(t.Rows) {
				break
			}
			if btn.Clicked(gtx) {
				t.OnRowClick(t.sortOrder[i])
			}
		}
	}
	// Action buttons
	sel := t.SelectedRows()
	for _, act := range t.actions {
		if act.click.Clicked(gtx) && act.OnClick != nil {
			act.OnClick(sel)
		}
	}
}

func (t *Table) selectedCount() int {
	n := 0
	for _, s := range t.selected {
		if s.Value {
			n++
		}
	}
	return n
}

func (t *Table) padCell() layout.Inset {
	th := t.th
	if t.Compact {
		return layout.Inset{Top: th.Space2, Bottom: th.Space2, Left: th.Space4, Right: th.Space4}
	}
	return layout.Inset{Top: th.Space3, Bottom: th.Space3, Left: th.Space4, Right: th.Space4}
}

func (t *Table) padCellV() unit.Dp {
	if t.Compact {
		return t.th.Space2
	}
	return t.th.Space3
}

// ─── Layout ─────────────────────────────────────────────────────────────────

func (t *Table) Layout(gtx layout.Context) layout.Dimensions {
	t.ensureState()
	t.processEvents(gtx)

	th := t.th
	radius := gtx.Dp(th.RoundedXl)

	if t.MaxHeight > 0 {
		if maxY := gtx.Dp(t.MaxHeight); maxY < gtx.Constraints.Max.Y {
			gtx.Constraints.Max.Y = maxY
		}
		gtx.Constraints.Min.Y = 0
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
			return t.renderAll(gtx)
		}),
	)
}

func (t *Table) renderAll(gtx layout.Context) layout.Dimensions {
	th := t.th
	nSel := t.selectedCount()

	// Top slot: action bar when rows are selected, otherwise the column header.
	topChild := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		if t.Selectable && nSel > 0 {
			return t.renderActionBar(gtx, nSel)
		}
		return t.renderHeader(gtx)
	})
	dividerChild := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return tableRowDivider(gtx, th)
	})

	if t.MaxHeight > 0 {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			topChild,
			dividerChild,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return t.list.Layout(gtx, len(t.Rows), t.renderRow)
			}),
		)
	}

	children := make([]layout.FlexChild, 0, 2+len(t.Rows)*2)
	children = append(children, topChild, dividerChild)
	for i := range t.Rows {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.renderRow(gtx, i)
		}))
		if i < len(t.Rows)-1 {
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return tableRowDivider(gtx, th)
			}))
		}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

// ─── Action bar ─────────────────────────────────────────────────────────────

// renderActionBar replaces the header row when rows are selected.
func (t *Table) renderActionBar(gtx layout.Context, nSel int) layout.Dimensions {
	th := t.th
	pad := t.padCell()

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: theme.Lerp(th.Base200, th.Primary, 0.08)}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return pad.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					// Selection count label
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := fmt.Sprintf("%d selected", nSel)
						return drawText(gtx, th, label, th.Primary, th.SmSize, font.SemiBold)
					}),
					// Action buttons
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if len(t.actions) == 0 {
							return layout.Dimensions{}
						}
						btns := make([]layout.FlexChild, 0, len(t.actions)*2+1)
						btns = append(btns, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Dimensions{}
						}))
						for _, act := range t.actions {
							act := act
							btns = append(btns,
								layout.Rigid(layout.Spacer{Width: th.Space2}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return NewButton(th, &act.click, act.Label).
										WithVariant(act.Variant).
										WithSize(BtnSm).
										Layout(gtx)
								}),
							)
						}
						return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, btns...)
					}),
				)
			})
		}),
	)
}

// ─── Header ─────────────────────────────────────────────────────────────────

func (t *Table) renderHeader(gtx layout.Context) layout.Dimensions {
	th := t.th

	children := make([]layout.FlexChild, 0, len(t.Cols)+1)
	if t.Selectable {
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.renderSelectAll(gtx)
		}))
	}
	for i, col := range t.Cols {
		children = append(children, layout.Flexed(col.Flex, func(gtx layout.Context) layout.Dimensions {
			return t.padCell().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return t.renderHeaderCell(gtx, i, col)
			})
		}))
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base200}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return layout.Flex{}.Layout(gtx, children...)
		}),
	)
}

func (t *Table) renderHeaderCell(gtx layout.Context, i int, col TableCol) layout.Dimensions {
	th := t.th
	label := col.Title
	fg := theme.Opacity(th.BaseContent, 0.7)

	if t.sortCol == i {
		if t.sortAsc {
			label += " ↑"
		} else {
			label += " ↓"
		}
	}

	if col.Sortable && i < len(t.colBtns) {
		btn := t.colBtns[i]
		return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			if btn.Hovered() {
				fg = th.Primary
			}
			return drawText(gtx, th, label, fg, th.SmSize, font.SemiBold)
		})
	}
	return drawText(gtx, th, label, fg, th.SmSize, font.SemiBold)
}

// renderSelectAll draws the "select all" checkbox in the header.
// Bug fix: widget.Bool.Layout already calls Update internally, so we cannot
// call Update again to detect changes. Instead we snapshot Value before
// rendering and compare after.
func (t *Table) renderSelectAll(gtx layout.Context) layout.Dimensions {
	th := t.th
	pad := layout.Inset{Top: t.padCellV(), Bottom: t.padCellV(), Left: th.Space4, Right: th.Space2}
	return pad.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		before := t.selectAll.Value
		dims := NewCheckbox(th, &t.selectAll, "").Layout(gtx)
		// Sync row checkboxes if select-all toggled this frame.
		if t.selectAll.Value != before {
			v := t.selectAll.Value
			for _, s := range t.selected {
				s.Value = v
			}
			t.prevSelectAll = v
		}
		return dims
	})
}

// ─── Rows ────────────────────────────────────────────────────────────────────

func (t *Table) renderRow(gtx layout.Context, idx int) layout.Dimensions {
	th := t.th
	if idx >= len(t.sortOrder) {
		return layout.Dimensions{}
	}
	origIdx := t.sortOrder[idx]
	if origIdx >= len(t.Rows) {
		return layout.Dimensions{}
	}
	row := t.Rows[origIdx]

	hovered := idx < len(t.rowBtns) && t.rowBtns[idx].Hovered()
	isSelected := origIdx < len(t.selected) && t.selected[origIdx].Value

	bg := th.Base100
	switch {
	case isSelected:
		bg = theme.Lerp(th.Base100, th.Primary, 0.12)
	case hovered:
		bg = th.Base200
	case t.Zebra && idx%2 == 1:
		bg = th.Base200
	}

	cells := make([]layout.FlexChild, 0, len(t.Cols)+1)
	if t.Selectable && origIdx < len(t.selected) {
		s := t.selected[origIdx]
		cells = append(cells, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			pad := layout.Inset{Top: t.padCellV(), Bottom: t.padCellV(), Left: th.Space4, Right: th.Space2}
			return pad.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return NewCheckbox(th, s, "").Layout(gtx)
			})
		}))
	}
	for j, col := range t.Cols {
		cells = append(cells, layout.Flexed(col.Flex, func(gtx layout.Context) layout.Dimensions {
			cell := ""
			if j < len(row) {
				cell = row[j]
			}
			return t.padCell().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, cell, th.BaseContent, th.SmSize, font.Normal)
			})
		}))
	}

	rowContent := func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return layout.Flex{}.Layout(gtx, cells...)
	}
	if idx < len(t.rowBtns) {
		btn := t.rowBtns[idx]
		rowContent = func(gtx layout.Context) layout.Dimensions {
			return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return layout.Flex{}.Layout(gtx, cells...)
			})
		}
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: bg}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(rowContent),
	)
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func tableRowDivider(gtx layout.Context, th *theme.Theme) layout.Dimensions {
	h := gtx.Dp(1)
	sz := image.Pt(gtx.Constraints.Max.X, h)
	defer clip.Rect{Max: sz}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: sz}
}
