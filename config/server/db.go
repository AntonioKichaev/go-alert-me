package server

import "go.uber.org/zap/zapcore"

type DBConfig struct {
	DatabaseDNS string `env:"DATABASE_DSN"`
}

func (db *DBConfig) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("DATABASE_DSN", db.DatabaseDNS)
	return nil
}

func NewDBConfig() *DBConfig {
	dbConf := &DBConfig{}

	return dbConf

}
