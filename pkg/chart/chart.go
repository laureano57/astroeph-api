package chart

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Chart represents an SVG natal chart generator
type Chart struct {
	Data1         *ChartData
	Data2         *ChartData // For composite charts
	Width         int
	Height        int
	CX            float64 // Center X
	CY            float64 // Center Y
	Config        Config
	MaxRadius     float64
	Margin        float64
	RingThickness float64
	FontSize      float64
	ScaleAdj      float64
	PosAdj        float64
}

// NewChart creates a new Chart instance
func NewChart(data1 *ChartData, width int, height *int, data2 *ChartData) *Chart {
	chart := &Chart{
		Data1:  data1,
		Data2:  data2,
		Width:  width,
		Config: data1.Config,
	}

	// Set height to width if not specified
	if height == nil {
		chart.Height = width
	} else {
		chart.Height = *height
	}

	// Calculate derived properties
	chart.CX = float64(chart.Width) / 2
	chart.CY = float64(chart.Height) / 2

	margin := math.Min(float64(chart.Width), float64(chart.Height)) * chart.Config.Chart.MarginFactor
	chart.MaxRadius = math.Min(float64(chart.Width)-margin, float64(chart.Height)-margin) / 2
	chart.Margin = margin
	chart.RingThickness = chart.MaxRadius * chart.Config.Chart.RingThicknessFraction
	chart.FontSize = chart.RingThickness * chart.Config.Chart.FontSizeFraction
	chart.ScaleAdj = float64(chart.Width) / chart.Config.Chart.ScaleAdjFactor
	chart.PosAdj = chart.FontSize / chart.Config.Chart.PosAdjFactor

	// Fix scale for symbol visibility - use a fixed scale that works well
	chart.ScaleAdj = 2.5

	return chart
}

// GenerateSVG generates the complete SVG representation of the chart
func (c *Chart) GenerateSVG() string {
	var elements []string

	// Generate all chart components
	elements = append(elements, c.generateSignWheel()...)
	elements = append(elements, c.generateHouseWheel()...)
	elements = append(elements, c.generateVertexWheel()...)
	elements = append(elements, c.generateSignWheelSymbols()...)
	elements = append(elements, c.generateOuterBodyWheel()...)

	if c.Data2 != nil {
		elements = append(elements, c.generateInnerBodyWheel()...)
	}

	elements = append(elements, c.generateOuterAspects()...)

	if c.Data2 != nil {
		elements = append(elements, c.generateInnerAspects()...)
	}

	// Wrap in SVG root element
	return c.svgRoot(elements)
}

// svgRoot generates the SVG root element
func (c *Chart) svgRoot(content []string) string {
	return fmt.Sprintf(`<svg height="%d" width="%d" font-family="%s" version="1.1" xmlns="http://www.w3.org/2000/svg">
%s
</svg>`, c.Height, c.Width, c.Config.Chart.Font, strings.Join(content, "\n"))
}

// generateSignWheel generates the zodiac sign wheel
func (c *Chart) generateSignWheel() []string {
	radius := c.MaxRadius
	theme := c.Config.GetTheme()
	var elements []string

	// Background circle
	elements = append(elements, c.backgroundCircle(radius, theme.Background))

	// Generate sectors for each sign
	for i := 0; i < 12; i++ {
		startDeg := c.Data1.Signs[i].NormalizedDegree
		endDeg := startDeg + 30

		bgColor := c.getBackgroundColor(i, theme)

		elements = append(elements, c.sector(
			radius,
			startDeg,
			endDeg,
			bgColor,
			theme.Foreground,
			float64(c.Config.Chart.StrokeWidth),
			c.Config.Chart.StrokeOpacity,
		))
	}

	return elements
}

