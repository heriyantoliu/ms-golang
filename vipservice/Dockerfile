FROM golang:alpine as buildervip

WORKDIR /app

COPY . .

ARG CGO_ENABLED=0

RUN go mod tidy
RUN go build -o vipservice-linux-amd64

FROM scratch

EXPOSE 6868

WORKDIR /app

COPY --from=buildervip /app/vipservice-linux-amd64 .

ENTRYPOINT [ "/app/vipservice-linux-amd64","-configServerUrl=http://configserver:8888", "-profile=test", "-configBranch=P9"]

