package ssdb_test

import (
	"context"
	"fmt"

	"github.com/ssdb-go/ssdb"
)

type ssdbHook struct{}

var _ ssdb.Hook = ssdbHook{}

func (ssdbHook) BeforeProcess(ctx context.Context, cmd ssdb.Cmder) (context.Context, error) {
	fmt.Printf("starting processing: <%s>\n", cmd)
	return ctx, nil
}

func (ssdbHook) AfterProcess(ctx context.Context, cmd ssdb.Cmder) error {
	fmt.Printf("finished processing: <%s>\n", cmd)
	return nil
}

func (ssdbHook) BeforeProcessPipeline(ctx context.Context, cmds []ssdb.Cmder) (context.Context, error) {
	fmt.Printf("pipeline starting processing: %v\n", cmds)
	return ctx, nil
}

func (ssdbHook) AfterProcessPipeline(ctx context.Context, cmds []ssdb.Cmder) error {
	fmt.Printf("pipeline finished processing: %v\n", cmds)
	return nil
}

func Example_instrumentation() {
	sdb := ssdb.NewClient(&ssdb.Options{
		Addr: ":8888",
	})
	sdb.AddHook(ssdbHook{})

	sdb.Ping(ctx)
	// Output: starting processing: <ping: >
	// finished processing: <ping: PONG>
}

func ExamplePipeline_instrumentation() {
	sdb := ssdb.NewClient(&ssdb.Options{
		Addr: ":8888",
	})
	sdb.AddHook(ssdbHook{})

	sdb.Pipelined(ctx, func(pipe ssdb.Pipeliner) error {
		pipe.Ping(ctx)
		pipe.Ping(ctx)
		return nil
	})
	// Output: pipeline starting processing: [ping:  ping: ]
	// pipeline finished processing: [ping: PONG ping: PONG]
}
