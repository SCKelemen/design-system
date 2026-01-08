package design

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SCKelemen/color"
)

// DesignTokens represents visual design configuration
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

	// Light/dark variant colors (if specified, override base colors based on mode)
	ColorLight      string
	ColorDark       string
	BackgroundLight string
	BackgroundDark  string
	AccentLight     string
	AccentDark      string

	// Radix UI theme tokens
	RadixAccentColor string // "pink", "blue", "green", etc.
	RadixGrayColor   string // "mauve", "slate", "gray", etc.
	RadixRadius      string // "none", "small", "medium", "large", "full"
	RadixScaling     string // "90%", "95%", "100%", "105%", "110%"

	// Layout configuration
	Layout *LayoutTokens
}

// LayoutTokens represents spacing and dimension configuration
type LayoutTokens struct {
	// Spacing scale (follows 8pt grid system)
	SpaceXS  int // 4px
	SpaceS   int // 8px
	SpaceM   int // 16px
	SpaceL   int // 20px
	SpaceXL  int // 24px
	Space2XL int // 32px

	// Card dimensions
	CardPaddingLeft   int // Horizontal padding inside cards
	CardPaddingRight  int
	CardPaddingTop    int // Vertical padding inside cards
	CardPaddingBottom int
	CardTitleHeight   int // Height reserved for title area
	CardIconWidth     int // Width of icon
	CardIconSpacing   int // Space between icon and title
	CardHeaderPadding int // Padding for header items

	// Component heights
	StatCardHeight      int // Height for stat cards without trend
	StatCardHeightTrend int // Height for stat cards with trend graph
	TrendGraphMinHeight int // Minimum height for trend graphs

	// Grid defaults
	DefaultGridGap     float64 // Default gap between grid items
	DefaultGridWidth   float64 // Default grid container width
	DefaultGridColumns int     // Default number of columns
}

// DefaultLayoutTokens returns the default layout token values
func DefaultLayoutTokens() *LayoutTokens {
	return &LayoutTokens{
		// Spacing scale
		SpaceXS:  4,
		SpaceS:   8,
		SpaceM:   16,
		SpaceL:   20,
		SpaceXL:  24,
		Space2XL: 32,

		// Card dimensions (migrated from components.go constants)
		CardPaddingLeft:   20,
		CardPaddingRight:  20,
		CardPaddingTop:    20,
		CardPaddingBottom: 20,
		CardTitleHeight:   50,
		CardIconWidth:     20,
		CardIconSpacing:   8,
		CardHeaderPadding: 10,

		// Component heights
		StatCardHeight:      70,
		StatCardHeightTrend: 84,
		TrendGraphMinHeight: 15,

		// Grid defaults
		DefaultGridGap:     8.0,
		DefaultGridWidth:   1000.0,
		DefaultGridColumns: 3,
	}
}

// MotionTokens represents animation configuration
type MotionTokens struct {
	Level      string // "none", "subtle", "regular", "loud"
	Durations  map[string]string
	Amplitudes map[string]float64
}

