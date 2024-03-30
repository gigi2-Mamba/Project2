package repository

import (
	"context"
	"project0/internal/domain"
	"project0/internal/repository/dao"
	"time"
)

/*
User: society-programmer
Date: 2024/2/23  周五
Time: 16:03
*/

type CronJobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, id int64) error
	UpdateTime(ctx context.Context, d int64) error
	UpdateNextTime(ctx context.Context, id int64, time time.Time) error
}

// mysql抢占
type PreemptRepository struct {
	dao dao.JobDAO
}

func (p *PreemptRepository) UpdateNextTime(ctx context.Context, id int64, time time.Time) error {

	return p.dao.UpdateNextTime(ctx, id, time)
}

func NewPreemptRepository(dao dao.JobDAO) *PreemptRepository {
	return &PreemptRepository{dao: dao}
}

func (p *PreemptRepository) UpdateTime(ctx context.Context, id int64) error {
	return p.dao.UpdateTime(ctx, id)
}
func (p *PreemptRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.dao.Preempt(ctx)

	if err != nil {
		return domain.Job{}, err
	}

	// 像简单拼凑可以这样写.  复杂才走domain 到entity，或者说开发初期。 就这样简单拼凑
	return domain.Job{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
	}, err
}

func (p *PreemptRepository) Release(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
