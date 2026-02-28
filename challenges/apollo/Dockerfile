FROM node:24 AS builder

WORKDIR /app

COPY challenges/apollo/package*.json ./
RUN npm ci

COPY challenges/apollo/ .

RUN npm run build

FROM node:24-slim

WORKDIR /app

COPY challenges/apollo/package*.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/dist ./dist

EXPOSE 4000

CMD [ "node", "dist/index.js" ]
