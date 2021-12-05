// Copyright 2018 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sub

import (
	"crypto/rand"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"

	"github.com/fatedier/frp/pkg/config"
)
var (
	addr string
	lp string
	rp string
	name string
	pwd string
)
func randName() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
func init() {
	RegisterCommonFlags(socks5Cmd)
	//socks5Cmd.PersistentFlags().StringVarP(&addr, "server", "s", "", "connect addr")
	//socks5Cmd.PersistentFlags().StringVarP(&lp, "lp", "", "", "local port")
	socks5Cmd.Flags().StringVarP(&rp, "rp", "r", "", "remote port")
	socks5Cmd.Flags().StringVarP(&name, "name", "n", randName(), "name")
	socks5Cmd.Flags().StringVarP(&user, "user", "u", "", "user")
	socks5Cmd.Flags().StringVarP(&pwd, "pwd", "p", "", "pwd")
	socks5Cmd.Flags().StringVarP(&token, "token", "t", "", "token")
	//socks5Cmd.PersistentFlags().StringVarP(&localIP, "local_ip", "i", "127.0.0.1", "local ip")
	//
	//socks5Cmd.PersistentFlags().BoolVarP(&useEncryption, "ue", "", false, "use encryption")
	//socks5Cmd.PersistentFlags().BoolVarP(&useCompression, "uc", "", false, "use compression")
	socks5Cmd.MarkFlagRequired("server_addr")
	socks5Cmd.MarkFlagRequired("rp")
	rootCmd.AddCommand(socks5Cmd)
}

var socks5Cmd = &cobra.Command{
	Use:   "socks",
	Short: "Run frpc with a single socks5 proxy\n\tExample: frpc socks -s 1.1.1.1:1234 -r 1234 [-t f86bc7ff68aff3ad] [-n zz] [-u z] [-p z] ",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipStr, portStr, err := net.SplitHostPort(serverAddr)
		ini := "[common]\nserver_addr = "+ipStr+
			"\nserver_port = "+portStr+"\n"
		if token != ""{
			ini += "token = "+token +"\n"
		}
		ini += "\n\n["+name+
			"]\ntype = tcp\n" +
			"remote_port ="+rp+"\n" +
			"plugin = socks5\n" +
			"use_encryption = true\n" +
			"use_compression = true\n"
		if user != ""{
			ini += "plugin_user = "+user+"\n"
		}
		if pwd != ""{
			ini += "plugin_passwd ="+pwd+"\n"
		}
		cfg, pxyCfgs, visitorCfgs, err := config.ParseClientConfig1(ini)
		if err != nil {
			return err
		}
		err = startService(cfg, pxyCfgs, visitorCfgs, "")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return nil
	},
}
