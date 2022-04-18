---
# Sprawozdanie - zadanie 1
### Autor: Sebastian Wiktor 
---
## CZĘŚĆ OBOWIĄZKOWA

### 1. Kod serwera

Serwer został napisany w języku `Go`. Kod programu wraz z komentarzami znajduje się w pliku [server.go](../main/server.go).

---

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
pakietów czy programów więc scratch jest optymalnym rozwiązaniem. Aby docker zbudował obraz jedynym wymaganiem było dodanie [certyfikatów](../main/ca-certificates.srt).    
Serwer uruchamiany jest na porcie 8082.  

---

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

> ⚠️ `UWAGA:` W logach zapisywanych przez serwer wyświetlana jest data i godzina zgodna dla domyślnego w kontenerze czasu wzorcowego UTC. Rozbieżność w godzinach wyświetlanych na głównej stronie i w logach wynika stąd, że adres IP dla którego wyświetlane są informacje znajduje się w strefie czasowej UTC+02:00.

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

---
## CZĘŚĆ DODATKOWA

### **DODATEK 1**

### 1. Wykorzystanie GitHub Actions, cache oraz Github Container Registry

Proces budowania naszego obrazu możemy w bardzo prosty i szybki sposób zautomatyzować za pomocą `GitHub Actions`. Aby rozpocząć korzystanie z tego narzędzia, w górnym panelu naszego repozytorium na `GitHub` wybieramy `Actions`.

