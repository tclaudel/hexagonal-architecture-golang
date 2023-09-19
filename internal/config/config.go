package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	mongoURIKey    = "MONGO_URI"
	mongoDBName    = "MONGO_DB_NAME"
	userServerPort = "USER_SERVER_PORT"
)

const (
	mongoURIDefault       = "mongodb://localhost:27017"
	mongoDBNameDefault    = "user"
	userServerPortDefault = "8080"
)

type Config struct {
	v *viper.Viper
}

func NewConfig() *Config {
	viperCfg := viper.New()

	viperCfg.AutomaticEnv()

	return &Config{
		v: viperCfg,
	}
}

func (c Config) MongoURI() string {
	mongoURI := c.v.GetString(mongoURIKey)
	if mongoURI == "" {
		log.Printf("Mongo uri not set, using default : %s", mongoURIDefault)

		return mongoURIDefault
	}

	log.Printf("Mongo uri set : %s", mongoURI)

	return mongoURI
}

func (c Config) MongoDBName() string {
	mongoDBName := c.v.GetString(mongoDBName)
	if mongoDBName == "" {
		log.Printf("Mongo database name not set, using default : %s", mongoDBNameDefault)

		return mongoDBNameDefault
	}

	log.Printf("Mongo database name set : %s", mongoDBName)

	return mongoDBName
}

func (c Config) UserServerPort() string {
	userServerPort := c.v.GetString(userServerPort)
	if userServerPort == "" {
		log.Printf("User server port not set, using default : %s", userServerPortDefault)

		return userServerPortDefault
	}

	log.Printf("User server port set : %s", userServerPort)

	return userServerPort
}