// generateSignWheelSymbols generates zodiac sign symbols
func (c *Chart) generateSignWheelSymbols() []string {
	var elements []string
	theme := c.Config.GetTheme()

	for i := 0; i < 12; i++ {
		startDeg := c.Data1.Signs[i].NormalizedDegree
		symbolRadius := c.MaxRadius - (c.RingThickness / 2)
		symbolAngle := math.Pi * (startDeg + 15) / 180 // Center of sector
		symbolX := c.CX - symbolRadius*math.Cos(symbolAngle) - c.PosAdj
		symbolY := c.CY + symbolRadius*math.Sin(symbolAngle) - c.PosAdj

		signMember := SIGN_MEMBERS[i]
		strokeColor := c.getColorForElement(signMember.Color, theme)

		// Get SVG path for sign symbol
		svgPath := GetSVGPath(signMember.Name)
		if svgPath == "" {
			continue
		}

		// Scale symbols to be similar size as house numbers
		symbolScale := (c.FontSize * 0.8) / 20.0 // 20 is the original SVG size
		elements = append(elements, fmt.Sprintf(
			`<g stroke="%s" stroke-width="%.1f" fill="none" transform="translate(%.1f, %.1f) scale(%.3f)">%s</g>`,
			strokeColor,
			float64(c.Config.Chart.StrokeWidth)*1.5,
			symbolX,
			symbolY,
			symbolScale,
			svgPath,
		))
	}

	return elements
}

// generateHouseWheel generates the house wheel
func (c *Chart) generateHouseWheel() []string {
	radius := c.MaxRadius - c.RingThickness
	theme := c.Config.GetTheme()
	var elements []string

	// Background circle
	elements = append(elements, c.backgroundCircle(radius, theme.Background))

	// Generate house sectors and numbers
	houseVertices := c.getHouseVertices()
	for i, vertex := range houseVertices {
		startDeg, endDeg := vertex[0], vertex[1]
		bgColor := c.getBackgroundColor(i, theme)

		// House sector
		elements = append(elements, c.sector(
			radius,
			startDeg,
			endDeg,
			bgColor,
			theme.Foreground,
			float64(c.Config.Chart.StrokeWidth),
			c.Config.Chart.StrokeOpacity,
		))

		// House number
		numberRadius := radius - (c.RingThickness / 2)
		numberAngle := math.Pi * (startDeg + (math.Mod(endDeg-startDeg, 360) / 2)) / 180
		numberX := c.CX - numberRadius*math.Cos(numberAngle)
		numberY := c.CY + numberRadius*math.Sin(numberAngle)

		numberColor := c.getColorForElement(SIGN_MEMBERS[i].Color, theme)
		fontSize := c.FontSize * 0.8

		elements = append(elements, fmt.Sprintf(
			`<text x="%.1f" y="%.1f" fill="%s" font-size="%.1f" text-anchor="middle" dominant-baseline="central">%d</text>`,
			numberX, numberY, numberColor, fontSize, i+1,
		))
	}

	return elements
}

// generateVertexWheel generates vertex lines
func (c *Chart) generateVertexWheel() []string {
	var elements []string
	theme := c.Config.GetTheme()

	vertexRadius := c.MaxRadius + c.Margin/2
	houseRadius := c.MaxRadius - 2*c.RingThickness
	bodyRadius := c.MaxRadius - 3*c.RingThickness

	// Background circles
	elements = append(elements, fmt.Sprintf(
		`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="%s" stroke-width="%d"/>`,
		c.CX, c.CY, houseRadius, theme.Background, theme.Foreground, c.Config.Chart.StrokeWidth,
	))

	elements = append(elements, fmt.Sprintf(
		`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#88888800" stroke="%s" stroke-width="%d"/>`,
		c.CX, c.CY, bodyRadius, theme.Dim, c.Config.Chart.StrokeWidth,
	))

	// Generate vertex lines
	for _, house := range c.Data1.Houses {
		radius := houseRadius
		strokeWidth := float64(c.Config.Chart.StrokeWidth)
		strokeColor := theme.Dim

		// Highlight major angles (1st, 4th, 7th, 10th houses)
		if house.Value == 1 || house.Value == 4 || house.Value == 7 || house.Value == 10 {
			radius = vertexRadius
			strokeColor = theme.Foreground
		}

		angle := math.Pi * house.NormalizedDegree / 180
		endX := c.CX - radius*math.Cos(angle)
		endY := c.CY + radius*math.Sin(angle)

		elements = append(elements, fmt.Sprintf(
			`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f" stroke-opacity="%.2f"/>`,
			c.CX, c.CY, endX, endY, strokeColor, strokeWidth, c.Config.Chart.StrokeOpacity,
		))
	}

	return elements
}

