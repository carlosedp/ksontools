# Build API and docs
ARG K8S_VERSION=1.14.7
FROM golang:1.13 AS builder

ARG K8S_VERSION
ENV SWAGGER_VERSION=$K8S_VERSION

RUN mkdir -p /go/src/github.com/bryanl/woowoo

WORKDIR /go/src/github.com/bryanl/woowoo
COPY . .

RUN go get github.com/bryanl/ksgen && \
    go install github.com/bryanl/woowoo/cmd/kslibdocgen

RUN ksgen -tag $SWAGGER_VERSION -output /tmp && \
    kslibdocgen -path /tmp/k8s.libsonnet -outPath /go/src/github.com/bryanl/woowoo/k8sdocs

# Build site
FROM node:8.16.1 as site

ARG K8S_VERSION
ENV SWAGGER_VERSION=$K8S_VERSION
ENV HUGO_VERSION 0.49.2
ENV HUGO_BINARY hugo_${HUGO_VERSION}_Linux-64bit.deb
ADD https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/${HUGO_BINARY} /tmp/hugo.deb

RUN dpkg -i /tmp/hugo.deb && \
    rm /tmp/hugo.deb && \
    mkdir -p /go/src/github.com/bryanl/woowoo

WORKDIR /go/src/github.com/bryanl/woowoo/k8sdocs/
COPY --from=builder /go/src/github.com/bryanl/woowoo/k8sdocs/ .

RUN npm install gulp-cli -g && \
    npm install gulp -D && \
    npm install @primer/octicons && \
    npm install

RUN gulp scss && \
    gulp images && \
    gulp js

RUN sed -i "s/SWAGGER_VERSION/$SWAGGER_VERSION/g" layouts/_default/baseof.html

# RUN hugo
RUN hugo

RUN mkdir -p public/octicons && \
    cp icons.data.svg.css public/octicons

# Build docs container
FROM nginx:1.17
LABEL MAINTAINER="carlosedp@gmail.com"

COPY --from=site /go/src/github.com/bryanl/woowoo/k8sdocs/public /usr/share/nginx/html
WORKDIR /usr/share/nginx/html