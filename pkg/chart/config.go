package chart

// ThemeType represents the available theme types
type ThemeType string

const (
	ThemeLight ThemeType = "light"
	ThemeDark  ThemeType = "dark"
	ThemeMono  ThemeType = "mono"
)

// HouseSystem represents different house systems
type HouseSystem string

const (
	HousePlacidus      HouseSystem = "Placidus"
	HouseKoch          HouseSystem = "Koch"
	HouseEqual         HouseSystem = "Equal"
	HouseCampanus      HouseSystem = "Campanus"
	HouseRegiomontanus HouseSystem = "Regiomontanus"
	HousePorphyry      HouseSystem = "Porphyry"
	HouseWholeSign     HouseSystem = "Whole_Sign"
)

// Orb configuration for aspects
type Orb struct {
	Conjunction int `json:"conjunction"`
	Opposition  int `json:"opposition"`
	Trine       int `json:"trine"`
	Square      int `json:"square"`
	Sextile     int `json:"sextile"`
	Quincunx    int `json:"quincunx"`
}

// DefaultOrb returns default orb values
func DefaultOrb() Orb {
	return Orb{
		Conjunction: 7,
		Opposition:  6,
		Trine:       6,
		Square:      6,
		Sextile:     5,
		Quincunx:    0,
	}
}

// Theme contains color definitions
type Theme struct {
	Fire         string  `json:"fire"`         // fire signs, square aspects, Asc
	Earth        string  `json:"earth"`        // earth signs, MC
	Air          string  `json:"air"`          // air signs, trine aspects
	Water        string  `json:"water"`        // water signs, opposition aspects
	Points       string  `json:"points"`       // lunar nodes, sextile aspects
	Asteroids    string  `json:"asteroids"`    // asteroids, quincunx aspects
	Positive     string  `json:"positive"`     // positive polarity
	Negative     string  `json:"negative"`     // negative polarity
	Others       string  `json:"others"`       // conjunction aspects
	Transparency float64 `json:"transparency"` // transparency level
	Foreground   string  `json:"foreground"`   // main foreground color
	Background   string  `json:"background"`   // main background color
	Dim          string  `json:"dim"`          // dimmed color for subtle elements
}

// LightTheme returns default light theme colors
func LightTheme() Theme {
	return Theme{
		Fire:         "#ef476f",
		Earth:        "#ffd166",
		Air:          "#06d6a0",
		Water:        "#81bce7",
		Points:       "#118ab2",
		Asteroids:    "#AA96DA",
		Positive:     "#FFC0CB",
		Negative:     "#AD8B73",
		Others:       "#FFA500",
		Transparency: 0.1,
		Foreground:   "#758492",
		Background:   "#FFFDF1",
		Dim:          "#A4BACD",
	}
}

// DarkTheme returns default dark theme colors
func DarkTheme() Theme {
	return Theme{
		Fire:         "#ef476f",
		Earth:        "#ffd166",
		Air:          "#06d6a0",
		Water:        "#81bce7",
		Points:       "#118ab2",
		Asteroids:    "#AA96DA",
		Positive:     "#FFC0CB",
		Negative:     "#AD8B73",
		Others:       "#FFA500",
		Transparency: 0.1,
		Foreground:   "#F7F3F0",
		Background:   "#343a40",
		Dim:          "#515860",
	}
}

// MonoTheme returns default monochrome theme colors
func MonoTheme() Theme {
	return Theme{
		Fire:         "#888888",
		Earth:        "#888888",
		Air:          "#888888",
		Water:        "#888888",
		Points:       "#888888",
		Asteroids:    "#888888",
		Positive:     "#888888",
		Negative:     "#888888",
		Others:       "#888888",
		Transparency: 0,
		Foreground:   "#888888",
		Background:   "#FFFFFF",
		Dim:          "#888888",
	}
}

