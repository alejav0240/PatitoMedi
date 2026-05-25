# Insomnia — Medical History (ejemplos)

Instrucciones rápidas:

- URL base (entorno de Insomnia): `http://localhost:8000/graphql/medical-history`
- Método: `POST`
- Header: `Content-Type: application/json`
- Body: tipo `JSON` con campo `query` y opcional `variables`.

Ejemplos listos para pegar en el body de una request en Insomnia.

1) Health check

Request body (JSON):

```json
{ "query": "query { healthGraphQL }" }
```

2) Crear registro clínico (mutation)

Request body (con `variables`):

```json
{
  "query": "mutation($input: ClinicalRecordInput!){ createRecord(input:$input){ id patientId bloodType notes createdAt } }",
  "variables": {
    "input": {
      "patientId": "00000000-0000-0000-0000-000000000001",
      "bloodType": "O+",
      "allergies": ["penicilina"],
      "chronicDiseases": ["hipertensión"],
      "medications": ["enalapril 10mg"],
      "notes": ["Paciente estable"]
    }
  }
}
```

3) Actualizar registro (mutation)

Request body (con `variables`):

```json
{
  "query": "mutation($id: ID!, $input: ClinicalRecordInput!){ updateRecord(id:$id, input:$input){ id updatedAt } }",
  "variables": {
    "id": "64a1b2c3d4e5f6a7b8c9d0e1",
    "input": {
      "patientId": "00000000-0000-0000-0000-000000000001",
      "medications": ["enalapril 20mg"],
      "notes": ["Dosis ajustada"]
    }
  }
}
```

4) Consultar registros por paciente (query)

Request body (con `variables`):

```json
{
  "query": "query($patientId: String!){ patientRecords(patientId:$patientId){ id patientId bloodType allergies chronicDiseases medications notes createdAt updatedAt } }",
  "variables": {
    "patientId": "00000000-0000-0000-0000-000000000001"
  }
}
```

5) Obtener registro por ID (query)

Request body (con `variables`):

```json
{
  "query": "query($id: ID!){ record(id:$id){ id patientId bloodType allergies chronicDiseases medications notes createdAt updatedAt } }",
  "variables": {
    "id": "64a1b2c3d4e5f6a7b8c9d0e1"
  }
}
```

Consejos para Insomnia:

- Crea un environment con la variable `base_url` y pon `http://localhost:8000/graphql/medical-history`.
- Crea una petición POST y usa `{{ base_url }}` como URL.
- Pega cualquiera de los bodies anteriores en la pestaña `Body -> JSON`.
- Para probar variantes, modifica `variables` en la pestaña `Body` sin tocar `query`.

Si quieres, puedo exportar esto como una colección de Insomnia `.json` lista para importar. ¿La genero? 
