FROM golang:latest as builder
ENV GOPATH=/go
ENV GO111MODULE=on
WORKDIR ${GOPATH}/src/github.com/Shikugawa/pingoo
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o ./dist/pingoo -i ./main.go

FROM alpine:latest
RUN apk add --update --no-cache ca-certificates tzdata && update-ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/Shikugawa/pingoo/dist .
RUN chmod +x ./pingoo
EXPOSE 3000
ENTRYPOINT [ "/app/pingoo" ]
