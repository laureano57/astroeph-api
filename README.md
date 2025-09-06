# AstroEph API

Un servicio en Go para cálculos astrológicos que genera cartas natales, sinastría, cartas compuestas, revoluciones y progresiones usando Swiss Ephemeris (swephgo) y genera gráficos en SVG junto con datos en JSON.

## Características

- ✨ **Cartas Natales**: Cálculo completo de posiciones planetarias, casas y aspectos
- 🔮 **Sinastría**: Análisis de compatibilidad entre dos cartas natales
- 🌟 **Cartas Compuestas**: Cálculo de cartas compuestas para relaciones
- ☀️ **Revolución Solar**: Cartas de revolución solar anuales
- 🌙 **Revolución Lunar**: Cartas de revolución lunar mensuales
- 📈 **Progresiones Secundarias**: Cálculo de progresiones
- 🎨 **Gráficos SVG**: Generación de gráficos visuales en múltiples temas
- 🤖 **Formato LLM**: Respuestas optimizadas para modelos de lenguaje
- 🌍 **Geocodificación**: Base de datos GeoNames embebida (223k+ ciudades)

## Arquitectura

El proyecto sigue una arquitectura limpia y modular:

```
/astroeph-api
├── cmd/
│   └── server/
│       └── main.go                 # Punto de entrada de la aplicación
│
├── internal/
│   ├── http/                       # Capa HTTP
│   │   ├── router.go               # Configuración de rutas
│   │   └── handlers/               # Manejadores HTTP específicos
│   │       ├── natal_handler.go
│   │       ├── synastry_handler.go
│   │       ├── composite_handler.go
│   │       ├── solar_return_handler.go
│   │       ├── lunar_return_handler.go
│   │       └── progressions_handler.go
│   │
│   ├── service/                    # Lógica de negocio
│   │   ├── natal_service.go
│   │   ├── synastry_service.go
│   │   ├── composite_service.go
│   │   ├── solar_return_service.go
│   │   ├── lunar_return_service.go
│   │   └── progressions_service.go
│   │
│   ├── domain/                     # Modelos de dominio
│   │   ├── chart.go                # Carta astrológica
│   │   ├── planet.go               # Planetas y cuerpos celestes
│   │   ├── aspect.go               # Aspectos astrológicos
│   │   ├── house.go                # Casas astrológicas
│   │   ├── time.go                 # Manejo de tiempo
│   │   ├── location.go             # Ubicaciones geográficas
│   │   └── utils.go                # Utilidades de dominio
│   │
│   ├── astro/                      # Capa de cálculos astrológicos
│   │   ├── ephemeris.go            # Wrapper sobre swephgo
│   │   ├── planets.go              # Cálculos planetarios
│   │   ├── houses.go               # Cálculos de casas
│   │   ├── aspects.go              # Cálculos de aspectos
│   │   ├── geocoding.go            # Geocodificación
│   │   └── chartdrawer.go          # Generación de gráficos SVG
│   │
│   ├── config/                     # Configuración
│   │   └── config.go
│   │
│   └── logging/                    # Sistema de logging
│       └── logger.go
│
├── pkg/                            # Paquetes públicos
│   ├── errors/                     # Manejo de errores
│   │   └── errors.go
│   ├── utils/                      # Utilidades generales
│   │   └── utils.go
│   └── chart/                      # Librería de generación de gráficos
│       └── [archivos existentes]
│
├── internal/astro/data/             # Datos embebidos de la aplicación
│   ├── cities500.txt               # Base de datos de ciudades (embebida)
│   └── readme.txt                  # Documentación de GeoNames
│
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

Todos los endpoints soportan respuestas JSON estructuradas y opcionalmente respuestas formateadas para LLM mediante el parámetro `"ai_response": true`.

### Cartas Natales
- `POST /api/v1/natal-chart` - Calcular carta natal

### Sinastría
- `POST /api/v1/synastry` - Calcular sinastría entre dos personas

### Cartas Compuestas
- `POST /api/v1/composite-chart` - Calcular carta compuesta

### Revoluciones Solares
- `POST /api/v1/solar-return` - Calcular revolución solar

### Revoluciones Lunares
- `POST /api/v1/lunar-return` - Calcular revolución lunar

### Progresiones
- `POST /api/v1/progressions` - Calcular progresiones secundarias

### Utilidades
- `GET /api/v1/house-systems` - Listar sistemas de casas disponibles
- `GET /health` - Verificar estado del servicio

### Parámetros Comunes

Todos los endpoints de cálculo astrológico soportan:
- `"ai_response": false` (default) - Respuesta JSON estructurada únicamente
- `"ai_response": true` - Incluye campo adicional `"ai_formatted_response"` optimizado para LLMs

## Instalación y Uso

### Prerrequisitos

- Go 1.21 o superior
- Git

### Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/astroeph-api.git
cd astroeph-api
```

### Instalar dependencias

```bash
go mod tidy
```

### Ejecutar el servidor

```bash
# Usando Make (recomendado - configura automáticamente Swiss Ephemeris)
make run

# O manualmente
go run cmd/server/main.go
```

El servidor se iniciará en `http://localhost:8080`

