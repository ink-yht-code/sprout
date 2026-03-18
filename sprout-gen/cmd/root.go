package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand 创建 sprout-gen 的根命令。
//
// use 为命令名（例如 sprout-gen 或 sprout gen）。
func NewRootCommand(use string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   use,
		Short: "Sprout 代码生成器",
		Long:  `Sprout 代码生成器 - 快速创建符合 DDD + Clean Architecture 的微服务`,
	}

	rootCmd.AddCommand(NewNewCommand())
	rootCmd.AddCommand(NewAPICommand())
	rootCmd.AddCommand(NewRPCCommand())
	rootCmd.AddCommand(NewRepoCommand())
	rootCmd.AddCommand(NewLintCommand())
	rootCmd.AddCommand(NewEntityCommand())
	rootCmd.AddCommand(NewModuleCommand())
	rootCmd.AddCommand(NewGenerateCommand())

	// 安装彩色帮助主题
	InstallHelpTheme(rootCmd)

	return rootCmd
}

// RootCommand 返回默认的根命令（命令名为 sprout-gen）。
func RootCommand() *cobra.Command {
	return NewRootCommand("sprout-gen")
}

// Execute 执行根命令。
func Execute() error {
	return RootCommand().Execute()
}
