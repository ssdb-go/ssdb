package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/ssdb-go/ssdb/internal/customvet/checks/setval"
)

func main() {
	multichecker.Main(
		setval.Analyzer,
	)
}
