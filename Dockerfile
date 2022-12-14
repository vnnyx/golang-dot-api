FROM golang:1.19-alpine as builder
WORKDIR /builder
COPY . .
RUN apk add upx
RUN go mod tidy
RUN go build \
    -ldflags "-s -w" \
    -o /builder/cmd/app/main /builder/cmd/app/main.go
RUN upx -9 /builder/cmd/app/main

FROM alpine:latest
ENV APP_PORT=9090
WORKDIR /app
COPY --from=builder /builder/cmd/app/main .
COPY --from=builder /builder/.env .
EXPOSE ${APP_PORT}
CMD /app/main