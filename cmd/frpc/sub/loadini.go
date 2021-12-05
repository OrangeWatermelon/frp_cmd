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
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/fatedier/frp/pkg/config"
)
var ini string
func init() {
	RegisterCommonFlags(loadiniCmd)
	//loadiniCmd.PersistentFlags().StringVarP(&addr, "server", "s", "", "connect addr")
	//loadiniCmd.PersistentFlags().StringVarP(&lp, "lp", "", "", "local port")
	loadiniCmd.Flags().StringVarP(&ini, "ini", "i", "", "base64 ini")

	loadiniCmd.MarkFlagRequired("ini")
	rootCmd.AddCommand(loadiniCmd)
}

var loadiniCmd = &cobra.Command{
	Use:   "loadini",
	Short: "Run frpc and load ini from base64\n\tExample: frpc loadini -i <base64ini>",
	RunE: func(cmd *cobra.Command, args []string) error {
		ini, err := base64.StdEncoding.DecodeString(ini)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfg, pxyCfgs, visitorCfgs, err := config.ParseClientConfig1(string(ini))
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
