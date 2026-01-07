package model

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type UserModel struct {
	Id            uuid.UUID
	Email         pgtype.Text
	EmailVerified bool
	UserName      pgtype.Text
	DisplayName   pgtype.Text
	AvatarURL     pgtype.Text
	Status        int8
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserProfileModel struct {
	Id            uuid.UUID
	Email         pgtype.Text
	EmailVerified bool
	UserName      pgtype.Text
	DisplayName   pgtype.Text
	AvatarURL     pgtype.Text
	Status        int8
	ActivePlan    UserPlanModel
}

type UserPlanModel struct {
	Code      pgtype.Text
	Name      pgtype.Text
	Status    pgtype.Text
	StartAt   time.Time
	ExpiresAt time.Time
}
