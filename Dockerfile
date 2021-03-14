FROM golang:alpine AS builder
WORKDIR /app
ADD ./ /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags musl -a -o ugpden

FROM alpine
COPY --from=builder /app/ugpden /usr/local/bin/ugpden
RUN addgroup -S ugpden && adduser -S -G ugpden ugpden
USER ugpden
LABEL maintainer="Alejandro Galue <agalue@opennms.org>" name="udpgen"
ENTRYPOINT [ "/usr/local/bin/ugpden" ]
