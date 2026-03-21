// Package theme provides a TailwindCSS/DaisyUI-inspired theming system for Gio UI.
//
// Color naming follows TailwindCSS conventions (e.g., Slate50, Blue500)
// with DaisyUI semantic tokens (Primary, Secondary, Accent, etc.).
package theme

import (
	"image/color"

	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"
)

// ---------- Tailwind Color Palette ----------

// Gray scales (Tailwind-style naming)
var (
	Slate50  = rgb(0xf8fafc)
	Slate100 = rgb(0xf1f5f9)
	Slate200 = rgb(0xe2e8f0)
	Slate300 = rgb(0xcbd5e1)
	Slate400 = rgb(0x94a3b8)
	Slate500 = rgb(0x64748b)
	Slate600 = rgb(0x475569)
	Slate700 = rgb(0x334155)
	Slate800 = rgb(0x1e293b)
	Slate900 = rgb(0x0f172a)
	Slate950 = rgb(0x020617)

	Gray50  = rgb(0xf9fafb)
	Gray100 = rgb(0xf3f4f6)
	Gray200 = rgb(0xe5e7eb)
	Gray300 = rgb(0xd1d5db)
	Gray400 = rgb(0x9ca3af)
	Gray500 = rgb(0x6b7280)
	Gray600 = rgb(0x4b5563)
	Gray700 = rgb(0x374151)
	Gray800 = rgb(0x1f2937)
	Gray900 = rgb(0x111827)
	Gray950 = rgb(0x030712)

	Zinc50  = rgb(0xfafafa)
	Zinc100 = rgb(0xf4f4f5)
	Zinc200 = rgb(0xe4e4e7)
	Zinc300 = rgb(0xd4d4d8)
	Zinc400 = rgb(0xa1a1aa)
	Zinc500 = rgb(0x71717a)
	Zinc600 = rgb(0x52525b)
	Zinc700 = rgb(0x3f3f46)
	Zinc800 = rgb(0x27272a)
	Zinc900 = rgb(0x18181b)
	Zinc950 = rgb(0x09090b)
)

