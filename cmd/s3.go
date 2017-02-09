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
	"fmt"

	"github.com/go-xweb/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tacogips/awstool/cmd/awstool"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "deal with s3",
	Long:  `deal with s3`,
}

var s3DLCmd = &cobra.Command{
	Use:   "dl",
	Short: "dl s3 files",
	Long:  `dl s3 files`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var s3ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list s3 files",
	Long:  `list s3 files`,
	Run: func(cmd *cobra.Command, args []string) {

		prefixFlag := cmd.Flag("pre")
		flagRegionFlag := cmd.Flag("r")
		//		outputFileFlag := cmd.Flag("o")
		bucketCmd := cmd.Flag("b")

		err := viper.ReadInConfig()
		if err != nil {
			log.Error(err)
			return
		}

		prefix := prefixFlag.Value.String()
		bucketName := bucketCmd.Value.String()
		if len(bucketName) == 0 {
			bucketName = viper.GetString("s3.bucket")
		}

		if len(bucketName) == 0 {
			log.Errorf("no bucket name")
			return
		}

		region := flagRegionFlag.Value.String()
		if len(region) == 0 {
			region = viper.GetString("region")
		}

		s3files, err := awstool.S3List(region, bucketName, prefix)
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Printf("%#v", s3files)
	},
}

func init() {

	// -- List cmd ---
	s3ListCmd.Flags().String("pre", "prefix(w/o bucket)", "output dir")

	// -- DL cmd ---
	s3DLCmd.Flags().String("dir", "s3file", "output dir")
	s3DLCmd.Flags().String("pre", "prefix(w/o bucket)", "output dir")

	s3Cmd.PersistentFlags().String("r", "", "region")
	s3Cmd.PersistentFlags().String("b", "", "bucket")
	s3Cmd.AddCommand(s3DLCmd)

	//	ebCmd.PersistentFlags().String("r", "", "region")
	//	ebCmd.AddCommand(ebListCmd)
	//	RootCmd.AddCommand(ebCmd)

	RootCmd.AddCommand(s3Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
