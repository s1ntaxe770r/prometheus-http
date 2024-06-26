FROM golang:latest as builder

WORKDIR /app
COPY ./go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o prom-http

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/prom-http /prom-http

EXPOSE 8080

ENTRYPOINT ["/prom-http"]
