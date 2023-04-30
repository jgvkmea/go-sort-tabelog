/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/jgvkmea/go-sort-tabelog/interface/gateway"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "指定したエリアとキーワードで食べログを検索する",
	Long:  `指定したエリアとキーワードで食べログを検索し、上位100件からレビューが高い順に5つ返す`,
	Run: func(cmd *cobra.Command, args []string) {
		l := logrus.New()
		area, err := cmd.Flags().GetString("area")
		if err != nil {
			l.Errorf("failed to get area: %s", err)
			return
		}
		keyword, err := cmd.Flags().GetString("keyword")
		if err != nil {
			l.Errorf("failed to get keyword: %s\n", err)
			return
		}
		l.Infof("area:%s keyword:%s で検索します", area, keyword)

		wd := gateway.NewWebDriver("--no-sandbox")
		shops, err := wd.GetShopList(area, keyword)
		if err != nil {
			l.Errorf("failed to get shop list: %s", err)
			return
		}
		shops.SortByRating()
		count := shops.GetOutputCount()
		shops = shops.GetTopShopList(count)
		if err != nil {
			l.Errorf("failed to get top shop list: %s", err)
			return
		}
		l.Infoln("shops: ", shops)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().String("area", "", "エリア名")
	searchCmd.Flags().String("keyword", "", "キーワード")
}
