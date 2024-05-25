# Builder state
FROM golang:1.22-alpine AS builder
RUN apk add --update make git curl

ARG MODULE_NAME=2024_1_kayros
WORKDIR /home/${MODULE_NAME}

# copy go.mod and go.sum in the root
COPY ../../../../go.mod ./
COPY ../../../../go.sum ./

# copy auth microservice directory (dirs cmd and internal)
COPY ../../../../microservices/session microservices/session
COPY ../../../../microservices/metrics  microservices/metrics 

# copy services (dbsql, s3, in-memory)
COPY ../../../../services/redis services/redis

# copy proto interfaces and data types
COPY ../../../../gen/go/session gen/go/session

# copy files configuration
COPY ../../../../config config

# copy internal dependences
COPY ../../../../internal/utils/myerrors internal/utils/myerrors
COPY ../../../../internal/middleware/grpc/server internal/middleware/grpc/server
COPY ../../../../internal/utils/constants internal/utils/constants

# building exe ile
RUN go build microservices/session/cmd/main.go

# Production state
FROM alpine:latest as production
WORKDIR /root/
ARG MODULE_NAME_BUILDER=2024_1_kayros

COPY --from=builder /home/${MODULE_NAME_BUILDER}/config/config.yaml config/config.yaml
COPY --from=builder /home/${MODULE_NAME_BUILDER}/main .

RUN chown root:root main

CMD ["./main"]