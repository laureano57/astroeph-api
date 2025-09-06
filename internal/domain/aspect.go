package domain

import "math"

// AspectType represents the type of aspect
type AspectType string

const (
	AspectConjunction  AspectType = "conjunction"
	AspectSextile      AspectType = "sextile"
	AspectSquare       AspectType = "square"
	AspectTrine        AspectType = "trine"
	AspectOpposition   AspectType = "opposition"
	AspectQuincunx     AspectType = "quincunx"
	AspectSemisextile  AspectType = "semisextile"
	AspectSemisquare   AspectType = "semisquare"
	AspectSesquisquare AspectType = "sesquisquare"
)

// Aspect represents an astrological aspect between two celestial bodies
type Aspect struct {
	Planet1    string     `json:"planet1"`
	Planet2    string     `json:"planet2"`
	Type       AspectType `json:"type"`
	Angle      float64    `json:"angle"`       // Exact angle between planets
	Orb        float64    `json:"orb"`         // Deviation from exact aspect
	IsApplying bool       `json:"is_applying"` // Whether aspect is applying or separating
	IsExact    bool       `json:"is_exact"`    // Whether aspect is exact (orb < 1 degree)
	Strength   float64    `json:"strength"`    // Aspect strength (0-1, based on orb)
	Nature     string     `json:"nature"`      // harmonious, challenging, neutral
}

// AspectDefinition defines the properties of an aspect type
type AspectDefinition struct {
	Type        AspectType
	Angle       float64
	BaseOrb     float64
	Symbol      string
	Nature      string // harmonious, challenging, neutral
	Description string
}

// GetAspectDefinitions returns all standard aspect definitions
func GetAspectDefinitions() []AspectDefinition {
	return []AspectDefinition{
		{
			Type:        AspectConjunction,
			Angle:       0,
			BaseOrb:     8,
			Symbol:      "☌",
			Nature:      "neutral",
			Description: "Unity, blending, intensity",
		},
		{
			Type:        AspectSextile,
			Angle:       60,
			BaseOrb:     4,
			Symbol:      "⚹",
			Nature:      "harmonious",
			Description: "Opportunity, cooperation, ease",
		},
		{
			Type:        AspectSquare,
			Angle:       90,
			BaseOrb:     6,
			Symbol:      "□",
			Nature:      "challenging",
			Description: "Tension, conflict, growth through challenge",
		},
		{
			Type:        AspectTrine,
			Angle:       120,
			BaseOrb:     7,
			Symbol:      "△",
			Nature:      "harmonious",
			Description: "Harmony, natural talent, ease",
		},
		{
			Type:        AspectOpposition,
			Angle:       180,
			BaseOrb:     8,
			Symbol:      "☍",
			Nature:      "challenging",
			Description: "Opposition, balance, projection",
		},
		{
			Type:        AspectQuincunx,
			Angle:       150,
			BaseOrb:     2,
			Symbol:      "⚻",
			Nature:      "neutral",
			Description: "Adjustment, adaptation, minor tension",
		},
		{
			Type:        AspectSemisextile,
			Angle:       30,
			BaseOrb:     1,
			Symbol:      "⚺",
			Nature:      "neutral",
			Description: "Mild connection, subtle influence",
		},
		{
			Type:        AspectSemisquare,
			Angle:       45,
			BaseOrb:     1,
			Symbol:      "∠",
			Nature:      "challenging",
			Description: "Minor tension, irritation",
		},
		{
			Type:        AspectSesquisquare,
			Angle:       135,
			BaseOrb:     1,
			Symbol:      "⚼",
			Nature:      "challenging",
			Description: "Minor tension, adjustment",
		},
	}
}

// GetAspectDefinition returns the definition for a specific aspect type
func GetAspectDefinition(aspectType AspectType) *AspectDefinition {
	definitions := GetAspectDefinitions()
	for _, def := range definitions {
		if def.Type == aspectType {
			return &def
		}
	}
	return nil
}

