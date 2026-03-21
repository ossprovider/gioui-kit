// Package component provides DaisyUI-inspired UI components for Gio.
//
// Components follow DaisyUI naming and variant patterns:
//
//	btn := component.Button{Text: "Click me", Variant: component.BtnPrimary}
//	card := component.Card{Title: "Hello"}
//	badge := component.Badge{Text: "NEW", Variant: component.BadgeAccent}
package component

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
// Button (DaisyUI btn / btn-primary / btn-sm / btn-outline)
// ============================================================

// BtnVariant defines button style variants.
type BtnVariant int

const (
	BtnDefault BtnVariant = iota
	BtnPrimary
	BtnSecondary
	BtnAccent
	BtnInfo
	BtnSuccess
	BtnWarning
	BtnError
	BtnGhost
	BtnLink
	BtnOutline
)

// BtnSize defines button sizes.
type BtnSize int

const (
	BtnMd BtnSize = iota
	BtnXs
	BtnSm
	BtnLg
)

// Button is a DaisyUI-style button component.
type Button struct {
	Text      string
	Variant   BtnVariant
	Size      BtnSize
	Disabled  bool
	Loading   bool
	FullWidth bool

	Clickable *widget.Clickable

	th *theme.Theme
}

// NewButton creates a new button with a theme.
func NewButton(th *theme.Theme, click *widget.Clickable, text string) *Button {
	return &Button{
		Text:      text,
		Variant:   BtnDefault,
		Size:      BtnMd,
		Clickable: click,
		th:        th,
	}
}

// WithVariant sets the button variant.
func (b *Button) WithVariant(v BtnVariant) *Button {
	b.Variant = v
	return b
}

// WithSize sets the button size.
func (b *Button) WithSize(s BtnSize) *Button {
	b.Size = s
	return b
}

func (b *Button) bgColor() (bg, fg color.NRGBA) {
	th := b.th
	switch b.Variant {
	case BtnPrimary:
		return th.Primary, th.PrimaryContent
	case BtnSecondary:
		return th.Secondary, th.SecondaryContent
	case BtnAccent:
		return th.Accent, th.AccentContent
	case BtnInfo:
		return th.Info, th.InfoContent
	case BtnSuccess:
		return th.Success, th.SuccessContent
	case BtnWarning:
		return th.Warning, th.WarningContent
	case BtnError:
		return th.Error, th.ErrorContent
	case BtnGhost:
		return theme.Transparent, th.BaseContent
	case BtnLink:
		return theme.Transparent, th.Primary
	case BtnOutline:
		return theme.Transparent, th.BaseContent
	default:
		return th.Neutral, th.NeutralContent
	}
}

func (b *Button) padding() layout.Inset {
	switch b.Size {
	case BtnXs:
		return layout.Inset{Top: 2, Bottom: 2, Left: 8, Right: 8}
	case BtnSm:
		return layout.Inset{Top: 4, Bottom: 4, Left: 12, Right: 12}
	case BtnLg:
		return layout.Inset{Top: 12, Bottom: 12, Left: 24, Right: 24}
	default: // BtnMd
		return layout.Inset{Top: 8, Bottom: 8, Left: 16, Right: 16}
	}
}

func (b *Button) fontSize() unit.Sp {
	switch b.Size {
	case BtnXs:
		return b.th.XsSize
	case BtnSm:
		return b.th.SmSize
	case BtnLg:
		return b.th.H4Size
	default:
		return b.th.SmSize
	}
}

