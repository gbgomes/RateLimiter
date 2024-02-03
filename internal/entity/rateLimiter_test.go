package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInclusaoDeChave(t *testing.T) {

	bd := NewRateLimiterTestRepository()
	assert.NotNil(t, bd)

	rl := NewRateLimiter(
		bd,
		5,
		5,
		10,
		5,
		5,
		10,
	)

	err := bd.IncluiChave("test", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bd.DataHM))

	err = bd.IncluiChave("test1", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(bd.DataHM))

}

func TestConsultaChave(t *testing.T) {

	bd := NewRateLimiterTestRepository()
	assert.NotNil(t, bd)

	rl := NewRateLimiter(
		bd,
		5,
		5,
		10,
		5,
		5,
		10,
	)

	err := bd.IncluiChave("test", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)

	found, count, err := bd.ConsultaChave("test")
	assert.Nil(t, err)
	assert.True(t, found)
	assert.Equal(t, int64(1), count)
}

func TestIncrementaChave(t *testing.T) {

	bd := NewRateLimiterTestRepository()
	assert.NotNil(t, bd)

	rl := NewRateLimiter(
		bd,
		5,
		5,
		10,
		5,
		5,
		10,
	)

	err := bd.IncluiChave("test1", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bd.DataHM))

	count, err := bd.Incrementa("test1")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
	count, err = bd.Incrementa("test1")
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)

	err = bd.IncluiChave("test2", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(bd.DataHM))

	count, err = bd.Incrementa("test2")
	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
	count, err = bd.Incrementa("test2")
	assert.Nil(t, err)
	assert.Equal(t, int64(3), count)
}

func TestExcluiChave(t *testing.T) {

	bd := NewRateLimiterTestRepository()
	assert.NotNil(t, bd)

	rl := NewRateLimiter(
		bd,
		5,
		5,
		10,
		5,
		5,
		10,
	)

	err := bd.IncluiChave("test1", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bd.DataHM))

	err = bd.IncluiChave("test2", int64(1), rl.IpTempoLimite)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(bd.DataHM))

	bd.ExcluiChave("test1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bd.DataHM))

	bd.ExcluiChave("test1")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(bd.DataHM))

	bd.ExcluiChave("test2")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(bd.DataHM))
}

func TestBloqueiaConsultaChave(t *testing.T) {

	bd := NewRateLimiterTestRepository()
	assert.NotNil(t, bd)

	rl := NewRateLimiter(
		bd,
		5,
		5,
		10,
		5,
		5,
		10,
	)

	err := bd.BloqueiaChave("test", rl.IpTempoBloqueio)
	assert.Nil(t, err)

	found, err := bd.ChaveBloqueada("test")
	assert.Nil(t, err)
	assert.True(t, found)
}
