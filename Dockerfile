FROM golang:alpine AS builder

WORKDIR /build

COPY src/go.mod .
COPY src/go.sum .

RUN go mod download

COPY src .

RUN go build -o ./server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/server .

EXPOSE 3003

CMD [ "./server" ]