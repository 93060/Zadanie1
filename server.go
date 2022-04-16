package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
)
//Struktura do obsługi odpowiedzi z API timezoneapi wygenerowana automatycznie
//za pomocą strony mholt.github.io/json-to-go/
type Response struct {
	Meta struct {
		Code          string `json:"code"`
		ExecutionTime string `json:"execution_time"`
	} `json:"meta"`
	Data struct {
		IP          string      `json:"ip"`
		City        interface{} `json:"city"`
		Postal      interface{} `json:"postal"`
		State       interface{} `json:"state"`
		StateCode   interface{} `json:"state_code"`
		Country     string      `json:"country"`
		CountryCode string      `json:"country_code"`
		Location    string      `json:"location"`
		Timezone    struct {
			ID                  string `json:"id"`
			Location            string `json:"location"`
			CountryCode         string `json:"country_code"`
			CountryName         string `json:"country_name"`
			Iso31661Alpha2      string `json:"iso3166_1_alpha_2"`
			Iso31661Alpha3      string `json:"iso3166_1_alpha_3"`
			UnM49Code           string `json:"un_m49_code"`
			Itu                 string `json:"itu"`
			Marc                string `json:"marc"`
			Wmo                 string `json:"wmo"`
			Ds                  string `json:"ds"`
			PhonePrefix         string `json:"phone_prefix"`
			Fifa                string `json:"fifa"`
			Fips                string `json:"fips"`
			Gual                string `json:"gual"`
			Ioc                 string `json:"ioc"`
			CurrencyAlphaCode   string `json:"currency_alpha_code"`
			CurrencyCountryName string `json:"currency_country_name"`
			CurrencyMinorUnit   string `json:"currency_minor_unit"`
			CurrencyName        string `json:"currency_name"`
			CurrencyCode        string `json:"currency_code"`
			Independent         string `json:"independent"`
			Capital             string `json:"capital"`
			Continent           string `json:"continent"`
			Tld                 string `json:"tld"`
			Languages           string `json:"languages"`
			GeonameID           string `json:"geoname_id"`
			Edgar               string `json:"edgar"`
		} `json:"timezone"`
		Datetime struct {
			Date          string `json:"date"`
			DateTime      string `json:"date_time"`
			DateTimeTxt   string `json:"date_time_txt"`
			DateTimeWti   string `json:"date_time_wti"`
			DateTimeYmd   string `json:"date_time_ymd"`
			Time          string `json:"time"`
			Month         string `json:"month"`
			MonthWilz     string `json:"month_wilz"`
			MonthAbbr     string `json:"month_abbr"`
			MonthFull     string `json:"month_full"`
			MonthDays     string `json:"month_days"`
			Day           string `json:"day"`
			DayWilz       string `json:"day_wilz"`
			DayAbbr       string `json:"day_abbr"`
			DayFull       string `json:"day_full"`
			Year          string `json:"year"`
			YearAbbr      string `json:"year_abbr"`
			Hour12Wolz    string `json:"hour_12_wolz"`
			Hour12Wilz    string `json:"hour_12_wilz"`
			Hour24Wolz    string `json:"hour_24_wolz"`
			Hour24Wilz    string `json:"hour_24_wilz"`
			HourAmPm      string `json:"hour_am_pm"`
			Minutes       string `json:"minutes"`
			Seconds       string `json:"seconds"`
			Week          string `json:"week"`
			OffsetSeconds string `json:"offset_seconds"`
			OffsetMinutes string `json:"offset_minutes"`
			OffsetHours   string `json:"offset_hours"`
			OffsetGmt     string `json:"offset_gmt"`
			OffsetTzid    string `json:"offset_tzid"`
			OffsetTzab    string `json:"offset_tzab"`
			OffsetTzfull  string `json:"offset_tzfull"`
			TzString      string `json:"tz_string"`
			Dst           string `json:"dst"`
			DstObserves   string `json:"dst_observes"`
			TimedaySpe    string `json:"timeday_spe"`
			TimedayGen    string `json:"timeday_gen"`
		} `json:"datetime"`
	} `json:"data"`
}

