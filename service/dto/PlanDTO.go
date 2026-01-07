package dto

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type subscriptionPlanResp struct {
	PlanCode pgtype.Text
	PlanName pgtype.Text
	StartAt  time.Time
	ExpireAt time.Time
}
