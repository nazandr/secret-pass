FROM golang:1.19-alpine AS builder
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN go build -o ./out/secret-pass 

FROM alpine:3.15

COPY --from=builder /build/out/secret-pass /app/secret-pass
COPY --from=builder /build/assets /app/assets
WORKDIR /app
ENTRYPOINT ["/app/secret-pass"]
