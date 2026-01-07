package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/model"
	"github.com/waterfish0129/Tabinote/utils"
	"strings"
	"time"
)

var userDao *UserDao

type UserDao struct {
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (m *UserDao) GetByGoogleSub(c context.Context,
	googleSub string) (*model.UserProfileModel, error) {
	var user model.UserProfileModel
	queryStr :=
		`
SELECT u.id, u.email, u.display_name, u.avatar_url
FROM users u
JOIN user_auth_identities uai ON u.id = uai.user_id
JOIN user_subscriptions us  ON u.id = us.user_id
WHERE uai.provider = 'google'
AND uai.provider_uid = $1 and us.status ='active'
LIMIT 1
`
	err := global.DBPool.QueryRow(c, queryStr, googleSub).
		Scan(&user.Id, &user.Email, &user.DisplayName, &user.AvatarURL)
	return &user, err
}

func (m *UserDao) CreateByGoogle(c context.Context, tx pgx.Tx,
	googleSub string,
	email string,
	name string,
	avatar *string) (*model.UserModel, error) {
	iUserModel := model.UserModel{
		Id:          uuid.New(),
		Email:       utils.ToPgText(utils.StringPtr(email)),
		UserName:    utils.ToPgText(utils.StringPtr(strings.Split(email, "@")[0])),
		DisplayName: utils.ToPgText(utils.StringPtr(name)),
		AvatarURL:   utils.ToPgText(avatar),
	}
	/*===========================================================================================*/
	//建立新user
	queryStr := `
INSERT INTO users (id, email,username, display_name, avatar_url, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
`
	_, err := tx.Exec(
		c,
		queryStr,
		iUserModel.Id,
		iUserModel.Email,
		iUserModel.UserName,
		iUserModel.DisplayName,
		iUserModel.AvatarURL,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}
	/*===========================================================================================*/
	//建立google identity
	queryStr = `
INSERT INTO user_auth_identities
(id, user_id, provider, provider_uid, password_hash, created_at)
VALUES ($1, $2, 'google', $3, $4, $5)
`
	_, err = tx.Exec(
		c,
		queryStr,
		uuid.New(),
		iUserModel.Id,
		googleSub,
		nil,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	return &iUserModel, nil
}

func (m *UserDao) CreateSubscription(c context.Context, tx pgx.Tx,
	userId uuid.UUID,
	planId uuid.UUID,
	startAt time.Time,
	expiresAt time.Time,
	provider pgtype.Text,
	providerRef pgtype.Text) error {
	iUserSubscriptionModel := model.UserSubscriptionModel{
		Id:          uuid.New(),
		UserId:      userId,
		PlanId:      planId,
		Status:      utils.ToPgText(utils.StringPtr("active")),
		StartAt:     startAt,
		ExpiresAt:   expiresAt,
		Provider:    provider,
		ProviderRef: providerRef,
		CreatedAt:   time.Now(),
	}
	queryStr := `
	INSERT INTO user_subscriptions
	(id, user_id, plan_id ,status, started_at,expires_at, provider,provider_ref, created_at)
	VALUES ($1, $2, $3, $4, $5,$6, $7, $8, $9)
`
	_, err := tx.Exec(
		c,
		queryStr,
		iUserSubscriptionModel.Id,
		iUserSubscriptionModel.UserId,
		iUserSubscriptionModel.PlanId,
		iUserSubscriptionModel.Status,
		iUserSubscriptionModel.StartAt,
		iUserSubscriptionModel.ExpiresAt,
		iUserSubscriptionModel.Provider,
		iUserSubscriptionModel.ProviderRef,
		iUserSubscriptionModel.CreatedAt,
	)
	return err
}

func (m *UserDao) Get(c context.Context, userId uuid.UUID) (*model.UserModel, error) {
	iUserModel := model.UserModel{}
	queryStr :=
		`
SELECT id, username, display_name ,avatar_url, email,email_verified
FROM  users
WHERE id=$1 AND status = 1
LIMIT 1
`
	err := global.DBPool.QueryRow(c, queryStr, userId).
		Scan(&iUserModel.Id, &iUserModel.UserName, &iUserModel.DisplayName, &iUserModel.AvatarURL, &iUserModel.Email, &iUserModel.EmailVerified)
	return &iUserModel, err
}
