package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func New(client *redis.Client) *Repository {
	return &Repository{client: client}
}

func newClient(url, port string, dataBase int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", url, port),
		DB:   dataBase,
	})
}
