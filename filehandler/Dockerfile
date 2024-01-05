# docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v8 -t uzrnem/filehandler:latest --push . --no-cache
# docker build -t notebook:0.0.2 -t notebook:latest . --no-cache

FROM golang:1.21.2-alpine3.18 AS build
WORKDIR /temp
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build main.go

FROM alpine:3.19.0 AS release
WORKDIR /app
COPY --from=build /temp/main .
EXPOSE 8080
ENTRYPOINT ["/app/main"]

# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping