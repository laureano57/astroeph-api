package domain

// HouseSystem represents different house systems
type HouseSystem string

const (
	HousePlacidus      HouseSystem = "Placidus"
	HouseKoch          HouseSystem = "Koch"
	HousePorphyrius    HouseSystem = "Porphyrius"
	HouseRegiomontanus HouseSystem = "Regiomontanus"
	HouseCampanus      HouseSystem = "Campanus"
	HouseEqual         HouseSystem = "Equal"
	HouseWholeSign     HouseSystem = "Whole Sign"
)

// House represents an astrological house
type House struct {
	Number    int     `json:"house"`
	Cusp      string  `json:"cusp"`       // Formatted degree string
	Sign      string  `json:"sign"`       // Zodiac sign on cusp
	CuspValue float64 `json:"cusp_value"` // Raw cusp degree value
	Size      float64 `json:"size"`       // House size in degrees
	Element   string  `json:"element"`    // Element of sign on cusp
	Modality  string  `json:"modality"`   // Modality of sign on cusp
	Ruler     string  `json:"ruler"`      // Traditional ruler of sign on cusp
}

// HouseInfo contains metadata about houses
type HouseInfo struct {
	Number      int
	Name        string
	Keywords    []string
	Description string
	Element     string // Associated element (angular, succedent, cadent)
	Type        string // angular, succedent, cadent
}

// GetHouseInfos returns information about all 12 houses
func GetHouseInfos() []HouseInfo {
	return []HouseInfo{
		{
			Number:      1,
			Name:        "First House",
			Keywords:    []string{"self", "identity", "appearance", "first impressions", "new beginnings"},
			Description: "House of Self and Identity",
			Element:     "fire",
			Type:        "angular",
		},
		{
			Number:      2,
			Name:        "Second House",
			Keywords:    []string{"money", "possessions", "values", "resources", "self-worth"},
			Description: "House of Money and Possessions",
			Element:     "earth",
			Type:        "succedent",
		},
		{
			Number:      3,
			Name:        "Third House",
			Keywords:    []string{"communication", "siblings", "short trips", "learning", "neighborhood"},
			Description: "House of Communication and Learning",
			Element:     "air",
			Type:        "cadent",
		},
		{
			Number:      4,
			Name:        "Fourth House",
			Keywords:    []string{"home", "family", "roots", "foundation", "mother"},
			Description: "House of Home and Family",
			Element:     "water",
			Type:        "angular",
		},
		{
			Number:      5,
			Name:        "Fifth House",
			Keywords:    []string{"creativity", "children", "romance", "fun", "self-expression"},
			Description: "House of Creativity and Romance",
			Element:     "fire",
			Type:        "succedent",
		},
		{
			Number:      6,
			Name:        "Sixth House",
			Keywords:    []string{"health", "work", "service", "daily routine", "pets"},
			Description: "House of Health and Service",
			Element:     "earth",
			Type:        "cadent",
		},
		{
			Number:      7,
			Name:        "Seventh House",
			Keywords:    []string{"partnerships", "marriage", "open enemies", "contracts"},
			Description: "House of Partnerships",
			Element:     "air",
			Type:        "angular",
		},
		{
			Number:      8,
			Name:        "Eighth House",
			Keywords:    []string{"transformation", "death", "other people's money", "occult", "sexuality"},
			Description: "House of Transformation and Shared Resources",
			Element:     "water",
			Type:        "succedent",
		},
		{
			Number:      9,
			Name:        "Ninth House",
			Keywords:    []string{"philosophy", "higher education", "foreign travel", "religion", "law"},
			Description: "House of Philosophy and Higher Learning",
			Element:     "fire",
			Type:        "cadent",
		},
		{
			Number:      10,
			Name:        "Tenth House",
			Keywords:    []string{"career", "reputation", "authority", "father", "public image"},
			Description: "House of Career and Reputation",
			Element:     "earth",
			Type:        "angular",
		},
		{
			Number:      11,
			Name:        "Eleventh House",
			Keywords:    []string{"friends", "groups", "hopes", "wishes", "humanitarian causes"},
			Description: "House of Friends and Aspirations",
			Element:     "air",
			Type:        "succedent",
		},
		{
			Number:      12,
			Name:        "Twelfth House",
			Keywords:    []string{"subconscious", "hidden enemies", "sacrifice", "spirituality", "karma"},
			Description: "House of Subconscious and Spirituality",
			Element:     "water",
			Type:        "cadent",
		},
	}
}

