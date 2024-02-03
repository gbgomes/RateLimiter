package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gbgomes/GoExpert/RateLimiter/configs"
	"github.com/gbgomes/GoExpert/RateLimiter/internal/entity"
	"github.com/gbgomes/GoExpert/RateLimiter/internal/infra/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handdler struct{}

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	ipMaxAccess, _ := strconv.ParseInt(configs.IPMaxNumberAccess, 10, 64)
	ipTimeLimit, _ := strconv.ParseInt(configs.IPTimeLimit, 10, 64)
	ipTimeBlocked, _ := strconv.ParseInt(configs.IPTimeBlock, 10, 64)
	tokenMaxAccess, _ := strconv.ParseInt(configs.IPMaxNumberAccess, 10, 64)
	tokenTimeLimit, _ := strconv.ParseInt(configs.IPTimeLimit, 10, 64)
	tokenTimeBlocked, _ := strconv.ParseInt(configs.IPTimeBlock, 10, 64)

	bd, err := database.Newdb(configs.BdType, configs.BdAddr, configs.BdPort)
	if bd == nil || err != nil {
		log.Fatalf("erro conectando ao %s: %s", configs.BdType, err)
	}

	rl := entity.NewRateLimiter(
		bd,
		ipMaxAccess,
		ipTimeLimit,
		ipTimeBlocked,
		tokenMaxAccess,
		tokenTimeLimit,
		tokenTimeBlocked,
	)

	bd.ExcluiListaTokens()
	readTokenLimits(configs.TokenFileLimits, bd)

	r := chi.NewRouter()
	r.Use(middleware.WithValue("rl", rl))
	r.Use(middleware.Recoverer)
	r.Use(ratelimit)
	r.Route("/", func(r chi.Router) {
		r.Get("/", httpHandler)
	})

	http.ListenAndServe(":8080", r)

}

func readTokenLimits(file string, bd database.RateLimiterInterfaceRepository) {
	jsonFile, err := os.Open(file)
	if err != nil {
		//arquivo não especificado
		log.Println("nenhum aquivo de tokens especificado")
		return
	}
	defer jsonFile.Close()
	var tokens []database.Token
	err = json.NewDecoder(jsonFile).Decode(&tokens)
	if err != nil {
		log.Println("erro decodificando json")
	}
	for _, token := range tokens {
		tk := make(map[string]interface{})
		tk["Token"] = token.Token
		tk["MaxNumberAccess"] = token.MaxNumberAccess
		tk["TimeLimit"] = token.TimeLimit
		tk["TimeBlock"] = token.TimeBlock
		bd.InsereHashMap("TkConfig:"+token.Token, tk)
	}
}

func httpHandler(w http.ResponseWriter, req *http.Request) {

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	dt := time.Now()
	strRet := "data e hora do acesso: " + dt.Format("2006-01-02 15:04:05")
	w.Write([]byte(strRet))
}

func ratelimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl := r.Context().Value("rl").(*entity.RateLimiter)
		ip, _ := getIP(r)
		limiteAtingido := rl.TrataRatelimit(r.Header.Get("API_KEY"), ip)

		if limiteAtingido {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		// o último IP da lista é o real do cliente.
		netIP := net.ParseIP(splitIps[len(splitIps)-1])
		if netIP != nil {
			return netIP.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	netIP := net.ParseIP(ip)
	if netIP != nil {
		ip := netIP.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("IP not found")
}
