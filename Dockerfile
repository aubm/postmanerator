FROM golang:latest
WORKDIR /go/src/github.com/aubm/postmanerator

COPY Gopkg.toml .
COPY Gopkg.lock .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only
COPY . /go/src/github.com/aubm/postmanerator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o postmanerator .


FROM alpine:3.6

ARG http_proxy
ARG https_proxy
ARG verify_ssl=y

ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
ENV verify_ssl=$verify_ssl

RUN apk update \
 && apk add ca-certificates wget git \
 && update-ca-certificates \
 && mkdir -p /root/.postmanerator/themes \
 && cd /root/.postmanerator/themes \
 && if [ "${verify_ssl}" = "n" ]; \
    then git config --global http.sslVerify "false"; \
    fi \
 && git clone https://github.com/aubm/postmanerator-default-theme.git default \
 && git clone https://github.com/zanaca/postmanerator-hu-theme.git hu \
 && git clone https://github.com/aubm/postmanerator-markdown-theme.git markdown

COPY --from=0 /go/src/github.com/srgrn/postmanerator/postmanerator /usr/bin/

ENTRYPOINT ["postmanerator"]
CMD ["-collection", "/usr/var/collection.json"]
