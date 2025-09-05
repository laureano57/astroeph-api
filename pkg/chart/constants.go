package chart

// Planet constants matching swephgo values
const (
	SUN      = 0
	MOON     = 1
	MERCURY  = 2
	VENUS    = 3
	MARS     = 4
	JUPITER  = 5
	SATURN   = 6
	URANUS   = 7
	NEPTUNE  = 8
	PLUTO    = 9
	ASC_NODE = 10
	CHIRON   = 15
)

// Vertex constants
const (
	ASC = "asc"
	IC  = "ic"
	DSC = "dsc"
	MC  = "mc"
)

// Aspect types
const (
	CONJUNCTION = "conjunction"
	OPPOSITION  = "opposition"
	TRINE       = "trine"
	SQUARE      = "square"
	SEXTILE     = "sextile"
	QUINCUNX    = "quincunx"
)

// Element and modality names
var (
	PLANET_NAMES   = []string{"sun", "moon", "mercury", "venus", "mars", "jupiter", "saturn", "uranus", "neptune", "pluto", "asc_node"}
	EXTRA_NAMES    = []string{"chiron", "ceres", "pallas", "juno", "vesta"}
	ELEMENT_NAMES  = []string{"fire", "earth", "air", "water"}
	MODALITY_NAMES = []string{"cardinal", "fixed", "mutable"}
	POLARITY_NAMES = []string{"positive", "negative"}
	SIGN_NAMES     = []string{"aries", "taurus", "gemini", "cancer", "leo", "virgo", "libra", "scorpio", "sagittarius", "capricorn", "aquarius", "pisces"}
	HOUSE_NAMES    = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve"}
	ASPECT_NAMES   = []string{"conjunction", "opposition", "trine", "square", "sextile", "quincunx"}
	VERTEX_NAMES   = []string{"asc", "ic", "dsc", "mc"}
)

// Body represents a celestial body with basic properties
type Body struct {
	Name   string
	Symbol string
	Value  int
	Color  string
}

// SignMember represents a zodiac sign with all its properties
type SignMember struct {
	Body
	Ruler            string
	Detriment        string
	Exaltation       string
	Fall             string
	ClassicRuler     string
	ClassicDetriment string
	Modality         string
	Element          string
	Polarity         string
}

// AspectMember represents an aspect type
type AspectMember struct {
	Body
}

// ElementMember represents an element
type ElementMember struct {
	Body
}

// ModalityMember represents a modality
type ModalityMember struct {
	Body
}

// PolarityMember represents polarity
type PolarityMember struct {
	Body
}

// HouseMember represents a house
type HouseMember struct {
	Body
}

// ExtraMember represents extra celestial bodies (asteroids, etc.)
type ExtraMember struct {
	Body
}

// VertexMember represents chart vertices (ASC, IC, DSC, MC)
type VertexMember struct {
	Body
}

// PLANET_MEMBERS contains all planet definitions
var PLANET_MEMBERS = []Body{
	{Name: "sun", Symbol: "☉", Value: 0, Color: "fire"},
	{Name: "moon", Symbol: "☽", Value: 1, Color: "water"},
	{Name: "mercury", Symbol: "☿", Value: 2, Color: "air"},
	{Name: "venus", Symbol: "♀", Value: 3, Color: "earth"},
	{Name: "mars", Symbol: "♂", Value: 4, Color: "fire"},
	{Name: "jupiter", Symbol: "♃", Value: 5, Color: "fire"},
	{Name: "saturn", Symbol: "♄", Value: 6, Color: "earth"},
	{Name: "uranus", Symbol: "♅", Value: 7, Color: "air"},
	{Name: "neptune", Symbol: "♆", Value: 8, Color: "water"},
	{Name: "pluto", Symbol: "♇", Value: 9, Color: "water"},
	{Name: "asc_node", Symbol: "☊", Value: 10, Color: "points"},
}

// ASPECT_MEMBERS contains all aspect definitions
var ASPECT_MEMBERS = []AspectMember{
	{Body: Body{Name: "conjunction", Symbol: "☌", Value: 0, Color: "others"}},
	{Body: Body{Name: "opposition", Symbol: "☍", Value: 180, Color: "water"}},
	{Body: Body{Name: "trine", Symbol: "△", Value: 120, Color: "air"}},
	{Body: Body{Name: "square", Symbol: "□", Value: 90, Color: "fire"}},
	{Body: Body{Name: "sextile", Symbol: "⚹", Value: 60, Color: "points"}},
	{Body: Body{Name: "quincunx", Symbol: "⚻", Value: 150, Color: "asteroids"}},
}

