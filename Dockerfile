FROM golang:1.20.5

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go get github.com/cespare/reflex
RUN go install github.com/cespare/reflex@latest

EXPOSE 5000

ENTRYPOINT ["make"]
CMD ["dev"]