![actions](https://user-images.githubusercontent.com/103113980/163727087-05235fcf-989a-43b1-9672-8057829ab559.png)

W zakładce `Actions` tworzymy nowy `Workflow`. Możemy wybrać przepływ sugerowany przez GitHub, który jest wstępnie skonfiguroway pod `Docker Image` lub utworzyć własny. 

![workflow](https://user-images.githubusercontent.com/103113980/163727808-1c7114e5-781c-4778-aa05-6925c2473d24.png)

---

![create_yml](https://user-images.githubusercontent.com/103113980/163727946-f091067d-b471-46a7-aae7-090f9dbf73ce.png)

Po prawej stronie znajduje się przydatne narzędzie `Marketplace`, w któym możemy znaleźć gotowe funkcje np. logowanie się do `Dockera`, czy konfiguracja `QEMU`.
Informacje o tym, jak stworzyć przepływ i zautomatyzować proces budowania i publikowania obrazu możemy znależć w dokumentacji Dockera pod tym [linkiem](https://docs.docker.com/ci-cd/github-actions/).

Na podstawie dokumentacji, informacji znalezionych w internecie oraz przepływu wykonanego na zajęciach stworzyłem `Workflow`, który zbuduje obraz na 3 wybrane platformy a następnie opublikuje go na `GitHub Containers Registry`. Dodałem również zapis do pamięci `cache` (ostatnie linijki kodu), jednak na razie na potrzeby przetestowania działania tej funkcji zostały one zakomentowane. Kod przepływu znajduje się [tutaj](../.github/workflows/workflow.yml)

```yml
name: GitHub Actions workflow with push to GHCR

#uruchom przeplyw po kazdym pushu na branch main z wyjatkiem aktualizacji pliku README.md
on:
  push:
    branches:
      - main
    paths-ignore:
      - '**/README.md'
      
# zadania do wykonania - build i push obrazu na maszynie z ubuntu
jobs:
  build-push-images:
    name: Build and push to GHCR
    runs-on: ubuntu-latest
    
#kroki do wykonania 
    steps:
        # sprawdzenie poprawnosci kodu
      - name: Checkout code
        uses: actions/checkout@v2
        # uruchomienie QEMU w celu zbudowania obrazow na rozne platformy 
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v1
        # uruchomienie buildera buildx
      - name: Buildx set-up
        id: buildx
        uses: docker/setup-buildx-action@v1
        
        # logowanie do github container registry z wykorzystaniem zmiennej srodowiskowej 
        # GHCR_PASSWORD umieszczonej w secrets, w ktorej zapisany jest token dostępu dostępowy 
      - name: Login to GitHub
        uses: docker/login-action@v1 
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GHCR_PASSWORD }}

        # budowa i publikacja obrazow na GHCR
      - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          context: ./
          file: ./Dockerfile
          platforms: linux/amd64,linux/arm64/v8,linux/arm/v7
          push: true
          tags: ghcr.io/93060/zadanie1:latest
          # konfiguracja cache typu gha
          #cache-from: type=gha
          #cache-to: type=gha,mode=max
```      
> ⚠️ `UWAGA:` Informacje o tym jak wygenerować token dostępowy znajduję się pod tym [linkiem](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token) w dokumentacji GitHub. [Tutaj](https://github.com/Azure/actions-workflow-samples/blob/master/assets/create-secrets-for-GitHub-workflows.md) znajdziemy informację o tym jak wygenerowany token dodać do zmiennej środowiskowej w `Secrets` 

Po utworzeniu pliku workflow.yml w `Actions` od razu uruchomił się przepływ.

![start](https://user-images.githubusercontent.com/103113980/163730764-96440991-885b-4f00-9109-20f5af000540.png)

Po chwili dostajemy informację, że przepływ wykonał się prawidłowo. 

![success](https://user-images.githubusercontent.com/103113980/163731132-e56a0d67-cdff-4aae-8643-32a0c788aabb.png)

Wchodząc na nasz profil `GitHub`, a następnie w zakładkę `Packages` widzimy paczkę zadanie1, która zawiera zbudowane przez nas obrazy. 

![package](https://user-images.githubusercontent.com/103113980/163731226-49ed5193-6be9-4e80-ae86-64929ef5df22.png)

Jak widać konfiguracja GitHub Container Registry jest bardzo prosta i wymaga niewielu modyfikacji w porównaniu gdy obraz publikowaliśmy na `DockerHub` podczas jednych z zajęć. Modyfikacji musimy poddać następujące fragmenty kodu: 

```yml
 - name: Login to GitHub
        uses: docker/login-action@v1 
        with:
          registry: ghcr.io
          username: ${{ github.repository_owner }}
          password: ${{ secrets.GHCR_PASSWORD }}  
```

W akcji `uses: docker/login-action@v1` należy dodać linijkę `registry: ghcr.io` oraz zamiast danych do logowania na `DockerHub` podać dane do logowania na `GitHub`.

```
 - name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          context: ./
          file: ./Dockerfile
          platforms: linux/amd64,linux/arm64/v8,linux/arm/v7
          push: true
          tags: ghcr.io/93060/zadanie1:latest
```

W akcji `docker/build-push-action@v2` należy zmodyfikować linijkę `tags`: tutaj zamiast repozytorium na `DockerHub` podajemy ściezkę do naszego `GHCR`. 

---

### Testowanie działania zapisu do pamięci podręcznej cache

Zaczniemy od zbudowania kolejnej wersji naszych obrazów nie uruchamiając jeszcze zapisu do pamięci cache. Modyfikacji poddamy pierwszą linijkę kodu w pliku workflow.yml 

![cachev1](https://user-images.githubusercontent.com/103113980/163731787-16363a02-58de-448f-bf33-163f53870afa.png)

Widzimy, że pomimo kosmetycznych zmian czas wykonania przepływu był zbliżony do tego pierwszego, co oznacza, że `GitHub Actions` automatycznie nie cachuje żadnych danych. Teraz odkomentujemy linijki i uruchomimy zapis do pamięci `cache`.

```yml
- name: Build and push
        id: docker_build
        uses: docker/build-push-action@v2
        with:
          context: ./
          file: ./Dockerfile
          platforms: linux/amd64,linux/arm64/v8,linux/arm/v7
          push: true
          tags: ghcr.io/93060/zadanie1:latest
          # konfiguracja cache typu gha
          cache-from: type=gha
          cache-to: type=gha,mode=max
```  

Wykonujemy kolejny przepływ - ten trwał wyraźnie dłużej niż budowanie poprzednich obrazów. 

![cachev2](https://user-images.githubusercontent.com/103113980/163732375-8dcaebd2-b4fa-4b72-aa6b-5701d177cfb8.png)

Wchodząc w logi widzimy, że na dodatkowy czas wpłynął eksport danych do pamięci cache. 

![cachev3](https://user-images.githubusercontent.com/103113980/163732378-e9247669-1c38-44f9-a257-c65c3fa0774b.jpg)

Uruchamiamy ostatni przepływ. Także tym razem wprowadzamy tylko kosmetyczną zmianę w pierwszej linijce pliku workflow.yml.
Widać ogromną różnicę - wykonanie przepływu trwało tylko `41 sekund`, około `6 razy szybciej` niż poprzednio. 

![cachev4](https://user-images.githubusercontent.com/103113980/163732626-b00afc4b-7822-4fe1-b668-e7a46f6ebe86.png)

Wchodząc w logi tego przepływu możemy znaleźć informację o wykorzystaniu pamięci `cache`.

![cached](https://user-images.githubusercontent.com/103113980/163732786-9545ca78-c7a7-4de1-978d-1a63f58ab566.png)

---

### **DODATEK 2**

### 1. Uruchomienie prywatnego rejestru 

**a.&ensp;Uruchomienie kontenera na porcie 6677** 

Uruchomienie kontenera z rejestrem. Flaga `--restart=always` powoduje, że kontener w razie zatrzymania automatycznie zostanie zrestartowany i uruchomiony ponownie. 
```
$ docker run -d -p 6677:5000 --restart=always --name private_registry registry
```
Kontener uruchomił się prawidłowo. Widzimy, że nasłuchuje on na porcie 6677 i jest widoczny na liście uruchomionych kontenerów. 

![registry](https://user-images.githubusercontent.com/103113980/163835977-70146c51-4277-48d0-8e06-583eaa09f9ae.png)

**b.&ensp;Pobranie najnowszego Ubuntu i wgranie go do utworzonego rejestru** 

```
$ docker pull ubuntu:latest
```

Zmiana nazwy obrazu - dodajemy tag do istniejącego obrazu. Gdy pierwsza część tagu zawiera nazwę hosta i port, `Docker` interpretuje ją jako lokalizację rejestru przy pushowaniu obrazu. 

```
$ docker tag ubuntu:latest localhost:6677/private_ubuntu
```

Wgrywamy obraz do naszego rejestru

```
$ docker push localhost:6677/private-ubuntu
```

Aby sprawdzić czy nasz rejestr rzeczywiście działa usuńmy oba obrazy Ubuntu i spróbujmy pobrać go z rejestru.

![ubuntu_push](https://user-images.githubusercontent.com/103113980/163839954-45f5428c-183a-4af5-a0da-fb118a4e7c10.png)

Udało pobrać się obraz Ubuntu z naszego prywatnego rejestru. 

---

### 2. Dodanie mechanizmu kontroli dostępu htpasswd. 

Zgodnie z tym, co napisano w [dokumentacji](https://docs.docker.com/registry/deploying/#native-basic-auth) Dockera, aby korzystać z uwierzytelniania należy wcześniej skonfigurować `TLS`. Aby to zrobić musimy wygenerować certyfikat SSL dla domeny localhost. 

W tym celu tworzymy folder `certs`, w którym zapiszemy wszystkie wygenerowane pliki. Folder znajduje się w katalogu domowym użytkownika `student`.

```
$ mkdir certs && cd certs
```

Stworzenie własnego urzędu certyfikacji CA

```
$ openssl req -x509 -nodes -new -sha256 -days 1024 -newkey rsa:2048 -keyout RootCA.key -out RootCA.pem -subj "/C=PL/CN=POLLUB-CA"
```
```
$ openssl x509 -outform pem -in RootCA.pem -out RootCA.crt
```

Tworzymy plik `domains.ext` który będzie zawierał informacje o naszej domenie

```
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = localhost
```

Generujemy pliki `localhost.key`, `localhost.crt` oraz `localhost.csr`. 

```
$ openssl req -new -nodes -newkey rsa:2048 -keyout localhost.key -out localhost.csr -subj "/C=PL/ST=Lubelskie/L=Lublin/O=PrivateLocalhostCert/CN=localhost.local"
```

```
$ openssl x509 -req -sha256 -days 1024 -in localhost.csr -CA RootCA.pem -CAkey RootCA.key -CAcreateserial -extfile domains.ext -out localhost.crt
```

Następnym krokiem jest utworzenie danych logowania użytkownika za pomocą `htpasswd`. Cały proces został opisany w podlinkowanej do zadania [dokumentacji](https://docs.docker.com/registry/deploying/#native-basic-auth) Dockera. Dane do logowania wygenerujemy w folderze `auth`, który również znajduje się w katalogu domowym użytkownika `student`.

```
$ sudo docker run --entrypoint htpasswd httpd:2 -Bbn testuser testpassword > auth/htpasswd
```

Po utworzeniu użytkownika `testuser` z hasłem `testpassword` zatrzymujemy kontener z rejestrem i usuwamy go. 

```
$ docker container stop private_registry && docker rm private_registry
```

Uruchamiamy kontener z uwierzytelnianiem. 

```
docker run -d \
  -p 6677:5000 \
  --restart=always \
  --name registry_private \
  -v "$(pwd)"/auth:/auth \
  -e "REGISTRY_AUTH=htpasswd" \
  -e "REGISTRY_AUTH_HTPASSWD_REALM=Registry Realm" \
  -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
  -v "$(pwd)"/certs:/certs \
  -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/localhost.crt \
  -e REGISTRY_HTTP_TLS_KEY=/certs/localhost.key \
  registry
  ```
  
![auth](https://user-images.githubusercontent.com/103113980/163878703-bcbac3b1-2bf8-44ce-b41c-97dbc55b5417.png)
  
Jak widać uruchomienie rejestru z uwierzytelnianiem powiodło się. Przy próbie wypchnięcia wcześniej utworzonego obrazu Ubuntu dostaliśmy informację o braku bycia uwierzytelnionym. Po zalogowaniu się na użytkownika `testuser` z hasłem `testpassword` próba przesłania obrazu na repozytorium powiodła się.

---
