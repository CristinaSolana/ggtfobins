FROM golang:1.13.0 as builder

ENV GO111MODULE=on
WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ggtfobins .

ENTRYPOINT ["/src/ggtfobins"]
