package server

import "go.uber.org/zap/zapcore"

type DBConfig struct {
	DatabaseDns string `env:"DATABASE_DSN"`
}

func (db *DBConfig) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("DATABASE_DSN", db.DatabaseDns)
	return nil
}

func NewDBConfig() *DBConfig {
	dbConf := &DBConfig{}

	return dbConf

}
