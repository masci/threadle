FROM golang:1.16-alpine AS build

ENV CGO_ENABLED=0
WORKDIR /go/src/github.com/masci/threadle
COPY . /go/src/github.com/masci/threadle/
RUN go build


FROM alpine:3.12

COPY --from=build /go/src/github.com/masci/threadle/threadle /usr/bin/threadle
COPY threadle.yaml.example /threadle.yaml
EXPOSE 3060
CMD ["threadle", "-c /threadle.yaml"]