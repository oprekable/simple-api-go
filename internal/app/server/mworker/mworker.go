package mworker

import (
	"errors"
	"fmt"
	appContext "simple-api-go/internal/app/context"
	"simple-api-go/internal/app/server"
	"simple-api-go/internal/pkg/utils/log"

	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/samber/lo"
	"github.com/segmentio/encoding/json"
)

const (
	MachineryCronJobTrigger = "MachineryCronJobTrigger"
)

type machineryWorker struct {
	App      appContext.IAppContext
	Repo     Repository
	Svc      Service
	TasksMap map[string]interface{}
}

var _ server.IServer = (*machineryWorker)(nil)

// NewMachineryWorker create object server
func NewMachineryWorker(appCtx appContext.IAppContext, repo Repository, svc Service) server.IServer {
	rd := &machineryWorker{
		App:  appCtx,
		Repo: repo,
		Svc:  svc,
	}

	rd.TasksMap = lo.Assign[string, interface{}](
		map[string]interface{}{
			fmt.Sprintf("%s-%s", MachineryCronJobTrigger, rd.App.GetEnvironment()): rd.Svc.TriggerJob,
		},
		rd.App.GetMachineryTaskMap(),
	)

	return rd
}

func (m *machineryWorker) AddHandlers(_ ...interface{}) {
}

func (m *machineryWorker) StartApp() {
	if !m.App.IsMachineryActive() {
		return
	}

	m.App.GetEg().Go(func() error {
		ctx := m.App.GetCtx()
		if m.App.GetMachineryServer() == nil {
			err := errors.New("MachineryServer was nil")
			log.AddErr(ctx, err)
			return err
		}

		registerPeriodicTask := func() (err error) {
			var args []tasks.Arg
			args = append(
				args,
				tasks.Arg{
					Type:  "string",
					Value: "",
				},
			)

			signatureUUID := fmt.Sprintf(
				"%s-%s",
				MachineryCronJobTrigger,
				m.App.GetEnvironment(),
			)

			signatureName := fmt.Sprintf(
				"%s-%s",
				MachineryCronJobTrigger,
				m.App.GetEnvironment(),
			)

			signature := &tasks.Signature{
				UUID:                        signatureUUID,
				Name:                        signatureName,
				Args:                        args,
				RetryCount:                  3,
				RetryTimeout:                5,
				IgnoreWhenTaskNotRegistered: true,
			}

			err = m.Repo.RegisterPeriodicTask("*/1 * * * ?", fmt.Sprintf("%s-%s-%s", "Job", MachineryCronJobTrigger, m.App.GetEnvironment()), signature)
			return
		}

		log.AddErr(ctx, m.Repo.RegisterTasks(m.TasksMap))
		log.AddErr(ctx, registerPeriodicTask())
		log.Msg(ctx, "[start] machinery server")

		worker := m.App.GetMachineryServer().NewCustomQueueWorker(
			fmt.Sprintf("%s-%s-%s-%d", m.App.GetEnvironment(), m.App.GetMachineID(), m.App.GetConfig().App.Host, m.App.GetConfig().App.Port),
			10000,
			m.App.GetDefaultQueueName(),
		)

		log.Msg(ctx, fmt.Sprintf("[start] machinery worker : %s", m.App.GetMachineID()))
		// Here we inject some custom code for error handling,
		// start and end of task hooks, useful for metrics for example.
		errorHandler := func(err error) {
			log.Err(ctx, "an error", err)
		}

		preTaskHandler := func(signature *tasks.Signature) {
		}

		postTaskHandler := func(signature *tasks.Signature) {
			byt, _ := json.Marshal(signature)
			log.Msg(ctx, string(byt))
		}

		worker.SetPostTaskHandler(postTaskHandler)
		worker.SetErrorHandler(errorHandler)
		worker.SetPreTaskHandler(preTaskHandler)
		err := worker.Launch()
		log.Msg(ctx, fmt.Sprintf("[shutdown] machinery worker [%v]", err.Error()))
		return err
	})
}

func (m *machineryWorker) StopApp() {
}
