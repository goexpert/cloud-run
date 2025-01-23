FROM golang:1.22.8 as builder
WORKDIR /app
COPY . .
RUN go build -o cloudrun cmd/api/main.go && chmod +x cloudrun

FROM debian:latest
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /app
COPY --from=builder /app/cloudrun .
ENTRYPOINT [ "./cloudrun" ]

