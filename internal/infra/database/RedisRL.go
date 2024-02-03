package database

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RedisRL struct {
	RD *redis.Client
}

func NewRedisRL(rd *redis.Client) *RedisRL {
	return &RedisRL{RD: rd}
}

func (r *RedisRL) IncluiChave(key string, value interface{}, limitTime int64) error {
	err := r.RD.Set(key, value, time.Duration(limitTime)*time.Second).Err()

	return err
}

func (r *RedisRL) ConsultaChave(key string) (bool, int64, error) {
	counter, err := r.RD.Get(key).Int64()
	if err == redis.Nil {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}

	return true, counter, nil
}

func (r *RedisRL) Incrementa(key string) (int64, error) {
	numAcessos, err := r.RD.Incr(key).Result()

	return numAcessos, err
}

func (r *RedisRL) ExcluiChave(key string) error {
	err := r.RD.Del(key).Err()
	return err
}

func (r *RedisRL) BloqueiaChave(key string, limitBloq int64) error {
	err := r.RD.Set(key+"_blocked", 1, time.Duration(limitBloq)*time.Second).Err()

	return err
}

func (r *RedisRL) ChaveBloqueada(key string) (bool, error) {
	err := r.RD.Get(key + "_blocked").Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisRL) InsereHashMap(key string, hm map[string]interface{}) error {

	r.RD.HMSet(key, hm)

	return nil
}

func (r *RedisRL) ExcluiListaTokens() error {

	tokens, err := r.RD.Keys("TkConfig:*").Result()
	if err != nil {
		return err
	}

	for _, token := range tokens {
		err := r.RD.Del(token).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisRL) ColetaToken(key string) (Token, error) {
	token := Token{}
	tk, err := r.RD.HGetAll("TkConfig:" + key).Result()
	if err == redis.Nil {
		return token, nil
	}
	if err != nil {
		return token, err
	}

	token.Token = tk["Token"]

	maxNumberAccess, err := strconv.ParseInt(tk["MaxNumberAccess"], 10, 64)
	if err != nil {
		maxNumberAccess = 0
	}
	timeLimit, err := strconv.ParseInt(tk["TimeLimit"], 10, 64)
	if err != nil {
		timeLimit = 0
	}
	timeBlock, err := strconv.ParseInt(tk["TimeBlock"], 10, 64)
	if err != nil {
		timeBlock = 0
	}

	token.MaxNumberAccess = maxNumberAccess
	token.TimeLimit = timeLimit
	token.TimeBlock = timeBlock

	return token, nil
}
