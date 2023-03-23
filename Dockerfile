FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

#copies everything
COPY . .

ENV PORT 8080

RUN go build

EXPOSE $PORT

CMD ["./saul"]