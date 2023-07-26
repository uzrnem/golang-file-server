# docker build . -t uzrnem/file-server:0.2 --no-cache
FROM golang:1.21rc2-alpine3.18 as dev

LABEL image_name="FileServer"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd cmd

RUN go build -o fileServer cmd/app/main.go

FROM alpine:3.18.2 as prod

WORKDIR /app

COPY --from=dev /app/fileServer /app/fileServer
COPY public /app/public

EXPOSE 9060

CMD [ "./fileServer" ]
