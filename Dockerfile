FROM golang:1.15-alpine AS builder
ADD . /go/
RUN addgroup -g 1010 api && adduser -D -G api -u 1010 api && \
    apk update && apk add --no-cache git && \
    go get -d -v && mkdir -p /app/config && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/pgapi .
COPY titanic.csv /app/migrations/


FROM scratch as app
LABEL app=pgapi
LABEL version=0.1
EXPOSE 8080:8080/tcp
USER 1010:1010
COPY --from=builder /etc/passwd /etc/group /etc/
COPY --from=builder /app /app
CMD ["/app/pgapi"]


