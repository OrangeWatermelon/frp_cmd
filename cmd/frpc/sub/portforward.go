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
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"

	"github.com/fatedier/frp/pkg/config"
)

var lip string
func init() {
	RegisterCommonFlags(pfCmd)
	//pfCmd.PersistentFlags().StringVarP(&addr, "server", "s", "", "connect addr")
	//pfCmd.PersistentFlags().StringVarP(&lp, "lp", "", "", "local port")
	pfCmd.Flags().StringVarP(&rp, "rp", "r", "", "remote port")
	pfCmd.Flags().StringVarP(&lp, "lp", "l", "", "local port")
	pfCmd.Flags().StringVarP(&lip, "lip", "", "127.0.0.1", "local ip")
	pfCmd.Flags().StringVarP(&name, "name", "n", randName(), "name")
	pfCmd.Flags().StringVarP(&token, "token", "t", "", "token")
	//pfCmd.PersistentFlags().StringVarP(&localIP, "local_ip", "i", "127.0.0.1", "local ip")
	//
	//pfCmd.PersistentFlags().BoolVarP(&useEncryption, "ue", "", false, "use encryption")
	//pfCmd.PersistentFlags().BoolVarP(&useCompression, "uc", "", false, "use compression")
	pfCmd.MarkFlagRequired("server_addr")
	pfCmd.MarkFlagRequired("rp")
	pfCmd.MarkFlagRequired("lp")
	rootCmd.AddCommand(pfCmd)
}

var pfCmd = &cobra.Command{
	Use:   "pf",
	Short: "Run frpc with a single portforward\n\tExample: frpc fp -s 1.1.1.1:1234 -r 1234 -l 8080 [--lip 192.168.1.3] [-t f86bc7ff68aff3ad] [-n zz]",
	RunE: func(cmd *cobra.Command, args []string) error {
		ipStr, portStr, err := net.SplitHostPort(serverAddr)
		ini := "[common]\nserver_addr = "+ipStr+
			"\nserver_port = "+portStr+"\n"
			//"tls_enable = true\n"
		if token != ""{
			ini += "token = "+token +"\n"
		}
		ini +=	"\n\n["+name+
			"]\ntype = tcp\n" +
			"remote_port ="+rp+"\n" +
			"local_port ="+lp+"\n"+
			"local_ip =" +lip+"\n"+
			"use_encryption = true\n" +
				"use_compression = true\n"
			//"disable_custom_tls_first_byte = true\n"
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
