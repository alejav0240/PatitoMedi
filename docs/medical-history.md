# Medical History Service

Service de historial médico (GraphQL) usado para almacenar y consultar registros clínicos.

**Stack:** Node.js · Express · Apollo Server · MongoDB · Mongoose · Kafka (opcional)

**Gateway base:** `/graphql/medical-history` (expuesto a través de Kong)

Estado: ✅ Implementado (servicio real disponible en Compose)

## Resumen

Este microservicio expone un endpoint GraphQL en `/graphql/medical-history` que ofrece consultas y mutaciones para manejar `ClinicalRecord`.

Dentro del contenedor el servidor escucha en el puerto `8080` y espera `POST /graphql/medical-history`.

Kong enruta la ruta pública `http://localhost:8000/graphql/medical-history` hacia el servicio. Importante: en la configuración actual `strip_path` está desactivado para que el upstream reciba el path completo.

## Esquema (resumen)

- Query
  - `healthGraphQL: String` — comprobación rápida del GraphQL
  - `patientRecords(patientId: String!): [ClinicalRecord!]!` — lista por paciente
  - `record(id: ID!): ClinicalRecord` — obtener por id

- Mutation
  - `createRecord(input: ClinicalRecordInput!): ClinicalRecord!` — crear registro
  - `updateRecord(id: ID!, input: ClinicalRecordInput!): ClinicalRecord!` — actualizar

## Modelo `ClinicalRecord`

- `id`, `patientId`, `bloodType`, `allergies`, `chronicDiseases`, `medications`, `notes`, `createdAt`, `updatedAt`

## Ejemplos (curl vía Kong)

Health check

```bash
curl -sS -X POST http://localhost:8000/graphql/medical-history \
  -H 'Content-Type: application/json' \
  -d '{"query":"query { healthGraphQL }"}'
```

Crear un registro (mutation)

```bash
curl -sS -X POST http://localhost:8000/graphql/medical-history \
  -H 'Content-Type: application/json' \
  -d '{"query":"mutation($input:ClinicalRecordInput!){ createRecord(input:$input){ id patientId bloodType notes createdAt } }","variables":{"input":{"patientId":"pat-1","bloodType":"O+","allergies":["penicillin"],"chronicDiseases":[],"medications":[],"notes":["nota"]}}}'
```

Consultar registros por paciente

```bash
curl -sS -X POST http://localhost:8000/graphql/medical-history \
  -H 'Content-Type: application/json' \
  -d '{"query":"query($patientId:String!){ patientRecords(patientId:$patientId){ id patientId bloodType notes createdAt } }","variables":{"patientId":"pat-1"}}'
```

## Ejecutar localmente

Opciones:

- Ejecutar dentro de Compose (recomendado, ya incluido): el servicio se levantará desde `services/medical-history` y se conectará a `medical-records-db`.

- Ejecutar localmente en tu máquina:

```bash
cd services/medical-history
npm install
export MONGO_URI='mongodb://localhost:27017/medical_records' # si mapear puerto 27017 a host
npm start
```

Después, probar las queries vía `http://localhost:4000/graphql/medical-history` (o el `PORT` expuesto si lo cambias).

## Notas y troubleshooting

- El servicio requiere Node >=20 para las dependencias (`@apollo/server`, `mongoose`, `mongodb`). El `Dockerfile` fue actualizado a `node:20-alpine`.
- Si ves errores relacionados con `crypto is not defined`, actualiza la versión de Node a 20 o superior.
- Asegúrate de que `medical-records-db` (Mongo) esté healthy en Compose antes de levantar el servicio.

## Integración

- Kong ruta declarada en `infra/kong/kong.yml` como `/graphql/medical-history`.
- El servicio publica/consume eventos en Kafka si `KAFKA` está activado; revisar `src` para puntos de integración.

---

Documentado por automatización: pruebas E2E realizadas vía Kong (health/mutation/query) con éxito.
