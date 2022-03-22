FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY api ./api
COPY proto ./proto

RUN go build ./api/main.go

CMD [ "./main" ]

EXPOSE 8080