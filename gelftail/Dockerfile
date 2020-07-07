FROM golang:alpine as builder

WORKDIR /app

COPY . .

ARG CGO_ENABLED=0

RUN go mod tidy
RUN go build -o gelftail-linux-amd64

FROM scratch

EXPOSE 12202/udp

WORKDIR /app

COPY --from=builder /app/gelftail-linux-amd64 .
COPY token.txt .

ENTRYPOINT [ "/app/gelftail-linux-amd64","-port=12202"]

