FROM golang:1.17.7-buster 

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /api

EXPOSE 8080

CMD ["/api"]