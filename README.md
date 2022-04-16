---
## Sprawozdanie - zadanie 1
### Autor: Sebastian Wiktor 
---
### CZĘŚĆ OBOWIĄZKOWA

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
[UPX](../main/upx-3.96-amd64_linux.tar.xz). Dzięki niemu plik wykonywalny zmniejszył się z 5.1MB do 1.58MB. 
UPX sam wykrywa architekturę na której jest używany, więc budowanie obrazów na platformy inne niż x86-64 nie sprawia żadnych problemów i możemy użyć tego samego pliku .tar (chociaż sama nazwa pliku może być myląca). Kolejna warstwa
to warstwa scratch. Serwer nie wymaga żadnych
pakietów czy programów więc scratch jest optymalnym rozwiązaniem. Aby docker zbudował obraz jedynym wymaganiem było dodanie [certyfikatów CA](../main/ca-certificates.srt).    
Serwer uruchamiany jest na porcie 8082.  

### 3. Polecenia
**a.&ensp;Zbudowanie opracowanego obrazu kontenera:** 
```
$ DOCKER_BUILDKIT=1 docker build -t serwer . 
```

**b.&ensp; Uruchomienie kontenera ze zbudowanym obrazem**

```
$ docker run -t --name serwer -p 8082:8082 serwer
```

**c.&ensp; Działanie serwera**

Wchodząc w przeglądarkę i wpisując adres `localhost:8082` ukazuje nam się działająca strona uruchomionego serwera

![server](https://user-images.githubusercontent.com/103113980/163670490-eb6ce6eb-f3ea-4246-9d50-03e403f70fc1.png)

Aby uzyskać dostęp do logów zapisywanych przez serwer należy użyć adresu `localhost:8082/log`

![logs](https://user-images.githubusercontent.com/103113980/163670599-77e18c3c-e21f-4d97-869d-ed2f2a3996a6.png)

**`UWAGA:`** W logach zapisywanych przez serwer wyświetlana jest data i godzina zgodna dla domyślnego w kontenerze czasu wzorcowego UTC. Rozbieżność w godzinach wyświetlanych na głównej stronie i w logach wynika stąd, że adres IP dla którego wyświetlane są informacje znajduje się w strefie czasowej UTC+02:00.

**d.&ensp; Sprawdzenie ilości warstw w zbudowanym obrazie**

```
$ docker image history serwer
```

![layers](https://user-images.githubusercontent.com/103113980/163673672-109c8f89-1470-450d-9733-70d08ec67d90.png)

**Alternatywny sposób sprawdzenia ilości warstw oraz uzyskania innych informacji na temat zbudowanego obrazu**

```
$ docker image inspect serwer
```
![inspect](https://user-images.githubusercontent.com/103113980/163675691-a389046c-5486-44d0-bd6a-49a6bae6b1e2.png)

### 4. Budowanie obrazów na różne architektury

Aby było możliwe zbudowanie obrazów na różne platformy sprzętowe musimy skorzystać z zasobów emulatora QEMU. Na potrzeby wykonania tego zadania zainstalujemy QEMU lokalnie, ale można to zrobić w alternatywny sposób z wykorzystaniem dedykowanego kontenera. Następnie do zbudowania obrazów wykorzystamy wraper buildx. 

**Instalacja zasobów QEMU**

```
$ sudo apt-get install qemu-user-static
```

**Utworzenie nowego buildera buildx oraz ustawienie go jako domyślnego**

```
$ docker buildx create --name serwerbuilder
```

```
$ docker buildx use serwerbuilder
```

**Zbudowanie obrazu serwera na 3 wybrane platformy i przesłanie ich na repozytorium DockerHub**

```
$ docker buildx build -t 93060/zadanie1:multiplatform --platform linux/amd64,linux/arm64/v8,linux/arm/v7 --push . 
```

**Potwierdzenie poprawnego zbudowania obrazów - repozytorium DockerHub**

![dockerhub](https://user-images.githubusercontent.com/103113980/163679177-ead3cd75-67c2-4d1f-a13b-ceddc1c003a2.png)

Zbudowane obrazy można znaleźć na repozytorium DockerHub do którego link znajduje się [tutaj](https://hub.docker.com/r/93060/zadanie1/tags).
