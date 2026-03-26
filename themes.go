package design

// DefaultTheme returns the default design tokens
func DefaultTheme() *DesignTokens {
	return &DesignTokens{
		Theme:      "default",
		Color:      "#E5E7EB",
		Background: "#020617",
		Accent:     "#1D4ED8",
		FontFamily: "system-ui",
		Radius:     16,
		Padding:    16,
		Density:    "comfortable",
		Mode:       "dark",

		FontSizeXS:       "xs",
		FontSizeS:        "s",
		FontSizeM:        "m",
		FontSizeL:        "l",
		FontSizeXL:       "xl",
		FontWeightNormal: "normal",
		FontWeightBold:   "bold",

		SelectionBackground: "#1D4ED880",
		SelectionForeground: "#E5E7EB",

		ToastInfoColor:    "#1D4ED8",
		ToastSuccessColor: "#10B981",
		ToastWarningColor: "#F59E0B",
		ToastErrorColor:   "#EF4444",
		Layout:            DefaultLayoutTokens(),
	}
}

// MidnightTheme returns the midnight theme (dark mode)
func MidnightTheme() *DesignTokens {
	return &DesignTokens{
		Theme:      "midnight",
		Color:      "#E5E7EB",
		Background: "#020617",
		Accent:     "#1D4ED8",
		FontFamily: "system-ui",
		Radius:     16,
		Padding:    16,
		Density:    "comfortable",
		Mode:       "dark",

		FontSizeXS:       "xs",
		FontSizeS:        "s",
		FontSizeM:        "m",
		FontSizeL:        "l",
		FontSizeXL:       "xl",
		FontWeightNormal: "normal",
		FontWeightBold:   "bold",

		SelectionBackground: "#1D4ED880",
		SelectionForeground: "#E5E7EB",

		ToastInfoColor:    "#1D4ED8",
		ToastSuccessColor: "#10B981",
		ToastWarningColor: "#F59E0B",
		ToastErrorColor:   "#EF4444",
		Layout:            DefaultLayoutTokens(),
	}
}

// NordTheme returns the Nord theme (dark mode)
func NordTheme() *DesignTokens {
	return &DesignTokens{
		Theme:      "nord",
		Color:      "#ECEFF4",
		Background: "#2E3440",
		Accent:     "#5E81AC",
		FontFamily: "system-ui",
		Radius:     16,
		Padding:    16,
		Density:    "comfortable",
		Mode:       "dark",

		FontSizeXS:       "xs",
		FontSizeS:        "s",
		FontSizeM:        "m",
		FontSizeL:        "l",
		FontSizeXL:       "xl",
		FontWeightNormal: "normal",
		FontWeightBold:   "bold",

		SelectionBackground: "#5E81AC80",
		SelectionForeground: "#ECEFF4",

		ToastInfoColor:    "#5E81AC",
		ToastSuccessColor: "#A3BE8C",
		ToastWarningColor: "#EBCB8B",
		ToastErrorColor:   "#BF616A",
		Layout:            DefaultLayoutTokens(),
	}
}

// PaperTheme returns the Paper theme (light mode)
func PaperTheme() *DesignTokens {
	return &DesignTokens{
		Theme:      "paper",
		Color:      "#1F2937",
		Background: "#F9FAFB",
		Accent:     "#3B82F6",
		FontFamily: "system-ui",
		Radius:     16,
		Padding:    16,
		Density:    "comfortable",
		Mode:       "light",

		FontSizeXS:       "xs",
		FontSizeS:        "s",
		FontSizeM:        "m",
		FontSizeL:        "l",
		FontSizeXL:       "xl",
		FontWeightNormal: "normal",
		FontWeightBold:   "bold",

		SelectionBackground: "#3B82F640",
		SelectionForeground: "#1F2937",

		ToastInfoColor:    "#3B82F6",
		ToastSuccessColor: "#10B981",
		ToastWarningColor: "#F59E0B",
		ToastErrorColor:   "#EF4444",
		Layout:            DefaultLayoutTokens(),
	}
}

// WrappedTheme returns the Wrapped theme (dark mode with special styling)
func WrappedTheme() *DesignTokens {
	return &DesignTokens{
		Theme:      "wrapped",
		Color:      "#EC4899",
		Background: "#020617",
		Accent:     "#7B58C9",
		FontFamily: "system-ui",
		Radius:     20, // Special larger radius for wrapped theme
		Padding:    16,
		Density:    "comfortable",
		Mode:       "dark",

		FontSizeXS:       "xs",
		FontSizeS:        "s",
		FontSizeM:        "m",
		FontSizeL:        "l",
		FontSizeXL:       "xl",
		FontWeightNormal: "normal",
		FontWeightBold:   "bold",

		SelectionBackground: "#7B58C980",
		SelectionForeground: "#EC4899",

		ToastInfoColor:    "#7B58C9",
		ToastSuccessColor: "#10B981",
		ToastWarningColor: "#F59E0B",
		ToastErrorColor:   "#EF4444",
		Layout:            DefaultLayoutTokens(),
	}
}

// CustomTheme creates a theme from query parameters
func CustomTheme(params map[string]string) *DesignTokens {
	return ResolveDesignTokens(params)
}

// LightMode returns a copy of the tokens with light mode applied
func (dt *DesignTokens) LightMode() *DesignTokens {
	lightTokens := *dt // Copy struct
	lightTokens.Mode = "light"

	// Apply light variants if available
	if dt.ColorLight != "" {
		lightTokens.Color = dt.ColorLight
	}
	if dt.BackgroundLight != "" {
		lightTokens.Background = dt.BackgroundLight
	}
	if dt.AccentLight != "" {
		lightTokens.Accent = dt.AccentLight
	}

	return &lightTokens
}

// DarkMode returns a copy of the tokens with dark mode applied
func (dt *DesignTokens) DarkMode() *DesignTokens {
	darkTokens := *dt // Copy struct
	darkTokens.Mode = "dark"

	// Apply dark variants if available
	if dt.ColorDark != "" {
		darkTokens.Color = dt.ColorDark
	}
	if dt.BackgroundDark != "" {
		darkTokens.Background = dt.BackgroundDark
	}
	if dt.AccentDark != "" {
		darkTokens.Accent = dt.AccentDark
	}

	return &darkTokens
}
