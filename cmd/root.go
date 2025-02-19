/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package cmd

import (
	"fmt"

	"github.com/konstructio/kubefirst-api/pkg/configs"
	"github.com/konstructio/kubefirst-api/pkg/progressPrinter"
	"github.com/konstructio/kubefirst/cmd/aws"
	"github.com/konstructio/kubefirst/cmd/civo"
	"github.com/konstructio/kubefirst/cmd/digitalocean"
	"github.com/konstructio/kubefirst/cmd/google"
	"github.com/konstructio/kubefirst/cmd/k3d"
	"github.com/konstructio/kubefirst/cmd/vultr"
	"github.com/konstructio/kubefirst/internal/common"
	"github.com/konstructio/kubefirst/internal/progress"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubefirst",
	Short: "kubefirst management cluster installer base command",
	Long: `kubefirst management cluster installer provisions an
	open source application delivery platform in under an hour.
	checkout the docs at https://kubefirst.konstruct.io/docs/.`,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		// wire viper config for flags for all commands
		return configs.InitializeViperConfig(cmd)
	},
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("To learn more about kubefirst, run:")
		fmt.Println("  kubefirst help")
		progress.Progress.Quit()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// This will allow all child commands to have informUser available for free.
	// Refers: https://github.com/konstructio/runtime/issues/525
	// Before removing next line, please read ticket above.
	common.CheckForVersionUpdate()
	progressPrinter.GetInstance()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error occurred during command execution:", err)
		fmt.Println("If a detailed error message was available, please make the necessary corrections before retrying.")
		fmt.Println("You can re-run the last command to try the operation again.")
		progress.Progress.Quit()
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(
		betaCmd,
		aws.NewCommand(),
		civo.NewCommand(),
		digitalocean.NewCommand(),
		k3d.NewCommand(),
		k3d.LocalCommandAlias(),
		google.NewCommand(),
		vultr.NewCommand(),
		LaunchCommand(),
		LetsEncryptCommand(),
		TerraformCommand(),
	)
}