// ResolveDesignTokens resolves design tokens from query parameters
// If auto_color_scheme is true, returns tokens that adapt to light/dark mode
func ResolveDesignTokens(queryParams map[string]string) *DesignTokens {
	tokens := &DesignTokens{
		Theme:      "default",
		Color:      "#E5E7EB",
		Background: "#020617",
		Accent:     "#1D4ED8",
		FontFamily: "system-ui",
		Radius:     16,
		Padding:    16,
		Density:    "comfortable",
		Mode:       "dark",
		Layout:     DefaultLayoutTokens(),
	}

	// Check for Radix UI theme tokens first
	if accentColor, ok := queryParams["accentColor"]; ok && accentColor != "" {
		tokens.RadixAccentColor = accentColor
	}
	if grayColor, ok := queryParams["grayColor"]; ok && grayColor != "" {
		tokens.RadixGrayColor = grayColor
	}
	if radius, ok := queryParams["radius"]; ok && radius != "" {
		// Check if it's a Radix radius token (string) vs numeric
		if radius == "none" || radius == "small" || radius == "medium" || radius == "large" || radius == "full" {
			tokens.RadixRadius = radius
		} else {
			// Try to parse as integer
			if r, err := strconv.Atoi(radius); err == nil {
				tokens.Radius = r
			}
		}
	}
	if scaling, ok := queryParams["scaling"]; ok && scaling != "" {
		tokens.RadixScaling = scaling
	}

	// Apply Radix theme if Radix tokens are present
	if tokens.RadixAccentColor != "" || tokens.RadixGrayColor != "" {
		applyRadixTheme(tokens)
	}

	// Apply theme if specified (and no Radix theme)
	if theme, ok := queryParams["theme"]; ok && theme != "" && tokens.RadixAccentColor == "" {
		applyTheme(tokens, theme)
	}

	// Helper function to parse color (supports single or light/dark format)
	parseColor := func(colorStr string) (string, string) {
		if colorStr == "" {
			return "", ""
		}

		// Check for dual color format: LIGHT/DARK
		if strings.Contains(colorStr, "/") {
			parts := strings.Split(colorStr, "/")
			if len(parts) == 2 {
				light := strings.TrimSpace(parts[0])
				dark := strings.TrimSpace(parts[1])

				// Validate both colors using the color package
				// Query params never have # prefix (it's a URL fragment delimiter)
				lightColor := light
				if !strings.HasPrefix(lightColor, "#") {
					lightColor = "#" + lightColor
				}
				darkColor := dark
				if !strings.HasPrefix(darkColor, "#") {
					darkColor = "#" + darkColor
				}

				// Validate parsing (but don't fail if invalid, just use as-is)
				if _, err := color.ParseColor(lightColor); err == nil {
					light = lightColor
				} else {
					light = lightColor // Use as-is even if parsing fails
				}
				if _, err := color.ParseColor(darkColor); err == nil {
					dark = darkColor
				} else {
					dark = darkColor // Use as-is even if parsing fails
				}

				return light, dark
			}
		}

		// Single color format - use for both modes
		// Query params never have # prefix
		singleColor := colorStr
		if !strings.HasPrefix(singleColor, "#") {
			singleColor = "#" + singleColor
		}

		// Validate parsing (but don't fail if invalid, just use as-is)
		if _, err := color.ParseColor(singleColor); err == nil {
			return singleColor, singleColor
		}
		return singleColor, singleColor // Use as-is even if parsing fails
	}

	// Override with individual parameters
	// Support both single color and light/dark variants (format: COLOR or LIGHT/DARK)
	if color, ok := queryParams["color"]; ok && color != "" {
		light, dark := parseColor(color)
		if light != "" {
			tokens.ColorLight = light
			tokens.ColorDark = dark
			// Set base color based on current mode
			if tokens.Mode == "light" {
				tokens.Color = light
			} else {
				tokens.Color = dark
			}
		}
	}
	// Backwards compatibility: still support color_light and color_dark
	if colorLight, ok := queryParams["color_light"]; ok && colorLight != "" {
		if !strings.HasPrefix(colorLight, "#") {
			colorLight = "#" + colorLight
		}
		tokens.ColorLight = colorLight
		if tokens.Mode == "light" {
			tokens.Color = colorLight
		}
	}
	if colorDark, ok := queryParams["color_dark"]; ok && colorDark != "" {
		if !strings.HasPrefix(colorDark, "#") {
			colorDark = "#" + colorDark
		}
		tokens.ColorDark = colorDark
		if tokens.Mode == "dark" {
			tokens.Color = colorDark
		}
	}

	if bg, ok := queryParams["background"]; ok && bg != "" {
		light, dark := parseColor(bg)
		if light != "" {
			tokens.BackgroundLight = light
			tokens.BackgroundDark = dark
			// Set base background based on current mode
			if tokens.Mode == "light" {
				tokens.Background = light
			} else {
				tokens.Background = dark
			}
		}
	}
	// Backwards compatibility: still support background_light and background_dark
	if bgLight, ok := queryParams["background_light"]; ok && bgLight != "" {
		if !strings.HasPrefix(bgLight, "#") {
			bgLight = "#" + bgLight
		}
		tokens.BackgroundLight = bgLight
		if tokens.Mode == "light" {
			tokens.Background = bgLight
		}
	}
	if bgDark, ok := queryParams["background_dark"]; ok && bgDark != "" {
		if !strings.HasPrefix(bgDark, "#") {
			bgDark = "#" + bgDark
		}
		tokens.BackgroundDark = bgDark
		if tokens.Mode == "dark" {
			tokens.Background = bgDark
		}
	}

	if accent, ok := queryParams["accent"]; ok && accent != "" {
		light, dark := parseColor(accent)
		if light != "" {
			tokens.AccentLight = light
			tokens.AccentDark = dark
			// Set base accent based on current mode
			if tokens.Mode == "light" {
				tokens.Accent = light
			} else {
				tokens.Accent = dark
			}
		}
	}
	// Backwards compatibility: still support accent_light and accent_dark
	if accentLight, ok := queryParams["accent_light"]; ok && accentLight != "" {
		if !strings.HasPrefix(accentLight, "#") {
			accentLight = "#" + accentLight
		}
		tokens.AccentLight = accentLight
		if tokens.Mode == "light" {
			tokens.Accent = accentLight
		}
	}
	if accentDark, ok := queryParams["accent_dark"]; ok && accentDark != "" {
		if !strings.HasPrefix(accentDark, "#") {
			accentDark = "#" + accentDark
		}
		tokens.AccentDark = accentDark
		if tokens.Mode == "dark" {
			tokens.Accent = accentDark
		}
	}

	if density, ok := queryParams["density"]; ok && density != "" {
		if density == "compact" || density == "comfortable" {
			tokens.Density = density
		}
	}

	// Handle mode - if not specified, try to infer from theme
	if mode, ok := queryParams["mode"]; ok && mode != "" {
		if mode == "light" || mode == "dark" {
			tokens.Mode = mode
		}
	} else {
		// If theme was specified without explicit mode, check if it has a mode suffix
		if theme, ok := queryParams["theme"]; ok && theme != "" {
			if strings.HasSuffix(theme, "-light") {
				tokens.Mode = "light"
			} else if strings.HasSuffix(theme, "-dark") {
				tokens.Mode = "dark"
			}
		}
	}

	// Apply Radix radius to numeric radius
	if tokens.RadixRadius != "" {
		tokens.Radius = radixRadiusToPixels(tokens.RadixRadius)
	}

	// Apply Radix scaling to padding and other spacing
	if tokens.RadixScaling != "" {
		scale := radixScalingToFloat(tokens.RadixScaling)
		tokens.Padding = int(float64(tokens.Padding) * scale)
		if tokens.Radius > 0 {
			tokens.Radius = int(float64(tokens.Radius) * scale)
		}
	}

	// Apply light/dark variant colors based on current mode
	// If variants are specified, they override the base colors
	// (This is already handled in the parsing above, but ensure consistency)
	if tokens.Mode == "light" {
		if tokens.ColorLight != "" {
			tokens.Color = tokens.ColorLight
		}
		if tokens.BackgroundLight != "" {
			tokens.Background = tokens.BackgroundLight
		}
		if tokens.AccentLight != "" {
			tokens.Accent = tokens.AccentLight
		}
	} else if tokens.Mode == "dark" {
		if tokens.ColorDark != "" {
			tokens.Color = tokens.ColorDark
		}
		if tokens.BackgroundDark != "" {
			tokens.Background = tokens.BackgroundDark
		}
		if tokens.AccentDark != "" {
			tokens.Accent = tokens.AccentDark
		}
	}

	return tokens
}

