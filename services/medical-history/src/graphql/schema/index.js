const typeDefs = `#graphql
  type ClinicalRecord {
    id: ID!
    patientId: String!
    bloodType: String
    allergies: [String!]
    chronicDiseases: [String!]
    medications: [String!]
    notes: [String!]
    createdAt: String
    updatedAt: String
  }

  input ClinicalRecordInput {
    patientId: String!
    bloodType: String
    allergies: [String!]
    chronicDiseases: [String!]
    medications: [String!]
    notes: [String!]
  }

  type Query {
    healthGraphQL: String
    patientRecords(patientId: String!): [ClinicalRecord!]!
    record(id: ID!): ClinicalRecord
  }

  type Mutation {
    createRecord(input: ClinicalRecordInput!): ClinicalRecord!
    updateRecord(id: ID!, input: ClinicalRecordInput!): ClinicalRecord!
  }
`;

export default typeDefs;