---
## Sprawozdanie - zadanie 1
### Autor: Sebastian Wiktor 
---
### 1. Kod serwera

Serwer został napisany w języku `Go`. Kod programu wraz z komentarzami znajduje się w pliku [server.go](../main/server.go).

### 2. Dockerfile
```dockerfile
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
```
Plik Dockerfile wykorzystuje wieloetapową metodę budowania obrazu. Pierwsza warstwa odpowiedzialna jest za zbudowanie pliku wykonywalnego serwera. Wykorzystywany
jest tutaj prosty [skrypt](../main/skrypt.sh) napisany w bashu, który umożliwi przyszłe zbudowanie obrazu kontenera na różne architektury `amd64`, `arm/v7` oraz `arm64/v8`.
Następnie aby zmniejszyć rozmiar pliku wykonywalnego (a co za tym idzie również obrazu kontenera) wykorzystywany jest program
[UPX](../main/upx-3.96-amd64_linux.tar.xz). Dzięki niemu plik wykonywalny zmniejszył się z 5.1MB do 1.58MB. Kolejna warstwa to warstwa scratch. Serwer nie wymaga żadnych
pakietów czy programów więc scratch jest optymalnym rozwiązaniem. Aby docker zbudował obraz wymagane było dodanie [certyfikatów CA](../main/ca-certificates.srt).
Serwer uruchamiany jest na porcie 8082.  

### 3. Polecenia
**a.&ensp;Zbudowanie opracowanego obrazu kontenera:** 

&ensp;`DOCKER_BUILDKIT=1 docker build -t serwer .` 

**b.&ensp; Uruchomienie kontenera ze zbudowanym obrazem**

&ensp; `docker run -t --name serwer -p 8082:8082 serwer`

