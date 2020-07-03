FROM golang:alpine as builderacc

WORKDIR /app

COPY /accountservice .

ARG CGO_ENABLED=0

RUN go mod tidy
RUN go build -o accountservice-linux-amd64

FROM scratch

EXPOSE 6767

WORKDIR /app

COPY --from=builderacc /app/accountservice-linux-amd64 .

ENTRYPOINT [ "/app/accountservice-linux-amd64"]

