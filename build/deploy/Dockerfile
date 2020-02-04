FROM golang:1.13-alpine3.10 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/go-postgres
COPY cmd ./cmd
COPY pkg ./pkg
COPY go.mod ./go.mod
COPY go.sum ./go.sum
RUN go mod vendor
RUN go build -o /go/bin/api ./cmd/api/main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["api"]
