package database

type Token struct {
	Token           string `json:"token"`
	MaxNumberAccess int64  `json:"maxNumberAccess"`
	TimeLimit       int64  `json:"timeLimit"`
	TimeBlock       int64  `json:"timeBlock"`
}

type RateLimiterInterfaceRepository interface {
	ConsultaChave(key string) (bool, int64, error)
	IncluiChave(key string, value interface{}, limitTime int64) error
	Incrementa(key string) (int64, error)
	ExcluiChave(key string) error
	BloqueiaChave(key string, limitBloq int64) error
	ChaveBloqueada(key string) (bool, error)
	InsereHashMap(key string, hm map[string]interface{}) error
	ExcluiListaTokens() error
	ColetaToken(key string) (Token, error)
}
