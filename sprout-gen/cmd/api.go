package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ink-yht-code/sprout/sprout-gen/generator"
	"github.com/ink-yht-code/sprout/sprout-gen/parser"
	"github.com/spf13/cobra"
)

// NewAPICommand 创建 api 子命令。
//
// 该命令读取 .sprout 文件，生成 HTTP 相关代码（types/handler/routes）。
func NewAPICommand() *cobra.Command {
	var sproutFileFlag string

	cmd := &cobra.Command{
		Use:   "api [service]",
		Short: "从 .sprout 文件生成 HTTP 代码",
		Long: `从 .sprout 文件生成 types_gen.go、handler_gen.go、handlers.go。

示例:
  sprout gen api                  # 在服务目录内运行，解析 .sprout
  sprout gen api --file my.sprout # 指定 .sprout 文件
  sprout gen api myservice        # 指定服务目录名`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var serviceName string
			if len(args) > 0 {
				serviceName = args[0]
			}

			sproutFile := sproutFileFlag
			if sproutFile == "" {
				// 优先查找当前目录的 .sprout 文件
				candidates := []string{".sprout", "service.sprout"}
				if serviceName != "" {
					candidates = append(candidates,
						filepath.Join(serviceName, ".sprout"),
						filepath.Join(serviceName, serviceName+".sprout"))
				}
				for _, c := range candidates {
					if _, err := os.Stat(c); err == nil {
						sproutFile = c
						break
					}
				}
			}

			if sproutFile == "" {
				color.Red("错误: 找不到 .sprout 文件")
				color.Yellow("请使用 --file 指定文件路径，或在服务目录内运行")
				os.Exit(1)
			}

			if _, err := os.Stat(sproutFile); os.IsNotExist(err) {
				color.Red("错误: .sprout 文件不存在: %s", sproutFile)
				os.Exit(1)
			}

			file, err := parser.ParseFile(sproutFile)
			if err != nil {
				color.Red("解析 .sprout 文件失败: %v", err)
				os.Exit(1)
			}

			if err := file.Validate(); err != nil {
				color.Red("验证失败: %v", err)
				os.Exit(1)
			}

			color.Green("解析成功:")
			fmt.Printf("  - 类型: %d 个\n", len(file.Types))
			fmt.Printf("  - 服务: %d 个\n", len(file.Services))
			for _, svc := range file.Services {
				fmt.Printf("    - %s: 公开 %d, 私有 %d\n", svc.Name, len(svc.Public), len(svc.Private))
			}

			// 确定服务目录和模块名
			servicePath := "."
			if serviceName != "" {
				servicePath = serviceName
			}

			moduleName := getModuleNameFromGoMod(servicePath)
			if moduleName == "" {
				color.Red("错误: 找不到 go.mod 或无法获取模块名")
				os.Exit(1)
			}

			gen := generator.New(moduleName, servicePath, file)
			if err := gen.Generate(); err != nil {
				color.Red("生成代码失败: %v", err)
				os.Exit(1)
			}

			color.Green("\n✓ 代码生成成功！")
			color.Yellow("\n生成的文件:")
			fmt.Println("  - internal/types/types_gen.go")
			fmt.Println("  - internal/web/handler_gen.go")
			fmt.Println("  - internal/web/handlers.go")
		},
	}

	cmd.Flags().StringVarP(&sproutFileFlag, "file", "f", "", ".sprout 文件路径")
	return cmd
}
