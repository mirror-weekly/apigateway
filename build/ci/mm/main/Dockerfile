FROM golang:1.15-buster AS build

WORKDIR /api-gateway

COPY . .

RUN apk add --update --no-cache make && \
    go get ./... && \
    make all

FROM debian:buster

WORKDIR /api-gateway

COPY --from=build /api-gateway/bin/ .

EXPOSE 8080
CMD ["./apigateway"]
