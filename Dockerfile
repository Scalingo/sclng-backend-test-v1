FROM golang:1.20
LABEL maintainer="Infrastructure Services Team <team-infrastructure-services@scalingo.com>"

RUN go install github.com/cespare/reflex@latest

WORKDIR $GOPATH/src/github.com/Scalingo/sclng-backend-test-v1

EXPOSE 5000

CMD $GOPATH/bin/sclng-backend-test-v1
