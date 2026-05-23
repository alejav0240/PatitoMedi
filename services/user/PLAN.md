# Plan De Desarrollo: User Service

## Objetivo

Implementar identidad, autenticacion, perfiles de pacientes y medicos, y emision de JWT para el resto de PatitoMedi.

## Fases

1. Crear proyecto Go con servidor HTTP, config por env y conexion PostgreSQL.
2. Implementar migraciones para pacientes, medicos, roles y sesiones.
3. Implementar registro de pacientes y medicos.
4. Implementar login, JWT, logout y endpoint `/me`.
5. Publicar eventos Kafka de usuario.
6. Agregar `/health`, `/metrics` y pruebas.

## Entregables Minimos

- Servicio reemplazando el stub en compose mediante `build: ./services/user`.
- Endpoints REST de registro, login y perfil.
- Password hashing seguro.
- JWT firmado por secreto local.
- Persistencia en `users-db`.
- Evento `user-registered`.

## Dependencias

- PostgreSQL `users-db`.
- Kafka para eventos.
- Kong para exposicion publica.

## Endpoints Esperados

- `POST /register/patient`
- `POST /register/doctor`
- `POST /login`
- `POST /logout`
- `GET /me`
- `PATCH /me`
- `GET /doctors`
- `GET /doctors/{id}`

Ruta gateway: `/api/users`.

## Eventos Kafka

Publica:

- `user-registered`
- `user-updated`
- `session-created`
- `session-ended`

Consume:

- Ninguno en MVP.

## Criterios De Aceptacion

- Un paciente y un medico pueden registrarse.
- Un usuario puede iniciar sesion y consultar `/me`.
- Los endpoints privados rechazan JWT ausente o invalido.
- El servicio publica `user-registered`.
- `/health` y `/metrics` responden.

## Orden Recomendado

1. Bootstrap Go.
2. Migraciones y repositorios.
3. Auth y JWT.
4. Eventos Kafka.
5. Observabilidad y tests.

## Estado Actual

- [x] Bootstrap Go con servidor HTTP.
- [x] Configuracion por variables de entorno.
- [x] Conexion PostgreSQL.
- [x] Schema inicial con pacientes, medicos, roles y sesiones.
- [x] Registro de pacientes y medicos.
- [x] Login, JWT, logout y `/me`.
- [x] Publicacion Kafka para eventos de usuario y sesion.
- [x] `/health` y `/metrics`.
- [x] Tests unitarios para JWT, validacion de registro y middleware de auth.
- [ ] Verificacion runtime pendiente con `docker compose build user-service` y `docker compose up` cuando Docker Hub y dependencias Go esten disponibles.
