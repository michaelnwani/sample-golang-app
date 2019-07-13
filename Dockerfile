FROM golang

MAINTAINER Michael Nwani <kmichael24@gmail.com>

ENV AWS_ACCESS_KEY_ID='<redacted>'
ENV AWS_SECRET_ACCESS_KEY='<redacted>'
ENV AWS_REGION='us-east-1'

COPY ./ /go/src/anon/
WORKDIR /go/src/anon/

RUN mkdir /var/log/anon/

RUN go get ./

RUN go build

EXPOSE 8080

CMD anon
