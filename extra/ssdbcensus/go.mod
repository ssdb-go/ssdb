module github.com/ssdb-go/ssdb/extra/ssdbcensus

go 1.15

replace github.com/ssdb-go/ssdb => ../..

replace github.com/go-ssdb/ssdb/extra/ssdbcmd => ../ssdbcmd

require (
	github.com/ssdb-go/ssdb/extra/ssdbcmd
	github.com/ssdb-go/ssdb
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	go.opencensus.io v0.23.0
)
