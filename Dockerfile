FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o browscap-go main.go

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/browscap-go .

RUN chown -R appuser:appgroup /app

USER appuser

ENTRYPOINT ["./browscap-go"]
