FROM golang:1.22

RUN mkdir /app

ADD . /app

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go build -o main ./cmd/server/main.go

EXPOSE 50051

CMD ["/app/main"]