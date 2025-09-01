# Astrological Calculation API (AstroEph-API)

An ultra-fast, self-contained, and scalable astrological API core built with Go. Provides real-time astrological calculations with the Swiss Ephemeris library, local SQLite geocoding service, and flexible output tailored for both human and LLM consumption.

## 🚀 Current Status

**Phase 1-3 Complete:** ✅ Core API structure with working natal chart calculations using mock data.

### Working Features
- ✅ Complete project structure with Go modules
- ✅ RESTful API with Gin framework
- ✅ Health check endpoint (`GET /health`)
- ✅ **Natal chart endpoint** (`POST /api/v1/natal-chart`) - **FULLY FUNCTIONAL**
- ✅ Data models for all endpoints
- ✅ Mock astrological calculations (placeholder for Swiss Ephemeris)
- ✅ House system support (Placidus, Koch, etc.)
- ✅ Aspect calculations between planets
- ✅ JSON and AI-response format support

### Planned Features
- 🔄 Swiss Ephemeris integration for accurate calculations
- 🔄 Transits calculations
- 🔄 Synastry (relationship compatibility) 
- 🔄 Composite charts
- 🔄 Progressions and returns
- 🔄 Local geocoding service with SQLite
- 🔄 Structured logging with zerolog

## 🏃‍♂️ Quick Start

### Prerequisites
- Go 1.22 or higher
- Git

### Installation & Running

```bash
# Clone the repository
git clone <repo-url>
cd astroeph-api

# Install dependencies
go mod tidy

# Build and run
go build
./astroeph-api

# Or run directly
go run main.go
```

The API will be available at `http://localhost:8080`

### Test the API

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Generate a Natal Chart:**
```bash
curl -X POST http://localhost:8080/api/v1/natal-chart \
  -H "Content-Type: application/json" \
  -d '{
    "day": 15,
    "month": 6,
    "year": 1990,
    "local_time": "14:30:00",
    "city": "New York",
    "house_system": "Placidus"
  }'
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### 1. Health Check
```
GET /health
```
**Response:** Service status and version info.

#### 2. Natal Chart (✅ Working)
```
POST /natal-chart
```

**Request Body:**
```json
{
  "day": 15,
  "month": 6, 
  "year": 1990,
  "local_time": "14:30:00",
  "city": "New York",
  "house_system": "Placidus",
  "ai_response": false
}
```

**Response:** Complete natal chart with planets, houses, aspects, and birth info.

#### 3. Transits (🔄 Coming Soon)
```
POST /transits
```
Calculate planetary transits for a specific date.

#### 4. Synastry (🔄 Coming Soon)
```
POST /synastry  
```
Calculate relationship compatibility between two charts.

#### 5. Additional Endpoints
- `POST /composite-chart` - Composite relationship chart
- `POST /progressions` - Secondary progressions
- `POST /solar-return` - Annual solar return chart
- `POST /lunar-return` - Monthly lunar return chart

## 🏗️ Architecture

```
├── main.go              # Application entry point
├── api/
│   └── routes.go        # HTTP route handlers
├── services/
│   └── astrology_service.go  # Core astrological calculations
├── models/
│   └── models.go        # Request/response data models
├── go.mod               # Go module definition
└── README.md           # This file
```

## 🔧 Configuration

### Supported House Systems
- Placidus (default)
- Koch
- Porphyrius  
- Regiomontanus
- Campanus
- Equal
- Whole Sign

### Supported Cities (Mock Data)
Currently includes coordinates for major cities: New York, London, Tokyo, Sydney, Los Angeles, Paris, Berlin, Moscow, Mumbai, São Paulo.

*Full geocoding service with SQLite database coming in Phase 4.*

## 🚧 Development Notes

### Current Implementation
This version uses **mock astronomical data** to provide a working API structure. The calculations return realistic-looking data for development and testing purposes.

### Next Steps (Phase 4)
1. **Swiss Ephemeris Integration**: Replace mock data with actual astronomical calculations
2. **Geocoding Service**: Local SQLite database for city coordinates and timezones
3. **Structured Logging**: Add comprehensive logging with zerolog
4. **LLM-Optimized Output**: Text formatting for AI applications

### Swiss Ephemeris Setup
To enable real astronomical calculations, install the Swiss Ephemeris C library:

```bash
# macOS
brew install swisseph

# Ubuntu/Debian
sudo apt-get install libswe-dev

# Then rebuild the project
go build
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable  
5. Submit a pull request

## 📄 License

See LICENSE file for details.

## 🔮 Sample Response

Here's what a natal chart response looks like:

```json
{
  "planets": [
    {
      "name": "Sun",
      "longitude": 120.5,
      "sign": "Leo", 
      "degree": 0.5,
      "house": 3
    }
    // ... more planets
  ],
  "houses": [
    {
      "house": 1,
      "cusp": 45,
      "sign": "Taurus"
    }
    // ... all 12 houses
  ],
  "aspects": [
    {
      "planet1": "Sun",
      "planet2": "Mars", 
      "type": "square",
      "angle": 89.6,
      "orb": 0.4
    }
    // ... all aspects
  ],
  "ascendant": 45,
  "midheaven": 315,
  "birth_info": {
    "date": "1990-06-15",
    "time": "14:30:00",
    "city": "New York",
    "latitude": 40.7128,
    "longitude": -74.006
  }
}
```