package astro

import (
	"astroeph-api/internal/domain"
)

// AspectCalculator handles aspect-related calculations
type AspectCalculator struct {
	// Configuration for aspect calculations
	aspectOrbs     map[domain.AspectType]float64
	planetOrbs     map[string]float64
	enabledAspects map[domain.AspectType]bool
}

// NewAspectCalculator creates a new aspect calculator with default settings
func NewAspectCalculator() *AspectCalculator {
	return &AspectCalculator{
		aspectOrbs:     getDefaultAspectOrbs(),
		planetOrbs:     getDefaultPlanetOrbs(),
		enabledAspects: getDefaultEnabledAspects(),
	}
}

// CalculateAspects calculates all aspects between a list of planets
func (ac *AspectCalculator) CalculateAspects(planets []domain.Planet) []domain.Aspect {
	var aspects []domain.Aspect

	// Compare each planet with every other planet
	for i, planet1 := range planets {
		for j, planet2 := range planets {
			if i >= j {
				continue // Avoid duplicate pairs and self-aspects
			}

			aspect := ac.calculateAspectBetweenPlanets(planet1, planet2)
			if aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}

	return aspects
}

// calculateAspectBetweenPlanets calculates the aspect between two planets
func (ac *AspectCalculator) calculateAspectBetweenPlanets(planet1, planet2 domain.Planet) *domain.Aspect {
	return domain.NewAspect(
		planet1.Name,
		planet2.Name,
		planet1.Longitude,
		planet2.Longitude,
		planet1.Speed,
		planet2.Speed,
	)
}

// CalculateAspectsBetweenCharts calculates aspects between planets from two different charts
func (ac *AspectCalculator) CalculateAspectsBetweenCharts(chart1Planets, chart2Planets []domain.Planet) []domain.Aspect {
	var aspects []domain.Aspect

	for _, planet1 := range chart1Planets {
		for _, planet2 := range chart2Planets {
			aspect := ac.calculateAspectBetweenPlanets(planet1, planet2)
			if aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}

	return aspects
}

// GetDynamicOrb calculates the effective orb between two planets for a specific aspect
func (ac *AspectCalculator) GetDynamicOrb(planet1Name, planet2Name string, aspectType domain.AspectType) float64 {
	baseOrb, exists := ac.aspectOrbs[aspectType]
	if !exists {
		return 0 // Unknown aspect
	}

	// Get planet adjustments
	adj1 := ac.planetOrbs[planet1Name]
	adj2 := ac.planetOrbs[planet2Name]

	// Calculate final orb: base_orb + adjustment_planet1 + adjustment_planet2
	return baseOrb + adj1 + adj2
}

// SetAspectOrb sets the orb for a specific aspect
func (ac *AspectCalculator) SetAspectOrb(aspectType domain.AspectType, orb float64) {
	ac.aspectOrbs[aspectType] = orb
}

// SetPlanetOrb sets the orb adjustment for a specific planet
func (ac *AspectCalculator) SetPlanetOrb(planetName string, orb float64) {
	ac.planetOrbs[planetName] = orb
}

// EnableAspect enables or disables an aspect type
func (ac *AspectCalculator) EnableAspect(aspectType domain.AspectType, enabled bool) {
	ac.enabledAspects[aspectType] = enabled
}

// IsAspectEnabled returns true if the aspect type is enabled
func (ac *AspectCalculator) IsAspectEnabled(aspectType domain.AspectType) bool {
	return ac.enabledAspects[aspectType]
}

// GetAspectOrb returns the base orb for an aspect type
func (ac *AspectCalculator) GetAspectOrb(aspectType domain.AspectType) float64 {
	return ac.aspectOrbs[aspectType]
}

// GetPlanetOrb returns the orb adjustment for a planet
func (ac *AspectCalculator) GetPlanetOrb(planetName string) float64 {
	return ac.planetOrbs[planetName]
}

// GetAspectStrength calculates the strength of an aspect based on its orb
func (ac *AspectCalculator) GetAspectStrength(aspect domain.Aspect) float64 {
	maxOrb := ac.GetDynamicOrb(aspect.Planet1, aspect.Planet2, aspect.Type)
	if maxOrb == 0 {
		return 0
	}

	strength := 1.0 - (aspect.Orb / maxOrb)
	if strength < 0 {
		strength = 0
	}

	return strength
}

// FilterAspectsByStrength filters aspects by minimum strength
func (ac *AspectCalculator) FilterAspectsByStrength(aspects []domain.Aspect, minStrength float64) []domain.Aspect {
	var filtered []domain.Aspect

	for _, aspect := range aspects {
		if ac.GetAspectStrength(aspect) >= minStrength {
			filtered = append(filtered, aspect)
		}
	}

	return filtered
}

// FilterMajorAspects returns only major aspects
func FilterMajorAspects(aspects []domain.Aspect) []domain.Aspect {
	var major []domain.Aspect

	for _, aspect := range aspects {
		if aspect.IsMajorAspect() {
			major = append(major, aspect)
		}
	}

	return major
}

// FilterMinorAspects returns only minor aspects
func FilterMinorAspects(aspects []domain.Aspect) []domain.Aspect {
	var minor []domain.Aspect

	for _, aspect := range aspects {
		if aspect.IsMinorAspect() {
			minor = append(minor, aspect)
		}
	}

	return minor
}

// FilterHarmoniousAspects returns only harmonious aspects
func FilterHarmoniousAspects(aspects []domain.Aspect) []domain.Aspect {
	var harmonious []domain.Aspect

	for _, aspect := range aspects {
		if aspect.IsHarmoniousAspect() {
			harmonious = append(harmonious, aspect)
		}
	}

	return harmonious
}

// FilterChallengingAspects returns only challenging aspects
func FilterChallengingAspects(aspects []domain.Aspect) []domain.Aspect {
	var challenging []domain.Aspect

	for _, aspect := range aspects {
		if aspect.IsChallengingAspect() {
			challenging = append(challenging, aspect)
		}
	}

	return challenging
}

// GroupAspectsByType groups aspects by their type
func GroupAspectsByType(aspects []domain.Aspect) map[domain.AspectType][]domain.Aspect {
	groups := make(map[domain.AspectType][]domain.Aspect)

	for _, aspect := range aspects {
		groups[aspect.Type] = append(groups[aspect.Type], aspect)
	}

	return groups
}

// CountAspectsByType counts aspects by their type
func CountAspectsByType(aspects []domain.Aspect) map[domain.AspectType]int {
	counts := make(map[domain.AspectType]int)

	for _, aspect := range aspects {
		counts[aspect.Type]++
	}

	return counts
}

// getDefaultAspectOrbs returns default orbs for each aspect type
func getDefaultAspectOrbs() map[domain.AspectType]float64 {
	return map[domain.AspectType]float64{
		domain.AspectConjunction:  8.0,
		domain.AspectOpposition:   8.0,
		domain.AspectTrine:        7.0,
		domain.AspectSquare:       6.0,
		domain.AspectSextile:      4.0,
		domain.AspectQuincunx:     2.0,
		domain.AspectSemisextile:  1.0,
		domain.AspectSemisquare:   1.0,
		domain.AspectSesquisquare: 1.0,
	}
}

// getDefaultPlanetOrbs returns default orb adjustments for planets
func getDefaultPlanetOrbs() map[string]float64 {
	return map[string]float64{
		"Sun":        1.0,
		"Moon":       1.0,
		"Mercury":    0.0,
		"Venus":      0.0,
		"Mars":       0.0,
		"Jupiter":    1.0,
		"Saturn":     1.0,
		"Uranus":     2.0,
		"Neptune":    2.0,
		"Pluto":      2.0,
		"North Node": 0.0,
		"Chiron":     0.0,
	}
}

// getDefaultEnabledAspects returns which aspects are enabled by default
func getDefaultEnabledAspects() map[domain.AspectType]bool {
	return map[domain.AspectType]bool{
		domain.AspectConjunction:  true,
		domain.AspectOpposition:   true,
		domain.AspectTrine:        true,
		domain.AspectSquare:       true,
		domain.AspectSextile:      true,
		domain.AspectQuincunx:     true,
		domain.AspectSemisextile:  false, // Minor aspects disabled by default
		domain.AspectSemisquare:   false,
		domain.AspectSesquisquare: false,
	}
}

// CalculateAspectPatterns identifies special aspect patterns
func (ac *AspectCalculator) CalculateAspectPatterns(aspects []domain.Aspect, planets []domain.Planet) []AspectPattern {
	var patterns []AspectPattern

	// Grand Trine pattern
	grandTrines := ac.findGrandTrines(aspects, planets)
	for _, gt := range grandTrines {
		patterns = append(patterns, AspectPattern{
			Type:    "Grand Trine",
			Planets: gt,
			Element: domain.GetElementForSign(domain.GetZodiacSign(planets[0].Longitude)), // Element of first planet
		})
	}

	// T-Square pattern
	tSquares := ac.findTSquares(aspects, planets)
	for _, ts := range tSquares {
		patterns = append(patterns, AspectPattern{
			Type:    "T-Square",
			Planets: ts,
			Element: "", // T-Squares don't have a single element
		})
	}

	// Stellium pattern (3+ planets in same sign)
	stelliums := ac.findStelliums(planets)
	for _, st := range stelliums {
		patterns = append(patterns, AspectPattern{
			Type:    "Stellium",
			Planets: st.Planets,
			Element: domain.GetElementForSign(st.Sign),
		})
	}

	return patterns
}

// AspectPattern represents a special aspect pattern
type AspectPattern struct {
	Type    string   `json:"type"`
	Planets []string `json:"planets"`
	Element string   `json:"element,omitempty"`
}

// Stellium represents a stellium pattern
type Stellium struct {
	Sign    string   `json:"sign"`
	Planets []string `json:"planets"`
}

// findGrandTrines finds Grand Trine patterns
func (ac *AspectCalculator) findGrandTrines(aspects []domain.Aspect, planets []domain.Planet) [][]string {
	// Implementation would find sets of three planets all forming trines with each other
	// This is a complex algorithm that would require more detailed implementation
	return [][]string{} // Placeholder
}

// findTSquares finds T-Square patterns
func (ac *AspectCalculator) findTSquares(aspects []domain.Aspect, planets []domain.Planet) [][]string {
	// Implementation would find patterns where two planets oppose each other and both square a third
	// This is a complex algorithm that would require more detailed implementation
	return [][]string{} // Placeholder
}

// findStelliums finds Stellium patterns (3+ planets in same sign)
func (ac *AspectCalculator) findStelliums(planets []domain.Planet) []Stellium {
	signGroups := make(map[string][]string)

	// Group planets by sign
	for _, planet := range planets {
		sign := planet.Sign
		signGroups[sign] = append(signGroups[sign], planet.Name)
	}

	var stelliums []Stellium
	for sign, planetNames := range signGroups {
		if len(planetNames) >= 3 {
			stelliums = append(stelliums, Stellium{
				Sign:    sign,
				Planets: planetNames,
			})
		}
	}

	return stelliums
}

// CalculateAspectGrid creates a grid showing all aspects between planets
func (ac *AspectCalculator) CalculateAspectGrid(planets []domain.Planet) map[string]map[string]*domain.Aspect {
	grid := make(map[string]map[string]*domain.Aspect)

	// Initialize grid
	for _, planet1 := range planets {
		grid[planet1.Name] = make(map[string]*domain.Aspect)
		for _, planet2 := range planets {
			grid[planet1.Name][planet2.Name] = nil
		}
	}

	// Calculate aspects
	aspects := ac.CalculateAspects(planets)

	// Fill grid
	for _, aspect := range aspects {
		grid[aspect.Planet1][aspect.Planet2] = &aspect
		grid[aspect.Planet2][aspect.Planet1] = &aspect // Symmetric
	}

	return grid
}
