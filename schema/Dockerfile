FROM golang:1.20

RUN go version

#Значение для переменной окружения
ENV GOPATH=/

# Копирование файлов проекта в контейнер
COPY ./ ./

# Установка зависимостей
RUN go mod download
RUN go build -o avito-tech-service ./cmd/avito-tech-service/main.go

# Команда, выполняемая при запуске контейнера
CMD ["./avito-tech-service"]