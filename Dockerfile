FROM golang:1.22.8 as builder
WORKDIR /app
COPY . .
RUN go build -o cloudrun cmd/api/main.go

FROM ubuntu
WORKDIR /app
COPY --from=builder /app/cloudrun /app/cloudrun
CMD [ "/app/cloudrun" ]