// NewHouse creates a new house with calculated properties
func NewHouse(number int, cuspDegree float64) House {
	sign := GetZodiacSign(cuspDegree)

	house := House{
		Number:    number,
		CuspValue: cuspDegree,
		Cusp:      FormatLongitude(cuspDegree),
		Sign:      sign,
		Element:   GetElementForSign(sign),
		Modality:  GetModalityForSign(sign),
		Ruler:     GetRulerForSign(sign),
	}

	return house
}

// CalculateHouseSizes calculates the size of each house given all cusps
func CalculateHouseSizes(cusps []float64) []float64 {
	if len(cusps) != 12 {
		return make([]float64, 12) // Return zeros if invalid input
	}

	sizes := make([]float64, 12)
	for i := 0; i < 12; i++ {
		nextIndex := (i + 1) % 12
		size := cusps[nextIndex] - cusps[i]

		// Handle wrap-around at 360 degrees
		if size < 0 {
			size += 360
		}

		sizes[i] = size
	}

	return sizes
}

// GetHouseInfo returns information about a specific house number
func GetHouseInfo(number int) *HouseInfo {
	infos := GetHouseInfos()
	if number < 1 || number > 12 {
		return nil
	}
	return &infos[number-1]
}

// IsAngularHouse returns true if the house is angular (1, 4, 7, 10)
func (h House) IsAngularHouse() bool {
	return h.Number == 1 || h.Number == 4 || h.Number == 7 || h.Number == 10
}

// IsSuccedentHouse returns true if the house is succedent (2, 5, 8, 11)
func (h House) IsSuccedentHouse() bool {
	return h.Number == 2 || h.Number == 5 || h.Number == 8 || h.Number == 11
}

// IsCadentHouse returns true if the house is cadent (3, 6, 9, 12)
func (h House) IsCadentHouse() bool {
	return h.Number == 3 || h.Number == 6 || h.Number == 9 || h.Number == 12
}

// GetHouseType returns the type of house (angular, succedent, cadent)
func (h House) GetHouseType() string {
	if h.IsAngularHouse() {
		return "angular"
	}
	if h.IsSuccedentHouse() {
		return "succedent"
	}
	return "cadent"
}

// IsFireHouse returns true if the house has fire association
func (h House) IsFireHouse() bool {
	return h.Number == 1 || h.Number == 5 || h.Number == 9
}

// IsEarthHouse returns true if the house has earth association
func (h House) IsEarthHouse() bool {
	return h.Number == 2 || h.Number == 6 || h.Number == 10
}

// IsAirHouse returns true if the house has air association
func (h House) IsAirHouse() bool {
	return h.Number == 3 || h.Number == 7 || h.Number == 11
}

// IsWaterHouse returns true if the house has water association
func (h House) IsWaterHouse() bool {
	return h.Number == 4 || h.Number == 8 || h.Number == 12
}

// GetHouseElement returns the elemental association of the house
func (h House) GetHouseElement() string {
	if h.IsFireHouse() {
		return "fire"
	}
	if h.IsEarthHouse() {
		return "earth"
	}
	if h.IsAirHouse() {
		return "air"
	}
	return "water"
}

// ContainsPlanet returns true if the given longitude is within this house
func (h House) ContainsPlanet(planetLongitude float64, nextHouseCusp float64) bool {
	// Normalize all values to 0-360 range
	start := normalizeAngle(h.CuspValue)
	end := normalizeAngle(nextHouseCusp)
	planet := normalizeAngle(planetLongitude)

	// Handle wrap-around case
	if start > end {
		return planet >= start || planet < end
	}

	return planet >= start && planet < end
}

// GetHouseKeywords returns keywords associated with this house
func (h House) GetHouseKeywords() []string {
	info := GetHouseInfo(h.Number)
	if info != nil {
		return info.Keywords
	}
	return []string{}
}

// GetHouseDescription returns a description of this house
func (h House) GetHouseDescription() string {
	info := GetHouseInfo(h.Number)
	if info != nil {
		return info.Description
	}
	return ""
}
