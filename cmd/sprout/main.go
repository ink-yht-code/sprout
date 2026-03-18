package main

import (
	"os"

	sproutgencmd "github.com/ink-yht-code/sprout/sprout-gen/cmd"
	registrycli "github.com/ink-yht-code/sprout/sprout-registry/cli"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "sprout",
		Short: "Sprout - Go 微服务开发框架工具集",
	}

	rootCmd.AddCommand(sproutgencmd.NewRootCommand("gen"))
	rootCmd.AddCommand(registrycli.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
