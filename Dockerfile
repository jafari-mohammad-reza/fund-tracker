FROM golang:1.20.5 as dev
ENV GOPROXY=https://goproxy.io,direct
ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/cespare/reflex@latest

EXPOSE 5000

ENTRYPOINT ["make"]
CMD ["dev"]
FROM golang:1.20.5 as prod
ENV GOPROXY=https://goproxy.io,direct
ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


EXPOSE 5000

ENTRYPOINT ["make"]
CMD ["run"]
