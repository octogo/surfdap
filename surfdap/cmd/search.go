package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/octogo/surfdap"
	"github.com/spf13/cobra"
)

var (
	searchScope  = ""
	searchFilter = ""
	searchAttrs  = ""
)

var scopeMap = map[string]surfdap.SearchScope{
	"base": surfdap.ScopeBase,
	"one":  surfdap.ScopeOne,
	"sub":  surfdap.ScopeSub,
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the LDAP directory tree",
	Run:   searchCmdRun,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().StringVarP(&searchScope, "scope", "s", "base", "LDAP search-scope")
	searchCmd.PersistentFlags().StringVarP(&searchFilter, "filter", "f", "(objectClass=*)", "LDAP search-fiter")
	searchCmd.PersistentFlags().StringVarP(&searchAttrs, "attrs", "a", "*", "filter only these attrs")
}

func searchCmdRun(ccmd *cobra.Command, args []string) {
	scope, found := scopeMap[searchScope]
	if !found {
		fmt.Println("Scope must be one of `base`, `one` or `sub`")
		os.Exit(1)
	}

	root := getRoot()

	nodes, err := root.Search(scope, searchFilter, strings.Split(searchAttrs, ","))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, node := range nodes {
		fmt.Println(node)
	}
}
