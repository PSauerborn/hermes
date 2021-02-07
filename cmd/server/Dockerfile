FROM golang:latest as build

ENV GO111MODULE=on

WORKDIR /app/server

# copy and install dependencies
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

# copy source code and build
COPY pkg ./pkg
COPY cmd ./cmd

RUN CGO_ENABLED=0 go build cmd/server/main.go

FROM alpine:latest as server

WORKDIR /app/server
# create directory for config settings
RUN mkdir /etc/hermes

COPY --from=build /app/server/main ./

RUN chmod +x ./main

CMD [ "./main" ]