// ELEMENT_MEMBERS contains all element definitions
var ELEMENT_MEMBERS = []ElementMember{
	{Body: Body{Name: "fire", Symbol: "🜂", Value: 0, Color: "fire"}},
	{Body: Body{Name: "earth", Symbol: "🜃", Value: 1, Color: "earth"}},
	{Body: Body{Name: "air", Symbol: "🜁", Value: 2, Color: "air"}},
	{Body: Body{Name: "water", Symbol: "🜄", Value: 3, Color: "water"}},
}

// MODALITY_MEMBERS contains all modality definitions
var MODALITY_MEMBERS = []ModalityMember{
	{Body: Body{Name: "cardinal", Symbol: "⟑", Value: 0, Color: "fire"}},
	{Body: Body{Name: "fixed", Symbol: "⊟", Value: 1, Color: "earth"}},
	{Body: Body{Name: "mutable", Symbol: "𛰣", Value: 2, Color: "air"}},
}

// POLARITY_MEMBERS contains polarity definitions
var POLARITY_MEMBERS = []PolarityMember{
	{Body: Body{Name: "positive", Symbol: "+", Value: 1, Color: "positive"}},
	{Body: Body{Name: "negative", Symbol: "-", Value: -1, Color: "negative"}},
}

// SIGN_MEMBERS contains all zodiac sign definitions
var SIGN_MEMBERS = []SignMember{
	{Body: Body{Name: "aries", Symbol: "♈", Value: 1, Color: "fire"}, Ruler: "mars", Detriment: "venus", Exaltation: "sun", Fall: "saturn", ClassicRuler: "mars", ClassicDetriment: "venus", Modality: "cardinal", Element: "fire", Polarity: "positive"},
	{Body: Body{Name: "taurus", Symbol: "♉", Value: 2, Color: "earth"}, Ruler: "venus", Detriment: "pluto", Exaltation: "moon", Fall: "", ClassicRuler: "venus", ClassicDetriment: "mars", Modality: "fixed", Element: "earth", Polarity: "negative"},
	{Body: Body{Name: "gemini", Symbol: "♊", Value: 3, Color: "air"}, Ruler: "mercury", Detriment: "jupiter", Exaltation: "", Fall: "", ClassicRuler: "mercury", ClassicDetriment: "jupiter", Modality: "mutable", Element: "air", Polarity: "positive"},
	{Body: Body{Name: "cancer", Symbol: "♋", Value: 4, Color: "water"}, Ruler: "moon", Detriment: "saturn", Exaltation: "jupiter", Fall: "mars", ClassicRuler: "moon", ClassicDetriment: "saturn", Modality: "cardinal", Element: "water", Polarity: "negative"},
	{Body: Body{Name: "leo", Symbol: "♌", Value: 5, Color: "fire"}, Ruler: "sun", Detriment: "uranus", Exaltation: "", Fall: "", ClassicRuler: "sun", ClassicDetriment: "saturn", Modality: "fixed", Element: "fire", Polarity: "positive"},
	{Body: Body{Name: "virgo", Symbol: "♍", Value: 6, Color: "earth"}, Ruler: "mercury", Detriment: "neptune", Exaltation: "mercury", Fall: "venus", ClassicRuler: "mercury", ClassicDetriment: "jupiter", Modality: "mutable", Element: "earth", Polarity: "negative"},
	{Body: Body{Name: "libra", Symbol: "♎", Value: 7, Color: "air"}, Ruler: "venus", Detriment: "mars", Exaltation: "saturn", Fall: "sun", ClassicRuler: "venus", ClassicDetriment: "mars", Modality: "cardinal", Element: "air", Polarity: "positive"},
	{Body: Body{Name: "scorpio", Symbol: "♏", Value: 8, Color: "water"}, Ruler: "pluto", Detriment: "venus", Exaltation: "", Fall: "moon", ClassicRuler: "mars", ClassicDetriment: "venus", Modality: "fixed", Element: "water", Polarity: "negative"},
	{Body: Body{Name: "sagittarius", Symbol: "♐", Value: 9, Color: "fire"}, Ruler: "jupiter", Detriment: "mercury", Exaltation: "", Fall: "", ClassicRuler: "jupiter", ClassicDetriment: "mercury", Modality: "mutable", Element: "fire", Polarity: "positive"},
	{Body: Body{Name: "capricorn", Symbol: "♑", Value: 10, Color: "earth"}, Ruler: "saturn", Detriment: "moon", Exaltation: "mars", Fall: "jupiter", ClassicRuler: "saturn", ClassicDetriment: "moon", Modality: "cardinal", Element: "earth", Polarity: "negative"},
	{Body: Body{Name: "aquarius", Symbol: "♒", Value: 11, Color: "air"}, Ruler: "uranus", Detriment: "sun", Exaltation: "", Fall: "", ClassicRuler: "saturn", ClassicDetriment: "sun", Modality: "fixed", Element: "air", Polarity: "positive"},
	{Body: Body{Name: "pisces", Symbol: "♓", Value: 12, Color: "water"}, Ruler: "neptune", Detriment: "mercury", Exaltation: "venus", Fall: "mercury", ClassicRuler: "jupiter", ClassicDetriment: "mercury", Modality: "mutable", Element: "water", Polarity: "negative"},
}

