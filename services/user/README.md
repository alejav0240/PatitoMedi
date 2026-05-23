# User Service

Servicio responsable de identidad, autenticación y perfiles de usuarios de PatitoMedi.

## Stack previsto

- Lenguaje: Golang.
- Protocolo público: REST.
- Base de datos: PostgreSQL `users-db`.
- Eventos: Apache Kafka.

## Responsabilidades

- Registrar pacientes y médicos.
- Autenticar usuarios y emitir tokens JWT.
- Gestionar perfiles, roles y sesiones.
- Exponer datos básicos de pacientes y médicos para otros servicios.
- Publicar eventos de ciclo de vida de usuarios.

## API esperada

Ruta base por gateway: `/api/users`.

- `POST /register/patient`: registrar paciente.
- `POST /register/doctor`: registrar médico.
- `POST /login`: autenticar usuario.
- `POST /logout`: cerrar sesión.
- `GET /me`: obtener perfil autenticado.
- `PATCH /me`: actualizar perfil autenticado.
- `GET /doctors`: listar médicos disponibles.
- `GET /doctors/{id}`: consultar médico.

## Datos

Tablas iniciales:

- `patients`
- `doctors`
- `sessions`

El servicio es dueño de su base de datos. Otros servicios no deben acceder directamente a sus tablas.

## Eventos Kafka

Publica:

- `user-registered`
- `user-updated`
- `session-created`
- `session-ended`

Consume:

- Eventos de negocio que requieran enriquecer auditoría o notificaciones futuras.

## Observabilidad

Debe exponer `/metrics` para Prometheus con métricas de latencia, errores, logins exitosos/fallidos y registros por tipo de usuario.

## Desarrollo local

El servicio ya tiene una implementacion MVP en Go y se construye desde Docker Compose:

```bash
docker compose build user-service
docker compose up -d users-db kafka user-service
```

Variables principales:

- `DATABASE_URL`: conexion PostgreSQL.
- `JWT_SECRET`: secreto HMAC para firmar tokens.
- `JWT_ISSUER`: issuer esperado en el token.
- `JWT_TTL`: duracion del token, por ejemplo `24h`.
- `KAFKA_BROKERS`: brokers Kafka separados por coma.
- `KAFKA_ENABLED`: permite desactivar publicacion de eventos en local.

Ejemplo de registro de paciente:

```bash
curl -X POST http://localhost/api/users/register/patient \
  -H 'Content-Type: application/json' \
  -d '{"email":"ana@example.com","password":"password123","full_name":"Ana Perez"}'
```

Ejemplo de registro de medico:

```bash
curl -X POST http://localhost/api/users/register/doctor \
  -H 'Content-Type: application/json' \
  -d '{"email":"dr@example.com","password":"password123","full_name":"Dr. Gomez","specialty":"Cardiologia"}'
```

Ejemplo de login:

```bash
curl -X POST http://localhost/api/users/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"ana@example.com","password":"password123"}'
```
