FROM golang:1.23.0

EXPOSE 8080

WORKDIR /usr/local/app

COPY . .

RUN go build main.go

CMD ["/usr/local/app/main"]