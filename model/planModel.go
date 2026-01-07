package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"time"
)

type SubscriptionPlanModel struct {
	Id        uuid.UUID
	Code      pgtype.Text
	Name      pgtype.Text
	Price     decimal.Decimal
	Currency  pgtype.Text
	Interval  pgtype.Text
	CreatedAt time.Time
}
type UserSubscriptionModel struct {
	Id          uuid.UUID
	UserId      uuid.UUID
	PlanId      uuid.UUID
	Status      pgtype.Text
	StartAt     time.Time
	ExpiresAt   time.Time
	Provider    pgtype.Text
	ProviderRef pgtype.Text
	CreatedAt   time.Time
}
