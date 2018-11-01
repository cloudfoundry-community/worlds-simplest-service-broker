FROM golang:alpine as build
WORKDIR /go/src/github.com/cloudfoundry-community/worlds-simplest-service-broker
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo github.com/cloudfoundry-community/worlds-simplest-service-broker/cmd/worlds-simplest-service-broker

FROM alpine:latest as final
EXPOSE 3000
ENV BASE_GUID=29140B3F-0E69-4C7E-8A35 \
    SERVICE_NAME=some-service-name \
    SERVICE_PLAN=shared \
    CREDENTIALS='{"port":"4000"}' \
    TAGS=shared,worlds-simplest-service-broker \
    FAKE_ASYNC=false \
    FAKE_STATEFUL=false \
    AUTH_USER=broker \
    AUTH_PASSWORD=broker

RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=build /go/bin/worlds-simplest-service-broker .
CMD ["./worlds-simplest-service-broker"]
