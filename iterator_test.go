package ssdb_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssdb-go/ssdb"
)

var _ = Describe("ScanIterator", func() {
	var client *ssdb.Client
	BeforeEach(func() {
		client = ssdb.NewClient(ssdbOptions())
	})

	AfterEach(func() {
		Expect(client.Close()).NotTo(HaveOccurred())
	})
})
