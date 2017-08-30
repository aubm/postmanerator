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
 && git clone https://github.com/aubm/postmanerator-markdown-theme.git markdown \
 && cd /usr/bin/ \
 && if [ "${verify_ssl}" = "n" ]; \
    then wget -O postmanerator https://github.com/aubm/postmanerator/releases/download/v0.8.0/postmanerator_linux_386 --no-check-certificate; \
    else wget -O postmanerator https://github.com/aubm/postmanerator/releases/download/v0.8.0/postmanerator_linux_386; \
    fi \
 && chmod +x postmanerator

ENTRYPOINT ["postmanerator"]
CMD ["-collection", "/usr/var/collection.json"]
