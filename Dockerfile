FROM golang:1.17.4-alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build ./framework/rest/server.go

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/server /app/
WORKDIR /app
CMD ["./server"]