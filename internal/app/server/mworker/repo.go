package mworker

import (
	"context"
	"simple-api-go/internal/pkg/utils/log"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/tasks"
)

type repo struct {
	MServer *machinery.Server
}

var _ Repository = (*repo)(nil)

// NewRepo ...
func NewRepo(mServer *machinery.Server) (returnData Repository) {
	return &repo{
		MServer: mServer,
	}
}

func (r *repo) RegisterTasks(tasksMap map[string]interface{}) (err error) {
	return r.MServer.RegisterTasks(tasksMap)
}

func (r *repo) RegisterPeriodicTask(spec, name string, signature *tasks.Signature) (err error) {
	return r.MServer.RegisterPeriodicTask(spec, name, signature)
}

func (r *repo) SendSingleTaskWithContext(ctx context.Context, taskSignature tasks.Signature) (err error) {
	_, e := r.MServer.SendTaskWithContext(ctx, &taskSignature)
	err = e
	log.AddErr(ctx, err)
	return
}

func (r *repo) SendGroupTaskWithContext(ctx context.Context, taskSignature []tasks.Signature) (err error) {
	countSlice := len(taskSignature)
	taskSlice := make([]*tasks.Signature, countSlice)

	for k := range taskSignature {
		taskSlice = append(taskSlice, &taskSignature[k])
	}

	group, err := tasks.NewGroup(taskSlice...)
	if err != nil {
		log.AddErr(ctx, err)
		return
	}
	_, err = r.MServer.SendGroupWithContext(ctx, group, countSlice)
	if err != nil {
		log.AddErr(ctx, err)
	}

	return
}