// Primary colors
var (
	Red50  = rgb(0xfef2f2)
	Red100 = rgb(0xfee2e2)
	Red200 = rgb(0xfecaca)
	Red300 = rgb(0xfca5a5)
	Red400 = rgb(0xf87171)
	Red500 = rgb(0xef4444)
	Red600 = rgb(0xdc2626)
	Red700 = rgb(0xb91c1c)
	Red800 = rgb(0x991b1b)
	Red900 = rgb(0x7f1d1d)

	Blue50  = rgb(0xeff6ff)
	Blue100 = rgb(0xdbeafe)
	Blue200 = rgb(0xbfdbfe)
	Blue300 = rgb(0x93c5fd)
	Blue400 = rgb(0x60a5fa)
	Blue500 = rgb(0x3b82f6)
	Blue600 = rgb(0x2563eb)
	Blue700 = rgb(0x1d4ed8)
	Blue800 = rgb(0x1e40af)
	Blue900 = rgb(0x1e3a8a)

	Green50  = rgb(0xf0fdf4)
	Green100 = rgb(0xdcfce7)
	Green200 = rgb(0xbbf7d0)
	Green300 = rgb(0x86efac)
	Green400 = rgb(0x4ade80)
	Green500 = rgb(0x22c55e)
	Green600 = rgb(0x16a34a)
	Green700 = rgb(0x15803d)
	Green800 = rgb(0x166534)
	Green900 = rgb(0x14532d)

	Yellow50  = rgb(0xfefce8)
	Yellow100 = rgb(0xfef9c3)
	Yellow200 = rgb(0xfef08a)
	Yellow300 = rgb(0xfde047)
	Yellow400 = rgb(0xfacc15)
	Yellow500 = rgb(0xeab308)
	Yellow600 = rgb(0xca8a04)
	Yellow700 = rgb(0xa16207)
	Yellow800 = rgb(0x854d0e)
	Yellow900 = rgb(0x713f12)

	Purple50  = rgb(0xfaf5ff)
	Purple100 = rgb(0xf3e8ff)
	Purple200 = rgb(0xe9d5ff)
	Purple300 = rgb(0xd8b4fe)
	Purple400 = rgb(0xc084fc)
	Purple500 = rgb(0xa855f7)
	Purple600 = rgb(0x9333ea)
	Purple700 = rgb(0x7e22ce)
	Purple800 = rgb(0x6b21a8)
	Purple900 = rgb(0x581c87)

	Indigo50  = rgb(0xeef2ff)
	Indigo100 = rgb(0xe0e7ff)
	Indigo200 = rgb(0xc7d2fe)
	Indigo300 = rgb(0xa5b4fc)
	Indigo400 = rgb(0x818cf8)
	Indigo500 = rgb(0x6366f1)
	Indigo600 = rgb(0x4f46e5)
	Indigo700 = rgb(0x4338ca)
	Indigo800 = rgb(0x3730a3)
	Indigo900 = rgb(0x312e81)

	Amber50  = rgb(0xfffbeb)
	Amber100 = rgb(0xfef3c7)
	Amber200 = rgb(0xfde68a)
	Amber300 = rgb(0xfcd34d)
	Amber400 = rgb(0xfbbf24)
	Amber500 = rgb(0xf59e0b)
	Amber600 = rgb(0xd97706)
	Amber700 = rgb(0xb45309)
	Amber800 = rgb(0x92400e)
	Amber900 = rgb(0x78350f)

	Cyan50  = rgb(0xecfeff)
	Cyan100 = rgb(0xcffafe)
	Cyan200 = rgb(0xa5f3fc)
	Cyan300 = rgb(0x67e8f9)
	Cyan400 = rgb(0x22d3ee)
	Cyan500 = rgb(0x06b6d4)
	Cyan600 = rgb(0x0891b2)
	Cyan700 = rgb(0x0e7490)
	Cyan800 = rgb(0x155e75)
	Cyan900 = rgb(0x164e63)

	Emerald50  = rgb(0xecfdf5)
	Emerald100 = rgb(0xd1fae5)
	Emerald200 = rgb(0xa7f3d0)
	Emerald300 = rgb(0x6ee7b7)
	Emerald400 = rgb(0x34d399)
	Emerald500 = rgb(0x10b981)
	Emerald600 = rgb(0x059669)
	Emerald700 = rgb(0x047857)
	Emerald800 = rgb(0x065f46)
	Emerald900 = rgb(0x064e3b)

	Rose50  = rgb(0xfff1f2)
	Rose100 = rgb(0xffe4e6)
	Rose200 = rgb(0xfecdd3)
	Rose300 = rgb(0xfda4af)
	Rose400 = rgb(0xfb7185)
	Rose500 = rgb(0xf43f5e)
	Rose600 = rgb(0xe11d48)
	Rose700 = rgb(0xbe123c)
	Rose800 = rgb(0x9f1239)
	Rose900 = rgb(0x881337)
)

// Special colors
var (
	White       = rgb(0xffffff)
	Black       = rgb(0x000000)
	Transparent = color.NRGBA{A: 0}
)

// ---------- DaisyUI Semantic Theme ----------

