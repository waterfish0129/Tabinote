package global

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	Logger *zap.Logger
	DBPool *pgxpool.Pool
)
