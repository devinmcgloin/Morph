FROM golang

WORKDIR /go/src/github.com/sprioc/conductor
COPY . /go/src/github.com/sprioc/conductor

RUN go get github.com/tools/godep
RUN godep restore

RUN go install ./...

ENTRYPOINT /go/bin/sprioc-serve

ENV PORT 80
EXPOSE 80
