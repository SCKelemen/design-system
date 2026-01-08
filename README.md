# Design System

Centralized design tokens, themes, and styling configuration for SCKelemen projects.

## Features

- **Design Tokens**: Visual configuration including colors, spacing, typography
- **Layout Tokens**: Spacing scale, card dimensions, component heights, grid defaults
- **Motion Tokens**: Animation configuration with levels (none, subtle, regular, loud)
- **Predefined Themes**: Default, Midnight, Nord, Paper, Wrapped
- **Radix UI Integration**: Support for Radix UI theme tokens
- **Light/Dark Mode**: Automatic mode switching and variant colors
- **Query Parameter Resolution**: Parse and resolve tokens from URL query parameters

## Installation

```bash
go get github.com/SCKelemen/design-system
```

## Usage

### Using Predefined Themes

```go
import design "github.com/SCKelemen/design-system"

// Get midnight theme
tokens := design.MidnightTheme()

// Get paper theme (light mode)
tokens := design.PaperTheme()

// Switch modes
lightTokens := tokens.LightMode()
darkTokens := tokens.DarkMode()
```

### Custom Themes from Query Parameters

```go
// Parse query parameters
params := map[string]string{
    "theme": "nord",
    "mode": "dark",
    "color": "ECEFF4",
    "background": "2E3440",
    "accent": "5E81AC",
}

tokens := design.ResolveDesignTokens(params)
```

### Dual Color Format (Light/Dark)

```go
// Specify different colors for light and dark modes
params := map[string]string{
    "color": "1F2937/E5E7EB",       // light/dark
    "background": "FFFFFF/020617",   // light/dark
    "accent": "2563EB/1D4ED8",       // light/dark
}

tokens := design.ResolveDesignTokens(params)
```

### Radix UI Themes

```go
params := map[string]string{
    "accentColor": "pink",
    "grayColor": "mauve",
    "radius": "large",
    "scaling": "105%",
}

tokens := design.ResolveDesignTokens(params)
```

### Layout Tokens

```go
layout := design.DefaultLayoutTokens()

// Spacing scale (8pt grid)
fmt.Println(layout.SpaceS)   // 8
fmt.Println(layout.SpaceM)   // 16
fmt.Println(layout.SpaceL)   // 20

// Card dimensions
fmt.Println(layout.CardPaddingLeft)  // 20
fmt.Println(layout.CardTitleHeight)  // 50

// Grid defaults
fmt.Println(layout.DefaultGridGap)     // 8.0
fmt.Println(layout.DefaultGridColumns) // 3
```

### Motion Tokens

```go
params := map[string]string{
    "motion": "subtle",
}

motion := design.ResolveMotionTokens(params)

fmt.Println(motion.Durations["fast"])          // "1.0s"
fmt.Println(motion.Amplitudes["scaleCard"])    // 0.02
```

## Available Themes

- **default**: Standard light/dark theme
- **midnight**: Dark blue theme with high contrast
- **nord**: Nordic-inspired theme with cool tones
- **paper**: Clean light theme with subtle colors
- **wrapped**: Special theme with pink accents and larger radius

## Token Structure

### DesignTokens

```go
type DesignTokens struct {
    Theme      string
    Color      string
    Background string
    Accent     string
    FontFamily string
    Radius     int
    Padding    int
    Density    string // "compact" or "comfortable"
    Mode       string // "light" or "dark"

    // Light/dark variants
    ColorLight      string
    ColorDark       string
    BackgroundLight string
    BackgroundDark  string
    AccentLight     string
    AccentDark      string

    // Radix UI tokens
    RadixAccentColor string
    RadixGrayColor   string
    RadixRadius      string
    RadixScaling     string

    Layout *LayoutTokens
}
```

### LayoutTokens

```go
type LayoutTokens struct {
    // Spacing scale
    SpaceXS, SpaceS, SpaceM, SpaceL, SpaceXL, Space2XL int

    // Card dimensions
    CardPaddingLeft, CardPaddingRight, CardPaddingTop, CardPaddingBottom int
    CardTitleHeight, CardIconWidth, CardIconSpacing, CardHeaderPadding int

    // Component heights
    StatCardHeight, StatCardHeightTrend, TrendGraphMinHeight int

    // Grid defaults
    DefaultGridGap, DefaultGridWidth float64
    DefaultGridColumns int
}
```

### MotionTokens

```go
type MotionTokens struct {
    Level      string // "none", "subtle", "regular", "loud"
    Durations  map[string]string
    Amplitudes map[string]float64
}
```

## Dependencies

- [github.com/SCKelemen/color](https://github.com/SCKelemen/color) - Color manipulation and parsing
