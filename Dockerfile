FROM golang:latest AS build
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/gophermart/bin/main ./cmd/gophermart/

FROM alpine:latest
WORKDIR /gophermart
RUN mkdir /gophermart/logs
COPY --from=build /build/cmd/gophermart/bin/main .
COPY --from=build /build/migrations /gophermart/migrations
CMD ["/gophermart/main"]