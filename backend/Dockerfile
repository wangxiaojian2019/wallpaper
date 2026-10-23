# Production runtime stage for Go backend
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

COPY wallpaper-server-linux /app/wallpaper-server
RUN chmod +x /app/wallpaper-server

EXPOSE 8080

CMD ["/app/wallpaper-server"]
