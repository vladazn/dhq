FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY swagger ./swagger

RUN go build ./swagger/main.go

CMD [ "./main" ]

EXPOSE 8080