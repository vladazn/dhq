FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY service ./service
COPY proto ./proto

RUN go build ./service/main.go

CMD [ "./main" ]

EXPOSE 8080