// generateOuterBodyWheel generates the outer body wheel
func (c *Chart) generateOuterBodyWheel() []string {
	radius := c.MaxRadius - 3*c.RingThickness
	data := c.Data2
	if data == nil {
		data = c.Data1
	}
	return c.generateBodyWheel(radius, data, c.Config.Chart.OuterMinDegree)
}

// generateInnerBodyWheel generates the inner body wheel for composite charts
func (c *Chart) generateInnerBodyWheel() []string {
	if c.Data2 == nil {
		return []string{}
	}
	radius := c.MaxRadius - 4*c.RingThickness
	return c.generateBodyWheel(radius, c.Data1, c.Config.Chart.InnerMinDegree)
}

// generateBodyWheel generates elements for a body wheel
func (c *Chart) generateBodyWheel(wheelRadius float64, data *ChartData, minDegree float64) []string {
	var elements []string
	theme := c.Config.GetTheme()

	// Sort bodies by normalized degree
	bodies := make([]MovableBody, len(data.Aspectables))
	copy(bodies, data.Aspectables)
	sort.Slice(bodies, func(i, j int) bool {
		return bodies[i].NormalizedDegree < bodies[j].NormalizedDegree
	})

	// Get normalized degrees
	degrees := make([]float64, len(bodies))
	for i, body := range bodies {
		degrees[i] = body.NormalizedDegree
	}

	// Adjust positions to avoid overlap
	adjustedDegrees := degrees
	if len(bodies) > 1 {
		adjustedDegrees = c.adjustedDegrees(degrees, minDegree)
	}

	// Generate body symbols and lines
	for i, body := range bodies {
		adjDeg := adjustedDegrees[i]

		strokeColor := c.getColorForElement(body.Color, theme)
		symbolRadius := wheelRadius + (c.RingThickness / 2)

		// Use original angle for line start position
		originalAngle := math.Pi * body.NormalizedDegree / 180
		degreeX := c.CX - wheelRadius*math.Cos(originalAngle)
		degreeY := c.CY + wheelRadius*math.Sin(originalAngle)

		// Use adjusted angle for symbol position
		adjustedAngle := math.Pi * adjDeg / 180
		symbolX := c.CX - symbolRadius*math.Cos(adjustedAngle)
		symbolY := c.CY + symbolRadius*math.Sin(adjustedAngle)

		// Inner line
		innerRadius := wheelRadius - c.RingThickness
		innerX := c.CX - innerRadius*math.Cos(originalAngle)
		innerY := c.CY + innerRadius*math.Sin(originalAngle)

		// Add connecting line
		elements = append(elements, fmt.Sprintf(
			`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f"/>`,
			degreeX, degreeY, symbolX, symbolY, strokeColor, float64(c.Config.Chart.StrokeWidth)/2,
		))

		// Add symbol background circle
		elements = append(elements, fmt.Sprintf(
			`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`,
			symbolX, symbolY, c.FontSize/2, theme.Background,
		))

		// Add dashed inner line
		elements = append(elements, fmt.Sprintf(
			`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f" stroke-dasharray="%.1f"/>`,
			degreeX, degreeY, innerX, innerY, theme.Dim, float64(c.Config.Chart.StrokeWidth)/2, c.RingThickness/11,
		))

		// Add symbol
		svgPath := GetSVGPath(body.Name)
		if svgPath != "" {
			// Default: use stroke for planets (like original implementation)
			fillAttr := "none"
			strokeAttr := strokeColor

			// Special handling for vertices (asc, ic, dsc, mc): use fill instead of stroke
			if body.Name == "asc" || body.Name == "ic" || body.Name == "dsc" || body.Name == "mc" {
				fillAttr = strokeColor
				strokeAttr = "none"
			}

			// Use same scale as sign symbols for consistency
			symbolScale := (c.FontSize * 0.8) / 20.0

			elements = append(elements, fmt.Sprintf(
				`<g fill="%s" stroke="%s" stroke-width="%.1f" transform="translate(%.1f, %.1f) scale(%.3f)">%s</g>`,
				fillAttr, strokeAttr, float64(c.Config.Chart.StrokeWidth)*1.5, symbolX-c.PosAdj, symbolY-c.PosAdj, symbolScale, svgPath,
			))
		}
	}

	return elements
}

