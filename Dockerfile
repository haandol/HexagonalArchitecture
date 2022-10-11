FROM golang:1.19.2 AS builder

# set environment path
ENV PATH /go/bin:$PATH
ENV GONOSUMDB github.com/haandol
ENV GOPRIVATE github.com/haandol

# create ssh directory
RUN mkdir ~/.ssh
RUN touch ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

WORKDIR /temp
COPY . /temp

ARG BUILD_TAG
ARG APP_NAME
RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.BuildTag=$BUILD_TAG -s" -o /go/bin/app ./cmd/${APP_NAME}

FROM golang:1.19.2-alpine AS server
ARG GIT_COMMIT=undefined
LABEL git_commit=$GIT_COMMIT

WORKDIR /
COPY --chown=0:0 --from=builder /go/bin/app /
COPY --chown=0:0 --from=builder /temp/xray.json /
COPY --chown=0:0 --from=builder /temp/docker-entrypoint.sh /
COPY --chown=0:0 --from=builder /temp/env/dev.env /.env

ARG APP_PORT
EXPOSE ${APP_PORT}

ENTRYPOINT ["/docker-entrypoint.sh"]