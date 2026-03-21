package component

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

// IconButton is a DaisyUI-style icon-only button.
type IconButton struct {
	Icon      *widget.Icon
	Variant   BtnVariant
	Size      BtnSize
	Disabled  bool
	Clickable *widget.Clickable
	th        *theme.Theme
}

// NewIconButton creates a new icon-only button.
func NewIconButton(th *theme.Theme, click *widget.Clickable, icon *widget.Icon) *IconButton {
	return &IconButton{
		Icon:      icon,
		Variant:   BtnDefault,
		Size:      BtnMd,
		Clickable: click,
		th:        th,
	}
}

// WithVariant sets the button variant.
func (b *IconButton) WithVariant(v BtnVariant) *IconButton {
	b.Variant = v
	return b
}

// WithSize sets the button size.
func (b *IconButton) WithSize(s BtnSize) *IconButton {
	b.Size = s
	return b
}

func (b *IconButton) iconColors() (bg, fg color.NRGBA) {
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
	case BtnOutline:
		return theme.Transparent, th.BaseContent
	default:
		return th.Neutral, th.NeutralContent
	}
}

func (b *IconButton) btnSize() (size unit.Dp, iconSize unit.Dp) {
	switch b.Size {
	case BtnXs:
		return 24, 14
	case BtnSm:
		return 32, 16
	case BtnLg:
		return 52, 28
	default: // BtnMd
		return 40, 20
	}
}

// Layout renders the icon button.
func (b *IconButton) Layout(gtx layout.Context) layout.Dimensions {
	th := b.th
	bg, fg := b.iconColors()
	btnDp, iconDp := b.btnSize()
	btnPx := gtx.Dp(btnDp)
	radius := btnPx / 2

	if b.Clickable.Hovered() && !b.Disabled {
		if bg.A > 0 {
			bg = theme.Lerp(bg, theme.Black, 0.1)
		} else {
			bg = theme.WithAlpha(th.BaseContent, 15)
		}
	}
	if b.Clickable.Pressed() && !b.Disabled {
		if bg.A > 0 {
			bg = theme.Lerp(bg, theme.Black, 0.2)
		} else {
			bg = theme.WithAlpha(th.BaseContent, 25)
		}
	}
	if b.Disabled {
		bg = theme.Opacity(bg, 0.5)
		fg = theme.Opacity(fg, 0.5)
	}

	sz := image.Pt(btnPx, btnPx)
	gtx.Constraints = layout.Exact(sz)

	return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				rect := image.Rectangle{Max: sz}
				if bg.A > 0 {
					defer clip.UniformRRect(rect, radius).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: bg}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
				}
				if b.Variant == BtnOutline {
					paint.FillShape(gtx.Ops, th.BaseContent,
						clip.Stroke{
							Path:  clip.UniformRRect(rect, radius).Path(gtx.Ops),
							Width: float32(gtx.Dp(1)),
						}.Op(),
					)
				}
				pointer.CursorPointer.Add(gtx.Ops)
				return layout.Dimensions{Size: sz}
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				iconPx := gtx.Dp(iconDp)
				gtx.Constraints = layout.Exact(image.Pt(iconPx, iconPx))
				if b.Icon != nil {
					return b.Icon.Layout(gtx, fg)
				}
				return layout.Dimensions{Size: image.Pt(iconPx, iconPx)}
			}),
		)
	})
}
