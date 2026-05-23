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
