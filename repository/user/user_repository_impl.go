package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"github.com/vnnyx/golang-dot-api/util"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewUserRepository(DB *gorm.DB, redis *redis.Client) UserRepository {
	return &UserRepositoryImpl{DB: DB, Redis: redis}
}

func (repository *UserRepositoryImpl) InsertUser(ctx context.Context, user entity.User) (entity.User, error) {
	fmt.Println("masuk create")
	err := repository.DB.WithContext(ctx).Create(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindUserByID(ctx context.Context, userId string) (user entity.User, err error) {
	err = repository.DB.WithContext(ctx).Where("id", userId).First(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindUserByUsername(ctx context.Context, username string) (user entity.User, err error) {
	err = repository.DB.WithContext(ctx).Where("username", username).First(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) FindAllUser(ctx context.Context, p *web.Pagination) (users []*entity.User, err error) {
	p.SetPagination(users, repository.DB)
	key := "user:limit:" + strconv.Itoa(p.Limit) + ":page:" + strconv.Itoa(p.Page)
	val, err := repository.Redis.Get(ctx, key).Result()
	if err != nil || val == "" {
		err := repository.DB.WithContext(ctx).Scopes(util.Paginate(users, p, repository.DB)).Find(&users).Error
		p.Rows = users
		if err != nil {
			return users, err
		}
		cache, _ := json.Marshal(users)
		err = repository.Redis.Set(ctx, key, cache, 10*time.Minute).Err()
		if err != nil {
			return users, err
		}
		return users, nil
	}
	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		fmt.Println(err)
		return users, err
	}
	return users, nil
}

func (repository *UserRepositoryImpl) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	err := repository.DB.WithContext(ctx).Where("id", user.UserID).Updates(&user).Error
	return user, err
}

func (repository *UserRepositoryImpl) DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error {
	return tx.WithContext(ctx).Where("id", userId).Delete(&entity.User{}).Error
}

func (repository *UserRepositoryImpl) DeleteAllUser(ctx context.Context) error {
	return repository.DB.WithContext(ctx).Exec("DELETE FROM users").Error
}

func (repository *UserRepositoryImpl) StoreToRedis(ctx context.Context, id string, token web.UserEmailVerification, user entity.User) error {
	err := repository.Redis.Set(ctx, token.UserID, token.OTP, token.ExpiredAt).Err()
	if err != nil {
		return err
	}
	kv, err := util.ToMap(user)
	if err != nil {
		return err
	}
	err = repository.Redis.HSet(ctx, id+"-hash", kv).Err()
	repository.Redis.Expire(ctx, id+"-hash", token.ExpiredAt)
	return err
}

func (repository *UserRepositoryImpl) GetDataToVerify(ctx context.Context, id string) (otp string, user entity.User, err error) {
	got, err := repository.Redis.HGetAll(ctx, id+"-hash").Result()
	if err != nil {
		return otp, user, err
	}
	otp, err = repository.Redis.Get(ctx, id).Result()
	if err != nil {
		return otp, user, err
	}
	b, err := json.Marshal(got)
	if err != nil {
		return otp, user, err
	}
	err = json.Unmarshal(b, &user)
	if err != nil {
		return otp, user, err
	}
	return otp, user, err
}

func (repository *UserRepositoryImpl) DeleteCache(ctx context.Context, id string) error {
	key := []string{id, id + "-hash"}
	return repository.Redis.Del(ctx, key...).Err()
}
