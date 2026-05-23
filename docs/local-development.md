# Desarrollo Local

## Preparar Variables

Copiar el archivo de ejemplo y ajustar valores si algun puerto local esta ocupado:

```bash
cp .env.example .env
```

Las credenciales son solo para desarrollo local. No deben usarse en produccion.

## Validar Docker Compose

Antes de levantar contenedores:

```bash
docker compose config --quiet
```

Si el comando no imprime salida, la configuracion es valida.

## Levantar Servicios

```bash
docker compose up -d
```

Ver estado:

```bash
docker compose ps
```

## Rutas Locales

- Nginx: `http://localhost`
- Kong proxy: `http://localhost:8000`
- Kong admin: `http://localhost:8001`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`

Credenciales locales de Grafana:

- Usuario: `admin`
- Password: `admin`

## Probar Gateway

Mientras los servicios reales no existan, el compose usa stubs `hashicorp/http-echo`.

```bash
curl http://localhost/api/users
curl http://localhost/api/appointments
curl http://localhost/api/payments
curl http://localhost/graphql/medical-history
curl http://localhost/ws/video
```

Tambien se pueden probar las rutas directamente contra Kong:

```bash
curl http://localhost:8000/api/users
curl http://localhost:8000/api/appointments
curl http://localhost:8000/api/payments
```

## Puertos De Infraestructura

- Users PostgreSQL: `localhost:5433`
- Appointments PostgreSQL: `localhost:5434`
- Payments PostgreSQL: `localhost:5435`
- MongoDB: `localhost:27017`
- Redis: `localhost:6379`
- Kafka: `localhost:9092`
- coturn: `localhost:3478`

## Apagar Servicios

```bash
docker compose down
```

Apagar y borrar volumenes locales:

```bash
docker compose down -v
```

Usar `down -v` solo si se quiere eliminar la informacion local de bases de datos, Kafka, Redis, Prometheus y Grafana.

## Reemplazar Stubs Por Servicios Reales

Cuando exista implementacion en una carpeta de `services/*`, reemplazar el bloque del stub en `docker-compose.yml` por un `build`, por ejemplo:

```yaml
user-service:
  build:
    context: ./services/user
  env_file:
    - .env
```

Mantener el mismo nombre de servicio en compose para no romper rutas Kong ni dependencias internas.
