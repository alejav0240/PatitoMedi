# Medical History Service

Servicio responsable del historial clínico, documentos médicos, recetas y registros asociados al paciente.

## Stack previsto

- Runtime: Node.js.
- Framework: Express.js.
- Protocolo público: GraphQL.
- Base de datos: MongoDB `medical-records-db`.
- Eventos: Apache Kafka.

## Responsabilidades

- Gestionar registros clínicos del paciente.
- Guardar diagnósticos, notas médicas, recetas y adjuntos.
- Exponer consultas GraphQL para el historial médico.
- Controlar acceso por rol entre paciente, médico y personal autorizado.
- Publicar eventos cuando el historial cambia.

## API esperada

Ruta por gateway: `/graphql/medical-history`.

Operaciones GraphQL previstas:

- `patientRecords(patientId)`: listar registros de un paciente.
- `record(id)`: consultar un registro clínico.
- `createRecord(input)`: crear registro clínico.
- `updateRecord(id, input)`: actualizar registro clínico.
- `addPrescription(recordId, input)`: agregar receta.
- `addAttachment(recordId, input)`: asociar documento.

## Datos

Colecciones previstas:

- `clinical_records`
- `prescriptions`
- `attachments`

MongoDB permite almacenar documentos clínicos flexibles, pero el esquema GraphQL debe mantener contratos claros hacia los clientes.

## Eventos Kafka

Publica:

- `record-created`
- `record-updated`
- `prescription-created`
- `attachment-added`

Consume:

- `appointment-confirmed`
- `call-ended`

## Observabilidad

Debe exponer `/metrics` para Prometheus con métricas de consultas GraphQL, mutaciones, errores de autorización y latencia por resolver.
