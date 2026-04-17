FROM golang:1.26.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o profile-service ./cmd/profile

EXPOSE 50051

CMD ["./profile-service"]
