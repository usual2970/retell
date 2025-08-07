package schedule

import (
	"github.com/usual2970/retell/pkg/logger"
	"github.com/usual2970/retell/pkg/schedule"
)

type OrderService interface {
	TakeOrderFees() error
}

type OrderFeeScheduler struct {
	service OrderService
}

func NewOrderFeeScheduler(service OrderService) {
	s := &OrderFeeScheduler{
		service: service,
	}

	// 每月 1 号 0 点 0 分执行
	schedule.AddJob("0 8 1 * *", s.Handle)
}

func (s *OrderFeeScheduler) Handle() {
	log := logger.WithField("module", "order_fee_scheduler")
	log.Info("Order fee scheduler started")
	if err := s.service.TakeOrderFees(); err != nil {
		// 处理错误
		log.WithField("error", err).Error("Failed to take order fees")
		return
	}
	log.Info("Order fee scheduler finished successfully")
}