// ResolveDesignTokensForBothModes resolves design tokens for both light and dark modes
// This is useful for generating adaptive SVGs that respond to color scheme
func ResolveDesignTokensForBothModes(queryParams map[string]string) (*DesignTokens, *DesignTokens) {
	// Create copies of params for light and dark
	lightParams := make(map[string]string)
	darkParams := make(map[string]string)
	for k, v := range queryParams {
		lightParams[k] = v
		darkParams[k] = v
	}

	// Force mode for each
	lightParams["mode"] = "light"
	darkParams["mode"] = "dark"

	// If theme is specified without mode suffix, apply it to both
	if theme, ok := queryParams["theme"]; ok && theme != "" {
		if !strings.HasSuffix(theme, "-light") && !strings.HasSuffix(theme, "-dark") {
			// Apply theme variant to both
			lightParams["theme"] = theme + "-light"
			darkParams["theme"] = theme + "-dark"
		} else {
			// Theme already has mode, use as-is
			lightParams["theme"] = theme
			darkParams["theme"] = theme
		}
	}

	// Handle color variants: parse LIGHT/DARK format or use single color
	parseColorForMode := func(colorStr string, mode string) string {
		if colorStr == "" {
			return ""
		}
		// Check for dual color format: LIGHT/DARK
		if strings.Contains(colorStr, "/") {
			parts := strings.Split(colorStr, "/")
			if len(parts) == 2 {
				var selectedColor string
				if mode == "light" {
					selectedColor = strings.TrimSpace(parts[0])
				} else {
					selectedColor = strings.TrimSpace(parts[1])
				}

				// Add # prefix if not present and validate with color package
				if !strings.HasPrefix(selectedColor, "#") {
					selectedColor = "#" + selectedColor
				}

				// Validate parsing (but don't fail if invalid, just use as-is)
				if _, err := color.ParseColor(selectedColor); err == nil {
					return selectedColor
				}
				return selectedColor // Use as-is even if parsing fails
			}
		}
		// Single color - use for both modes
		singleColor := colorStr
		if !strings.HasPrefix(singleColor, "#") {
			singleColor = "#" + singleColor
		}

		// Validate parsing (but don't fail if invalid, just use as-is)
		if _, err := color.ParseColor(singleColor); err == nil {
			return singleColor
		}
		return singleColor // Use as-is even if parsing fails
	}

	if color, ok := queryParams["color"]; ok && color != "" {
		lightColor := parseColorForMode(color, "light")
		darkColor := parseColorForMode(color, "dark")
		if lightColor != "" {
			lightParams["color"] = lightColor
		}
		if darkColor != "" {
			darkParams["color"] = darkColor
		}
	}

	if bg, ok := queryParams["background"]; ok && bg != "" {
		lightBg := parseColorForMode(bg, "light")
		darkBg := parseColorForMode(bg, "dark")
		if lightBg != "" {
			lightParams["background"] = lightBg
		}
		if darkBg != "" {
			darkParams["background"] = darkBg
		}
	}

	if accent, ok := queryParams["accent"]; ok && accent != "" {
		lightAccent := parseColorForMode(accent, "light")
		darkAccent := parseColorForMode(accent, "dark")
		if lightAccent != "" {
			lightParams["accent"] = lightAccent
		}
		if darkAccent != "" {
			darkParams["accent"] = darkAccent
		}
	}

	lightTokens := ResolveDesignTokens(lightParams)
	darkTokens := ResolveDesignTokens(darkParams)

	return lightTokens, darkTokens
}

