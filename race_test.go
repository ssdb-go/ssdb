package ssdb_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssdb-go/ssdb"
)

var _ = Describe("races", func() {
	var client *ssdb.Client
	var C, N int

	BeforeEach(func() {
		client = ssdb.NewClient(ssdbOptions())

		C, N = 10, 1000
		if testing.Short() {
			C = 4
			N = 100
		}
	})

	AfterEach(func() {
		err := client.Close()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should handle big vals in Set", func() {
		C, N = 4, 100

		bigVal := bigVal()
		perform(C, func(id int) {
			for i := 0; i < N; i++ {
				err := client.Set(ctx, "key", bigVal, 0).Err()
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	It("should select db", func() {
		err := client.Set(ctx, "db", 1, 0).Err()
		Expect(err).NotTo(HaveOccurred())

		perform(C, func(id int) {
			opt := ssdbOptions()
			opt.DB = id
			client := ssdb.NewClient(opt)
			for i := 0; i < N; i++ {
				err := client.Set(ctx, "db", id, 0).Err()
				Expect(err).NotTo(HaveOccurred())

				n, err := client.Get(ctx, "db").Int64()
				Expect(err).NotTo(HaveOccurred())
				Expect(n).To(Equal(int64(id)))
			}
			err := client.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		n, err := client.Get(ctx, "db").Int64()
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(int64(1)))
	})

	It("should select DB with read timeout", func() {
		perform(C, func(id int) {
			opt := ssdbOptions()
			opt.DB = id
			opt.ReadTimeout = time.Nanosecond
			client := ssdb.NewClient(opt)

			perform(C, func(id int) {
				err := client.Ping(ctx).Err()
				Expect(err).To(HaveOccurred())
				Expect(err.(net.Error).Timeout()).To(BeTrue())
			})

			err := client.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func bigVal() []byte {
	return bytes.Repeat([]byte{'*'}, 1<<17) // 128kb
}
