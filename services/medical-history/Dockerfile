FROM node:20-alpine

WORKDIR /app

# Copiar package.json primero para aprovechar cache de capas
COPY package*.json ./

RUN npm ci --omit=dev || true

COPY . .

CMD ["npm", "start"]
