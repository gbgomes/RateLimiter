package database

import "github.com/go-redis/redis"

func Newdb(tipo string, addr string, port string) (RateLimiterInterfaceRepository, error) {

	var bd RateLimiterInterfaceRepository
	if tipo == "Redis" {
		BdClient := redis.NewClient(&redis.Options{
			Addr:     addr + ":" + port,
			Password: "",
			DB:       0,
		})
		_, err := BdClient.Ping().Result()
		if err != nil {
			return nil, err
		}
		bd = NewRedisRL(BdClient)
	}

	return bd, nil
}
