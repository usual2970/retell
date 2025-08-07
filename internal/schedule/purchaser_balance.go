package schedule

import (
	"context"

	"github.com/usual2970/retell/pkg/logger"
	"github.com/usual2970/retell/pkg/schedule"
)

type PurchaserService interface {
	LowBalanceNotify(ctx context.Context) error
}

type PurchaserBalanceSchedule struct {
	service PurchaserService
}

func NewPurchaserBalanceSchedule(service PurchaserService) {
	s := &PurchaserBalanceSchedule{
		service: service,
	}

	// 每 日 0 点 0 分执行
	schedule.AddJob("0 0 * * *", s.Handle)
}

func (s *PurchaserBalanceSchedule) Handle() {
	log := logger.WithField("module", "purchaser_balance_scheduler")
	ctx := context.Background()
	if err := s.service.LowBalanceNotify(ctx); err != nil {
		log.WithField("err", err).Error("failed to notify low balance")
	} else {
		log.Info("Purchaser balance scheduler finished successfully")
	}

}
