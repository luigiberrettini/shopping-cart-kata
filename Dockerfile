FROM golang:alpine
WORKDIR $GOPATH/src/shopping-cart-kata
COPY . .
RUN set -ex; \
    apk update; \
    apk add --no-cache git; \
    go get -d -v ./... && \
    CGO_ENABLED=0 GOOS=linux go test ./... -timeout=60s -parallel=4 && \
    CGO_ENABLED=0 GOOS=linux go install -v ./... && \
    mkdir /cartTools && \
    cp $GOPATH/bin/* /cartTools

FROM scratch
COPY --from=0 /cartTools/* /
EXPOSE 8000
WORKDIR /
CMD ["/cartsvc"]