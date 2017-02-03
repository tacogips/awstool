// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ebCmd represents the eb command
var ebCmd = &cobra.Command{
	Use:   "eb",
	Short: "util for elastic beanstalk",
	Long:  `util for elastic beanstalk`,
}

var ebListCmd = &cobra.Command{
	Use:   "list",
	Short: "util for elastic beanstalk",
	Long:  `util for elastic beanstalk`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Get("eb")
		//awstool.ListEB(region string, filterAppNames []*string)
	},
}

func init() {
	ebCmd.AddCommand(ebListCmd)
	RootCmd.AddCommand(ebCmd)

}
