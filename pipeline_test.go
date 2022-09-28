package ssdb_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssdb-go/ssdb"
)

var _ = Describe("pipelining", func() {
	var client *ssdb.Client
	var pipe *ssdb.Pipeline

	BeforeEach(func() {
		client = ssdb.NewClient(ssdbOptions())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	It("supports block style", func() {
		var get *ssdb.Cmd
		cmds, err := client.Pipelined(ctx, func(pipe ssdb.Pipeliner) error {
			get = pipe.Get(ctx, "foo")
			return nil
		})
		Expect(err).To(Equal(ssdb.Nil))
		Expect(cmds).To(HaveLen(1))
		Expect(cmds[0]).To(Equal(get))
		Expect(get.Err()).To(Equal(ssdb.Nil))
		Expect(get.Val()).To(Equal(""))
	})

	assertPipeline := func() {
		It("returns no errors when there are no commands", func() {
			_, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("discards queued commands", func() {
			pipe.Get(ctx, "key")
			pipe.Discard()
			cmds, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmds).To(BeNil())
		})

		It("handles val/err", func() {
			err := client.Set(ctx, "key", "value", 0).Err()
			Expect(err).NotTo(HaveOccurred())

			get := pipe.Get(ctx, "key")
			cmds, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmds).To(HaveLen(1))

			val, err := get.Result()
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal("value"))
		})

		It("supports custom command", func() {
			pipe.Do(ctx, "ping")
			cmds, err := pipe.Exec(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(cmds).To(HaveLen(1))
		})

		It("handles large pipelines", func() {
			for callCount := 1; callCount < 16; callCount++ {
				for i := 1; i <= callCount; i++ {
					pipe.SetNX(ctx, strconv.Itoa(i)+"_key", strconv.Itoa(i)+"_value")
				}

				cmds, err := pipe.Exec(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(cmds).To(HaveLen(callCount))
				for _, cmd := range cmds {
					Expect(cmd).To(BeAssignableToTypeOf(&ssdb.Cmd{}))
				}
			}
		})
	}

	Describe("Pipeline", func() {
		BeforeEach(func() {
			pipe = client.Pipeline().(*ssdb.Pipeline)
		})

		assertPipeline()
	})
})
