package ssdbcensus

import (
	"context"

	"go.opencensus.io/trace"

	"github.com/ssdb-go/ssdb"
	"github.com/ssdb-go/ssdb/extra/ssdbcmd"
)

type TracingHook struct{}

var _ ssdb.Hook = (*TracingHook)(nil)

func NewTracingHook() *TracingHook {
	return new(TracingHook)
}

func (TracingHook) BeforeProcess(ctx context.Context, cmd ssdb.Cmder) (context.Context, error) {
	ctx, span := trace.StartSpan(ctx, cmd.FullName())
	span.AddAttributes(trace.StringAttribute("db.system", "ssdb"),
		trace.StringAttribute("ssdb.cmd", ssdbcmd.CmdString(cmd)))

	return ctx, nil
}

func (TracingHook) AfterProcess(ctx context.Context, cmd ssdb.Cmder) error {
	span := trace.FromContext(ctx)
	if err := cmd.Err(); err != nil {
		recordErrorOnOCSpan(ctx, span, err)
	}
	span.End()
	return nil
}

func (TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []ssdb.Cmder) (context.Context, error) {
	return ctx, nil
}

func (TracingHook) AfterProcessPipeline(ctx context.Context, cmds []ssdb.Cmder) error {
	return nil
}

func recordErrorOnOCSpan(ctx context.Context, span *trace.Span, err error) {
	if err != ssdb.Nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate([]trace.Attribute{trace.StringAttribute("Error", "ssdb error")}, err.Error())
	}
}
