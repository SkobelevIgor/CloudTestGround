FROM golang:latest

WORKDIR /go/src/app

ADD main.go .

RUN go build -o main . 

CMD ["./main"]

EXPOSE 8080
