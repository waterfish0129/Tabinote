package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/waterfish0129/Tabinote/global"
	"github.com/waterfish0129/Tabinote/model"
)

type PlanDao struct {
}

func NewPlanDao() *PlanDao {
	return &PlanDao{}
}

func (m *PlanDao) FindFreeSubscriptionPlan(c context.Context) (*model.SubscriptionPlanModel, error) {
	var plan model.SubscriptionPlanModel
	queryStr :=
		`
	SELECT id, code, name, price, currency, interval, created_at
	FROM subscription_plans
	WHERE code = $1
	`
	err := global.DBPool.QueryRow(c, queryStr, "FREE").Scan(
		&plan.Id,
		&plan.Code,
		&plan.Name,
		&plan.Price,
		&plan.Currency,
		&plan.Interval,
		&plan.CreatedAt)

	return &plan, err
}

func (m *PlanDao) FindAllSubscriptionPlans(c context.Context) ([]*model.SubscriptionPlanModel, error) {
	queryStr :=
		`
SELECT id, code, name, price, currency, interval, created_at
FROM subscription_plans
`
	rows, err := global.DBPool.Query(c, queryStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]*model.SubscriptionPlanModel, 0)

	for rows.Next() {
		var plan model.SubscriptionPlanModel
		if err := rows.Scan(
			&plan.Id,
			&plan.Code,
			&plan.Name,
			&plan.Price,
			&plan.Currency,
			&plan.Interval,
			&plan.CreatedAt,
		); err != nil {
			return nil, err
		}

		plans = append(plans, &plan)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return plans, nil
}

func (m *PlanDao) FindUserPlanByUserId(c context.Context, userId uuid.UUID, status string) (*model.UserPlanModel, error) {
	plan := model.UserPlanModel{}
	queryStr :=
		`
		SELECT sp.code, sp.name,us.status, us.started_at, us.expires_at
		FROM user_subscriptions us
		join subscription_plans sp on sp.id = us.plan_id
		where us.user_id =$1 AND us.status =$2
		LIMIT 1
		`

	err := global.DBPool.QueryRow(c, queryStr, userId, status).Scan(
		&plan.Code,
		&plan.Name,
		&plan.Status,
		&plan.StartAt,
		&plan.ExpiresAt,
	)
	return &plan, err
}
