package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/octogo/surfdap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	ldapHost   string
	ldapPort   uint16
	ldapBaseDN string
)

var rootCmd = &cobra.Command{
	Use:              "surfdap",
	Short:            "A command-line tool for surfing an LDAP directory tree.",
	PersistentPreRun: rootCmdPersistentPreRun,
	Run:              rootCmdRun,
}

func init() {
	viper.SetDefault("garbage-collect", false)
	viper.SetEnvPrefix("surfdap")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path/to/config.yml")
	rootCmd.PersistentFlags().StringVarP(&ldapHost, "host", "H", "", "LDAP server URL.")
	rootCmd.PersistentFlags().Uint16VarP(&ldapPort, "port", "P", 389, "Port of LDAP server.")
	rootCmd.PersistentFlags().StringVarP(&ldapBaseDN, "base", "B", "", "search-base DN")

	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("base", rootCmd.PersistentFlags().Lookup("base"))
}

// Execute serves as main entrypoint for cobra.
func Execute() {
	rootCmd.Execute()
}

func rootCmdPersistentPreRun(ccmd *cobra.Command, args []string) {
	if cfgFile != "" {
		abs, err := filepath.Abs(cfgFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		path := filepath.Dir(abs)
		base := filepath.Base(abs)
		nameParts := strings.Split(base, ".")
		name := strings.Join(nameParts[:len(nameParts)-1], ".")

		viper.SetConfigName(name)
		viper.AddConfigPath(path)
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func rootCmdRun(ccmd *cobra.Command, args []string) {
	root := getRoot()
	fmt.Println(root.DN())
}

func getConfig() surfdap.Config {
	host := viper.GetString("host")
	if host == "" {
		host = envy.Get("SURFDAP_HOST", "localhost")
	}

	port := viper.GetInt("port")
	if port == 0 {
		var err error
		port, err = strconv.Atoi(envy.Get("SURFDAP_PORT", "389"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	base := viper.GetString("base")
	if base == "" {
		base = envy.Get("SURFDAP_BASE", "")
	}

	return surfdap.Config{
		Host:   host,
		Port:   uint16(port),
		BaseDN: base,
	}
}

func getRoot() surfdap.Node {
	root, err := surfdap.New(getConfig())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return root
}
