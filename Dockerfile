FROM golang:1.15

ENV GO111MODULE=on

WORKDIR /app/server
COPY go.mod .
COPY go.sum .

RUN go get github.com/cespare/reflex

RUN go mod download
COPY . .

RUN go build 
# CMD ["./server"]

#CMD make watch