// Theme holds DaisyUI-style semantic color tokens and typography settings.
type Theme struct {
	// DaisyUI semantic colors
	Primary        color.NRGBA
	PrimaryContent color.NRGBA
	Secondary        color.NRGBA
	SecondaryContent color.NRGBA
	Accent        color.NRGBA
	AccentContent color.NRGBA
	Neutral        color.NRGBA
	NeutralContent color.NRGBA

	// State colors
	Info        color.NRGBA
	InfoContent color.NRGBA
	Success        color.NRGBA
	SuccessContent color.NRGBA
	Warning        color.NRGBA
	WarningContent color.NRGBA
	Error        color.NRGBA
	ErrorContent color.NRGBA

	// Surface colors
	Base100 color.NRGBA // bg base
	Base200 color.NRGBA // bg slightly darker
	Base300 color.NRGBA // bg more darker
	BaseContent color.NRGBA // text on base

	// Typography
	FontSize    unit.Sp
	H1Size      unit.Sp
	H2Size      unit.Sp
	H3Size      unit.Sp
	H4Size      unit.Sp
	SmSize      unit.Sp
	XsSize      unit.Sp

	// Font collection for Gio shaper
	Shaper *text.Shaper

	// Spacing scale (Tailwind-like)
	Space0  unit.Dp // 0
	Space1  unit.Dp // 4
	Space2  unit.Dp // 8
	Space3  unit.Dp // 12
	Space4  unit.Dp // 16
	Space5  unit.Dp // 20
	Space6  unit.Dp // 24
	Space8  unit.Dp // 32
	Space10 unit.Dp // 40
	Space12 unit.Dp // 48
	Space16 unit.Dp // 64

	// Border radius
	RoundedNone unit.Dp
	RoundedSm   unit.Dp
	RoundedMd   unit.Dp
	RoundedLg   unit.Dp
	RoundedXl   unit.Dp
	Rounded2xl  unit.Dp
	RoundedFull unit.Dp
}

// ---------- Preset Themes (DaisyUI-inspired) ----------

// Light returns the default light theme (similar to DaisyUI "light").
func Light() *Theme {
	return &Theme{
		Primary:          Indigo500,
		PrimaryContent:   White,
		Secondary:        Purple500,
		SecondaryContent: White,
		Accent:           Cyan500,
		AccentContent:    White,
		Neutral:          Gray700,
		NeutralContent:   Gray100,

		Info:           Blue500,
		InfoContent:    White,
		Success:        Emerald500,
		SuccessContent: White,
		Warning:        Amber500,
		WarningContent: White,
		Error:          Red500,
		ErrorContent:   White,

		Base100:     White,
		Base200:     Gray50,
		Base300:     Gray100,
		BaseContent: Gray900,

		FontSize: 16,
		H1Size:   30,
		H2Size:   24,
		H3Size:   20,
		H4Size:   18,
		SmSize:   14,
		XsSize:   12,

		Space0:  0,
		Space1:  4,
		Space2:  8,
		Space3:  12,
		Space4:  16,
		Space5:  20,
		Space6:  24,
		Space8:  32,
		Space10: 40,
		Space12: 48,
		Space16: 64,

		RoundedNone: 0,
		RoundedSm:   2,
		RoundedMd:   6,
		RoundedLg:   8,
		RoundedXl:   12,
		Rounded2xl:  16,
		RoundedFull: 9999,
	}
}

// Dark returns a dark theme (similar to DaisyUI "dark").
func Dark() *Theme {
	return &Theme{
		Primary:          Indigo400,
		PrimaryContent:   White,
		Secondary:        Purple400,
		SecondaryContent: White,
		Accent:           Cyan400,
		AccentContent:    Slate900,
		Neutral:          Slate600,
		NeutralContent:   Slate200,

		Info:           Blue400,
		InfoContent:    Slate900,
		Success:        Emerald400,
		SuccessContent: Slate900,
		Warning:        Amber400,
		WarningContent: Slate900,
		Error:          Rose400,
		ErrorContent:   White,

		Base100:     Slate900,
		Base200:     Slate800,
		Base300:     Slate700,
		BaseContent: Slate100,

		FontSize: 16,
		H1Size:   30,
		H2Size:   24,
		H3Size:   20,
		H4Size:   18,
		SmSize:   14,
		XsSize:   12,

		Space0:  0,
		Space1:  4,
		Space2:  8,
		Space3:  12,
		Space4:  16,
		Space5:  20,
		Space6:  24,
		Space8:  32,
		Space10: 40,
		Space12: 48,
		Space16: 64,

		RoundedNone: 0,
		RoundedSm:   2,
		RoundedMd:   6,
		RoundedLg:   8,
		RoundedXl:   12,
		Rounded2xl:  16,
		RoundedFull: 9999,
	}
}

