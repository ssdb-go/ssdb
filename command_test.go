package ssdb_test

import (
	"github.com/ssdb-go/ssdb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cmd", func() {
	var client *ssdb.Client

	BeforeEach(func() {
		client = ssdb.NewClient(ssdbOptions())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})

	It("dbsize", func() {
		set := client.DBSize(ctx)
		Expect(set.String()).To(Equal("set foo bar: OK"))

		get := client.Get(ctx, "foo")
		Expect(get.String()).To(Equal("get foo: bar"))
	})
})
