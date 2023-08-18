#обычный докер
#FROM golang:1.19-alpine
#WORKDIR /app
#COPY . .
#RUN go build -o main main.go
#RUN apk add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
#COPY app/migrate.linux-amd64 ./migrate

#COPY db/migrate ./migrate

#EXPOSE 8080
#CMD [ "/app/main" ]
#ENTRYPOINT [ "/app/start.sh" ]

#легковесный докер фаил
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./db/migration


EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]