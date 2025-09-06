package domain

// PlanetName represents planetary names as constants
type PlanetName string

const (
	Sun       PlanetName = "Sun"
	Moon      PlanetName = "Moon"
	Mercury   PlanetName = "Mercury"
	Venus     PlanetName = "Venus"
	Mars      PlanetName = "Mars"
	Jupiter   PlanetName = "Jupiter"
	Saturn    PlanetName = "Saturn"
	Uranus    PlanetName = "Uranus"
	Neptune   PlanetName = "Neptune"
	Pluto     PlanetName = "Pluto"
	NorthNode PlanetName = "North Node"
	SouthNode PlanetName = "South Node"
	Chiron    PlanetName = "Chiron"
	Ceres     PlanetName = "Ceres"
	Pallas    PlanetName = "Pallas"
	Juno      PlanetName = "Juno"
	Vesta     PlanetName = "Vesta"
)

// Planet represents a celestial body in an astrological chart
type Planet struct {
	Name         string  `json:"name"`
	Sign         string  `json:"sign"`
	Degree       string  `json:"degree"`
	House        int     `json:"house"`
	Longitude    float64 `json:"longitude"` // Raw longitude in degrees
	Latitude     float64 `json:"latitude"`  // Raw latitude in degrees (for some calculations)
	Speed        float64 `json:"speed"`     // Daily motion in degrees
	IsRetrograde bool    `json:"is_retrograde"`
	Element      string  `json:"element"`     // fire, earth, air, water
	Modality     string  `json:"modality"`    // cardinal, fixed, mutable
	PlanetType   string  `json:"planet_type"` // personal, social, transpersonal, etc.
}

// PlanetType represents the classification of planets
type PlanetType string

const (
	TypePersonal      PlanetType = "personal"      // Sun, Moon, Mercury, Venus, Mars
	TypeSocial        PlanetType = "social"        // Jupiter, Saturn
	TypeTranspersonal PlanetType = "transpersonal" // Uranus, Neptune, Pluto
	TypeLunar         PlanetType = "lunar"         // Nodes
	TypeAsteroid      PlanetType = "asteroid"      // Chiron, Ceres, etc.
)

// NewPlanet creates a new planet with calculated properties
func NewPlanet(name string, longitude, latitude, speed float64, houseNumber int) Planet {
	planet := Planet{
		Name:         name,
		Longitude:    longitude,
		Latitude:     latitude,
		Speed:        speed,
		House:        houseNumber,
		Sign:         GetZodiacSign(longitude),
		Degree:       FormatDegreeInSign(longitude),
		IsRetrograde: speed < 0,
		Element:      GetElementForSign(GetZodiacSign(longitude)),
		Modality:     GetModalityForSign(GetZodiacSign(longitude)),
		PlanetType:   string(GetPlanetType(name)),
	}

	return planet
}

// GetPlanetType returns the type classification for a planet
func GetPlanetType(planetName string) PlanetType {
	switch planetName {
	case string(Sun), string(Moon), string(Mercury), string(Venus), string(Mars):
		return TypePersonal
	case string(Jupiter), string(Saturn):
		return TypeSocial
	case string(Uranus), string(Neptune), string(Pluto):
		return TypeTranspersonal
	case string(NorthNode), string(SouthNode):
		return TypeLunar
	case string(Chiron), string(Ceres), string(Pallas), string(Juno), string(Vesta):
		return TypeAsteroid
	default:
		return TypePersonal
	}
}

// IsPersonalPlanet returns true if the planet is a personal planet
func (p Planet) IsPersonalPlanet() bool {
	return GetPlanetType(p.Name) == TypePersonal
}

// IsSocialPlanet returns true if the planet is a social planet
func (p Planet) IsSocialPlanet() bool {
	return GetPlanetType(p.Name) == TypeSocial
}

// IsTranspersonalPlanet returns true if the planet is a transpersonal planet
func (p Planet) IsTranspersonalPlanet() bool {
	return GetPlanetType(p.Name) == TypeTranspersonal
}

// IsRetrogradePlanet returns true if the planet is moving retrograde
func (p Planet) IsRetrogradePlanet() bool {
	return p.IsRetrograde
}

// GetDegreeInSign returns the degree within the zodiac sign (0-29.999...)
func (p Planet) GetDegreeInSign() float64 {
	return GetDegreeInSign(p.Longitude)
}

// GetSignNumber returns the sign number (0-11 for Aries-Pisces)
func (p Planet) GetSignNumber() int {
	return int(p.Longitude / 30.0)
}

// IsCombust returns true if the planet is combust (too close to the Sun)
func (p Planet) IsCombust(sunLongitude float64) bool {
	if p.Name == string(Sun) {
		return false
	}

	// Calculate angular distance to Sun
	distance := AngularDistance(p.Longitude, sunLongitude)

	// Combustion orbs vary by planet
	var combustOrb float64
	switch p.Name {
	case string(Moon):
		combustOrb = 12.0 // Moon has wider orb
	case string(Mercury):
		combustOrb = 8.5
	case string(Venus):
		combustOrb = 8.0
	case string(Mars):
		combustOrb = 17.0
	case string(Jupiter):
		combustOrb = 11.0
	case string(Saturn):
		combustOrb = 15.0
	default:
		combustOrb = 8.0
	}

	return distance <= combustOrb
}

// IsInDetriment returns true if the planet is in its detriment sign
func (p Planet) IsInDetriment() bool {
	detrimentSigns := map[string][]string{
		string(Sun):     {"Aquarius"},
		string(Moon):    {"Capricorn"},
		string(Mercury): {"Sagittarius", "Pisces"},
		string(Venus):   {"Aries", "Scorpio"},
		string(Mars):    {"Libra", "Taurus"},
		string(Jupiter): {"Gemini", "Virgo"},
		string(Saturn):  {"Cancer", "Leo"},
		string(Uranus):  {"Leo"},
		string(Neptune): {"Virgo"},
		string(Pluto):   {"Taurus"},
	}

	if signs, exists := detrimentSigns[p.Name]; exists {
		for _, sign := range signs {
			if p.Sign == sign {
				return true
			}
		}
	}
	return false
}

// IsInExaltation returns true if the planet is in its exaltation sign
func (p Planet) IsInExaltation() bool {
	exaltationSigns := map[string]string{
		string(Sun):     "Aries",
		string(Moon):    "Taurus",
		string(Mercury): "Virgo",
		string(Venus):   "Pisces",
		string(Mars):    "Capricorn",
		string(Jupiter): "Cancer",
		string(Saturn):  "Libra",
		string(Uranus):  "Scorpio",
		string(Neptune): "Aquarius",
		string(Pluto):   "Aries",
	}

	if sign, exists := exaltationSigns[p.Name]; exists {
		return p.Sign == sign
	}
	return false
}

// IsInFall returns true if the planet is in its fall sign
func (p Planet) IsInFall() bool {
	fallSigns := map[string]string{
		string(Sun):     "Libra",
		string(Moon):    "Scorpio",
		string(Mercury): "Pisces",
		string(Venus):   "Virgo",
		string(Mars):    "Cancer",
		string(Jupiter): "Capricorn",
		string(Saturn):  "Aries",
		string(Uranus):  "Taurus",
		string(Neptune): "Leo",
		string(Pluto):   "Libra",
	}

	if sign, exists := fallSigns[p.Name]; exists {
		return p.Sign == sign
	}
	return false
}
