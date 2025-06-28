#!/bin/sh

echo "🌱 Iniciando Hospital API..."

# Esperar a que la base de datos esté lista
echo "⏳ Esperando a que la base de datos esté ready..."
sleep 10

# Verificar conexión a la base de datos
echo "🔍 Verificando conexión a la base de datos..."
while ! nc -z $DB_HOST $DB_PORT; do
  echo "⏳ Esperando conexión a la base de datos..."
  sleep 2
done

echo "✅ Base de datos conectada!"

# Ejecutar las semillas si AUTO_SEED está habilitado
if [ "$AUTO_SEED" = "true" ]; then
  echo "🌱 Ejecutando semillas automáticamente..."
  ./seed -clean
  if [ $? -eq 0 ]; then
    echo "✅ Semillas ejecutadas correctamente!"
    echo "📋 Datos disponibles:"
    echo "   🏥 5 Hospitales (usuario: admin@hospitalcentral.com / admin123)"
    echo "   👥 15 Pacientes"
    echo "   📊 12+ Historiales clínicos con datos geográficos"
  else
    echo "❌ Error ejecutando semillas, pero continuando..."
  fi
else
  echo "⚠️  AUTO_SEED deshabilitado, saltando seeding"
fi

# Iniciar el servidor principal
echo "🚀 Iniciando servidor Hospital API en puerto $PORT..."
exec ./main
