FROM golang:stretch AS build-env

WORKDIR /go/src/github.com/gridiron-zone/gridiron

RUN apt update
RUN apt install git -y

COPY . .

RUN make build

FROM golang:stretch

RUN apt update
RUN apt install ca-certificates jq -y

WORKDIR /root

COPY --from=build-env /go/src/github.com/gridiron-zone/gridiron/build/gridirond /usr/bin/gridirond

EXPOSE 26656 26657 1317 9090

CMD ["gridirond"]