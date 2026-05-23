FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/user-service ./cmd/server

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
COPY --from=build /bin/user-service /bin/user-service
USER app

EXPOSE 8080

ENTRYPOINT ["/bin/user-service"]
