# GioUI Kit

> A **TailwindCSS + DaisyUI** inspired component library and scaffold for [Gio UI](https://gioui.org) — Go's immediate-mode cross-platform GUI framework.

**Naming follows Tailwind / DaisyUI conventions** so that if you know Tailwind's `flex flex-col gap-4 p-6 rounded-lg` or DaisyUI's `btn btn-primary`, you'll feel right at home in Go.

---

## Architecture

```
gioui-kit/
├── pkg/
│   ├── theme/        # 🎨 Color palette + semantic tokens + typography
│   │   └── theme.go  #    Tailwind colors (Slate50..950, Blue500, etc.)
│   │                  #    DaisyUI tokens (Primary, Secondary, Base100, etc.)
│   │                  #    Preset themes: Light(), Dark(), Cupcake(), Nord()
│   │
│   ├── layout/       # 📐 Flexbox, Grid, Container, Spacing utilities
│   │   └── layout.go #    FlexRow, FlexCol with Gap support
│   │                  #    Grid{Cols: 3, Gap: 16}
│   │                  #    P(), Px(), Py(), W(), H(), MinW(), MaxW()
│   │                  #    Container, Box, DividerH, ScrollY/ScrollX
│   │
│   ├── modifier/     # ✨ Visual decorators (shadow, gradient, ring, bg)
│   │   └── modifier.go#   Shadow (ShadowSm..Shadow2xl)
│   │                  #    LinearGradient, Bg, Rounded, Ring, Opacity
│   │
│   ├── component/    # 🧩 DaisyUI components
│   │   └── component.go#  Button (variants + sizes)
│   │                  #    Badge, Card, Input, Toggle, Alert
│   │                  #    Avatar, Progress, Tabs, Chip, Skeleton, Text
│   │
│   └── scaffold/     # 🏗️ App-level layout scaffolds
│       └── scaffold.go#   AppShell (Navbar + Sidebar + Content)
│                      #    Navbar, Sidebar, Modal, Drawer
│                      #    BottomNav, Toast, Breadcrumb
│
├── cmd/demo/         # 🚀 Example application
│   └── main.go
│
├── go.mod
└── README.md
```

---

## Quick Start

```go
package main

import (
    "gioui.org/app"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/widget"

    "github.com/hongshengjie/gioui-kit/pkg/component"
    kit "github.com/hongshengjie/gioui-kit/pkg/layout"
    "github.com/hongshengjie/gioui-kit/pkg/theme"
)

func main() {
    th := theme.Light()
    var btnClick widget.Clickable

    // ... inside your frame loop:
    // A card with a button inside a padded flex column
    kit.FlexCol{Gap: 16}.Layout(gtx,
        kit.Rigid(func(gtx layout.Context) layout.Dimensions {
            return component.NewCard(th).WithBorder().CardWithHeader(gtx, "Welcome",
                func(gtx layout.Context) layout.Dimensions {
                    return component.NewButton(th, &btnClick, "Get Started").
                        WithVariant(component.BtnPrimary).
                        Layout(gtx)
                },
            )
        }),
    )
}
```

---

## API Reference

### theme — Color & Typography System

#### Tailwind Color Palette

All Tailwind colors are exported as `color.NRGBA` variables:

```go
// Gray scales
theme.Slate50 ... theme.Slate950
theme.Gray50  ... theme.Gray950
theme.Zinc50  ... theme.Zinc950

// Chromatic colors
theme.Red50    ... theme.Red900
theme.Blue50   ... theme.Blue900
theme.Green50  ... theme.Green900
theme.Yellow50 ... theme.Yellow900
theme.Purple50 ... theme.Purple900
theme.Indigo50 ... theme.Indigo900
theme.Amber50  ... theme.Amber900
theme.Cyan50   ... theme.Cyan900
theme.Emerald50... theme.Emerald900
theme.Rose50   ... theme.Rose900

// Utilities
theme.White, theme.Black, theme.Transparent
theme.WithAlpha(c, 128)    // set alpha
theme.Opacity(c, 0.5)      // multiply alpha
theme.Lerp(a, b, 0.5)      // interpolate
theme.RGB(0xff6600)         // create from hex
```

#### DaisyUI Semantic Theme

```go
th := theme.Light()   // or Dark(), Cupcake(), Nord()

th.Primary            // Main brand color
th.PrimaryContent     // Text on Primary
th.Secondary          // Secondary brand
th.Accent             // Accent / highlight
th.Neutral            // Neutral dark
th.Info / Success / Warning / Error  // State colors

th.Base100            // Background (lightest)
th.Base200            // Background (medium)
th.Base300            // Background (darkest) / border
th.BaseContent        // Text on base backgrounds

// Typography scale
th.FontSize           // 16sp base
th.H1Size .. th.XsSize

// Spacing scale (Tailwind-like)
th.Space0 (0dp) .. th.Space16 (64dp)

// Border radius
th.RoundedNone (0) .. th.RoundedFull (9999)
```

---

### layout — Flexbox, Grid & Spacing

#### FlexRow / FlexCol

```go
// Horizontal flex with gap (like `flex flex-row gap-4 items-center`)
kit.FlexRow{
    Gap:       16,                     // gap-4
    Alignment: kit.ItemsCenter,        // items-center
    Spacing:   kit.JustifyBetween,     // justify-between
}.Layout(gtx,
    kit.Rigid(widgetA),                // flex child
    kit.Grow(widgetB),                 // flex-1
    kit.Flexed(0.5, widgetC),          // flex-grow: 0.5
)

// Vertical flex
kit.FlexCol{Gap: 8}.Layout(gtx, children...)
```

#### Alignment Constants (Tailwind naming)

| Tailwind Class   | GioUI Kit Constant     |
|------------------|------------------------|
| `items-start`    | `kit.ItemsStart`       |
| `items-center`   | `kit.ItemsCenter`      |
| `items-end`      | `kit.ItemsEnd`         |
| `justify-start`  | `kit.JustifyStart`     |
| `justify-center` | `kit.JustifyCenter`    |
| `justify-between`| `kit.JustifyBetween`   |
| `justify-around` | `kit.JustifyAround`    |

#### Grid

```go
// 3-column grid with gap (like `grid grid-cols-3 gap-4`)
kit.Grid{Cols: 3, Gap: 16}.Layout(gtx,
    widgetA, widgetB, widgetC,
    widgetD, widgetE,
)
```

#### Box / Container

```go
// Box (like <div class="p-4 bg-white rounded-lg border">)
kit.Box{
    Padding:    kit.P(16),
    Background: th.Base100,
    Radius:     th.RoundedLg,
    Border:     kit.Border{Color: th.Base300, Width: 1},
    MaxWidth:   400,
}.Layout(gtx, content)

// Centered max-width container (like `container mx-auto px-4`)
kit.Container{MaxWidth: 1280, Padding: 16}.Layout(gtx, content)
```

#### Spacing Utilities

```go
kit.P(16)              // p-4      uniform padding
kit.Px(16)             // px-4     horizontal padding
kit.Py(8)              // py-2     vertical padding
kit.Pt(12)             // pt-3     top padding
kit.Pb(12)             // pb-3     bottom padding
kit.Pl(8)              // pl-2     left padding
kit.Pr(8)              // pr-2     right padding
kit.Inset4(t, r, b, l) // custom all sides

kit.SpaceH(16)         // horizontal spacer widget
kit.SpaceV(16)         // vertical spacer widget
```

#### Size Utilities

```go
kit.W(200)(gtx, widget)       // w-[200px] fixed width
kit.H(100)(gtx, widget)       // h-[100px] fixed height
kit.MinW(300)(gtx, widget)    // min-w-[300px]
kit.MaxW(600)(gtx, widget)    // max-w-[600px]
kit.WFull(gtx, widget)        // w-full
```

#### Scrollable Lists

```go
scroll := kit.NewScrollY()
scroll.Layout(gtx, len(items), func(gtx layout.Context, i int) layout.Dimensions {
    return renderItem(gtx, items[i])
})
```

---

### modifier — Visual Decorators

```go
// Shadow (like `shadow-lg rounded-xl`)
modifier.Shadow{
    Style:  modifier.ShadowLg,
    Radius: th.RoundedXl,
}.Layout(gtx, content)

// Available presets: ShadowSm, ShadowMd, ShadowLg, ShadowXl, Shadow2xl

// Background with radius
modifier.Bg{Color: th.Primary, Radius: th.RoundedLg}.Layout(gtx, content)

// Clip to rounded shape
modifier.Rounded{Radius: th.RoundedFull}.Layout(gtx, content)

// Focus ring (like `ring-2 ring-blue-500 ring-offset-2`)
modifier.Ring{
    Width:  2,
    Color:  th.Primary,
    Offset: 2,
    Radius: th.RoundedLg,
}.Layout(gtx, content)

// Linear gradient background
modifier.LinearGradient{
    From:   th.Primary,
    To:     th.Secondary,
    Dir:    modifier.GradientToRight,
    Radius: th.RoundedLg,
}.Layout(gtx, content)
```

---

### component — DaisyUI Components

#### Button

```go
btn := component.NewButton(th, &clickable, "Click Me").
    WithVariant(component.BtnPrimary).   // BtnSecondary, BtnAccent, BtnGhost, ...
    WithSize(component.BtnLg)            // BtnXs, BtnSm, BtnMd, BtnLg
btn.Disabled = true
btn.Loading = true
btn.Layout(gtx)
```

| Variant          | Description            |
|------------------|------------------------|
| `BtnDefault`     | Neutral background     |
| `BtnPrimary`     | Brand primary color    |
| `BtnSecondary`   | Secondary color        |
| `BtnAccent`      | Accent color           |
| `BtnGhost`       | Transparent bg         |
| `BtnLink`        | Text link style        |
| `BtnOutline`     | Bordered, no fill      |
| `BtnInfo/Success/Warning/Error` | State colors |

#### Badge

```go
component.NewBadge(th, "NEW").WithVariant(component.BadgePrimary).Layout(gtx)
// Variants: BadgeDefault, BadgePrimary, BadgeSecondary, BadgeAccent,
//           BadgeSuccess, BadgeWarning, BadgeError, BadgeGhost
```

#### Card

```go
card := component.NewCard(th).WithBorder().WithCompact()
card.Layout(gtx, content)
card.CardWithHeader(gtx, "Title", content)
```

#### Alert

```go
component.NewAlert(th, "Success!", component.AlertSuccess).Layout(gtx)
// Variants: AlertInfo, AlertSuccess, AlertWarning, AlertError
```

#### Input

```go
input := component.NewInput(th, &editor, "Placeholder...").
    WithLabel("Email").
    WithVariant(component.InputPrimary)
input.Layout(gtx)
// Variants: InputBordered, InputGhost, InputPrimary, InputError, ...
// Sizes: InputXs, InputSm, InputMd, InputLg
```

#### Toggle

```go
component.NewToggle(th, &boolState, "Enable feature").Layout(gtx)
```

#### Avatar

```go
avatar := component.NewAvatar(th, "HS")
avatar.Size = component.AvatarLg
avatar.Layout(gtx)
```

#### Progress

```go
p := component.NewProgress(th, 0.65)
p.Variant = component.ProgressSuccess
p.Layout(gtx)
```

#### Tabs

```go
tabs := component.NewTabs(th, []string{"Tab 1", "Tab 2", "Tab 3"})
tabs.Layout(gtx)
selected := tabs.Selected
```

#### Text

```go
component.NewText(th, "Hello World").H1().Bold().WithColor(th.Primary).Layout(gtx)
// Sizes: .H1(), .H2(), .H3(), .H4(), .Sm(), .Xs()
// Weight: .Bold()
```

#### Chip / Skeleton

```go
component.NewChip(th, "Go").Layout(gtx)
component.NewSkeleton(th).Layout(gtx)
```

---

### scaffold — App-Level Layouts

#### AppShell

```go
shell := scaffold.NewAppShell(th).
    WithNavbar(navbarWidget).
    WithSidebar(sidebarWidget, 256).
    WithContent(contentWidget)
shell.Layout(gtx)
```

#### Navbar

```go
navbar := scaffold.NewNavbar(th)
navbar.Layout(gtx,
    startWidget,   // left section (brand/logo)
    centerWidget,  // center section
    endWidget,     // right section (actions)
)
```

#### Sidebar

```go
items := []scaffold.SidebarItem{
    {Label: "Dashboard", Icon: "◉", Active: true},
    {Label: "Settings",  Icon: "⚙"},
}
sidebar := scaffold.NewSidebar(th, items)
sidebar.OnSelect = func(i int) { /* handle */ }
sidebar.Layout(gtx)
```

#### Modal

```go
modal := scaffold.NewModal(th)
modal.Show()   // or modal.Hide()
modal.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
    return dialogContent(gtx)
})
```

#### Drawer

```go
drawer := scaffold.NewDrawer(th)
drawer.Side = scaffold.DrawerLeft  // or DrawerRight
drawer.Open()                      // .Close(), .Toggle()
drawer.Layout(gtx, drawerContent)
```

#### BottomNav (Mobile)

```go
nav := scaffold.NewBottomNav(th, []scaffold.BottomNavItem{
    {Label: "Home", Icon: "⌂", Active: true},
    {Label: "Search", Icon: "🔍"},
})
nav.Layout(gtx)
```

#### Toast

```go
toast := scaffold.NewToast(th)
toast.Show("Saved!")
toast.Position = scaffold.ToastBottomRight
toast.Layout(gtx)
```

#### Breadcrumb

```go
scaffold.NewBreadcrumb(th, "Home", "Products", "Detail").Layout(gtx)
```

---

## Tailwind ↔ GioUI Kit Cheat Sheet

| Tailwind CSS                    | GioUI Kit                                        |
|---------------------------------|--------------------------------------------------|
| `flex flex-row gap-4`           | `kit.FlexRow{Gap: 16}`                           |
| `flex flex-col gap-2`           | `kit.FlexCol{Gap: 8}`                            |
| `items-center`                  | `Alignment: kit.ItemsCenter`                     |
| `justify-between`               | `Spacing: kit.JustifyBetween`                    |
| `flex-1`                        | `kit.Grow(widget)`                               |
| `grid grid-cols-3 gap-4`       | `kit.Grid{Cols: 3, Gap: 16}`                     |
| `p-4`                           | `kit.P(16)`                                      |
| `px-6 py-2`                     | `kit.Px(24)` + `kit.Py(8)`                       |
| `w-full`                        | `kit.WFull`                                      |
| `max-w-lg`                      | `kit.MaxW(512)`                                  |
| `bg-white rounded-lg shadow-md` | `modifier.Shadow{ShadowMd, RoundedLg} + Bg{}`   |
| `text-sm font-bold text-gray-500`| `NewText(th, "...").Sm().Bold().WithColor(Gray500)` |
| `btn btn-primary btn-lg`       | `NewButton(th, &c, "Go").WithVariant(BtnPrimary).WithSize(BtnLg)` |
| `badge badge-accent`           | `NewBadge(th, "Tag").WithVariant(BadgeAccent)`   |
| `card card-bordered`           | `NewCard(th).WithBorder()`                       |
| `alert alert-success`          | `NewAlert(th, "OK!", AlertSuccess)`              |
| `input input-bordered`         | `NewInput(th, &ed, "...").WithVariant(InputBordered)` |
| `toggle`                       | `NewToggle(th, &b, "Label")`                     |

---

## Run the Demo

```bash
cd cmd/demo
go run .
```

---

## License

MIT