// Layout renders the button.
func (b *Button) Layout(gtx layout.Context) layout.Dimensions {
	th := b.th
	bg, fg := b.bgColor()
	radius := gtx.Dp(th.RoundedLg)

	// Hover state
	if b.Clickable.Hovered() && !b.Disabled {
		bg = theme.Lerp(bg, theme.Black, 0.1)
	}
	// Pressed state
	if b.Clickable.Pressed() && !b.Disabled {
		bg = theme.Lerp(bg, theme.Black, 0.2)
	}
	// Disabled state
	if b.Disabled {
		bg = theme.Opacity(bg, 0.5)
		fg = theme.Opacity(fg, 0.5)
	}

	return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return b.padding().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					sz := gtx.Constraints.Min
					rrect := clip.UniformRRect(image.Rectangle{Max: sz}, radius)
					defer rrect.Push(gtx.Ops).Pop()

					// Background fill
					if bg.A > 0 {
						paint.ColorOp{Color: bg}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
					}

					// Outline variant border
					if b.Variant == BtnOutline {
						paint.FillShape(gtx.Ops, th.BaseContent,
							clip.Stroke{
								Path:  clip.UniformRRect(image.Rectangle{Max: sz}, radius).Path(gtx.Ops),
								Width: float32(gtx.Dp(1)),
							}.Op(),
						)
					}

					// Cursor
					pointer.CursorPointer.Add(gtx.Ops)
					return layout.Dimensions{Size: sz}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					if b.Loading {
						return drawSpinner(gtx, fg, b.fontSize())
					}
					return drawText(gtx, th, b.Text, fg, b.fontSize(), font.SemiBold)
				}),
			)
		})
	})
}

// ============================================================
// Badge (DaisyUI badge / badge-primary)
// ============================================================

type BadgeVariant int

const (
	BadgeDefault BadgeVariant = iota
	BadgePrimary
	BadgeSecondary
	BadgeAccent
	BadgeInfo
	BadgeSuccess
	BadgeWarning
	BadgeError
	BadgeOutline
	BadgeGhost
)

// Badge is a DaisyUI-style badge/tag component.
type Badge struct {
	Text    string
	Variant BadgeVariant
	th      *theme.Theme
}

func NewBadge(th *theme.Theme, text string) *Badge {
	return &Badge{Text: text, th: th}
}

func (b *Badge) WithVariant(v BadgeVariant) *Badge {
	b.Variant = v
	return b
}

func (b *Badge) colors() (bg, fg color.NRGBA) {
	th := b.th
	switch b.Variant {
	case BadgePrimary:
		return th.Primary, th.PrimaryContent
	case BadgeSecondary:
		return th.Secondary, th.SecondaryContent
	case BadgeAccent:
		return th.Accent, th.AccentContent
	case BadgeInfo:
		return th.Info, th.InfoContent
	case BadgeSuccess:
		return th.Success, th.SuccessContent
	case BadgeWarning:
		return th.Warning, th.WarningContent
	case BadgeError:
		return th.Error, th.ErrorContent
	case BadgeGhost:
		return th.Base200, th.BaseContent
	default:
		return th.Neutral, th.NeutralContent
	}
}

func (b *Badge) Layout(gtx layout.Context) layout.Dimensions {
	bg, fg := b.colors()
	padding := layout.Inset{Top: 2, Bottom: 2, Left: 8, Right: 8}
	radius := gtx.Dp(b.th.RoundedFull)

	return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				sz := gtx.Constraints.Min
				rrect := clip.UniformRRect(image.Rectangle{Max: sz}, radius)
				defer rrect.Push(gtx.Ops).Pop()
				paint.ColorOp{Color: bg}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, b.th, b.Text, fg, b.th.XsSize, font.SemiBold)
			}),
		)
	})
}

// ============================================================
// Card (DaisyUI card / card-bordered)
// ============================================================

// Card is a DaisyUI-style card container.
type Card struct {
	Bordered bool
	Compact  bool
	th       *theme.Theme
}

func NewCard(th *theme.Theme) *Card {
	return &Card{th: th}
}

func (c *Card) WithBorder() *Card {
	c.Bordered = true
	return c
}

func (c *Card) WithCompact() *Card {
	c.Compact = true
	return c
}

func (c *Card) Layout(gtx layout.Context, body layout.Widget) layout.Dimensions {
	th := c.th
	radius := gtx.Dp(th.RoundedXl)
	padding := layout.UniformInset(th.Space6)
	if c.Compact {
		padding = layout.UniformInset(th.Space4)
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			rrect := clip.UniformRRect(image.Rectangle{Max: sz}, radius)
			defer rrect.Push(gtx.Ops).Pop()

			// Background
			paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Border
			if c.Bordered {
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
			return padding.Layout(gtx, body)
		}),
	)
}

