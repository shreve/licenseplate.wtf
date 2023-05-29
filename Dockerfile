FROM golang:1.20 AS build

ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini-static /tini
RUN chmod +x /tini

WORKDIR /app

RUN apt-get update && apt-get install -y tini musl-dev

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /lp ./cmd/lp/

EXPOSE 8081
ENTRYPOINT [ "/tini", "--" ]
CMD [ "/lp", "serve" ]
