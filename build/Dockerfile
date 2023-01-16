FROM golang:1.18-alpine AS build

WORKDIR /xyauth

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./

RUN go build -o /server ./cmd/server/*.go \
 && go build -o /database ./cmd/database/*.go


FROM scratch

ENV DOCKER_RUNNING=1
WORKDIR /

COPY --from=build /server /server
COPY --from=build /database /database

EXPOSE 8443

ENTRYPOINT [ "/server" ]

