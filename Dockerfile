# service backend
FROM golang:1.22 AS building
RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /app/main cmd/main/main.go
CMD ["/app/main"]

#FROM alpine:3.13
#WORKDIR /usr/bin
#COPY --from=build /app/config/config.yaml /go/bin/config/coni
#COPY --from=build /app/bin /go/bin
#EXPOSE 8000
#ENTRYPOINT /go/bin/main --port 8000