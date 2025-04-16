# syntax=docker/dockerfile:1
FROM alpine:3.21
LABEL author=sorc@sction.org
ARG TARGETOS
ARG TARGETARCH
COPY ./dist/minas_linux_${TARGETARCH} /usr/bin/minas
COPY ./rclone/rclone-v1.68.2-linux-${TARGETARCH}/rclone /usr/bin/rclone
COPY ./rclone/ca.crt /usr/local/share/ca-certificates/myrootca.crt
#RUN apk add ncat openssh-client
RUN apk add tzdata ca-certificates
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" >/etc/timezone && apk del tzdata && mkdir /app && chmod +x /usr/bin/minas && chmod +x /usr/bin/rclone && update-ca-certificates
WORKDIR /app
CMD ["minas","server"]