// ResolveMotionTokens resolves motion tokens from query parameters
func ResolveMotionTokens(queryParams map[string]string) *MotionTokens {
	tokens := &MotionTokens{
		Level: "subtle",
		Durations: map[string]string{
			"fast":   "0.7s",
			"normal": "1.6s",
			"slow":   "2.8s",
		},
		Amplitudes: map[string]float64{
			"scaleCard":  0.03,
			"ledBreathe": 0.06,
		},
	}

	if motion, ok := queryParams["motion"]; ok && motion != "" {
		switch motion {
		case "none", "subtle", "regular", "loud":
			tokens.Level = motion
		}
	}

	// Adjust durations and amplitudes based on level
	switch tokens.Level {
	case "none":
		tokens.Durations["fast"] = "0s"
		tokens.Durations["normal"] = "0s"
		tokens.Durations["slow"] = "0s"
	case "subtle":
		tokens.Durations["fast"] = "1.0s"
		tokens.Durations["normal"] = "2.4s"
		tokens.Durations["slow"] = "4.0s"
		tokens.Amplitudes["scaleCard"] = 0.02
		tokens.Amplitudes["ledBreathe"] = 0.04
	case "regular":
		tokens.Durations["fast"] = "0.7s"
		tokens.Durations["normal"] = "1.6s"
		tokens.Durations["slow"] = "2.8s"
		tokens.Amplitudes["scaleCard"] = 0.03
		tokens.Amplitudes["ledBreathe"] = 0.06
	case "loud":
		tokens.Durations["fast"] = "0.5s"
		tokens.Durations["normal"] = "1.2s"
		tokens.Durations["slow"] = "2.0s"
		tokens.Amplitudes["scaleCard"] = 0.05
		tokens.Amplitudes["ledBreathe"] = 0.10
	}

	return tokens
}

