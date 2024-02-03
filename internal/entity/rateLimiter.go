package entity

import (
	"log"

	"github.com/gbgomes/GoExpert/RateLimiter/internal/infra/database"
)

type RateLimiter struct {
	BD                 database.RateLimiterInterfaceRepository
	IpNumMaxAcessos    int64
	IpTempoLimite      int64
	IpTempoBloqueio    int64
	TokenNumMaxAcessos int64
	TokenTempoLimite   int64
	TokenTempoBloqueio int64
}

func NewRateLimiter(bd database.RateLimiterInterfaceRepository,
	ipNumMaxAcessos, ipTempoLimite, ipTempoBloqueio,
	tokenMaxAccess, tokenTimeLimit, tokenTimeBlocked int64) *RateLimiter {
	return &RateLimiter{
		BD:                 bd,
		IpNumMaxAcessos:    ipNumMaxAcessos,
		IpTempoLimite:      ipTempoLimite,
		IpTempoBloqueio:    ipTempoBloqueio,
		TokenNumMaxAcessos: tokenMaxAccess,
		TokenTempoLimite:   tokenTimeLimit,
		TokenTempoBloqueio: tokenTimeBlocked,
	}
}

func (rl *RateLimiter) TrataRatelimit(tk, ip string) bool {
	var numMaxAccess,
		timeLimit,
		timeBlocked int64

	var key string

	// caso não tenha token no request usa o ip
	if len(tk) == 0 {
		key = ip
		numMaxAccess = rl.IpNumMaxAcessos
		timeLimit = rl.IpTempoLimite
		timeBlocked = rl.IpTempoBloqueio
	} else {
		key = tk
		token, _ := rl.BD.ColetaToken(key)
		if len(token.Token) > 0 {
			numMaxAccess = token.MaxNumberAccess
			timeLimit = token.TimeLimit
			timeBlocked = token.TimeBlock
		} else {
			numMaxAccess = rl.TokenNumMaxAcessos
			timeLimit = rl.TokenTempoLimite
			timeBlocked = rl.IpTempoBloqueio
		}
	}

	//verifica se a chave está bloqueada
	bBloq, _ := rl.BD.ChaveBloqueada(key)
	if bBloq {
		//se estiver bloqueado, remove o contador de acessos
		_ = rl.excluiAcesso(key)
		log.Printf("token/ip %s bloqueado", key)
		return true
	} else {
		// tenta realizar a coleta pela chave.
		found, acessos, err := rl.verificaSeTemAcesso(key)
		if !found {
			// caso não exista registro para a chave, inclui inicializando com 1 acesso
			// e configurando o tempo de expiração para o tempo de verificação de limite de acessos
			err = rl.registraPrimeiroAcesso(key, timeLimit)
			if err != nil {
				log.Println(err)
			}
			log.Printf("numero de acessos do token/ip %s: %d", key, 1)
			return false
		} else if err != nil {
			// loga erro genério de acesso ao Redis
			log.Println(err)
			return false
		} else if acessos >= numMaxAccess {
			// se atingir o limite de acesso, inclui registro para indicar que está bloqueado
			// configurando a expiração para o tempo de bloqueio de acesso
			err = rl.bloqueiAcesso(key, timeBlocked)
			if err != nil {
				log.Println(err)
			}
			log.Printf("limite de acessos alcançado. o token/ip %s foi bloqueado por %d segundos", key, timeBlocked)
			return true
		} else {
			// incrementa a quantidade de acessos
			numAcessos, err := rl.registraAcesso(key)
			if err != nil {
				log.Println(err)
			}
			log.Printf("numero de acessos do token/ip %s: %d", key, numAcessos)
			return false
		}
	}
}

func (rl *RateLimiter) verificaSeTemAcesso(key string) (bool, int64, error) {
	found, acessos, err := rl.BD.ConsultaChave(key)
	return found, acessos, err
}

func (rl *RateLimiter) registraPrimeiroAcesso(key string, timeLimit int64) error {
	return rl.BD.IncluiChave(key, 1, timeLimit)
}

func (rl *RateLimiter) bloqueiAcesso(key string, timeBlocked int64) error {
	return rl.BD.BloqueiaChave(key, timeBlocked)
}

func (rl *RateLimiter) registraAcesso(key string) (int64, error) {
	numAcessos, err := rl.BD.Incrementa(key)
	return numAcessos, err
}

func (rl *RateLimiter) excluiAcesso(key string) error {
	return rl.BD.ExcluiChave(key)
}
