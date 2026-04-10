package component

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/hongshengjie/gioui-kit/theme"
)

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
			borderCol := inp.borderColor()

			field := func(gtx layout.Context) layout.Dimensions {
				sz := image.Pt(gtx.Constraints.Max.X, h)
				gtx.Constraints = layout.Exact(sz)

				return layout.Stack{}.Layout(gtx,
					layout.Expanded(func(gtx layout.Context) layout.Dimensions {
						sz := gtx.Constraints.Min
						defer clip.UniformRRect(image.Rectangle{Max: sz}, gtx.Dp(th.RoundedLg)).Push(gtx.Ops).Pop()
						paint.ColorOp{Color: th.Base100}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
						return layout.Dimensions{Size: sz}
					}),
					layout.Stacked(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						gtx.Constraints.Min.Y = h
						return layout.Inset{
							Left: th.Space3, Right: th.Space3,
							Top: th.Space2, Bottom: th.Space2,
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							if th.Shaper == nil {
								th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(defaultFonts()))
							}
							textMat := op.Record(gtx.Ops)
							paint.ColorOp{Color: th.BaseContent}.Add(gtx.Ops)
							textCallOp := textMat.Stop()
							selMat := op.Record(gtx.Ops)
							paint.ColorOp{Color: theme.WithAlpha(th.Primary, 80)}.Add(gtx.Ops)
							selCallOp := selMat.Stop()
							if inp.Editor.Text() == "" && inp.Placeholder != "" {
								drawText(gtx, th, inp.Placeholder,
									theme.Opacity(th.BaseContent, 0.4), th.FontSize, font.Normal)
							}
							return inp.Editor.Layout(gtx, th.Shaper, font.Font{}, th.FontSize, textCallOp, selCallOp)
						})
					}),
				)
			}

			if inp.Variant != InputGhost {
				return widget.Border{Color: borderCol, CornerRadius: th.RoundedLg, Width: 1}.Layout(gtx, field)
			}
			return field(gtx)
		}),
	)
}
