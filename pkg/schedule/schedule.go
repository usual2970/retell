package schedule

import (
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

var scheduler gocron.Scheduler
var schedulerOnce sync.Once

func GetSchedule() gocron.Scheduler {
	schedulerOnce.Do(func() {
		s, err := gocron.NewScheduler(gocron.WithLocation(time.UTC))
		if err != nil {
			panic(err)
		}
		scheduler = s
	})

	return scheduler
}

func StartSchedule() {
	s := GetSchedule()
	s.Start()
}

func StopSchedule() {
	s := GetSchedule()
	s.Shutdown()
}

func AddJob(spec string, jobFunc func()) (uuid.UUID, error) {
	s := GetSchedule()
	job, err := s.NewJob(
		gocron.CronJob(spec, false),
		gocron.NewTask(jobFunc),
	)
	if err != nil {
		return uuid.Nil, err
	}

	return job.ID(), nil
}
