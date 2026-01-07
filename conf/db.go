package conf

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"time"
)

var Pool *pgxpool.Pool

func InitPostgres() error {
	connStr := viper.GetString("DATABASE_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}

	/*===========================================================================================*/
	//連線池設定

	cfg.MaxConns = viper.GetInt32("db.maxConnects")
	cfg.MinConns = viper.GetInt32("db.minConnects")
	cfg.MaxConnLifetime = viper.GetDuration("db.connectLifetime_min")

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return err
	}

	/*===========================================================================================*/
	//測試連線
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	Pool = pool
	return nil
}
