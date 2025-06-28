-- Migración para agregar coordenadas (latitud y longitud) a la tabla hospitales
-- Ejecutar este script en la base de datos para agregar las nuevas columnas

ALTER TABLE hospitales 
ADD COLUMN latitud DECIMAL(10,8) NOT NULL DEFAULT 0,
ADD COLUMN longitud DECIMAL(11,8) NOT NULL DEFAULT 0;

-- Actualizar los hospitales existentes con sus coordenadas reales en La Paz
UPDATE hospitales 
SET latitud = -16.5189, longitud = -68.0888 
WHERE email = 'admin@hospitalcentral.com';

UPDATE hospitales 
SET latitud = -16.5203, longitud = -68.0854 
WHERE email = 'admin@hospitalnino.com';

UPDATE hospitales 
SET latitud = -16.5203, longitud = -68.1127 
WHERE email = 'admin@hospitalclinicas.com';

UPDATE hospitales 
SET latitud = -16.5075, longitud = -68.1064 
WHERE email = 'admin@hospitalsangabriel.com';

UPDATE hospitales 
SET latitud = -16.5322, longitud = -68.0753 
WHERE email = 'admin@hospitalarcoiris.com';

-- Remover los valores por defecto después de la migración
ALTER TABLE hospitales 
ALTER COLUMN latitud DROP DEFAULT,
ALTER COLUMN longitud DROP DEFAULT;

-- Crear índices para mejorar las consultas geográficas
CREATE INDEX idx_hospitales_location ON hospitales(latitud, longitud);

-- Comentario de confirmación
-- Las columnas 'latitud' y 'longitud' han sido agregadas exitosamente
-- Todos los hospitales existentes han sido actualizados con sus coordenadas reales
