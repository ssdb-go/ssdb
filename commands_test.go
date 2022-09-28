package ssdb_test

import (
	"context"
	"encoding/json"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssdb-go/ssdb"
)

var _ = Describe("Commands", func() {
	ctx := context.TODO()
	var client *ssdb.Client

	BeforeEach(func() {
		client = ssdb.NewClient(ssdbOptions())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	Describe("server", func() {
		It("should Auth", func() {
			cmds, err := client.Pipelined(ctx, func(pipe ssdb.Pipeliner) error {
				pipe.Auth(ctx, "password")
				pipe.Auth(ctx, "")
				return nil
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("ERR AUTH"))
			Expect(cmds[0].Err().Error()).To(ContainSubstring("ERR AUTH"))
			Expect(cmds[1].Err().Error()).To(ContainSubstring("ERR AUTH"))

			stats := client.PoolStats()
			Expect(stats.Hits).To(Equal(uint32(1)))
			Expect(stats.Misses).To(Equal(uint32(1)))
			Expect(stats.Timeouts).To(Equal(uint32(0)))
			Expect(stats.TotalConns).To(Equal(uint32(1)))
			Expect(stats.IdleConns).To(Equal(uint32(1)))
		})
	})
})

type numberStruct struct {
	Number int
}

func (s *numberStruct) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *numberStruct) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, s)
}

func deref(viface interface{}) interface{} {
	v := reflect.ValueOf(viface)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface()
}
