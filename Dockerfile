# stage 1 # - build
FROM golang:1.17.2-alpine3.14 as ws_build

ENV GO111MODULES=on

WORKDIR /unity-web-service

COPY . .

RUN go mod tidy && go build -o unity-ws

# stage 2 # - exe
FROM alpine:3.14

COPY --from=ws_build /unity-web-service/unity-ws .

EXPOSE 8080

CMD ["./unity-ws"]