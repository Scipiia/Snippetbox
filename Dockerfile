#обычный докер
FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD [ "/app/main" ]

#легковесный докер фаил
# FROM golang:1.19-alpine AS builder
# WORKDIR /app
# COPY . .
# RUN go build -o main main.go

# FROM alpine
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY app.env .


# EXPOSE 8080
# CMD [ "/app/main" ]
