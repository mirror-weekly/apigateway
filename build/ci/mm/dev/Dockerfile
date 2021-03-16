FROM golang:1.15.6-alpine3.12 AS build

WORKDIR /api-gateway

COPY . .

RUN apk add --update --no-cache make && \
    go get ./... && \
    make all

FROM alpine:3.12

WORKDIR /api-gateway

COPY --from=build /api-gateway/bin/ .
COPY firebaseCredential.json .
COPY config.yaml .

EXPOSE 8080
CMD ["./apigateway"]
