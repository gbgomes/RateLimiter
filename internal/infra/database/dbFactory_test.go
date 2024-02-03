package database

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDataBaseType(t *testing.T) {
	bd, err := Newdb("Redis", "localhost", "6379")
	assert.NotNil(t, err)
	assert.NotNil(t, bd)

	assert.Equal(t, "*database.RedisRL", reflect.TypeOf(bd).String())

}
