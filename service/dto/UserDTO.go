package dto

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserProfileResp struct {
	Id          uuid.UUID
	Email       pgtype.Text
	UserName    pgtype.Text
	DisplayName pgtype.Text
	AvatarURL   pgtype.Text
	Role        pgtype.Text
	Plan        subscriptionPlanResp
}