### Compilar para producción

```bash
go build -o astroeph-api cmd/server/main.go
./astroeph-api
```

## Ejemplo de Uso

### Carta Natal (JSON estructurado)

```bash
curl -X POST http://localhost:8080/api/v1/natal-chart \
  -H "Content-Type: application/json" \
  -d '{
    "day": 15,
    "month": 3,
    "year": 1990,
    "local_time": "14:30",
    "city": "Madrid",
    "house_system": "Placidus",
    "draw_chart": true,
    "svg_theme": "light",
    "ai_response": false
  }'
```

### Carta Natal con respuesta optimizada para LLM

```bash
curl -X POST http://localhost:8080/api/v1/natal-chart \
  -H "Content-Type: application/json" \
  -d '{
    "day": 15,
    "month": 3,
    "year": 1990,
    "local_time": "14:30",
    "city": "Madrid",
    "house_system": "Placidus",
    "ai_response": true
  }'
```

### Sinastría

```bash
curl -X POST http://localhost:8080/api/v1/synastry \
  -H "Content-Type: application/json" \
  -d '{
    "person1": {
      "name": "Persona 1",
      "day": 15,
      "month": 3,
      "year": 1990,
      "local_time": "14:30",
      "city": "Madrid"
    },
    "person2": {
      "name": "Persona 2",
      "day": 22,
      "month": 7,
      "year": 1988,
      "local_time": "09:15",
      "city": "Barcelona"
    },
    "ai_response": true
  }'
```

### Comandos Make para pruebas rápidas

```bash
make health     # Verificar estado del servidor
make natal      # Probar endpoint de carta natal
make natal-ai   # Probar carta natal con respuesta AI
make synastry   # Probar endpoint de sinastría
make test-all   # Ejecutar todas las pruebas
```

## Configuración

La aplicación puede configurarse mediante variables de entorno:

- `PORT`: Puerto del servidor (default: 8080)
- `LOG_LEVEL`: Nivel de logging (default: info)
- `LOG_FORMAT`: Formato de logs (default: console)

## Sistemas de Casas Soportados

- Placidus (por defecto)
- Koch
- Porphyrius  
- Regiomontanus
- Campanus
- Equal (Casas Iguales)
- Whole Sign (Signos Completos)

## Temas de Gráficos Disponibles

- `light`: Tema claro
- `dark`: Tema oscuro  
- `mono`: Tema monocromático

## Tecnologías Utilizadas

- **Go**: Lenguaje de programación principal
- **Gin**: Framework web HTTP
- **Swiss Ephemeris (swephgo)**: Cálculos astronómicos precisos
- **SQLite**: Base de datos embebida para geocodificación
- **Zerolog**: Logging estructurado
- **SVG**: Generación de gráficos vectoriales
- **GeoNames**: Base de datos geográfica mundial embebida

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-caracteristica`)
3. Commit tus cambios (`git commit -am 'Agregar nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)
5. Crea un Pull Request

## Licencias y Créditos

### Licencia del Proyecto

Este proyecto está bajo la **Licencia MIT**. Ver el archivo `LICENSE` para más detalles.

### Datos Geográficos - GeoNames

Este proyecto utiliza datos geográficos provenientes del **[GeoNames Gazetteer](https://www.geonames.org/)**, los cuales están licenciados bajo la **[Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/)**.

**Créditos de GeoNames:**
- **Fuente**: GeoNames Gazetteer (https://www.geonames.org/)
- **Licencia**: Creative Commons Attribution 4.0 International License
- **Archivo utilizado**: `cities500.txt` - Ciudades con población > 500 habitantes
- **Formato**: Los datos están embebidos en el binario para mejorar la portabilidad
- **Propósito**: Geocodificación y resolución de coordenadas de ciudades mundiales

**Aviso de Licencia GeoNames:**
```
This work is licensed under a Creative Commons Attribution 4.0 License.
See https://creativecommons.org/licenses/by/4.0/
The Data is provided "as is" without warranty or any representation of accuracy, timeliness or completeness.
```

### Otros Componentes de Terceros

- **[Swiss Ephemeris](https://www.astro.com/swisseph/)**: Cálculos astronómicos precisos (GNU GPL v2 para uso no comercial, licencia comercial disponible)
- **[swephgo](https://github.com/mshafiee/swephgo)**: Wrapper Go para Swiss Ephemeris
- **[Gin Web Framework](https://github.com/gin-gonic/gin)**: Framework HTTP (MIT License)
- **[Zerolog](https://github.com/rs/zerolog)**: Biblioteca de logging (MIT License)
- **[modernc.org/sqlite](https://gitlab.com/cznic/sqlite)**: Driver SQLite puro Go (BSD-3-Clause)

### Agradecimientos

- **GeoNames.org** por proporcionar una base de datos geográfica mundial completa y accesible
- **Astrodienst** por el desarrollo y mantenimiento de Swiss Ephemeris
- **La comunidad de Go** por las excelentes bibliotecas y herramientas
- **Todos los contribuidores** que hacen posible el ecosistema de software libre

## Soporte

Para reportar bugs o solicitar características, por favor abre un issue en GitHub.
