FROM golang:1.20

RUN go version

# Value for environment variable
ENV GOPATH=/

# Copying Project Files to a Container
COPY ./ ./

# Installing dependencies
RUN go mod download
RUN go build -o avito-tech-service ./cmd/avito-tech-service/main.go

# The command to be executed when the container is starte
CMD ["./avito-tech-service"]