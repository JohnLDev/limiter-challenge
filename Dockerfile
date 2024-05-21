FROM golang:1.22-alpine as builder

WORKDIR /go/app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main -ldflags="-w -s" ./cmd/main.go

FROM scratch as runner

WORKDIR /go/app

COPY --from=builder /go/app/main .
COPY --from=builder /go/app/.env .
COPY tokens.json .

ENTRYPOINT [ "./main" ]

CMD [""]