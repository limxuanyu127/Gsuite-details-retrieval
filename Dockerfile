FROM golang:latest

RUN mkdir -p /go/src/GsuiteRetrieval
RUN mkdir -p /go/src/GsuiteRetrieval/controllers

COPY main.go /go/src/GsuiteRetrieval/
COPY controllers /go/src/GsuiteRetrieval/controllers

WORKDIR /go/src/GsuiteRetrieval

RUN go get -u google.golang.org/api/admin/reports/v1
RUN go get -u google.golang.org/api/sheets/v4
RUN go get -u golang.org/x/oauth2/google
RUN go get -u google.golang.org/api/admin/directory/v1
RUN go get github.com/sendgrid/sendgrid-go
RUN go build -o main .

CMD ["./main"]