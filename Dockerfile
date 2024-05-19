FROM golang:alpine

WORKDIR /app

COPY . .

WORKDIR /app/telemetry

RUN CGO_ENABLED=0 GOOS=linux go build -o qr

CMD ["./qr"]
