package ssdb_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ssdb-go/ssdb"
)

const (
	ssdbPort          = "8888"
	ssdbAddr          = ":" + ssdbPort
	ssdbSecondaryPort = "6381"
)

const (
	ringShard1Port = "6390"
	ringShard2Port = "6391"
	ringShard3Port = "6392"
)

const (
	sentinelName       = "mymaster"
	sentinelMasterPort = "9123"
	sentinelSlave1Port = "9124"
	sentinelSlave2Port = "9125"
	sentinelPort1      = "9126"
	sentinelPort2      = "9127"
	sentinelPort3      = "9128"
)

var (
	sentinelAddrs = []string{":" + sentinelPort1, ":" + sentinelPort2, ":" + sentinelPort3}

	processes map[string]*ssdbProcess

	ssdbMain                                       *ssdbProcess
	ringShard1, ringShard2, ringShard3             *ssdbProcess
	sentinelMaster, sentinelSlave1, sentinelSlave2 *ssdbProcess
	sentinel1, sentinel2, sentinel3                *ssdbProcess
)

func registerProcess(port string, p *ssdbProcess) {
	if processes == nil {
		processes = make(map[string]*ssdbProcess)
	}
	processes[port] = p
}

var _ = BeforeSuite(func() {
	var err error

	ssdbMain, err = startSsdb(ssdbPort)
	Expect(err).NotTo(HaveOccurred())

	ringShard1, err = startSsdb(ringShard1Port)
	Expect(err).NotTo(HaveOccurred())

	ringShard2, err = startSsdb(ringShard2Port)
	Expect(err).NotTo(HaveOccurred())

	ringShard3, err = startSsdb(ringShard3Port)
	Expect(err).NotTo(HaveOccurred())

	sentinelMaster, err = startSsdb(sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())

	Expect(err).NotTo(HaveOccurred())

	sentinelSlave1, err = startSsdb(
		sentinelSlave1Port, "--slaveof", "127.0.0.1", sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())

	sentinelSlave2, err = startSsdb(
		sentinelSlave2Port, "--slaveof", "127.0.0.1", sentinelMasterPort)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	for _, p := range processes {
		Expect(p.Close()).NotTo(HaveOccurred())
	}
	processes = nil
})

func TestGinkgoSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-ssdb")
}

//------------------------------------------------------------------------------

func ssdbOptions() *ssdb.Options {
	return &ssdb.Options{
		Addr: ssdbAddr,
		DB:   15,

		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,

		MaxRetries: -1,

		PoolSize:        10,
		PoolTimeout:     30 * time.Second,
		ConnMaxIdleTime: time.Minute,
	}
}

func performAsync(n int, cbs ...func(int)) *sync.WaitGroup {
	var wg sync.WaitGroup
	for _, cb := range cbs {
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func(cb func(int), i int) {
				defer GinkgoRecover()
				defer wg.Done()

				cb(i)
			}(cb, i)
		}
	}
	return &wg
}

func perform(n int, cbs ...func(int)) {
	wg := performAsync(n, cbs...)
	wg.Wait()
}

