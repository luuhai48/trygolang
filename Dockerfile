ARG GO_VERSION=1.19
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN go build ./src/*.go -o /build


FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /build /app/build

EXPOSE 3333
ENTRYPOINT [ "/app/build" ]
