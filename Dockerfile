FROM golang:1.20.5

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go get -u github.com/cespare/reflex
EXPOSE 5000

ENTRYPOINT ["make"]
CMD ["dev"]