// generateOuterAspects generates aspect lines for single charts
func (c *Chart) generateOuterAspects() []string {
	if c.Data2 != nil {
		return []string{}
	}
	aspectRadius := c.MaxRadius - 3*c.RingThickness
	return c.generateAspectLines(aspectRadius, c.Data1.Aspects)
}

// generateInnerAspects generates aspect lines for composite charts
func (c *Chart) generateInnerAspects() []string {
	if c.Data2 == nil {
		return []string{}
	}
	aspectRadius := c.MaxRadius - 4*c.RingThickness

	// Calculate composite aspects between the two charts
	compositeAspects := calculateCompositeAspects(c.Data1.Aspectables, c.Data2.Aspectables, c.Config)

	return c.generateAspectLines(aspectRadius, compositeAspects)
}

// generateAspectLines generates aspect lines between bodies
func (c *Chart) generateAspectLines(radius float64, aspects []Aspect) []string {
	var elements []string
	theme := c.Config.GetTheme()

	// Background circle
	elements = append(elements, fmt.Sprintf(
		`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s" stroke="%s" stroke-width="%d"/>`,
		c.CX, c.CY, radius, theme.Background, theme.Dim, c.Config.Chart.StrokeWidth,
	))

	// Generate aspect lines
	for _, aspect := range aspects {
		startAngle := math.Pi * aspect.Body1.NormalizedDegree / 180
		endAngle := math.Pi * aspect.Body2.NormalizedDegree / 180

		orbConfig := c.Config.GetOrbForAspect(aspect.AspectMember.Name)
		if orbConfig == 0 {
			continue
		}

		orb := 1.0
		if aspect.Orb != nil {
			orb = *aspect.Orb
		}

		orbFraction := 1 - orb/float64(orbConfig)
		opacityFactor := orbFraction
		if aspect.AspectMember.Name == "conjunction" {
			opacityFactor = 1.0
		}

		strokeColor := c.getColorForElement(aspect.AspectMember.Color, theme)

		elements = append(elements, fmt.Sprintf(
			`<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f" stroke-opacity="%.3f"/>`,
			c.CX-radius*math.Cos(startAngle),
			c.CY+radius*math.Sin(startAngle),
			c.CX-radius*math.Cos(endAngle),
			c.CY+radius*math.Sin(endAngle),
			strokeColor,
			float64(c.Config.Chart.StrokeWidth)/2,
			c.Config.Chart.StrokeOpacity*opacityFactor,
		))
	}

	return elements
}

// Helper methods

// sector creates an SVG sector (pie slice)
func (c *Chart) sector(radius, startDeg, endDeg float64, fill, stroke string, strokeWidth, strokeOpacity float64) string {
	startRad := math.Pi * startDeg / 180
	endRad := math.Pi * endDeg / 180
	startX := c.CX - radius*math.Cos(startRad)
	startY := c.CY + radius*math.Sin(startRad)
	endX := c.CX - radius*math.Cos(endRad)
	endY := c.CY + radius*math.Sin(endRad)

	// Round coordinates
	startX = math.Round(startX*100) / 100
	startY = math.Round(startY*100) / 100
	endX = math.Round(endX*100) / 100
	endY = math.Round(endY*100) / 100

	pathData := fmt.Sprintf("M%.2f %.2f L%.2f %.2f A%.2f %.2f 0 0 0 %.2f %.2f Z",
		c.CX, c.CY, startX, startY, radius, radius, endX, endY)

	return fmt.Sprintf(
		`<path d="%s" fill="%s" fill-opacity="0.3" stroke="%s" stroke-width="%.1f" stroke-opacity="%.2f"/>`,
		pathData, fill, stroke, strokeWidth, strokeOpacity,
	)
}

// backgroundCircle creates a background circle
func (c *Chart) backgroundCircle(radius float64, fill string) string {
	return fmt.Sprintf(`<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`, c.CX, c.CY, radius, fill)
}