// CardWithHeader renders a card with a separate title section.
func (c *Card) CardWithHeader(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
	th := c.th
	return c.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: th.Space3}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawText(gtx, th, title, th.BaseContent, th.H3Size, font.Bold)
				})
			}),
			layout.Rigid(body),
		)
	})
}

// ============================================================
// Input (DaisyUI input / input-bordered)
// ============================================================

type InputVariant int

const (
	InputDefault InputVariant = iota
	InputBordered
	InputGhost
	InputPrimary
	InputSecondary
	InputAccent
	InputInfo
	InputSuccess
	InputWarning
	InputError
)

type InputSize int

const (
	InputMd InputSize = iota
	InputXs
	InputSm
	InputLg
)

// Input is a DaisyUI-style text input.
type Input struct {
	Editor      *widget.Editor
	Placeholder string
	Label       string
	Variant     InputVariant
	Size        InputSize
	th          *theme.Theme
}

func NewInput(th *theme.Theme, editor *widget.Editor, placeholder string) *Input {
	return &Input{
		Editor:      editor,
		Placeholder: placeholder,
		Variant:     InputBordered,
		Size:        InputMd,
		th:          th,
	}
}

func (inp *Input) WithLabel(label string) *Input {
	inp.Label = label
	return inp
}

func (inp *Input) WithVariant(v InputVariant) *Input {
	inp.Variant = v
	return inp
}

func (inp *Input) borderColor() color.NRGBA {
	th := inp.th
	switch inp.Variant {
	case InputPrimary:
		return th.Primary
	case InputSecondary:
		return th.Secondary
	case InputAccent:
		return th.Accent
	case InputInfo:
		return th.Info
	case InputSuccess:
		return th.Success
	case InputWarning:
		return th.Warning
	case InputError:
		return th.Error
	case InputGhost:
		return theme.Transparent
	default:
		return th.Base300
	}
}

func (inp *Input) height() unit.Dp {
	switch inp.Size {
	case InputXs:
		return 24
	case InputSm:
		return 32
	case InputLg:
		return 48
	default:
		return 40
	}
}

func (inp *Input) Layout(gtx layout.Context) layout.Dimensions {
	th := inp.th

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Label
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if inp.Label == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Bottom: th.Space1}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, inp.Label, th.BaseContent, th.SmSize, font.Medium)
			})
		}),
		// Input field
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			h := gtx.Dp(inp.height())
			radius := gtx.Dp(th.RoundedLg)
			borderCol := inp.borderColor()

			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					sz := image.Pt(gtx.Constraints.Max.X, h)
					rrect := clip.UniformRRect(image.Rectangle{Max: sz}, radius)
					defer rrect.Push(gtx.Ops).Pop()

					// Background
					paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)

					// Border
					if inp.Variant != InputGhost {
						paint.FillShape(gtx.Ops, borderCol,
							clip.Stroke{
								Path:  clip.UniformRRect(image.Rectangle{Max: sz}, radius).Path(gtx.Ops),
								Width: float32(gtx.Dp(1)),
							}.Op(),
						)
					}
					return layout.Dimensions{Size: sz}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.Y = h
					return layout.Inset{
						Left: th.Space3, Right: th.Space3,
						Top: th.Space2, Bottom: th.Space2,
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if th.Shaper == nil {
							th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(defaultFonts()))
						}
						return inp.Editor.Layout(gtx, th.Shaper, font.Font{}, th.FontSize, op.CallOp{}, op.CallOp{})
					})
				}),
			)
		}),
	)
}

// ============================================================
// Toggle / Checkbox (DaisyUI toggle)
// ============================================================

// Toggle is a DaisyUI-style toggle switch.
type Toggle struct {
	Bool    *widget.Bool
	Label   string
	Variant BtnVariant
	th      *theme.Theme
}

func NewToggle(th *theme.Theme, b *widget.Bool, label string) *Toggle {
	return &Toggle{
		Bool:    b,
		Label:   label,
		Variant: BtnPrimary,
		th:      th,
	}
}

