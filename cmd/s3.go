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
	"os"
	"text/template"
	"time"

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

		prefixFlag := cmd.Flag("pre")
		flagRegionFlag := cmd.Flag("r")
		outputDirFlag := cmd.Flag("dir")
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

		outputDir := outputDirFlag.Value.String()
		if len(outputDir) == 0 {
			outputDir = time.Now().Format(fmt.Sprintf("./s3_pre_%s_2006_01_02_15_04_05", prefix))
		}
		if _, err := os.Stat(outputDir); err != nil {
			log.Errorf("dir exists %s", outputDir)
			return
		} else {
			err := os.MkdirAll(outputDir, 0774)
			if err != nil {
				log.Errorf("failed to create dir %s :%s", outputDir, err.Error())
				return
			}
		}

		err = awstool.S3DownloadPrefix(region, bucketName, prefix, outputDir)
		if err != nil {
			log.Error(err)
			return
		}

	},
}

var s3ListCmd = &cobra.Command{
	Use:   "list",
	Short: "list s3 files",
	Long:  `list s3 files`,
	Run: func(cmd *cobra.Command, args []string) {

		prefixFlag := cmd.Flag("pre")
		flagRegionFlag := cmd.Flag("r")
		outputFileFlag := cmd.Flag("o")
		bucketCmd := cmd.Flag("b")

		err := viper.ReadInConfig()
		if err != nil {
			log.Error(err)
			return
		}

		outputFile := outputFileFlag.Value.String()
		if len(outputFile) == 0 {
			log.Errorf("invalid output path [%s]", outputFile)
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

		tmpl := viper.GetString("s3.template")

		o, err := os.Create(outputFile)
		if err != nil {
			log.Error(err)
			return
		}
		defer o.Close()

		useTemplate := len(tmpl) != 0

		var t *template.Template
		if useTemplate {
			funcMap := template.FuncMap{
				"ToJstFormatFunc": awstool.ToJstFormatFunc(time.RFC3339),
				"AsKiB":           awstool.AsKiB,
				"AsMiB":           awstool.AsMiB,
				"AsGiB":           awstool.AsGiB,
			}
			t = template.New("s3listtmpl").Funcs(funcMap)

			_, err := t.Parse(tmpl)
			if err != nil {
				log.Error(err)
				return
			}
		}

		if useTemplate {

			err := t.Execute(o, s3files)
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			fmt.Fprintf(o, "%#v", s3files)
		}

	},
}

func init() {

	// -- List cmd ---
	s3ListCmd.Flags().String("o", "s3list.out", "output")

	// -- DL cmd ---
	s3DLCmd.Flags().String("dir", "", "output dir(must be new dir)")

	// -- s3 cmd ---
	s3Cmd.PersistentFlags().String("r", "", "regionCmd")
	s3Cmd.PersistentFlags().String("b", "", "bucket")
	s3Cmd.PersistentFlags().String("pre", "", "prefix(w/o bucket)")

	s3Cmd.AddCommand(s3DLCmd)
	s3Cmd.AddCommand(s3ListCmd)

	//	ebCmd.PersistentFlags().String("r", "", "region")
	//	ebCmd.AddCommand(ebListCmd)
	//	RootCmd.AddCommand(ebCmd)

	RootCmd.AddCommand(s3Cmd)

}
