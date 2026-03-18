package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewLintCommand 创建 lint 子命令。
//
// 用于检查生成项目的分层约束（通过简单的 import 规则扫描）。
func NewLintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint <path>",
		Short: "检查分层约束",
		Long: `检查代码是否符合分层约束规则。

规则:
  - web 禁止 import repository/**
  - server 禁止 import gorm.io/*, github.com/redis/*
  - domain 禁止 import 任何基础设施库

示例:
  sprout gen lint user
  sprout gen lint .`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			if err := lintArch(path); err != nil {
				color.Red("错误: %v", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

// Rule 表示一条分层约束规则。
type Rule struct {
	Package     string
	Forbidden   []string
	Description string
}

var rules = []Rule{
	{
		Package:     "web",
		Forbidden:   []string{"repository/", "gorm.io/", "github.com/redis/"},
		Description: "web 禁止直接依赖 repository 或基础设施库",
	},
	{
		Package:     "server",
		Forbidden:   []string{"gorm.io/", "github.com/redis/"},
		Description: "server 禁止直接依赖 gorm/redis",
	},
	{
		Package:     "domain",
		Forbidden:   []string{"gorm.io/", "github.com/redis/", "github.com/gin-gonic/", "database/sql"},
		Description: "domain 禁止依赖任何基础设施库",
	},
	{
		Package:     "domain/entity",
		Forbidden:   []string{"gorm.io/", "github.com/redis/", "github.com/gin-gonic/", "database/sql"},
		Description: "domain/entity 禁止依赖任何基础设施库",
	},
	{
		Package:     "domain/port",
		Forbidden:   []string{"gorm.io/", "github.com/redis/", "github.com/gin-gonic/"},
		Description: "domain/port 禁止依赖具体实现",
	},
}

func lintArch(root string) error {
	errors := 0

	for _, rule := range rules {
		pkgPath := filepath.Join(root, "internal", strings.ReplaceAll(rule.Package, "/", string(filepath.Separator)))

		if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()

				// 检查 import
				if strings.Contains(line, "import") {
					for _, forbidden := range rule.Forbidden {
						if strings.Contains(line, forbidden) {
							relPath, _ := filepath.Rel(root, path)
							color.Red("❌ %s:%d: %s", relPath, lineNum, rule.Description)
							fmt.Printf("   发现违规 import: %s\n\n", forbidden)
							errors++
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if errors > 0 {
		return fmt.Errorf("发现 %d 个分层约束违规", errors)
	}

	color.Green("✓ 分层约束检查通过")
	return nil
}
