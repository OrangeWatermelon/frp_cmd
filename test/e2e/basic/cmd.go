package basic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatedier/frp/test/e2e/framework"
	"github.com/fatedier/frp/test/e2e/pkg/request"

	. "github.com/onsi/ginkgo"
)

const (
	ConfigValidStr = "syntax is ok"
)

var _ = Describe("[Feature: Cmd]", func() {
	f := framework.NewDefaultFramework()

	Describe("Verify", func() {
		It("frps valid", func() {
			path := f.GenerateConfigFile(`
			[common]
			bind_addr = 0.0.0.0
			bind_port = 7000
			`)
			_, output, err := f.RunFrps("verify", "-c", path)
			framework.ExpectNoError(err)
			framework.ExpectTrue(strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
		It("frps invalid", func() {
			path := f.GenerateConfigFile(`
			[common]
			bind_addr = 0.0.0.0
			bind_port = 70000
			`)
			_, output, err := f.RunFrps("verify", "-c", path)
			framework.ExpectNoError(err)
			framework.ExpectTrue(!strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
		It("frpc valid", func() {
			path := f.GenerateConfigFile(`
			[common]
			server_addr = 0.0.0.0
			server_port = 7000
			`)
			_, output, err := f.RunFrpc("verify", "-c", path)
			framework.ExpectNoError(err)
			framework.ExpectTrue(strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
		It("frpc invalid", func() {
			path := f.GenerateConfigFile(`
			[common]
			server_addr = 0.0.0.0
			server_port = 7000
			protocol = invalid
			`)
			_, output, err := f.RunFrpc("verify", "-c", path)
			framework.ExpectNoError(err)
			framework.ExpectTrue(!strings.Contains(output, ConfigValidStr), "output: %s", output)
		})
	})

	Describe("Single proxy", func() {
		It("TCP", func() {
			serverPort := f.AllocPort()
			_, _, err := f.RunFrps("-t", "123", "-p", strconv.Itoa(serverPort))
			framework.ExpectNoError(err)

			localPort := f.PortByName(framework.TCPEchoServerPort)
			remotePort := f.AllocPort()
			_, _, err = f.RunFrpc("tcp", "-s", fmt.Sprintf("127.0.0.1:%d", serverPort), "-t", "123", "-u", "test",
				"-l", strconv.Itoa(localPort), "-r", strconv.Itoa(remotePort), "-n", "tcp_test")
			framework.ExpectNoError(err)

			framework.NewRequestExpect(f).Port(remotePort).Ensure()
		})

		It("UDP", func() {
			serverPort := f.AllocPort()
			_, _, err := f.RunFrps("-t", "123", "-p", strconv.Itoa(serverPort))
			framework.ExpectNoError(err)

			localPort := f.PortByName(framework.UDPEchoServerPort)
			remotePort := f.AllocPort()
			_, _, err = f.RunFrpc("udp", "-s", fmt.Sprintf("127.0.0.1:%d", serverPort), "-t", "123", "-u", "test",
				"-l", strconv.Itoa(localPort), "-r", strconv.Itoa(remotePort), "-n", "udp_test")
			framework.ExpectNoError(err)

			framework.NewRequestExpect(f).Protocol("udp").
				Port(remotePort).Ensure()
		})

		It("HTTP", func() {
			serverPort := f.AllocPort()
			vhostHTTPPort := f.AllocPort()
			_, _, err := f.RunFrps("-t", "123", "-p", strconv.Itoa(serverPort), "--vhost_http_port", strconv.Itoa(vhostHTTPPort))
			framework.ExpectNoError(err)

			_, _, err = f.RunFrpc("http", "-s", "127.0.0.1:"+strconv.Itoa(serverPort), "-t", "123", "-u", "test",
				"-n", "udp_test", "-l", strconv.Itoa(f.PortByName(framework.HTTPSimpleServerPort)),
				"--custom_domain", "test.example.com")
			framework.ExpectNoError(err)

			framework.NewRequestExpect(f).Port(vhostHTTPPort).
				RequestModify(func(r *request.Request) {
					r.HTTP().HTTPHost("test.example.com")
				}).
				Ensure()
		})
	})
})