func (t *Toggle) Layout(gtx layout.Context) layout.Dimensions {
	th := t.th

	return layout.Flex{
		Alignment: layout.Middle,
	}.Layout(gtx,
		// Toggle track
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			trackW := gtx.Dp(44)
			trackH := gtx.Dp(24)
			thumbSize := gtx.Dp(20)
			radius := trackH / 2

			return t.Bool.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Track
				trackColor := th.Base300
				if t.Bool.Value {
					trackColor = th.Primary
				}
				trackRect := image.Rect(0, 0, trackW, trackH)
				defer clip.UniformRRect(trackRect, radius).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: trackColor}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)

				// Thumb
				thumbX := gtx.Dp(2)
				if t.Bool.Value {
					thumbX = trackW - thumbSize - gtx.Dp(2)
				}
				thumbY := (trackH - thumbSize) / 2
				thumbRect := image.Rect(thumbX, thumbY, thumbX+thumbSize, thumbY+thumbSize)
				defer clip.UniformRRect(thumbRect, thumbSize/2).Push(gtx.Ops).Pop()
				paint.ColorOp{Color: theme.White}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)

				pointer.CursorPointer.Add(gtx.Ops)
				return layout.Dimensions{Size: image.Pt(trackW, trackH)}
			})
		}),
		// Label
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if t.Label == "" {
				return layout.Dimensions{}
			}
			return layout.Inset{Left: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, t.Label, th.BaseContent, th.FontSize, font.Normal)
			})
		}),
	)
}

// ============================================================
// Alert (DaisyUI alert / alert-info / alert-success)
// ============================================================

type AlertVariant int

const (
	AlertInfo AlertVariant = iota
	AlertSuccess
	AlertWarning
	AlertError
)

// Alert is a DaisyUI-style alert banner.
type Alert struct {
	Text    string
	Variant AlertVariant
	th      *theme.Theme
}

func NewAlert(th *theme.Theme, text string, variant AlertVariant) *Alert {
	return &Alert{Text: text, Variant: variant, th: th}
}

func (a *Alert) colors() (bg, fg, icon color.NRGBA) {
	th := a.th
	switch a.Variant {
	case AlertSuccess:
		return theme.WithAlpha(th.Success, 30), th.Success, th.Success
	case AlertWarning:
		return theme.WithAlpha(th.Warning, 30), th.Warning, th.Warning
	case AlertError:
		return theme.WithAlpha(th.Error, 30), th.Error, th.Error
	default:
		return theme.WithAlpha(th.Info, 30), th.Info, th.Info
	}
}

func (a *Alert) icon() string {
	switch a.Variant {
	case AlertSuccess:
		return "✓"
	case AlertWarning:
		return "⚠"
	case AlertError:
		return "✕"
	default:
		return "ℹ"
	}
}

func (a *Alert) Layout(gtx layout.Context) layout.Dimensions {
	th := a.th
	bg, fg, _ := a.colors()
	radius := gtx.Dp(th.RoundedLg)
	padding := layout.Inset{
		Top: th.Space3, Bottom: th.Space3,
		Left: th.Space4, Right: th.Space4,
	}

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			rrect := clip.UniformRRect(image.Rectangle{Max: sz}, radius)
			defer rrect.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: bg}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			return padding.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Right: th.Space2}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return drawText(gtx, th, a.icon(), fg, th.H3Size, font.Bold)
						})
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return drawText(gtx, th, a.Text, fg, th.SmSize, font.Medium)
					}),
				)
			})
		}),
	)
}

// ============================================================
// Avatar (DaisyUI avatar)
// ============================================================

type AvatarSize int

const (
	AvatarMd AvatarSize = iota
	AvatarXs
	AvatarSm
	AvatarLg
)

// Avatar renders a circular avatar placeholder.
type Avatar struct {
	Initials string
	Size     AvatarSize
	Online   bool
	th       *theme.Theme
}

func NewAvatar(th *theme.Theme, initials string) *Avatar {
	return &Avatar{Initials: initials, Size: AvatarMd, th: th}
}

