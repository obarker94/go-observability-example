FROM golang:alpine

WORKDIR /app

COPY . .

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -o s2

CMD ["./s2"]
