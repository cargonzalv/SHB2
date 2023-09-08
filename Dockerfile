FROM golang:alpine AS build

WORKDIR /go/src/srv
COPY . /go/src/srv

RUN apk --no-cache add openssh=~9 git=~2

RUN --mount=type=ssh set -x \
    && mkdir 0600 ~/.ssh \
    && ssh-keyscan github.com >> ~/.ssh/known_hosts

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

RUN --mount=type=ssh go get -d -v ./...

RUN go build -o /go/bin/srv ./cmd/app/main.go \
    && mkdir /logs

# Now copy it into our base image.
# FROM gcr.io/distroless/base-debian11
FROM alpine:3

ARG cert_location=/usr/local/share/ca-certificates
RUN apk --no-cache add ca-certificates=20230506-r0 openssl=~3 && update-ca-certificates

# Get certificate from "proxy.golang.org"
RUN sh -c ' \
    set -euxo pipefail && \
    openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/storage.googleapis.com.crt &&\
    update-ca-certificates &&\
    set -x; \
    adduser -u 1001 --disabled-password --home /app app \
'

# Update certificates
# RUN update-ca-certificates &&\
    # set -x; \
 # adduser -u 1001 --disabled-password --home /app app

USER app
WORKDIR /app

COPY --chown=1001:1001 --from=build /go/bin/srv .
COPY --chown=1001:1001 --from=build /go/src/srv/config/config.yml ./config/config.yml
COPY --chown=1001:1001 --from=build /go/src/srv/.env ./.env
COPY --chown=1001:1001 --from=build /logs ./logs

# Docker build arguments to store information about
# the Erlang build (version of OTP, git branch, etc.)
# into the Docker image.
ARG DELIVERY_BUILD_DATE
ARG DELIVERY_GIT_BRANCH
ARG DELIVERY_GIT_TAG
ARG DELIVERY_GIT_PRIV_TAG
ARG DELIVERY_GIT_HASH
ARG DELIVERY_HELM_TAG
ARG DELIVERY_ADDITIONAL_TAGS

ARG SERVICE_NAME
ENV BUILD_DATE $DELIVERY_BUILD_DATE
ENV GIT_BRANCH $DELIVERY_GIT_BRANCH
ENV GIT_TAG $DELIVERY_GIT_TAG
ENV GIT_PRIV_TAG $DELIVERY_GIT_PRIV_TAG
ENV GIT_HASH $DELIVERY_GIT_HASH
ENV HELM_TAG $DELIVERY_HELM_TAG
ENV ADDITIONAL_TAGS $DELIVERY_ADDITIONA
ENV SERVICE_NAME $SERVICE_NAME
ENV HB_APP_NAME $DELIVERY_HB_APP_NAME
ENV HB_APP_VERSION $DELIVERY_GIT_TAG

LABEL git-hash=${GIT_HASH} \
      git-tag=${GIT_TAG} \
      git-priv-tag=${GIT_PRIV_TAG} \
      build-date=${BUILD_DATE} \
      branch-name=${GIT_BRANCH} \
      helm-tag=${HELM_TAG} \
      additional-tags=${ADDITIONAL_TAGS} \
      appname=${SERVICE_NAME} \
      maintainer=sre@adgear.com

CMD ["/app/srv"]