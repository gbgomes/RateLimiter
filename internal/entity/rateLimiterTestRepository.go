package entity

import "github.com/gbgomes/GoExpert/RateLimiter/internal/infra/database"

type Data struct {
	Key   string
	Value any
}

type RateLimiterRepository struct {
	DataHM map[string]Data
}

func NewRateLimiterTestRepository() *RateLimiterRepository {
	return &RateLimiterRepository{make(map[string]Data)}
}

func (r *RateLimiterRepository) ConsultaChave(key string) (bool, int64, error) {
	item := r.DataHM[key]
	if item.Key == "" {
		return false, 0, nil
	}
	return true, item.Value.(int64), nil
}

//type Token struct {
//	Token           string `json:"token"`
//	MaxNumberAccess int64  `json:"maxNumberAccess"`
//	TimeLimit       int64  `json:"timeLimit"`
//	TimeBlock       int64  `json:"timeBlock"`
//	Value           int64  `json:"value"`
//}

func (r *RateLimiterRepository) IncluiChave(key string, value interface{}, limitTime int64) error {
	item := Data{key, value.(int64)}
	r.DataHM[item.Key] = item
	return nil
}

func (r *RateLimiterRepository) Incrementa(key string) (int64, error) {
	item := r.DataHM[key]
	count := item.Value.(int64)
	count += 1
	item.Value = count

	r.DataHM[item.Key] = item

	return item.Value.(int64), nil
}

func (r *RateLimiterRepository) ExcluiChave(key string) error {
	delete(r.DataHM, key)
	return nil
}

func (r *RateLimiterRepository) BloqueiaChave(key string, limitBloq int64) error {
	err := r.IncluiChave(key+"_Block", int64(1), limitBloq)
	return err
}

func (r *RateLimiterRepository) ChaveBloqueada(key string) (bool, error) {
	found, _, err := r.ConsultaChave("test_Block")
	return found, err
}

func (r *RateLimiterRepository) InsereHashMap(key string, hm map[string]interface{}) error {
	item := Data{key, hm}
	r.DataHM[item.Key] = item
	return nil
}

func (r *RateLimiterRepository) ExcluiListaTokens() error {
	return nil
}

func (r *RateLimiterRepository) ColetaToken(key string) (database.Token, error) {
	return database.Token{}, nil
}
