FROM golang:1.18-alpine AS build

WORKDIR /xyauth

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /server ./cmd/server/*.go \
 && CGO_ENABLED=0 go build -ldflags="-w -s" -o /database ./cmd/database/*.go

FROM scratch

ENV DOCKER_RUNNING=1
WORKDIR /

COPY --from=build /server /database /
COPY --from=build /xyauth/configs /configs
COPY --from=build /xyauth/web /web

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "/server" ]
