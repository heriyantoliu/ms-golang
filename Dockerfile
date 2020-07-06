FROM golang:alpine as builderacc

WORKDIR /app

COPY /accountservice .

ARG CGO_ENABLED=0

RUN go mod tidy
RUN go build -o accountservice-linux-amd64

FROM golang:alpine as builderhealth

WORKDIR /app

COPY /healthchecker .

ARG CGO_ENABLED=0

RUN go mod tidy
RUN go build -o healthchecker-linux-amd64

FROM scratch

EXPOSE 6767

WORKDIR /app

COPY --from=builderacc /app/accountservice-linux-amd64 .
COPY --from=builderhealth /app/healthchecker-linux-amd64 .

HEALTHCHECK --interval=1s --timeout=3s --start-period=40s CMD ["./healthchecker-linux-amd64", "-port=6767"] || exit 1

ENTRYPOINT [ "/app/accountservice-linux-amd64","-configServerUrl=http://configserver:8888", "-profile=test", "-configBranch=P8"]