// HOUSE_MEMBERS contains all house definitions
var HOUSE_MEMBERS = []HouseMember{
	{Body: Body{Name: "one", Symbol: "1", Value: 1, Color: "fire"}},
	{Body: Body{Name: "two", Symbol: "2", Value: 2, Color: "earth"}},
	{Body: Body{Name: "three", Symbol: "3", Value: 3, Color: "air"}},
	{Body: Body{Name: "four", Symbol: "4", Value: 4, Color: "water"}},
	{Body: Body{Name: "five", Symbol: "5", Value: 5, Color: "fire"}},
	{Body: Body{Name: "six", Symbol: "6", Value: 6, Color: "earth"}},
	{Body: Body{Name: "seven", Symbol: "7", Value: 7, Color: "air"}},
	{Body: Body{Name: "eight", Symbol: "8", Value: 8, Color: "water"}},
	{Body: Body{Name: "nine", Symbol: "9", Value: 9, Color: "fire"}},
	{Body: Body{Name: "ten", Symbol: "10", Value: 10, Color: "earth"}},
	{Body: Body{Name: "eleven", Symbol: "11", Value: 11, Color: "air"}},
	{Body: Body{Name: "twelve", Symbol: "12", Value: 12, Color: "water"}},
}

// EXTRA_MEMBERS contains extra celestial body definitions
var EXTRA_MEMBERS = []ExtraMember{
	{Body: Body{Name: "chiron", Symbol: "⚷", Value: 15, Color: "asteroids"}},
	{Body: Body{Name: "ceres", Symbol: "⚳", Value: 17, Color: "asteroids"}},
	{Body: Body{Name: "pallas", Symbol: "⚴", Value: 18, Color: "asteroids"}},
	{Body: Body{Name: "juno", Symbol: "⚵", Value: 19, Color: "asteroids"}},
	{Body: Body{Name: "vesta", Symbol: "⚶", Value: 20, Color: "asteroids"}},
}

// VERTEX_MEMBERS contains vertex definitions
var VERTEX_MEMBERS = []VertexMember{
	{Body: Body{Name: "asc", Symbol: "Asc", Value: 1, Color: "fire"}},
	{Body: Body{Name: "ic", Symbol: "IC", Value: 4, Color: "water"}},
	{Body: Body{Name: "dsc", Symbol: "Dsc", Value: 7, Color: "air"}},
	{Body: Body{Name: "mc", Symbol: "MC", Value: 10, Color: "earth"}},
}

// GetSignMember returns the sign member for a given index
func GetSignMember(index int) SignMember {
	if index < 0 || index >= len(SIGN_MEMBERS) {
		return SIGN_MEMBERS[0] // Default to Aries
	}
	return SIGN_MEMBERS[index]
}

// GetAspectMember returns the aspect member by name
func GetAspectMember(name string) *AspectMember {
	for _, aspect := range ASPECT_MEMBERS {
		if aspect.Name == name {
			return &aspect
		}
	}
	return nil
}