func main() {
	//ustawienie portu na ktorym bedzie dzialal serwer
	PORT := "8082"
	//Wyswietlenie wiadomosci powitalnej na konsoli
	fmt.Println("Hello! I'm running on port: "+ PORT)
	//lokalizacja pliku z logami
	LOG_FILE := "./app.log"
	//utworzenie pliku z logami
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE,0644)
	//jezeli blad
	if err != nil {
        log.Panic(err)
    }
	defer logFile.Close()
	//ustawienie wyjscia dla zapisu logow
	log.SetOutput(logFile)
	//ustawienie flag dla logow
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//obsluga wywolania strony serwera
	http.HandleFunc("/", server)
	//obsluga wywolania logow
	http.HandleFunc("/log", logShow)
	//log w razie bledu
	log.Fatal(http.ListenAndServe(":" + PORT, nil))
	log.Println("Running on port " + PORT)
}

//funkcja obslugujaca wyswietlenie informacji o polaczeniu z serwerem
func server(w http.ResponseWriter, r *http.Request) {
	//Pobranie adresu IP klienta
	IPAddress := r.Header.Get("X-Real-Ip")
	//zmienna przechowujaca tresc wyswietlanej wiadomosci
	msg := ""
	//zmienna przechowujaca dane o autorze
	name := ""
	//jezeli naglowek X-Real-Ip jest pusty to sprawdz X-Forwarded-For
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	//jezeli adres IP nadal pusty to ustaw adres zdalny
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	//podzielenie adresu na adres IP i port
	tmpStr := strings.Split(IPAddress, ":")
	//przypisanie adresu IP z tablicy tmpStr(index 0 to adres IP, index 1 to port)
	IPAddress = tmpStr[0]
	name = "Autor: Sebastian Wiktor. Serwer wykonany na potrzeby Zadania nr 1 z przedmiotu Technologie Chmurowe.\n\n"
	msg = "Adres IP klienta: " + IPAddress
	//Jezeli adres IP jest prywatny
	if net.ParseIP(IPAddress).IsPrivate(){
		realIP := IPAddress
		IPAddress = "5.173.14.19"
		strPriv := "Polaczenie z prywatnego adresu IP (" + realIP + ") - ustawiono przykladowy adres publiczny. "
		msg = strPriv + "\nAdres IP: " + IPAddress
		log.Println(strPriv + "IP: " + IPAddress +" Autor: Sebastian Wiktor")
	}
	//Jezeli adres IP to loopback
	if net.ParseIP(IPAddress).IsLoopback(){
		realIP := IPAddress
		IPAddress = "5.173.14.19"
		strLoop := "Polaczenie z adresu loopback (" + realIP + ") - ustawiono przykladowy adres publiczny. "
		msg = strLoop+ "\nAdres IP: " + IPAddress
		log.Println(strLoop + "IP: " + IPAddress +" Autor: Sebastian Wiktor")
	}
	//URL wysylany do API timezoneapi w celu uzyskania informacji nt. strefy czasowej klienta
	url := "https://timezoneapi.io/api/ip/?" + IPAddress + "&token=aefZDpAVRzhAMLcCimEO"
	//wyslanie zapytania 
	response, err := http.Get(url)
	//jezeli blad to zwroc error
	if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }
    	//zmienna przechowujaca odpowiedz
	responseData, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
    }
    	//utworzenie struktury response
	var result Response
	//wypakowanie otrzymanego jsona do struktury response
	json.Unmarshal(responseData, &result)
	//przekazanie daty i godziny do wyswietlanej wiadomosci
	msg += "\nData i godzina polaczenia: " + result.Data.Datetime.Date + " " + result.Data.Datetime.Time
	//wyswietlenie danych o autorze
	fmt.Fprint(w, name)
	//wyswietlenie wiadomosci o polaczeniu
	fmt.Fprint(w, msg)
	
}
//funkcja obslugujaca wyswietlenie pliku z logami
func logShow(w http.ResponseWriter, r *http.Request){
	//do zmiennej content wczytaj zawartosc pliku app.log
	content, err := ioutil.ReadFile("app.log")
	//jezeli blad to zwroc error
    if err != nil {
        log.Fatal(err)
    }	
	//wyswietlenie zawartosci content na strone
	fmt.Fprint(w, string(content))
}
