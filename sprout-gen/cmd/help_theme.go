package cmd

import (
	"bytes"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// InstallHelpTheme 为 Cobra 命令安装彩色帮助主题。
func InstallHelpTheme(root *cobra.Command) {
	root.SetHelpTemplate(colorizeHelpTemplate(root.HelpTemplate()))
	root.SetUsageTemplate(colorizeHelpTemplate(root.UsageTemplate()))

	root.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		tpl := colorizeHelpTemplate(cmd.HelpTemplate())
		_ = cmd.HelpTemplate()
		cmd.SetHelpTemplate(tpl)
		_ = cmd.Help()
	})

	root.SetUsageFunc(func(cmd *cobra.Command) error {
		tpl := colorizeHelpTemplate(cmd.UsageTemplate())
		cmd.SetUsageTemplate(tpl)
		return cmd.Usage()
	})
}

func colorizeHelpTemplate(tpl string) string {
	b := bytes.NewBuffer(nil)
	lines := strings.Split(tpl, "\n")

	head := color.New(color.FgCyan, color.Bold).SprintFunc()
	cmdName := color.New(color.FgGreen, color.Bold).SprintFunc()
	flag := color.New(color.FgYellow).SprintFunc()

	for _, line := range lines {
		trim := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trim, "Usage:"):
			line = strings.Replace(line, "Usage:", head("Usage:"), 1)
		case strings.HasPrefix(trim, "Available Commands:"):
			line = strings.Replace(line, "Available Commands:", head("Available Commands:"), 1)
		case strings.HasPrefix(trim, "Flags:"):
			line = strings.Replace(line, "Flags:", head("Flags:"), 1)
		case strings.HasPrefix(trim, "Global Flags:"):
			line = strings.Replace(line, "Global Flags:", head("Global Flags:"), 1)
		}

		// 高亮命令列表
		if strings.HasPrefix(line, "  ") && len(strings.TrimLeft(line, " ")) > 0 {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				first := fields[0]
				if !strings.HasPrefix(first, "-") && first == strings.ToLower(first) {
					idx := strings.Index(line, first)
					if idx >= 0 {
						line = line[:idx] + cmdName(first) + line[idx+len(first):]
					}
				}
			}
		}

		// 高亮 flags
		line = highlightFlags(line, flag)

		b.WriteString(line)
		b.WriteByte('\n')
	}

	return b.String()
}

func highlightFlags(line string, paint func(a ...any) string) string {
	replacer := func(token string) string {
		if strings.HasPrefix(token, "--") || (strings.HasPrefix(token, "-") && len(token) <= 3) {
			return paint(token)
		}
		return token
	}

	parts := strings.Split(line, " ")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = replacer(p)
	}
	return strings.Join(parts, " ")
}
