import express from "express";
import dotenv from "dotenv";
import { connectDB } from "./config/db.js";
import typeDefs from "./graphql/schema/index.js";

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

const startServer = async () => {
    try {
        await connectDB();

        app.listen(PORT, () => {
            console.log(`Medical History Service running on port ${PORT}`);
        });
    } catch (error) {
        console.error("Failed to start server:", error.message);
        process.exit(1);
    }
};

startServer();