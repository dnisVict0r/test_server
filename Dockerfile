FROM golang:latest AS builder

RUN mkdir -p /go/src/server

#RUN apt-get update && apt-get install -y curl

WORKDIR /go/src/server
#ведение логов
#ARG LOG_DIR=/go/src/server/logs

#RUN mkdir -p ${LOG_DIR}

#ENV LOG_FILE_PATH=${LOG_DIR}/server.Log

COPY go.mod go.sum ./
#сторонние библиотеки
RUN go mod download 
COPY server.go ./
#оптимизация сборки go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .
#RUN CGO_ENABLED=0 GOOS=Linux go build -a -installsuffix cgo -o server .
#этап сборки контейнера на alpine
FROM alpine:latest
# установка сертификатов 
RUN apk --no-cache add ca-certificates
#этап сборки scratch
#FROM scratch
# добавляем сертификаты 
#ADD ca-certificates.crt /etc/ssl/certs/
WORKDIR /root/

COPY --from=builder /go/src/server .

EXPOSE 5000

CMD ["./server"]
