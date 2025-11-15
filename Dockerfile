FROM golang:1.25.4 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ekko-bot cmd/bot/main.go

FROM alpine:3.22.0 AS run

WORKDIR /app
COPY --from=build /app/ekko-bot .
COPY db/ ./db/

CMD ["./ekko-bot"]
