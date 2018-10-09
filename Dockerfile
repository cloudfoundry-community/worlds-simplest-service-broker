FROM golang:alpine
WORKDIR /go/src/worlds-simplest-service-broker
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worlds-simplest-service-broker .

FROM alpine:latest
EXPOSE 3000
ENV BASE_GUID=29140B3F-0E69-4C7E-8A35 \
    SERVICE_NAME=some-service-name \
    SERVICE_PLAN=shared \
    CREDENTIALS='{"port":"4000"}' \
    FAKE_ASYNC=false

RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=0 /go/src/worlds-simplest-service-broker/worlds-simplest-service-broker .
CMD ["./worlds-simplest-service-broker"]
