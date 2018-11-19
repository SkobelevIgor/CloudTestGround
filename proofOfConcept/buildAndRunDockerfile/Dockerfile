FROM golang:latest

WORKDIR /go/src/app

ADD ./toBuild/main.go .

RUN go build -o main . 

CMD ["./main"]

EXPOSE 8080
