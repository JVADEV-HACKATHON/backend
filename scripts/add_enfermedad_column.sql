-- Migración para agregar el campo 'enfermedad' a la tabla historial_clinico
-- Ejecutar este script en la base de datos para agregar la nueva columna

ALTER TABLE historial_clinico 
ADD COLUMN enfermedad VARCHAR(150) NOT NULL DEFAULT 'Sin especificar';

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
ALTER TABLE historial_clinico 
ALTER COLUMN enfermedad DROP DEFAULT;

-- Comentario de confirmación
-- La columna 'enfermedad' ha sido agregada exitosamente
-- Todos los registros existentes han sido actualizados
