FROM golang:1.18-alpine AS build

WORKDIR /xyauth

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN apk add -U --no-cache ca-certificates \
 && CGO_ENABLED=0 go build -ldflags="-w -s" -o /server ./cmd/server/*.go

FROM scratch

ENV DOCKER_RUNNING=1
WORKDIR /

COPY --from=build /server /
COPY --from=build /xyauth/configs /configs
COPY --from=build /xyauth/web /web
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "/server" ]