// Cupcake returns a pastel-friendly theme (DaisyUI "cupcake").
func Cupcake() *Theme {
	t := Light()
	t.Primary = rgb(0x65c3c8)
	t.PrimaryContent = rgb(0x223D40)
	t.Secondary = rgb(0xef9fbc)
	t.SecondaryContent = rgb(0x49242e)
	t.Accent = rgb(0xeeaf3a)
	t.AccentContent = rgb(0x452f10)
	t.Neutral = rgb(0x291334)
	t.NeutralContent = rgb(0xe8d8f0)
	t.Base100 = rgb(0xfaf7f5)
	t.Base200 = rgb(0xefeae6)
	t.Base300 = rgb(0xe7e2df)
	t.BaseContent = rgb(0x291334)
	return t
}

// Nord returns a Nord-palette theme.
func Nord() *Theme {
	t := Light()
	t.Primary = rgb(0x5E81AC)
	t.PrimaryContent = rgb(0xECEFF4)
	t.Secondary = rgb(0x81A1C1)
	t.SecondaryContent = rgb(0x2E3440)
	t.Accent = rgb(0x88C0D0)
	t.AccentContent = rgb(0x2E3440)
	t.Neutral = rgb(0x4C566A)
	t.NeutralContent = rgb(0xD8DEE9)
	t.Info = rgb(0x88C0D0)
	t.Success = rgb(0xA3BE8C)
	t.Warning = rgb(0xEBCB8B)
	t.Error = rgb(0xBF616A)
	t.Base100 = rgb(0xECEFF4)
	t.Base200 = rgb(0xE5E9F0)
	t.Base300 = rgb(0xD8DEE9)
	t.BaseContent = rgb(0x2E3440)
	return t
}

// ---------- Font helpers ----------

// FontFace returns a font.Font with the given weight.
func FontFace(weight font.Weight) font.Font {
	return font.Font{Weight: weight}
}

var (
	FontThin       = FontFace(font.Thin)
	FontLight      = FontFace(font.Light)
	FontNormal     = FontFace(font.Normal)
	FontMedium     = FontFace(font.Medium)
	FontSemiBold   = FontFace(font.SemiBold)
	FontBold       = FontFace(font.Bold)
	FontExtraBold  = FontFace(font.ExtraBold)
)

// ---------- Color utilities ----------

func rgb(hex uint32) color.NRGBA {
	return color.NRGBA{
		R: uint8(hex >> 16),
		G: uint8(hex >> 8),
		B: uint8(hex),
		A: 0xff,
	}
}

// WithAlpha returns a copy of c with the given alpha (0–255).
func WithAlpha(c color.NRGBA, a uint8) color.NRGBA {
	c.A = a
	return c
}

// Opacity returns a color with the given opacity fraction (0.0–1.0).
func Opacity(c color.NRGBA, opacity float32) color.NRGBA {
	c.A = uint8(float32(c.A) * opacity)
	return c
}

// Lerp linearly interpolates between two colors.
func Lerp(a, b color.NRGBA, t float32) color.NRGBA {
	return color.NRGBA{
		R: uint8(float32(a.R)*(1-t) + float32(b.R)*t),
		G: uint8(float32(a.G)*(1-t) + float32(b.G)*t),
		B: uint8(float32(a.B)*(1-t) + float32(b.B)*t),
		A: uint8(float32(a.A)*(1-t) + float32(b.A)*t),
	}
}

// RGB creates a color from hex value.
func RGB(hex uint32) color.NRGBA {
	return rgb(hex)
}
