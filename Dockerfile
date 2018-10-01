FROM golang:1.11
COPY . /go/src/github.com/andreymgn/RSOI
WORKDIR /go/src/github.com/andreymgn/RSOI
RUN go install ./cmd/...