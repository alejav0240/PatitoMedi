FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download || true

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /bin/video-call ./cmd

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app
COPY --from=build /bin/video-call /bin/video-call
USER app

EXPOSE 8080

ENTRYPOINT ["/bin/video-call"]
