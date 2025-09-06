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
- 🌍 **Geocodificación**: Base de datos integrada de ciudades mundiales

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
├── data/                           # Datos de la aplicación
│   └── geocoding/
│       └── cities500.txt           # Base de datos de ciudades
│
├── go.mod
├── go.sum
└── README.md
```

## API Endpoints

### Cartas Natales
- `POST /api/v1/natal-chart` - Calcular carta natal
- `POST /api/v1/natal-chart/formatted` - Obtener carta natal formateada para LLM

### Sinastría
- `POST /api/v1/synastry` - Calcular sinastría entre dos personas
- `POST /api/v1/synastry/formatted` - Obtener sinastría formateada para LLM

### Cartas Compuestas
- `POST /api/v1/composite-chart` - Calcular carta compuesta
- `POST /api/v1/composite-chart/formatted` - Obtener carta compuesta formateada

### Revoluciones Solares
- `POST /api/v1/solar-return` - Calcular revolución solar
- `POST /api/v1/solar-return/formatted` - Obtener revolución solar formateada

### Revoluciones Lunares
- `POST /api/v1/lunar-return` - Calcular revolución lunar
- `POST /api/v1/lunar-return/formatted` - Obtener revolución lunar formateada

### Progresiones
- `POST /api/v1/progressions` - Calcular progresiones secundarias
- `POST /api/v1/progressions/formatted` - Obtener progresiones formateadas

### Utilidades
- `GET /api/v1/house-systems` - Listar sistemas de casas disponibles
- `GET /health` - Verificar estado del servicio

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
go run cmd/server/main.go
```

El servidor se iniciará en `http://localhost:8080`

### Compilar para producción

```bash
go build -o astroeph-api cmd/server/main.go
./astroeph-api
```

## Ejemplo de Uso

### Carta Natal

```bash
curl -X POST http://localhost:8080/api/v1/natal-chart \
  -H "Content-Type: application/json" \
  -d '{
    "day": 15,
    "month": 3,
    "year": 1990,
    "local_time": "14:30:00",
    "city": "Madrid",
    "house_system": "Placidus",
    "draw_chart": true,
    "svg_theme": "light"
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
      "local_time": "14:30:00",
      "city": "Madrid"
    },
    "person2": {
      "name": "Persona 2",
      "day": 22,
      "month": 7,
      "year": 1988,
      "local_time": "09:15:00",
      "city": "Barcelona"
    }
  }'
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
- **SQLite**: Base de datos para geocodificación
- **Zerolog**: Logging estructurado
- **SVG**: Generación de gráficos vectoriales

## Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-caracteristica`)
3. Commit tus cambios (`git commit -am 'Agregar nueva característica'`)
4. Push a la rama (`git push origin feature/nueva-caracteristica`)
5. Crea un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Agradecimientos

- Swiss Ephemeris por los cálculos astronómicos precisos
- GeoNames por la base de datos de ciudades
- La comunidad de Go por las excelentes librerías

## Soporte

Para reportar bugs o solicitar características, por favor abre un issue en GitHub.

## Changelog

### v1.0.0 (Refactorización Arquitectónica)
- ✨ Arquitectura limpia con separación de capas
- 🔧 Inyección de dependencias
- 📝 Sistema de logging mejorado
- 🛠️ Manejo de errores robusto
- 🎯 Servicios especializados por tipo de carta
- 🌐 Handlers HTTP separados por funcionalidad
- 📊 Modelos de dominio bien definidos
- ⚡ Capa astro optimizada
- 🔍 Mejor organización del código