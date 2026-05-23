import express from "express";
import dotenv from "dotenv";
import { ApolloServer } from "@apollo/server";
import { expressMiddleware } from "@as-integrations/express4";
import cors from "cors";

import { connectDB } from "./config/db.js";
import typeDefs from "./graphql/schema/index.js";
import resolvers from "./graphql/resolvers/index.js";

dotenv.config();

const app = express();
const PORT = process.env.PORT || 4000;

async function startServer() {
    try {
        // Conectar a MongoDB
        await connectDB();

        // Crear Apollo Server
        const server = new ApolloServer({
            typeDefs,
            resolvers,
        });

        await server.start();

        // Middlewares
        app.use(cors());
        app.use(express.json());

        // Health check
        app.get("/health", (req, res) => {
            res.status(200).json({
                status: "ok",
                service: "medical-history",
            });
        });

        // GraphQL endpoint
        app.use(
            "/graphql/medical-history",
            expressMiddleware(server, {
                context: async ({ req }) => ({ req }),
            })
        );

        // Arrancar servidor HTTP
        app.listen(PORT, () => {
            console.log(`Medical History Service running on port ${PORT}`);
            console.log(`GraphQL ready at /graphql/medical-history`);
        });
    } catch (error) {
        console.error("Failed to start server:", error);
        process.exit(1);
    }
}

startServer();