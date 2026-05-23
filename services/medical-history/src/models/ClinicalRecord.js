import mongoose from "mongoose";

const clinicalRecordSchema = new mongoose.Schema(
    {
        patientId: {
            type: String,
            required: true
        },
        bloodType: {
            type: String
        },
        allergies: {
            type: [String],
            default: []
        },
        chronicDiseases: {
            type: [String],
            default: []
        },
        medications: {
            type: [String],
            default: []
        },
        notes: {
            type: [String],
            default: []
        }
    },
    {
        timestamps: true
    }
);

const ClinicalRecord = mongoose.model("ClinicalRecord", clinicalRecordSchema);

export default ClinicalRecord;