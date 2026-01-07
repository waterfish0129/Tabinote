package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/waterfish0129/Tabinote/dao"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/model"
	"github.com/waterfish0129/Tabinote/service/dto"
	"github.com/waterfish0129/Tabinote/utils"
	"go.uber.org/zap"
	"time"
)

type AuthService struct {
	userDao *dao.UserDao
	planDao *dao.PlanDao
}

func NewAuthService(
	userDao *dao.UserDao,
	planDao *dao.PlanDao,
) *AuthService {
	return &AuthService{
		userDao: userDao,
		planDao: planDao,
	}
}

func (m *AuthService) GoogleLogin(c context.Context, idToken string) (*dto.GoogleLoginResponseDTO, error) {
	/*===========================================================================================*/
	//驗證google 的id token是否有效
	googleUser, err := utils.VerifyGoogleIDToken(c, idToken)
	if err != nil {
		err = errors.New("invalid google token")
		return nil, err
	}

	/*===========================================================================================*/
	// 檢查該使用者是否已經存在資料庫
	user, err := m.userDao.GetByGoogleSub(c, googleUser.Sub)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		//如果這個錯誤不是查無資料
		return nil, err
	}
	//如果這邊錯誤不空  代表是資料庫查無此用戶 要新增用戶
	if err != nil {
		user, err = m.RegisterUserByGoogle(c, googleUser.Sub, googleUser.Email, googleUser.Name, &googleUser.Picture)
	}
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.Id, "user", user.ActivePlan.Code.String)
	return &dto.GoogleLoginResponseDTO{
		AccessToken: token,
		ExpiresIn:   int(utils.GetExpiresIn().Seconds()),
	}, nil
}

func (m *AuthService) RegisterUserByGoogle(c context.Context, googleSub string,
	email string,
	name string,
	avatar *string) (*model.UserProfileModel, error) {
	/*===========================================================================================*/
	//取得免費方案
	freeplan, err := m.planDao.FindFreeSubscriptionPlan(c)
	if err != nil {
		return nil, errors.New("FREE plan not found")
	}

	/*===========================================================================================*/
	//建立transaction
	tx, err := global.DBPool.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		// 檢查錯誤，如果錯誤是 "tx is closed" (pgx.ErrTxClosed)，代表已經成功 Commit 了，直接忽略即可
		if err != nil && err != pgx.ErrTxClosed {
			global.Logger.Error("Database transaction rollback failed",
				zap.Error(err),
				zap.String("location", "auth_service.RegisterUserByGoogle"))
		}
	}(tx, c)

	/*===========================================================================================*/
	//新建使用者
	user, err := m.userDao.CreateByGoogle(
		c,
		tx,
		googleSub,
		email,
		name,
		avatar,
	)
	if err != nil {
		return nil, err
	}

	/*===========================================================================================*/
	//創建使用者的方案
	err = m.userDao.CreateSubscription(c, tx,
		user.Id,
		freeplan.Id,
		time.Now(),
		time.Now().Add(time.Hour*24*365*100),
		utils.ToPgText(utils.StringPtr("default")),
		utils.ToPgText(nil),
	)
	if err != nil {
		return nil, err
	}
	/*===========================================================================================*/
	//Commit
	if err := tx.Commit(c); err != nil {
		return nil, err
	}

	iUserProfileModel := model.UserProfileModel{
		Id:            user.Id,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		UserName:      user.UserName,
		DisplayName:   user.DisplayName,
		AvatarURL:     user.AvatarURL,
		Status:        user.Status,
		ActivePlan: model.UserPlanModel{
			Code:      freeplan.Code,
			Name:      freeplan.Name,
			Status:    utils.ToPgText(utils.StringPtr("active")),
			StartAt:   time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 100),
		},
	}

	return &iUserProfileModel, nil
}

func (m *AuthService) Me(c context.Context, userId uuid.UUID) (*dto.UserProfileResp, error) {
	iUserProfileResp := dto.UserProfileResp{}
	user, err := m.userDao.Get(c, userId)
	if err != nil {
		return nil, err
	}
	userPlan, err := m.planDao.FindUserPlanByUserId(c, userId, "active")
	if err != nil {
		return nil, err

	}
	iUserProfileResp.Id = user.Id
	iUserProfileResp.UserName = user.UserName
	iUserProfileResp.DisplayName = user.DisplayName
	iUserProfileResp.Email = user.Email
	iUserProfileResp.AvatarURL = user.AvatarURL
	iUserProfileResp.Role = utils.ToPgText(utils.StringPtr("user"))
	iUserProfileResp.Plan.PlanCode = userPlan.Code
	iUserProfileResp.Plan.PlanName = userPlan.Name
	iUserProfileResp.Plan.StartAt = userPlan.StartAt
	iUserProfileResp.Plan.ExpireAt = userPlan.ExpiresAt

	return &iUserProfileResp, nil
}