// Display configuration for celestial bodies
type Display struct {
	Sun     bool `json:"sun"`
	Moon    bool `json:"moon"`
	Mercury bool `json:"mercury"`
	Venus   bool `json:"venus"`
	Mars    bool `json:"mars"`
	Jupiter bool `json:"jupiter"`
	Saturn  bool `json:"saturn"`
	Uranus  bool `json:"uranus"`
	Neptune bool `json:"neptune"`
	Pluto   bool `json:"pluto"`
	AscNode bool `json:"asc_node"`
	Chiron  bool `json:"chiron"`
	Ceres   bool `json:"ceres"`
	Pallas  bool `json:"pallas"`
	Juno    bool `json:"juno"`
	Vesta   bool `json:"vesta"`
	Asc     bool `json:"asc"`
	IC      bool `json:"ic"`
	Dsc     bool `json:"dsc"`
	MC      bool `json:"mc"`
}

// DefaultDisplay returns default display configuration
func DefaultDisplay() Display {
	return Display{
		Sun:     true,
		Moon:    true,
		Mercury: true,
		Venus:   true,
		Mars:    true,
		Jupiter: true,
		Saturn:  true,
		Uranus:  true,
		Neptune: true,
		Pluto:   true,
		AscNode: true,
		Chiron:  false,
		Ceres:   false,
		Pallas:  false,
		Juno:    false,
		Vesta:   false,
		Asc:     true,
		IC:      false,
		Dsc:     false,
		MC:      true,
	}
}

// ChartConfig contains chart rendering configuration
type ChartConfig struct {
	StrokeWidth           int     `json:"stroke_width"`
	StrokeOpacity         float64 `json:"stroke_opacity"`
	Font                  string  `json:"font"`
	FontSizeFraction      float64 `json:"font_size_fraction"`
	InnerMinDegree        float64 `json:"inner_min_degree"`
	OuterMinDegree        float64 `json:"outer_min_degree"`
	MarginFactor          float64 `json:"margin_factor"`
	RingThicknessFraction float64 `json:"ring_thickness_fraction"`
	ScaleAdjFactor        float64 `json:"scale_adj_factor"`
	PosAdjFactor          float64 `json:"pos_adj_factor"`
}

// DefaultChartConfig returns default chart configuration
func DefaultChartConfig() ChartConfig {
	return ChartConfig{
		StrokeWidth:           1,
		StrokeOpacity:         1.0,
		Font:                  "sans-serif",
		FontSizeFraction:      0.55,
		InnerMinDegree:        9.0,
		OuterMinDegree:        8.0,
		MarginFactor:          0.04,
		RingThicknessFraction: 0.15,
		ScaleAdjFactor:        600.0,
		PosAdjFactor:          2.2,
	}
}

// Config is the main configuration struct
type Config struct {
	ThemeType   ThemeType   `json:"theme_type"`
	HouseSystem HouseSystem `json:"house_system"`
	Orb         Orb         `json:"orb"`
	LightTheme  Theme       `json:"light_theme"`
	DarkTheme   Theme       `json:"dark_theme"`
	Display     Display     `json:"display"`
	Chart       ChartConfig `json:"chart"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		ThemeType:   ThemeDark,
		HouseSystem: HousePlacidus,
		Orb:         DefaultOrb(),
		LightTheme:  LightTheme(),
		DarkTheme:   DarkTheme(),
		Display:     DefaultDisplay(),
		Chart:       DefaultChartConfig(),
	}
}

// GetTheme returns the current theme based on theme type
func (c *Config) GetTheme() Theme {
	switch c.ThemeType {
	case ThemeLight:
		return c.LightTheme
	case ThemeDark:
		return c.DarkTheme
	case ThemeMono:
		return MonoTheme()
	default:
		return c.DarkTheme
	}
}

// GetOrbForAspect returns the orb value for a specific aspect
func (c *Config) GetOrbForAspect(aspectName string) int {
	switch aspectName {
	case "conjunction":
		return c.Orb.Conjunction
	case "opposition":
		return c.Orb.Opposition
	case "trine":
		return c.Orb.Trine
	case "square":
		return c.Orb.Square
	case "sextile":
		return c.Orb.Sextile
	case "quincunx":
		return c.Orb.Quincunx
	default:
		return 0
	}
}
