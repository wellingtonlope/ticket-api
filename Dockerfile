FROM golang:1.25-alpine AS builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build ./cmd/main.go

FROM alpine:3
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
