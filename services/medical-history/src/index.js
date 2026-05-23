import express from "express";
import dotenv from "dotenv";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 4000;

app.use(express.json());

app.get("/health", (req, res) => {
    res.status(200).json({
        status: "ok",
        service: "medical-history"
    });
});

app.listen(PORT, () => {
    console.log(`Medical History Service running on port ${PORT}`);
});