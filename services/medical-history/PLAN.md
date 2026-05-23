# Plan De Desarrollo: Medical History Service

## Objetivo

Implementar historial clinico con GraphQL, documentos flexibles en MongoDB y control de acceso por usuario.

## Fases

1. Crear proyecto Express.js con GraphQL.
2. Definir schema para registros, recetas y adjuntos.
3. Implementar persistencia MongoDB.
4. Implementar queries y mutations principales.
5. Publicar eventos de cambios de historial.
6. Agregar autorizacion, observabilidad y pruebas.

## Entregables Minimos

- Servicio reemplazando el stub en compose mediante `build: ./services/medical-history`.
- Endpoint GraphQL.
- Colecciones para registros clinicos, recetas y adjuntos.
- Mutaciones para crear y actualizar historial.
- Eventos `record-created` y `record-updated`.
- `/health` y `/metrics`.

## Dependencias

- MongoDB `medical-records-db`.
- Kafka.
- User Service para identidad y roles.
- Appointments Service para contexto de cita.

## Interfaces Esperadas

Ruta gateway: `/graphql/medical-history`.

Operaciones GraphQL:

- `patientRecords(patientId)`
- `record(id)`
- `createRecord(input)`
- `updateRecord(id, input)`
- `addPrescription(recordId, input)`
- `addAttachment(recordId, input)`

## Eventos Kafka

Publica:

- `record-created`
- `record-updated`
- `prescription-created`
- `attachment-added`

Consume:

- `appointment-confirmed`
- `call-ended`

## Criterios De Aceptacion

- Un medico autorizado puede crear un registro.
- Un paciente puede consultar su propio historial.
- Usuarios no autorizados reciben error.
- Cambios clinicos publican eventos.
- `/health` y `/metrics` responden.

## Orden Recomendado

1. Bootstrap Express y GraphQL.
2. Schema y modelos MongoDB.
3. Autorizacion.
4. Eventos Kafka.
5. Observabilidad y tests.
