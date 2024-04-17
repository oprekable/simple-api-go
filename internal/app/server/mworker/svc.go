package mworker

import (
	"context"
	"fmt"
	"simple-api-go/internal/app/config"
	"time"

	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/robfig/cron/v3"
)

type svc struct {
	MachineryRepo Repository
	Config        *config.Data
	Environment   string
	MachineID     string
}

var _ Service = (*svc)(nil)

func NewSvc(env string, machineID string, config *config.Data, machineryRepo Repository) (returnData Service) {
	return &svc{
		Environment:   env,
		MachineID:     machineID,
		Config:        config,
		MachineryRepo: machineryRepo,
	}
}

func (s *svc) TriggerJob(str string) (err error) {
	return s.TriggerJobWithContext(context.Background(), str)
}

func (s *svc) TriggerJobWithContext(ctx context.Context, _ string) (err error) {
	for _, v := range s.Config.CronJob {
		if v.IsEnabled {
			var args []tasks.Arg

			for _, vv := range v.Signature.TaskArgs {
				args = append(
					args,
					tasks.Arg{
						Type:  "string",
						Value: vv.Value,
					},
				)
			}

			//check spec
			schedule, er := cron.ParseStandard(v.Spec)

			if er != nil {
				err = er
				return
			}

			nextEta := schedule.Next(time.Now())
			signature := tasks.Signature{
				UUID:                        fmt.Sprintf("%s-%s", v.Signature.UUID, s.Environment),
				Name:                        fmt.Sprintf("%s-%s", v.Signature.Name, s.Environment),
				Args:                        args,
				ETA:                         &nextEta,
				RetryCount:                  v.Signature.RetryCount,
				RetryTimeout:                v.Signature.RetryTimeout,
				IgnoreWhenTaskNotRegistered: true,
			}

			err = s.MachineryRepo.SendSingleTaskWithContext(ctx, signature)

			if err != nil {
				return
			}
		}
	}

	return
}
