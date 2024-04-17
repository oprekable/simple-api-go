package mworker

import (
	"context"

	"github.com/RichardKnop/machinery/v2/tasks"
)

type Repository interface {
	RegisterTasks(tasksMap map[string]interface{}) (err error)
	RegisterPeriodicTask(spec, name string, signature *tasks.Signature) (err error)
	SendSingleTaskWithContext(ctx context.Context, taskSignature tasks.Signature) (err error)
	SendGroupTaskWithContext(ctx context.Context, taskSignature []tasks.Signature) (err error)
}

type Service interface {
	TriggerJob(str string) (err error)
	TriggerJobWithContext(ctx context.Context, str string) (err error)
}
