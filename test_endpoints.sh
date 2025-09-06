#!/bin/bash

# AstroEph API - Test All Endpoints
# Make sure the server is running on localhost:8080

echo "🏥 Testing Health Check..."
curl -s http://localhost:8080/health | jq .

echo -e "\n🌟 Testing Natal Chart (JSON Only)..."
curl -X POST http://localhost:8080/api/v1/natal-chart \
  -H "Content-Type: application/json" \
  -d '{
    "day": 15,
    "month": 6,
    "year": 1990,
    "local_time": "14:30",
    "city": "London",
    "house_system": "Placidus",
    "draw_chart": false,
    "ai_response": false
  }' | jq '.birth_info'

echo -e "\n🌟 Testing Natal Chart (With AI Response)..."
curl -X POST http://localhost:8080/api/v1/natal-chart \
  -H "Content-Type: application/json" \
  -d '{
    "day": 15,
    "month": 6,
    "year": 1990,
    "local_time": "14:30",
    "city": "London",
    "house_system": "Placidus",
    "draw_chart": false,
    "ai_response": true
  }' | jq 'has("ai_formatted_response")'

echo -e "\n💕 Testing Synastry (With AI Response)..."
curl -X POST http://localhost:8080/api/v1/synastry \
  -H "Content-Type: application/json" \
  -d '{
    "person1": {
      "day": 15,
      "month": 6,
      "year": 1990,
      "local_time": "14:30",
      "city": "London",
      "name": "Person 1"
    },
    "person2": {
      "day": 22,
      "month": 3,
      "year": 1992,
      "local_time": "10:15",
      "city": "Paris",
      "name": "Person 2"
    },
    "draw_chart": false,
    "ai_response": true
  }' | jq 'has("ai_formatted_response")'

echo -e "\n🔗 Testing Composite Chart..."
curl -X POST http://localhost:8080/api/v1/composite \
  -H "Content-Type: application/json" \
  -d '{
    "person1": {
      "day": 15,
      "month": 6,
      "year": 1990,
      "local_time": "14:30",
      "city": "London",
      "name": "Person 1"
    },
    "person2": {
      "day": 22,
      "month": 3,
      "year": 1992,
      "local_time": "10:15",
      "city": "Paris",
      "name": "Person 2"
    },
    "draw_chart": false,
    "ai_response": true
  }' | jq 'has("ai_formatted_response")'

echo -e "\n☀️ Testing Solar Return..."
curl -X POST http://localhost:8080/api/v1/solar-return \
  -H "Content-Type: application/json" \
  -d '{
    "birth_day": 15,
    "birth_month": 6,
    "birth_year": 1990,
    "birth_time": "14:30",
    "birth_city": "London",
    "return_year": 2024,
    "return_city": "New York",
    "draw_chart": false,
    "ai_response": true
  }' | jq '.return_date'

echo -e "\n🌙 Testing Lunar Return..."
curl -X POST http://localhost:8080/api/v1/lunar-return \
  -H "Content-Type: application/json" \
  -d '{
    "birth_day": 15,
    "birth_month": 6,
    "birth_year": 1990,
    "birth_time": "14:30",
    "birth_city": "London",
    "return_month": 12,
    "return_year": 2024,
    "return_city": "Madrid",
    "draw_chart": false,
    "ai_response": true
  }' | jq '.return_date'

echo -e "\n📈 Testing Progressions..."
curl -X POST http://localhost:8080/api/v1/progressions \
  -H "Content-Type: application/json" \
  -d '{
    "birth_day": 15,
    "birth_month": 6,
    "birth_year": 1990,
    "birth_time": "14:30",
    "birth_city": "London",
    "progression_day": 15,
    "progression_month": 6,
    "progression_year": 2024,
    "draw_chart": false,
    "ai_response": true
  }' | jq '.years_progressed'

echo -e "\n✅ All endpoint tests completed!"
