package doctor

import (
	"context"
	"time"

	"multienv/common/doctor/model"
)

func Measure(ctx context.Context, id, group string, run func(context.Context) (string, error)) model.Check {
	start := time.Now()
	message, err := run(ctx)

	check := model.Check{ID: id, Group: group, Status: model.OK, Message: message}
	if err != nil {
		check.Status = model.Fail
		check.Message = err.Error()
	}

	check.DurationMS = time.Since(start).Milliseconds()

	return check
}
