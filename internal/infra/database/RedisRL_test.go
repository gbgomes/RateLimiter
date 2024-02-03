package database

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestInclusaoDeChave(t *testing.T) {
	BdClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	bd.IncluiChave("test", 1, 5)
	counter, err := BdClient.Get("test").Int64()
	assert.Nil(t, err)
	assert.NotEqual(t, counter, 0)
	assert.Nil(t, err)
}

func TestConsultaDeChave(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	bd.IncluiChave("test", 1, 0)
	found, counter, err := bd.ConsultaChave("test")
	assert.True(t, found)
	assert.Equal(t, counter, int64(1))
	assert.Nil(t, err)
}

func TestExpiracaoDeChave(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	bd.IncluiChave("test", 1, 5)
	found, counter, err := bd.ConsultaChave("test")
	assert.True(t, found)
	assert.NotEqual(t, counter, 0)
	assert.Nil(t, err)

	time.Sleep(5 * time.Second)

	found, counter, err = bd.ConsultaChave("test")
	assert.False(t, found)
	assert.Equal(t, counter, int64(0))
	assert.Nil(t, err)

}

func TestIncrementaChave(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	bd.IncluiChave("test", 1, 5)
	found, counter, err := bd.ConsultaChave("test")
	assert.True(t, found)
	assert.Equal(t, counter, int64(1))
	assert.Nil(t, err)

	counter, err = bd.Incrementa("test")
	assert.Nil(t, err)
	assert.Equal(t, counter, int64(2))
}

func TestExclusaoDeChave(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	bd.IncluiChave("test", 1, 0)
	found, counter, err := bd.ConsultaChave("test")
	assert.True(t, found)
	assert.NotEqual(t, counter, 0)
	assert.Nil(t, err)

	err = bd.ExcluiChave("test")
	assert.Nil(t, err)

	found, counter, err = bd.ConsultaChave("test")
	assert.False(t, found)
	assert.Equal(t, counter, int64(0))
	assert.Nil(t, err)

}

func TestBloqueiaChave(t *testing.T) {
	BdClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	err = bd.BloqueiaChave("test", 5)
	assert.Nil(t, err)

	err = BdClient.Get("test" + "_blocked").Err()
	assert.Nil(t, err)

	time.Sleep(5 * time.Second)

	err = BdClient.Get("test" + "_blocked").Err()
	assert.NotNil(t, err)
}

func TestVerificaChaveBloqueada(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	bd.BloqueiaChave("test", 1)
	found, err := bd.ChaveBloqueada("test")
	assert.Nil(t, err)
	assert.True(t, found)
}

func TestListaTokens(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	tokens := []Token{
		{"tk1", 1, 1, 10},
		{"tk2", 2, 2, 20},
	}

	for _, token := range tokens {
		tk := make(map[string]interface{})
		tk["Token"] = token.Token
		tk["MaxNumberAccess"] = token.MaxNumberAccess
		tk["TimeLimit"] = token.TimeLimit
		tk["TimeBlock"] = token.TimeBlock
		err := bd.InsereHashMap(token.Token, tk)
		assert.Nil(t, err)
	}

	tk, err := bd.ColetaToken("tk1")
	assert.Nil(t, err)
	assert.NotNil(t, tk)
	assert.Equal(t, int64(1), tk.MaxNumberAccess)
	assert.Equal(t, int64(1), tk.TimeLimit)
	assert.Equal(t, int64(10), tk.TimeBlock)

	tk, err = bd.ColetaToken("tk2")
	assert.Nil(t, err)
	assert.NotNil(t, tk)
	assert.Equal(t, int64(2), tk.MaxNumberAccess)
	assert.Equal(t, int64(2), tk.TimeLimit)
	assert.Equal(t, int64(20), tk.TimeBlock)

	err = bd.ExcluiListaTokens()
	assert.Nil(t, err)

	tk, err = bd.ColetaToken("tk1")
	assert.Nil(t, err)
	assert.NotNil(t, tk)
	assert.Equal(t, int64(0), tk.MaxNumberAccess)
	assert.Equal(t, int64(0), tk.TimeLimit)
	assert.Equal(t, int64(0), tk.TimeBlock)
}
