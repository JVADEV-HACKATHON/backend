-- Script de migración completo para agregar coordenadas a hospitales y enfermedad a historial clínico
-- Ejecutar este script para aplicar todas las modificaciones necesarias

-- ============================================
-- PARTE 1: Agregar coordenadas a hospitales
-- ============================================

-- Verificar si las columnas ya existen antes de agregarlas
DO $$ 
BEGIN
    -- Agregar latitud si no existe
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hospitales' AND column_name = 'latitud'
    ) THEN
        ALTER TABLE hospitales ADD COLUMN latitud DECIMAL(10,8) NOT NULL DEFAULT 0;
    END IF;
    
    -- Agregar longitud si no existe
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'hospitales' AND column_name = 'longitud'
    ) THEN
        ALTER TABLE hospitales ADD COLUMN longitud DECIMAL(11,8) NOT NULL DEFAULT 0;
    END IF;
END $$;

-- Actualizar hospitales existentes con coordenadas reales de La Paz
UPDATE hospitales 
SET latitud = -16.5189, longitud = -68.0888 
WHERE email = 'admin@hospitalcentral.com' AND (latitud = 0 OR latitud IS NULL);

UPDATE hospitales 
SET latitud = -16.5203, longitud = -68.0854 
WHERE email = 'admin@hospitalnino.com' AND (latitud = 0 OR latitud IS NULL);

UPDATE hospitales 
SET latitud = -16.5203, longitud = -68.1127 
WHERE email = 'admin@hospitalclinicas.com' AND (latitud = 0 OR latitud IS NULL);

UPDATE hospitales 
SET latitud = -16.5075, longitud = -68.1064 
WHERE email = 'admin@hospitalsangabriel.com' AND (latitud = 0 OR latitud IS NULL);

UPDATE hospitales 
SET latitud = -16.5322, longitud = -68.0753 
WHERE email = 'admin@hospitalarcoiris.com' AND (latitud = 0 OR latitud IS NULL);

-- Remover valores por defecto
ALTER TABLE hospitales ALTER COLUMN latitud DROP DEFAULT;
ALTER TABLE hospitales ALTER COLUMN longitud DROP DEFAULT;

-- ============================================
-- PARTE 2: Agregar enfermedad a historial clínico
-- ============================================

-- Verificar si la columna enfermedad ya existe
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'historial_clinico' AND column_name = 'enfermedad'
    ) THEN
        ALTER TABLE historial_clinico ADD COLUMN enfermedad VARCHAR(150) NOT NULL DEFAULT 'Sin especificar';
    END IF;
END $$;

-- Actualizar registros existentes con valores por defecto basados en el diagnóstico
UPDATE historial_clinico 
SET enfermedad = CASE 
    WHEN diagnostico LIKE '%Migraña%' THEN 'Migraña'
    WHEN diagnostico LIKE '%Bronquitis%' THEN 'Bronquitis aguda'
    WHEN diagnostico LIKE '%Gastritis%' THEN 'Gastritis'
    WHEN diagnostico LIKE '%Dermatitis%' THEN 'Dermatitis alérgica'
    WHEN diagnostico LIKE '%gripal%' OR diagnostico LIKE '%Influenza%' THEN 'Influenza'
    WHEN diagnostico LIKE '%Esguince%' THEN 'Esguince de tobillo'
    WHEN diagnostico LIKE '%Conjuntivitis%' THEN 'Conjuntivitis viral'
    WHEN diagnostico LIKE '%hipertensiv%' THEN 'Hipertensión arterial'
    WHEN diagnostico LIKE '%Gastroenteritis%' THEN 'Gastroenteritis viral'
    WHEN diagnostico LIKE '%Embarazo%' THEN 'Embarazo normal'
    WHEN diagnostico LIKE '%COVID%' THEN 'COVID-19'
    WHEN diagnostico LIKE '%Intoxicación%' THEN 'Intoxicación alimentaria'
    ELSE 'Sin especificar'
END
WHERE enfermedad = 'Sin especificar' OR enfermedad IS NULL;

-- Remover el valor por defecto después de la migración
ALTER TABLE historial_clinico ALTER COLUMN enfermedad DROP DEFAULT;

-- ============================================
-- PARTE 3: Crear índices para optimización
-- ============================================

-- Índice para búsquedas geográficas de hospitales
CREATE INDEX IF NOT EXISTS idx_hospitales_location ON hospitales(latitud, longitud);

-- Índice para búsquedas por enfermedad
CREATE INDEX IF NOT EXISTS idx_historial_enfermedad ON historial_clinico(enfermedad);

-- Índice combinado para análisis epidemiológicos
CREATE INDEX IF NOT EXISTS idx_historial_epidemiologico ON historial_clinico(enfermedad, is_contagious, consultation_date);

-- ============================================
-- CONFIRMACIÓN
-- ============================================

-- Mostrar resumen de cambios
SELECT 
    'Hospitales con coordenadas' as tabla,
    COUNT(*) as total_registros,
    COUNT(CASE WHEN latitud != 0 AND longitud != 0 THEN 1 END) as con_coordenadas
FROM hospitales

UNION ALL

SELECT 
    'Historiales con enfermedad' as tabla,
    COUNT(*) as total_registros,
    COUNT(CASE WHEN enfermedad != 'Sin especificar' THEN 1 END) as con_enfermedad
FROM historial_clinico;

-- Comentario final
-- ✅ Migración completada exitosamente
-- ✅ Coordenadas agregadas a hospitales
-- ✅ Campo enfermedad agregado a historial clínico
-- ✅ Índices optimizados creados