func eventually(fn func() error, timeout time.Duration) error {
	errCh := make(chan error, 1)
	done := make(chan struct{})
	exit := make(chan struct{})

	go func() {
		for {
			err := fn()
			if err == nil {
				close(done)
				return
			}

			select {
			case errCh <- err:
			default:
			}

			select {
			case <-exit:
				return
			case <-time.After(timeout / 100):
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		close(exit)
		select {
		case err := <-errCh:
			return err
		default:
			return fmt.Errorf("timeout after %s without an error", timeout)
		}
	}
}

func execCmd(name string, args ...string) (*os.Process, error) {
	cmd := exec.Command(name, args...)
	if testing.Verbose() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Process, cmd.Start()
}

func connectTo(port string) (*ssdb.Client, error) {
	client := ssdb.NewClient(&ssdb.Options{
		Addr:       ":" + port,
		MaxRetries: -1,
	})

	err := eventually(func() error {
		return client.Ping(ctx).Err()
	}, 30*time.Second)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type ssdbProcess struct {
	*os.Process
	*ssdb.Client
}

func (p *ssdbProcess) Close() error {
	if err := p.Kill(); err != nil {
		return err
	}

	err := eventually(func() error {
		if err := p.Client.Ping(ctx).Err(); err != nil {
			return nil
		}
		return errors.New("client is not shutdown")
	}, 10*time.Second)
	if err != nil {
		return err
	}

	p.Client.Close()
	return nil
}

var (
	ssdbServerBin, _    = filepath.Abs(filepath.Join("testdata", "ssdb", "src", "ssdb-server"))
	ssdbServerConf, _   = filepath.Abs(filepath.Join("testdata", "ssdb", "ssdb.conf"))
	ssdbSentinelConf, _ = filepath.Abs(filepath.Join("testdata", "ssdb", "sentinel.conf"))
)

func ssdbDir(port string) (string, error) {
	dir, err := filepath.Abs(filepath.Join("testdata", "instances", port))
	if err != nil {
		return "", err
	}
	if err := os.RemoveAll(dir); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o775); err != nil {
		return "", err
	}
	return dir, nil
}

func startSsdb(port string, args ...string) (*ssdbProcess, error) {
	dir, err := ssdbDir(port)
	if err != nil {
		return nil, err
	}

	if err := exec.Command("cp", "-f", ssdbServerConf, dir).Run(); err != nil {
		return nil, err
	}

	baseArgs := []string{filepath.Join(dir, "ssdb.conf"), "--port", port, "--dir", dir}
	process, err := execCmd(ssdbServerBin, append(baseArgs, args...)...)
	if err != nil {
		return nil, err
	}

	client, err := connectTo(port)
	if err != nil {
		process.Kill()
		return nil, err
	}

	p := &ssdbProcess{process, client}
	registerProcess(port, p)
	return p, nil
}

//------------------------------------------------------------------------------

type badConnError string

func (e badConnError) Error() string   { return string(e) }
func (e badConnError) Timeout() bool   { return true }
func (e badConnError) Temporary() bool { return false }

type badConn struct {
	net.TCPConn

	readDelay, writeDelay time.Duration
	readErr, writeErr     error
}

var _ net.Conn = &badConn{}

func (cn *badConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (cn *badConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (cn *badConn) Read([]byte) (int, error) {
	if cn.readDelay != 0 {
		time.Sleep(cn.readDelay)
	}
	if cn.readErr != nil {
		return 0, cn.readErr
	}
	return 0, badConnError("bad connection")
}

func (cn *badConn) Write([]byte) (int, error) {
	if cn.writeDelay != 0 {
		time.Sleep(cn.writeDelay)
	}
	if cn.writeErr != nil {
		return 0, cn.writeErr
	}
	return 0, badConnError("bad connection")
}

//------------------------------------------------------------------------------

type hook struct {
	beforeProcess func(ctx context.Context, cmd ssdb.Cmder) (context.Context, error)
	afterProcess  func(ctx context.Context, cmd ssdb.Cmder) error

	beforeProcessPipeline func(ctx context.Context, cmds []ssdb.Cmder) (context.Context, error)
	afterProcessPipeline  func(ctx context.Context, cmds []ssdb.Cmder) error
}

func (h *hook) BeforeProcess(ctx context.Context, cmd ssdb.Cmder) (context.Context, error) {
	if h.beforeProcess != nil {
		return h.beforeProcess(ctx, cmd)
	}
	return ctx, nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd ssdb.Cmder) error {
	if h.afterProcess != nil {
		return h.afterProcess(ctx, cmd)
	}
	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmds []ssdb.Cmder) (context.Context, error) {
	if h.beforeProcessPipeline != nil {
		return h.beforeProcessPipeline(ctx, cmds)
	}
	return ctx, nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmds []ssdb.Cmder) error {
	if h.afterProcessPipeline != nil {
		return h.afterProcessPipeline(ctx, cmds)
	}
	return nil
}