func (a *Avatar) sizeDp() unit.Dp {
	switch a.Size {
	case AvatarXs:
		return 24
	case AvatarSm:
		return 32
	case AvatarLg:
		return 64
	default:
		return 48
	}
}

func (a *Avatar) Layout(gtx layout.Context) layout.Dimensions {
	th := a.th
	sz := gtx.Dp(a.sizeDp())
	rect := image.Rect(0, 0, sz, sz)

	// Circle background
	defer clip.UniformRRect(rect, sz/2).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: th.Primary}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// Initials
	textSize := unit.Sp(float32(sz) * 0.4)
	macro := op.Record(gtx.Ops)
	dims := drawText(gtx, th, a.Initials, th.PrimaryContent, textSize, font.Bold)
	call := macro.Stop()

	offX := (sz - dims.Size.X) / 2
	offY := (sz - dims.Size.Y) / 2
	defer op.Offset(image.Pt(offX, offY)).Push(gtx.Ops).Pop()
	call.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Pt(sz, sz)}
}

// ============================================================
// Progress (DaisyUI progress)
// ============================================================

type ProgressVariant int

const (
	ProgressPrimary ProgressVariant = iota
	ProgressSecondary
	ProgressAccent
	ProgressInfo
	ProgressSuccess
	ProgressWarning
	ProgressError
)

// Progress renders a progress bar.
type Progress struct {
	Value   float32 // 0.0 to 1.0
	Variant ProgressVariant
	th      *theme.Theme
}

func NewProgress(th *theme.Theme, value float32) *Progress {
	return &Progress{Value: value, th: th}
}

func (p *Progress) color() color.NRGBA {
	th := p.th
	switch p.Variant {
	case ProgressSecondary:
		return th.Secondary
	case ProgressAccent:
		return th.Accent
	case ProgressInfo:
		return th.Info
	case ProgressSuccess:
		return th.Success
	case ProgressWarning:
		return th.Warning
	case ProgressError:
		return th.Error
	default:
		return th.Primary
	}
}

func (p *Progress) Layout(gtx layout.Context) layout.Dimensions {
	th := p.th
	h := gtx.Dp(8)
	w := gtx.Constraints.Max.X
	radius := h / 2

	// Track
	trackRect := image.Rect(0, 0, w, h)
	defer clip.UniformRRect(trackRect, radius).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: th.Base300}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	// Fill
	fillW := int(float32(w) * p.Value)
	if fillW > 0 {
		fillRect := image.Rect(0, 0, fillW, h)
		defer clip.UniformRRect(fillRect, radius).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: p.color()}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
	}

	return layout.Dimensions{Size: image.Pt(w, h)}
}

// ============================================================
// Tooltip
// ============================================================

// Tooltip wraps content with a tooltip text.
type Tooltip struct {
	Text string
	th   *theme.Theme
}

func NewTooltip(th *theme.Theme, text string) *Tooltip {
	return &Tooltip{Text: text, th: th}
}

// ============================================================
// Tabs (DaisyUI tabs)
// ============================================================

type TabVariant int

const (
	TabBoxed TabVariant = iota
	TabBordered
	TabLifted
)

// Tabs manages a tabbed interface.
type Tabs struct {
	Items    []string
	Selected int
	Variant  TabVariant
	clicks   []widget.Clickable
	children []layout.FlexChild // reused across frames to avoid per-frame alloc
	th       *theme.Theme
}

func NewTabs(th *theme.Theme, items []string) *Tabs {
	return &Tabs{
		Items:  items,
		clicks: make([]widget.Clickable, len(items)),
		th:     th,
	}
}

