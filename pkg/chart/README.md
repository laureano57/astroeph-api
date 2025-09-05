# Chart Package

Go port of SVG chart rendering functionality from the [natal](https://github.com/hoishing/natal) Python library.

**Original Library:** [hoishing/natal](https://github.com/hoishing/natal)  
**License:** MIT License Copyright (c) 2022 Kelvin Ng  
**Credits:** This package is a Go port of the chart rendering capabilities from the natal library created by Kelvin Ng.

## Features

- Natal chart generation with SVG output
- Multiple theme support (light, dark, monochrome)
- Astrological symbol rendering
- Aspect calculation and visualization
- Synastry and composite chart support
- Transit chart capabilities
- Secondary progressions framework

## Project Structure

### Core Files

- `chart.go` - Main SVG chart generator
- `chart_data.go` - Astrological chart data structures
- `config.go` - Configuration and themes
- `constants.go` - Astrological constants and symbols
- `generator.go` - High-level chart generation functions
- `svg_loader.go` - SVG symbol loader with embedded assets
- `raw_chart_data.go` - Raw data processing and aspect calculations

### Directories

- `svg_paths/` - SVG symbols for planets, signs, aspects, etc.

## Basic Usage

```go
package main

import (
    "astroeph-api/pkg/chart"
    "astroeph-api/models"
)

func main() {
    // Create configuration
    config := chart.DefaultConfig()
    
    // Generate natal chart
    response, err := chart.GenerateNatalChartSVG(natalData, 600, nil, &config)
    if err != nil {
        panic(err)
    }
    
    // SVG output is in response.SVG
    fmt.Println(response.SVG)
}
```

## Chart Types

### Natal Charts
```go
response, err := chart.GenerateNatalChartSVG(natalData, 600, nil, &config)
```

### Synastry Charts (Relationship Compatibility)
```go
response, err := chart.GenerateSynastryChartSVG(person1Chart, person2Chart, 600, nil, &config)
```

### Transit Charts
```go
response, err := chart.GenerateTransitChartSVG(natalChart, currentPositions, 600, nil, &config)
```

### Composite Charts
```go
response, err := chart.GenerateCompositeChartSVG(chart1, chart2, 600, nil, &config)
```

## Configuration

### Available Themes

- `light` - Light theme (default)
- `dark` - Dark theme  
- `mono` - Monochrome theme

### House Systems

- `placidus` - Placidus (default)
- `koch` - Koch
- `equal` - Equal houses
- `campanus` - Campanus
- `regiomontanus` - Regiomontanus
- `porphyry` - Porphyry
- `whole_sign` - Whole sign houses

## Supported Aspects

- Conjunction (0°)
- Opposition (180°)
- Trine (120°)
- Square (90°)
- Sextile (60°)
- Quincunx (150°)

## SVG Structure

### Visual Layers (from outside to inside)

1. **Outer wheel** - Zodiac signs
2. **Second wheel** - Astrological houses with numbers
3. **Third wheel** - House cusps and vertex lines
4. **Fourth wheel** - Planets (outer wheel)
5. **Fifth wheel** - Planets (inner wheel, composite charts only)
6. **Center** - Aspect lines

### Included Symbols

#### Planets
- Classical: Sun, Moon, Mercury, Venus, Mars, Jupiter, Saturn
- Modern: Uranus, Neptune, Pluto
- Other: Lunar Nodes, Chiron
- Asteroids: Ceres, Pallas, Juno, Vesta

#### Zodiac Signs
- Fire: Aries, Leo, Sagittarius
- Earth: Taurus, Virgo, Capricorn
- Air: Gemini, Libra, Aquarius
- Water: Cancer, Scorpio, Pisces

#### Chart Points
- Ascendant (ASC), Midheaven (MC)
- Descendant (DSC), Imum Coeli (IC)

## Customization

### Theme Colors by Element

```go
config := chart.Config{
    ThemeType: chart.ThemeLight,
    Theme: chart.Theme{
        Fire:  "#ef476f", // Fire signs (Aries, Leo, Sagittarius)
        Earth: "#ffd166", // Earth signs (Taurus, Virgo, Capricorn)
        Air:   "#06d6a0", // Air signs (Gemini, Libra, Aquarius)
        Water: "#81bce7", // Water signs (Cancer, Scorpio, Pisces)
    },
}
```

### Aspect Orbs

```go
config.Orbs = chart.Orb{
    Conjunction: 8,
    Opposition:  8,
    Trine:       6,
    Square:      6,
    Sextile:     5,
    Quincunx:    3,
}
```

## Main Functions

### Chart Generation

- `GenerateNatalChartSVG()` - Individual natal chart
- `GenerateSynastryChartSVG()` - Relationship compatibility chart
- `GenerateTransitChartSVG()` - Current planetary positions vs natal
- `GenerateCompositeChartSVG()` - Dual chart display
- `PrepareProgressionData()` - Secondary progressions data preparation

### Utilities

- `DefaultConfig()` - Default configuration
- `GetAvailableThemes()` - Available theme types
- `GetAvailableHouseSystems()` - Available house systems
- `ValidateChartRequest()` - Request validation
- `calculateCompositeAspects()` - Cross-chart aspect calculation

## Technical Considerations

### Performance

- SVG symbols embedded in binary using Go's `embed` package
- Optimized for charts up to 600x600px
- Efficient aspect calculations with configurable orbs
- Anti-overlap algorithms for planet positioning

### Accuracy

- Planetary positions with arc-minute precision
- Configurable aspect orbs
- Automatic degree normalization (0-360°)
- Support for retrograde motion indicators

## Architecture

### Data Flow

1. **Input**: Astrological calculation results (from swephgo)
2. **Processing**: Convert to internal chart data structures
3. **Calculation**: Aspects, house positions, symbol placement
4. **Rendering**: SVG generation with themes and styling
5. **Output**: Complete SVG chart ready for display

### Composite Chart Structure

- **Outer wheel** - Primary chart (person 1 or natal)
- **Inner wheel** - Secondary chart (person 2 or transiting planets)
- **Outer aspects** - Internal aspects of primary chart
- **Inner aspects** - Cross-chart aspects between the two charts

## API Integration

This package integrates seamlessly with the astroeph-api endpoints:

- `/api/v1/natal-chart` - Natal chart generation
- `/api/v1/natal-chart/svg` - SVG-specific natal charts
- Future endpoints for synastry, transits, and progressions

## Acknowledgments

This Go package is a port of the excellent [natal](https://github.com/hoishing/natal) Python library created by Kelvin Ng. The original library provided the foundation for the chart rendering logic, SVG structure, and astrological calculations implemented here.

**Original Library Features Ported:**
- SVG chart rendering engine
- Symbol positioning algorithms
- Theme and color systems
- Aspect calculation methods
- Multi-chart support architecture

## License

This package maintains compatibility with the original natal library's MIT License while being part of the larger astroeph-api project.