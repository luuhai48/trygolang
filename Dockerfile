ARG GO_VERSION=1.19
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN go build -o build ./src/*.go


FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir /app
COPY --from=builder /app/build /app/build
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder --chmod=755 /app/entrypoint.sh /app/entrypoint.sh

WORKDIR /app
EXPOSE 3333
ENTRYPOINT [ "/app/entrypoint.sh" ]
