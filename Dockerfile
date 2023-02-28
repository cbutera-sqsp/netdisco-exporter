# builder image
FROM golang:1.19 AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/netdisco-exporter 
RUN pwd
RUN ls

# container image
FROM alpine:3.17
RUN adduser -D netdisco-exporter
RUN apk --no-cache add ca-certificates

COPY --from=builder /build/netdisco-exporter /bin/netdisco-exporter
USER netdisco-exporter
CMD ["/bin/netdisco-exporter"]