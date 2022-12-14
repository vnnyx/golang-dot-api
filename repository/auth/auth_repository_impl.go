package auth

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vnnyx/golang-dot-api/model"
)

type AuthRepositoryImpl struct {
	Redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) AuthRepository {
	return &AuthRepositoryImpl{Redis: redis}
}

func (repository *AuthRepositoryImpl) StoreToken(ctx context.Context, details model.TokenDetails) error {
	return repository.Redis.Set(ctx, details.AccessUUID, details.AccessToken, time.Unix(details.AtExpires, 0).Sub(time.Now())).Err()
}

func (repository *AuthRepositoryImpl) DeleteToken(ctx context.Context, accessUuid string) error {
	return repository.Redis.Del(ctx, accessUuid).Err()
}

func (repository *AuthRepositoryImpl) GetToken(ctx context.Context, accessUuid string) (access string, err error) {
	return repository.Redis.Get(ctx, accessUuid).Result()
}