func (t *Tabs) Layout(gtx layout.Context) layout.Dimensions {
	th := t.th

	// Check clicks
	for i := range t.clicks {
		if t.clicks[i].Clicked(gtx) {
			t.Selected = i
		}
	}

	if cap(t.children) < len(t.Items) {
		t.children = make([]layout.FlexChild, len(t.Items))
	}
	t.children = t.children[:len(t.Items)]
	for i, item := range t.Items {
		i, item := i, item
		t.children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			isActive := i == t.Selected
			return t.clicks[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
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
							// Active indicator
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

// ============================================================
// Chip (tag-like removable badges)
// ============================================================

type Chip struct {
	Text     string
	Closable bool
	close    widget.Clickable
	th       *theme.Theme
}

func NewChip(th *theme.Theme, text string) *Chip {
	return &Chip{Text: text, th: th}
}

func (c *Chip) Layout(gtx layout.Context) layout.Dimensions {
	th := c.th
	radius := gtx.Dp(th.RoundedFull)

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			sz := gtx.Constraints.Min
			defer clip.UniformRRect(image.Rectangle{Max: sz}, radius).Push(gtx.Ops).Pop()
			paint.ColorOp{Color: th.Base200}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: sz}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			inset := layout.Inset{Top: 4, Bottom: 4, Left: 12, Right: 12}
			return inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawText(gtx, th, c.Text, th.BaseContent, th.SmSize, font.Medium)
			})
		}),
	)
}

// ============================================================
// Skeleton (loading placeholder)
// ============================================================

// Skeleton renders a loading skeleton placeholder.
type Skeleton struct {
	Width  unit.Dp
	Height unit.Dp
	Radius unit.Dp
	th     *theme.Theme
}

func NewSkeleton(th *theme.Theme) *Skeleton {
	return &Skeleton{Width: 200, Height: 20, Radius: th.RoundedMd, th: th}
}

func (s *Skeleton) Layout(gtx layout.Context) layout.Dimensions {
	w := gtx.Dp(s.Width)
	h := gtx.Dp(s.Height)
	rr := gtx.Dp(s.Radius)

	rect := image.Rect(0, 0, w, h)
	defer clip.UniformRRect(rect, rr).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: s.th.Base300}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: image.Pt(w, h)}
}

// ============================================================
// Text helpers
// ============================================================

// Text renders themed text (like Tailwind text-sm, text-lg, font-bold, etc.).
type Text struct {
	Content string
	Color   color.NRGBA
	Size    unit.Sp
	Weight  font.Weight
	th      *theme.Theme
}

func NewText(th *theme.Theme, content string) *Text {
	return &Text{
		Content: content,
		Color:   th.BaseContent,
		Size:    th.FontSize,
		Weight:  font.Normal,
		th:      th,
	}
}

func (t *Text) H1() *Text   { t.Size = t.th.H1Size; t.Weight = font.Bold; return t }
func (t *Text) H2() *Text   { t.Size = t.th.H2Size; t.Weight = font.Bold; return t }
func (t *Text) H3() *Text   { t.Size = t.th.H3Size; t.Weight = font.SemiBold; return t }
func (t *Text) H4() *Text   { t.Size = t.th.H4Size; t.Weight = font.SemiBold; return t }
func (t *Text) Sm() *Text   { t.Size = t.th.SmSize; return t }
func (t *Text) Xs() *Text   { t.Size = t.th.XsSize; return t }
func (t *Text) Bold() *Text { t.Weight = font.Bold; return t }
func (t *Text) WithColor(c color.NRGBA) *Text { t.Color = c; return t }

func (t *Text) Layout(gtx layout.Context) layout.Dimensions {
	return drawText(gtx, t.th, t.Content, t.Color, t.Size, t.Weight)
}

// drawText is a shared text drawing utility.
func drawText(gtx layout.Context, th *theme.Theme, txt string, col color.NRGBA, size unit.Sp, weight font.Weight) layout.Dimensions {
	lbl := widget.Label{MaxLines: 0}
	f := font.Font{Weight: weight}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	if th.Shaper == nil {
		th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(defaultFonts()))
	}
	return lbl.Layout(gtx, th.Shaper, f, size, txt, op.CallOp{})
}

// drawSpinner draws a simple loading spinner placeholder.
func drawSpinner(gtx layout.Context, col color.NRGBA, size unit.Sp) layout.Dimensions {
	sz := gtx.Sp(size)
	return layout.Dimensions{Size: image.Pt(int(sz), int(sz))}
}

// defaultFonts returns the built-in Go font collection.
func defaultFonts() []font.FontFace {
	return gofont.Collection()
}
