FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

COPY cert/fullchain.pem /app/cert/fullchain.pem
COPY cert/privkey.pem /app/cert/privkey.pem

RUN go build -o main ./main.go

RUN chmod +x main

EXPOSE 434

CMD [ "./main" ]