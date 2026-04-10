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
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

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
	radius := th.RoundedLg

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

	inner := func(gtx layout.Context) layout.Dimensions {
		return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{Alignment: layout.Center}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					sz := gtx.Constraints.Min
					defer clip.UniformRRect(image.Rectangle{Max: sz}, gtx.Dp(radius)).Push(gtx.Ops).Pop()
					if bg.A > 0 {
						paint.ColorOp{Color: bg}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
					}
					pointer.CursorPointer.Add(gtx.Ops)
					return layout.Dimensions{Size: sz}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return b.padding().Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						if b.Loading {
							return drawSpinner(gtx, fg, b.fontSize())
						}
						return drawText(gtx, th, b.Text, fg, b.fontSize(), font.SemiBold)
					})
				}),
			)
		})
	}

	if b.Variant == BtnOutline {
		return widget.Border{Color: th.BaseContent, CornerRadius: radius, Width: 1}.Layout(gtx, inner)
	}
	return inner(gtx)
}
