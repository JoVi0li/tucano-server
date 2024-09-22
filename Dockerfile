FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main ./cmd/server/main.go

RUN chmod +x main

EXPOSE 443

CMD [ "./main" ]