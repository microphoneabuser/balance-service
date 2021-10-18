FROM golang:1.17

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client
# install netcat
RUN apt-get -y install netcat

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-it.sh

# build go app
RUN go mod download
RUN go build -o balance-service ./cmd/main.go

CMD ["./balance-service"]