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

	"github.com/go-xweb/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tacogips/awstool/cmd/awstool"
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

		flagRegionFlag := cmd.Flag("r")
		outputFileFlag := cmd.Flag("o")

		outputFile := outputFileFlag.Value.String()
		if len(outputFile) == 0 {
			log.Errorf("invalid output path [%s]", outputFile)
			return
		}

		err := viper.ReadInConfig()
		if err != nil {
			log.Error(err)
			return
		}

		region := flagRegionFlag.Value.String()
		if len(region) == 0 {
			region = viper.GetString("region")
		}
		filterAppNames := viper.GetStringSlice("eb.filter_app_names")

		apps, err := awstool.ListEB(region, filterAppNames)
		if err != nil {
			log.Error(err)
			return
		}
		if len(apps) == 0 {
			log.Errorf("no applications found region:%s filter-app-names:%#v", region, filterAppNames)
			return
		}

		tmpl := viper.GetString("eb.template")

		o, err := os.Create(outputFile)
		if err != nil {
			log.Error(err)
			return
		}
		defer o.Close()

		useTemplate := len(tmpl) != 0

		var t *template.Template
		if useTemplate {
			t = template.New("eblisttmpl")

			_, err := t.Parse(tmpl)
			if err != nil {
				log.Error(err)
				return
			}
		}

		for _, app := range apps {
			if useTemplate {
				err := t.Execute(o, app)
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				fmt.Fprintf(o, "%#v", app)
			}
		}

	},
}

func init() {

	ebListCmd.Flags().String("o", "eblist.out", "output")

	ebCmd.PersistentFlags().String("r", "", "region")
	ebCmd.AddCommand(ebListCmd)
	RootCmd.AddCommand(ebCmd)

}
