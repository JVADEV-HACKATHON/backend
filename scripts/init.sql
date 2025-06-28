-- Hospital API Database Initialization Script

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_paciente_nombre ON pacientes(nombre);
CREATE INDEX IF NOT EXISTS idx_historial_fecha_ingreso ON historial_clinico(fecha_ingreso);
CREATE INDEX IF NOT EXISTS idx_historial_paciente ON historial_clinico(id_paciente);
CREATE INDEX IF NOT EXISTS idx_historial_hospital ON historial_clinico(id_hospital);
CREATE INDEX IF NOT EXISTS idx_historial_geolocation ON historial_clinico(patient_latitude, patient_longitude);
CREATE INDEX IF NOT EXISTS idx_historial_consultation_date ON historial_clinico(consultation_date);

-- Insert default hospital for testing
-- Note: Password is 'admin123' hashed with bcrypt
INSERT INTO hospitales (nombre, direccion, ciudad, telefono, email, password) 
VALUES (
    'Hospital General Central',
    'Avenida Principal 123, Centro',
    'Ciudad Capital',
    '+1234567890',
    'admin@hospitalcentral.com',
    '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'
) ON CONFLICT (email) DO NOTHING;
