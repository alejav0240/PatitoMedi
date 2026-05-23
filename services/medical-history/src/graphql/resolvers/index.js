import ClinicalRecord from "../../models/ClinicalRecord.js";

const resolvers = {
    Query: {
        healthGraphQL: () => "GraphQL is working",
        patientRecords: async (_, { patientId }) => {
            return await ClinicalRecord.find({ patientId });
        },
        record: async (_, { id }) => {
            return await ClinicalRecord.findById(id);
        },
    },
    Mutation: {
        createRecord: async (_, { input }) => {
            return await ClinicalRecord.create(input);
        },
        updateRecord: async (_, { id, input }) => {
            return await ClinicalRecord.findByIdAndUpdate(id, input, {
                new: true,
            });
        },
    },
};

export default resolvers;