FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

#COPY ./docs .
COPY . .
#COPY .env .
#COPY ./config/config.yml .
RUN go build -o main .
EXPOSE 3000
# Запускаем приложение
CMD ["/app/main"]