// getHouseVertices calculates house vertices (start and end degrees)
func (c *Chart) getHouseVertices() [][2]float64 {
	vertices := make([][2]float64, 12)

	for i := 0; i < 12; i++ {
		nextI := (i + 1) % 12
		startDeg := c.Data1.Houses[i].NormalizedDegree
		endDeg := c.Data1.Houses[nextI].NormalizedDegree

		// Handle wrap-around case
		if endDeg < startDeg {
			endDeg += 360
		}

		vertices[i] = [2]float64{startDeg, endDeg}
	}

	return vertices
}

// getBackgroundColor returns background color for a given index
func (c *Chart) getBackgroundColor(index int, theme Theme) string {
	transparency := theme.Transparency
	signMember := SIGN_MEMBERS[index%4] // Cycle through first 4 signs for color pattern

	baseColor := c.getColorForElement(signMember.Color, theme)
	// For simplicity, return the base color - in a full implementation you'd blend with background
	return c.blendColorWithBackground(baseColor, theme.Background, transparency)
}

// getColorForElement returns the color for a given element
func (c *Chart) getColorForElement(element string, theme Theme) string {
	switch element {
	case "fire":
		return theme.Fire
	case "earth":
		return theme.Earth
	case "air":
		return theme.Air
	case "water":
		return theme.Water
	case "points":
		return theme.Points
	case "asteroids":
		return theme.Asteroids
	case "positive":
		return theme.Positive
	case "negative":
		return theme.Negative
	case "others":
		return theme.Others
	default:
		return theme.Foreground
	}
}

// blendColorWithBackground blends a color with background (simplified)
func (c *Chart) blendColorWithBackground(color, background string, transparency float64) string {
	// Simplified - just return the color for now
	// In a full implementation, you'd do proper color blending
	return color
}

// adjustedDegrees adjusts spacing between bodies to avoid overlap
func (c *Chart) adjustedDegrees(degrees []float64, minDegree float64) []float64 {
	step := minDegree + 0.1
	n := len(degrees)

	fwdDegs := make([]float64, n)
	copy(fwdDegs, degrees)

	bwdDegs := make([]float64, n)
	copy(bwdDegs, degrees)
	// Reverse the backward array
	for i := 0; i < n/2; i++ {
		bwdDegs[i], bwdDegs[n-1-i] = bwdDegs[n-1-i], bwdDegs[i]
	}

	// Forward adjustment
	changed := true
	for changed {
		changed = false
		for i := 0; i < n; i++ {
			var prevDeg float64
			if i == 0 {
				prevDeg = fwdDegs[n-1] - 360
			} else {
				prevDeg = fwdDegs[i-1]
			}

			delta := fwdDegs[i] - prevDeg
			diff := math.Min(delta, 360-delta)

			if fwdDegs[i] < prevDeg || diff < minDegree {
				fwdDegs[i] = prevDeg + step
				changed = true
			}
		}
	}

	// Backward adjustment
	changed = true
	for changed {
		changed = false
		for i := 0; i < n; i++ {
			var prevDeg float64
			if i == 0 {
				prevDeg = bwdDegs[n-1] + 360
			} else {
				prevDeg = bwdDegs[i-1]
			}

			delta := prevDeg - bwdDegs[i]
			diff := math.Min(delta, 360-delta)

			if prevDeg < bwdDegs[i] || diff < minDegree {
				bwdDegs[i] = prevDeg - step
				changed = true
			}
		}
	}

	// Reverse backward array back
	for i := 0; i < n/2; i++ {
		bwdDegs[i], bwdDegs[n-1-i] = bwdDegs[n-1-i], bwdDegs[i]
	}

	// Average forward and backward adjustments
	avgAdj := make([]float64, n)
	for i := 0; i < n; i++ {
		fwd := math.Mod(fwdDegs[i], 360)
		bwd := math.Mod(bwdDegs[i], 360)

		if math.Abs(fwd-bwd) < 180 {
			avgAdj[i] = (fwd + bwd) / 2
		} else {
			avgAdj[i] = math.Mod((fwd+bwd+360)/2, 360)
		}
	}

	return avgAdj
}
