package ssdb_test

import (
	"context"
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssdb-go/ssdb"
)

type ssdbHookError struct {
	ssdb.Hook
}

var _ ssdb.Hook = ssdbHookError{}

func (ssdbHookError) BeforeProcess(ctx context.Context, cmd ssdb.Cmder) (context.Context, error) {
	return ctx, nil
}

func (ssdbHookError) AfterProcess(ctx context.Context, cmd ssdb.Cmder) error {
	return nil
}

func TestHookError(t *testing.T) {
	sdb := ssdb.NewClient(&ssdb.Options{
		Addr: ":8888",
	})
	sdb.AddHook(ssdbHookError{})

	err := sdb.Ping(ctx).Err()
	if err != nil {
		t.Fatalf("got ann error %v", err)
	}
	//fmt.Println(sdb.Ping(ctx).String())
}

//------------------------------------------------------------------------------

var _ = Describe("Client", func() {
	var client *ssdb.Client

	AfterEach(func() {
		client.Close()
	})

	It("should Stringer", func() {
		Expect(client.String()).To(Equal("Ssdb<:8888 db:15>"))
	})

	It("supports context", func() {
		ctx, cancel := context.WithCancel(ctx)
		cancel()

		err := client.Ping(ctx).Err()
		Expect(err).To(MatchError("context canceled"))
	})
})

var _ = Describe("Client timeout", func() {
	var opt *ssdb.Options
	var client *ssdb.Client

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	testTimeout := func() {
		It("Ping timeouts", func() {
			err := client.Ping(ctx).Err()
			Expect(err).To(HaveOccurred())
			Expect(err.(net.Error).Timeout()).To(BeTrue())
		})
	}

	Context("read timeout", func() {
		BeforeEach(func() {
			opt = ssdbOptions()
			opt.ReadTimeout = time.Nanosecond
			opt.WriteTimeout = -1
			client = ssdb.NewClient(opt)
		})

		testTimeout()
	})

	Context("write timeout", func() {
		BeforeEach(func() {
			opt = ssdbOptions()
			opt.ReadTimeout = -1
			opt.WriteTimeout = time.Nanosecond
			client = ssdb.NewClient(opt)
		})

		testTimeout()
	})
})
