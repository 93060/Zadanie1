FROM golang:1.18 as gobuilder
WORKDIR /app
COPY server.go ./
COPY go.mod ./
COPY setup.sh ./
COPY upx-3.96-amd64_linux.tar.xz ./
RUN bash setup.sh && \
CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server && \
apt-get update && \
apt-get install xz-utils && \
tar -C /usr/local -xf upx-3.96-amd64_linux.tar.xz && \
/usr/local/upx-3.96-amd64_linux/upx --ultra-brute --overlay=strip ./server

FROM scratch as main
LABEL Autor: "Sebastian Wiktor"
COPY --from=gobuilder /app/server /
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8082
ENTRYPOINT [ "/server" ]
