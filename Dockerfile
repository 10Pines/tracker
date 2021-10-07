FROM golang:1.16-alpine3.14 as builder
RUN apk add --no-cache gcc musl-dev make
WORKDIR /var/build
COPY . .
RUN ["make", "test"]
RUN ["make", "build"]

FROM alpine:3.14
WORKDIR /var/tracker
ENV GIN_MODE=release
COPY --from=builder /var/build/server .
ENTRYPOINT ["./server"]
