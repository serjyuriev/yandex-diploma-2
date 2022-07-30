FROM golang:1.18-alpine

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go build -o main ./cmd/go-keeper-server

CMD [ "/app/main", "-c", "dev_srv_config.yaml" ]