// NewAspect creates a new aspect between two planets
func NewAspect(planet1, planet2 string, planet1Lon, planet2Lon, planet1Speed, planet2Speed float64) *Aspect {
	// Calculate the angular difference
	angle := AngularDistance(planet1Lon, planet2Lon)

	// Find the best matching aspect
	aspectType, exactAngle, orb := FindBestAspect(angle)
	if aspectType == "" {
		return nil // No aspect found within orb
	}

	// Determine if aspect is applying or separating
	// Applying means faster planet is moving toward exact aspect
	isApplying := IsAspectApplying(planet1Lon, planet2Lon, planet1Speed, planet2Speed, exactAngle)

	// Calculate aspect strength based on orb
	def := GetAspectDefinition(aspectType)
	strength := 1.0
	if def != nil && def.BaseOrb > 0 {
		strength = 1.0 - (orb / def.BaseOrb)
		if strength < 0 {
			strength = 0
		}
	}

	aspect := &Aspect{
		Planet1:    planet1,
		Planet2:    planet2,
		Type:       aspectType,
		Angle:      angle,
		Orb:        orb,
		IsApplying: isApplying,
		IsExact:    orb < 1.0,
		Strength:   strength,
		Nature:     def.Nature,
	}

	return aspect
}

// FindBestAspect finds the best matching aspect for a given angle
func FindBestAspect(angle float64) (AspectType, float64, float64) {
	definitions := GetAspectDefinitions()

	bestAspect := AspectType("")
	bestExactAngle := 0.0
	smallestOrb := 999.0

	for _, def := range definitions {
		orb := math.Abs(angle - def.Angle)
		// Handle the wrap-around case for angles near 0/360
		if orb > 180 {
			orb = 360 - orb
		}

		// Check if this aspect is within orb and better than current best
		if orb <= def.BaseOrb && orb < smallestOrb {
			bestAspect = def.Type
			bestExactAngle = def.Angle
			smallestOrb = orb
		}
	}

	if bestAspect == "" {
		return "", 0, 0
	}

	return bestAspect, bestExactAngle, smallestOrb
}

// IsAspectApplying determines if an aspect is applying (getting closer) or separating
func IsAspectApplying(lon1, lon2, speed1, speed2, exactAngle float64) bool {
	// Calculate the current angular difference
	currentAngle := AngularDistance(lon1, lon2)

	// Calculate what the angle will be in one day
	futureAngle := AngularDistance(lon1+speed1, lon2+speed2)

	// If future angle is closer to exact aspect angle, it's applying
	currentDiff := math.Abs(currentAngle - exactAngle)
	futureDiff := math.Abs(futureAngle - exactAngle)

	return futureDiff < currentDiff
}

// IsHarmoniousAspect returns true if the aspect is generally harmonious
func (a Aspect) IsHarmoniousAspect() bool {
	return a.Nature == "harmonious"
}

// IsChallengingAspect returns true if the aspect is generally challenging
func (a Aspect) IsChallengingAspect() bool {
	return a.Nature == "challenging"
}

// IsNeutralAspect returns true if the aspect is neutral
func (a Aspect) IsNeutralAspect() bool {
	return a.Nature == "neutral"
}

// IsMajorAspect returns true if the aspect is considered a major aspect
func (a Aspect) IsMajorAspect() bool {
	majorAspects := []AspectType{
		AspectConjunction,
		AspectSextile,
		AspectSquare,
		AspectTrine,
		AspectOpposition,
	}

	for _, major := range majorAspects {
		if a.Type == major {
			return true
		}
	}
	return false
}

// IsMinorAspect returns true if the aspect is considered a minor aspect
func (a Aspect) IsMinorAspect() bool {
	return !a.IsMajorAspect()
}

// GetSymbol returns the symbol for the aspect
func (a Aspect) GetSymbol() string {
	def := GetAspectDefinition(a.Type)
	if def != nil {
		return def.Symbol
	}
	return ""
}

// GetDescription returns a description of the aspect
func (a Aspect) GetDescription() string {
	def := GetAspectDefinition(a.Type)
	if def != nil {
		return def.Description
	}
	return ""
}

// CalculateAspects calculates all aspects between a list of planets
func CalculateAspects(planets []Planet) []Aspect {
	var aspects []Aspect

	// Compare each planet with every other planet
	for i, planet1 := range planets {
		for j, planet2 := range planets {
			if i >= j {
				continue // Avoid duplicate pairs and self-aspects
			}

			aspect := NewAspect(
				planet1.Name,
				planet2.Name,
				planet1.Longitude,
				planet2.Longitude,
				planet1.Speed,
				planet2.Speed,
			)

			if aspect != nil {
				aspects = append(aspects, *aspect)
			}
		}
	}

	return aspects
}
