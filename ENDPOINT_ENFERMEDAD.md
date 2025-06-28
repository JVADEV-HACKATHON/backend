# Endpoint: Buscar Historiales por Enfermedad

## Descripción
Endpoint que permite buscar historiales clínicos por nombre de enfermedad y devuelve los datos en el formato específico solicitado.

## URL
```
GET /api/v1/historial/enfermedad
```

## Autenticación
- Requiere token JWT de hospital autenticado
- Header: `Authorization: Bearer <token>`

## Parámetros de Consulta

| Parámetro  | Tipo   | Requerido | Descripción                    | Ejemplo     |
|------------|--------|-----------|--------------------------------|-------------|
| enfermedad | string | Sí        | Nombre de la enfermedad        | `Dengue`    |
| page       | int    | No        | Número de página (default: 1)  | `1`         |
| limit      | int    | No        | Elementos por página (max: 100, default: 10) | `10` |

## Ejemplo de Solicitud
```bash
curl -X GET "http://localhost:8080/api/v1/historial/enfermedad?enfermedad=Dengue&page=1&limit=10" \
     -H "Authorization: Bearer your-jwt-token"
```

## Ejemplo de Respuesta
```json
{
  "message": "Datos obtenidos exitosamente",
  "total": 2,
  "data": [
    {
      "id": 1,
      "fecha_ingreso": "2024-06-15T10:30:00Z",
      "motivo_consulta": "Fiebre alta y dolor de cabeza intenso",
      "diagnostico": "Dengue clásico sin signos de alarma",
      "tratamiento": "Reposo absoluto e hidratación oral",
      "medicamentos": "Paracetamol 500mg cada 6 horas",
      "observaciones": "Paciente en vigilancia epidemiológica",
      "patient_latitude": -17.783300,
      "patient_longitude": -63.182100,
      "patient_address": "Av. San Martín 3456, Equipetrol",
      "patient_district": "Equipetrol",
      "patient_neighborhood": "Equipetrol Norte",
      "consultation_date": "2024-06-15T00:00:00Z",
      "symptoms_start_date": "2024-06-10T00:00:00Z",
      "is_contagious": true,
      "created_at": "2024-06-15T10:35:00Z",
      "paciente": {
        "id": 1,
        "nombre": "Carlos",
        "apellido": "Suárez Mendoza",
        "edad": 45,
        "sexo": "M"
      },
      "hospital": {
        "id": 1,
        "nombre": "Hospital Universitario Japonés",
        "direccion": "Av. Japón s/n, 3er Anillo Externo",
        "hospital_latitude": -17.7833,
        "hospital_longitude": -63.1821,
        "ciudad": "Santa Cruz de la Sierra",
        "telefono": "+591-3-3460101"
      }
    }
  ]
}
```

## Enfermedades Disponibles
El seeder genera datos para las siguientes enfermedades:

1. **Dengue** (contagiosa)
2. **Sarampión** (contagiosa)
3. **Zika** (contagiosa)
4. **Influenza** (contagiosa)
5. **Gripe AH1N1** (contagiosa)
6. **Bronquitis** (no contagiosa)

## Códigos de Respuesta

| Código | Descripción |
|--------|-------------|
| 200    | Éxito - Datos obtenidos correctamente |
| 400    | Error - Parámetro 'enfermedad' faltante |
| 401    | No autenticado - Token JWT inválido o faltante |
| 500    | Error interno del servidor |

## Características Especiales

- **Búsqueda case-insensitive**: La búsqueda no distingue entre mayúsculas y minúsculas
- **Paginación**: Soporte completo para paginación de resultados
- **Datos completos**: Incluye información del paciente y hospital relacionados
- **Información geográfica**: Coordenadas reales para mapas de calor epidemiológicos
- **Separación de nombres**: Los nombres completos se dividen automáticamente en nombre y apellido

## Notas Técnicas

- La búsqueda es exacta por nombre de enfermedad
- Los datos del paciente incluyen la edad calculada automáticamente
- Las coordenadas están en formato decimal (WGS84)
- Los datos de Santa Cruz de la Sierra incluyen distritos y barrios reales
- Todas las fechas están en formato ISO 8601 UTC
