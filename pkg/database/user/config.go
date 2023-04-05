package database

import (
	"errors"
	"os"
	"strconv"
)

const (
	REDIS_ADDRESS  = "REDIS_ADDRESS"
	REDIS_USERNAME = "REDIS_USERNAME"
	REDIS_PASSWORD = "REDIS_PASSWORD"
	REDIS_DATABASE = "REDIS_DATABASE"
)

const (
	defaultRedisUsername = "default"
	defaultRedisInstance = 0
)

type RedisConfig interface {
	GetAddress() string
	GetUsername() string
	GetPassword() string
	GetDbInstance() int
}

type redisConfig struct {
	address    string
	username   string
	password   string
	dbInstance int
}

func NewRedisConfig() (RedisConfig, error) {

	address, ok := os.LookupEnv(REDIS_ADDRESS)
	if !ok {
		return redisConfig{}, errors.New("no adress env var")
	}

	password, ok := os.LookupEnv(REDIS_PASSWORD)
	if !ok {
		return redisConfig{}, errors.New("no password env var")
	}

	dbInstance, err := strconv.Atoi(getEnv(REDIS_DATABASE, ""))
	if err != nil {
		dbInstance = defaultRedisInstance
	}

	username := getEnv(REDIS_USERNAME, defaultRedisUsername)

	return redisConfig{
		address:    address,
		username:   username,
		password:   password,
		dbInstance: dbInstance,
	}, nil
}

func (c redisConfig) GetAddress() string {
	return c.address
}

func (c redisConfig) GetUsername() string {
	return c.username
}

func (c redisConfig) GetPassword() string {
	return c.password
}
func (c redisConfig) GetDbInstance() int {
	return c.dbInstance
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
