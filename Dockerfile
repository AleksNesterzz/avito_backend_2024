FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x ./scripts/wait-for-postgres.sh

RUN go mod download
RUN go build -o go-app ./cmd/main.go
CMD ["./go-app"]