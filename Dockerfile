FROM golang:1.22.8 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o cloud-run

FROM scratch
WORKDIR /app
COPY --from=builder /app/cloud-run .
ENTRYPOINT [ "./cloud-run" ]