// applyTheme applies a named theme to design tokens
// Supports theme variants: "nord", "nord-light", "nord-dark", etc.
func applyTheme(tokens *DesignTokens, theme string) {
	// Parse theme name and mode
	themeName := theme
	mode := tokens.Mode // Use existing mode or default

	// Check for explicit mode suffix (e.g., "nord-light", "nord-dark")
	if strings.HasSuffix(theme, "-light") {
		themeName = strings.TrimSuffix(theme, "-light")
		mode = "light"
	} else if strings.HasSuffix(theme, "-dark") {
		themeName = strings.TrimSuffix(theme, "-dark")
		mode = "dark"
	}

	// Theme definitions with light/dark variants
	themes := map[string]map[string]map[string]string{
		"nord": {
			"light": {
				"color":      "#2E3440",
				"background": "#ECEFF4",
				"accent":     "#5E81AC",
			},
			"dark": {
				"color":      "#ECEFF4",
				"background": "#2E3440",
				"accent":     "#5E81AC",
			},
		},
		"midnight": {
			"light": {
				"color":      "#1F2937",
				"background": "#F9FAFB",
				"accent":     "#2563EB",
			},
			"dark": {
				"color":      "#E5E7EB",
				"background": "#020617",
				"accent":     "#1D4ED8",
			},
		},
		"paper": {
			"light": {
				"color":      "#1F2937",
				"background": "#F9FAFB",
				"accent":     "#3B82F6",
			},
			"dark": {
				"color":      "#E5E7EB",
				"background": "#1F2937",
				"accent":     "#60A5FA",
			},
		},
		"wrapped": {
			"light": {
				"color":      "#1F2937",
				"background": "#FDF2F8",
				"accent":     "#EC4899",
			},
			"dark": {
				"color":      "#EC4899",
				"background": "#020617",
				"accent":     "#7B58C9",
			},
		},
		"default": {
			"light": {
				"color":      "#1F2937",
				"background": "#FFFFFF",
				"accent":     "#2563EB",
			},
			"dark": {
				"color":      "#E5E7EB",
				"background": "#020617",
				"accent":     "#1D4ED8",
			},
		},
	}

	// Apply theme colors based on mode
	if themeMap, ok := themes[themeName]; ok {
		if modeMap, ok := themeMap[mode]; ok {
			tokens.Color = modeMap["color"]
			tokens.Background = modeMap["background"]
			tokens.Accent = modeMap["accent"]
			tokens.Mode = mode
		} else {
			// Fallback to dark if mode not found
			if darkMap, ok := themeMap["dark"]; ok {
				tokens.Color = darkMap["color"]
				tokens.Background = darkMap["background"]
				tokens.Accent = darkMap["accent"]
				tokens.Mode = "dark"
			}
		}

		// Special handling for wrapped theme
		if themeName == "wrapped" {
			tokens.Radius = 20
		}
	}
}

