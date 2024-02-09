package redis

import (
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	client *redis.Client
}

func New(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}