// applyRadixTheme applies Radix UI theme tokens
func applyRadixTheme(tokens *DesignTokens) {
	// Radix accent colors (approximate mappings)
	accentColors := map[string]map[string]string{
		"pink":   {"light": "#EC4899", "dark": "#F472B6"},
		"blue":   {"light": "#3B82F6", "dark": "#60A5FA"},
		"green":  {"light": "#10B981", "dark": "#34D399"},
		"purple": {"light": "#8B5CF6", "dark": "#A78BFA"},
		"red":    {"light": "#EF4444", "dark": "#F87171"},
		"orange": {"light": "#F97316", "dark": "#FB923C"},
		"yellow": {"light": "#EAB308", "dark": "#FCD34D"},
		"cyan":   {"light": "#06B6D4", "dark": "#22D3EE"},
		"violet": {"light": "#7C3AED", "dark": "#8B5CF6"},
		"indigo": {"light": "#6366F1", "dark": "#818CF8"},
	}

	// Radix gray colors (approximate mappings)
	grayColors := map[string]map[string]map[string]string{
		"mauve": {
			"light": {"bg": "#FDFCFD", "fg": "#1A1523", "border": "#E9E4ED"},
			"dark":  {"bg": "#1A1523", "fg": "#EDE9FE", "border": "#2F2655"},
		},
		"slate": {
			"light": {"bg": "#FBFCFD", "fg": "#1E293B", "border": "#E2E8F0"},
			"dark":  {"bg": "#0F172A", "fg": "#F1F5F9", "border": "#1E293B"},
		},
		"gray": {
			"light": {"bg": "#FBFBFB", "fg": "#1C1C1F", "border": "#E4E4E7"},
			"dark":  {"bg": "#111113", "fg": "#E4E4E7", "border": "#2A2A2B"},
		},
		"sage": {
			"light": {"bg": "#FBFDFC", "fg": "#1C211C", "border": "#E8EDE8"},
			"dark":  {"bg": "#141716", "fg": "#ECEDEC", "border": "#272D27"},
		},
		"olive": {
			"light": {"bg": "#FCFDFC", "fg": "#1C211C", "border": "#E8EDE8"},
			"dark":  {"bg": "#181B18", "fg": "#ECEDEC", "border": "#2A2E2A"},
		},
		"sand": {
			"light": {"bg": "#FAF9F6", "fg": "#1C1C1A", "border": "#E8E6E1"},
			"dark":  {"bg": "#161615", "fg": "#E8E6E1", "border": "#282826"},
		},
	}

	// Apply accent color
	if tokens.RadixAccentColor != "" {
		if colors, ok := accentColors[tokens.RadixAccentColor]; ok {
			if tokens.Mode == "light" {
				tokens.Accent = colors["light"]
			} else {
				tokens.Accent = colors["dark"]
			}
		}
	}

	// Apply gray color
	if tokens.RadixGrayColor != "" {
		if grayMap, ok := grayColors[tokens.RadixGrayColor]; ok {
			if modeMap, ok := grayMap[tokens.Mode]; ok {
				tokens.Background = modeMap["bg"]
				tokens.Color = modeMap["fg"]
			}
		}
	}
}

// radixRadiusToPixels converts Radix radius token to pixels
func radixRadiusToPixels(radius string) int {
	switch radius {
	case "none":
		return 0
	case "small":
		return 4
	case "medium":
		return 8
	case "large":
		return 16
	case "full":
		return 9999 // Very large for full rounding
	default:
		return 8 // Default to medium
	}
}

// radixScalingToFloat converts Radix scaling percentage to float multiplier
func radixScalingToFloat(scaling string) float64 {
	// Remove % if present
	scaling = strings.TrimSuffix(scaling, "%")

	var scale float64
	fmt.Sscanf(scaling, "%f", &scale)
	if scale == 0 {
		return 1.0
	}
	return scale / 100.0
}

// ToCSS converts design tokens to CSS string for SVG
func (dt *DesignTokens) ToCSS() string {
	return fmt.Sprintf(`
		:root {
			--color: %s;
			--background: %s;
			--accent: %s;
			--font-family: %s;
			--radius: %dpx;
			--padding: %dpx;
		}
	`, dt.Color, dt.Background, dt.Accent, dt.FontFamily, dt.Radius, dt.Padding)
}
