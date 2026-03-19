package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ink-yht-code/sprout/sprout-gen/scaffold"
	"github.com/ink-yht-code/sprout/sprout-gen/scaffoldtmpl"
	"github.com/spf13/cobra"
)

// NewGenerateCommand 创建 generate 子命令。
//
// 该命令使用 legacy 模板生成服务骨架，主要用于兼容旧项目结构。
func NewGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <name>",
		Short: "生成服务骨架（legacy 模板模式）",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if strings.Contains(name, "-") {
				return fmt.Errorf("invalid service name %q: legacy templates use it as Go identifier; please use letters/numbers/underscore only (e.g. demosvc)", name)
			}
			serviceID, _ := cmd.Flags().GetInt("service-id")
			hasHTTP, _ := cmd.Flags().GetBool("http")
			hasRPC, _ := cmd.Flags().GetBool("rpc")
			tidy, _ := cmd.Flags().GetBool("tidy")
			if serviceID == 0 {
				return fmt.Errorf("service-id is required (non-zero) for generate")
			}
			if err := generateSkeleton(name, serviceID, hasHTTP, hasRPC); err != nil {
				return err
			}
			if tidy {
				c := exec.Command("go", "mod", "tidy")
				c.Dir = name
				out, err := c.CombinedOutput()
				if err != nil {
					return fmt.Errorf("go mod tidy failed in %s: %w\n%s", name, err, string(out))
				}
			}
			fmt.Printf("\n✓ generated service: %s\n\n", name)
			fmt.Printf("next steps:\n")
			fmt.Printf("  cd %s\n", name)
			if !tidy {
				fmt.Printf("  go mod tidy\n")
			}
			fmt.Printf("  go test ./...\n")
			fmt.Printf("  go run ./cmd\n")
			return nil
		},
	}

	cmd.Flags().Int("service-id", 0, "Service ID（必须非 0）")
	cmd.Flags().Bool("http", true, "生成 HTTP 相关文件")
	cmd.Flags().Bool("rpc", false, "生成 gRPC 相关文件")
	cmd.Flags().Bool("tidy", false, "生成后自动执行 go mod tidy")
	return cmd
}

func generateSkeleton(name string, serviceID int, hasHTTP, hasRPC bool) error {
	data := scaffold.ServiceData{
		Name:           name,
		NameUpper:      title(name),
		ServiceID:      serviceID,
		HasHTTP:        hasHTTP,
		HasRPC:         hasRPC,
		SproutVersion:  SproutVersion,
		SproutxVersion: SproutxVersion,
		GintVersion:    "v1.0.0",
		GintxVersion:   "v1.0.0",
	}

	// .sprout file
	if hasHTTP {
		content, err := scaffold.ExecuteTemplate(scaffoldtmpl.SproutTmpl, data)
		if err != nil {
			return err
		}
		if err := scaffold.GenerateFile(filepath.Join(name, name+".sprout"), content); err != nil {
			return err
		}
	}

	// configs
	cfgContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.ConfigYamlTmpl, data)
	if err != nil {
		return err
	}
	cfgContent = rewriteLegacyRefs(cfgContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "configs", name+".yaml"), cfgContent); err != nil {
		return err
	}

	// main
	mainContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.MainTmpl, data)
	if err != nil {
		return err
	}
	mainContent = rewriteLegacyRefs(mainContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "cmd", "main.go"), mainContent); err != nil {
		return err
	}

	// config.go
	configGo := rewriteLegacyRefs(scaffoldtmpl.ConfigGoTmpl)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "config", "config.go"), configGo); err != nil {
		return err
	}

	// errs
	codesContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.CodesTmpl, data)
	if err != nil {
		return err
	}
	codesContent = rewriteLegacyRefs(codesContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "errs", "codes.go"), codesContent); err != nil {
		return err
	}
	errorGo := rewriteLegacyRefs(scaffoldtmpl.ErrorTmpl)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "errs", "error.go"), errorGo); err != nil {
		return err
	}

	// wiring
	wiringContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.WiringTmpl, data)
	if err != nil {
		return err
	}
	wiringContent = rewriteLegacyRefs(wiringContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "wiring", "wire.go"), wiringContent); err != nil {
		return err
	}
	wireGenContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.WireGenTmpl, data)
	if err != nil {
		return err
	}
	wireGenContent = rewriteLegacyRefs(wireGenContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "wire_gen.go"), wireGenContent); err != nil {
		return err
	}

	if hasHTTP {
		// web layer
		handlerContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.HTTPTmpl, data)
		if err != nil {
			return err
		}
		handlerContent = rewriteLegacyRefs(handlerContent)
		if err := scaffold.GenerateFile(filepath.Join(name, "internal", "web", "handler.go"), handlerContent); err != nil {
			return err
		}
		voContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.TypesTmpl, data)
		if err != nil {
			return err
		}
		voContent = rewriteLegacyRefs(voContent)
		if err := scaffold.GenerateFile(filepath.Join(name, "internal", "web", "vo.go"), voContent); err != nil {
			return err
		}
		resultContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.ResultTmpl, data)
		if err != nil {
			return err
		}
		resultContent = rewriteLegacyRefs(resultContent)
		if err := scaffold.GenerateFile(filepath.Join(name, "internal", "web", "result.go"), resultContent); err != nil {
			return err
		}
	}

	// service
	svcContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.ServiceTmpl, data)
	if err != nil {
		return err
	}
	svcContent = rewriteLegacyRefs(svcContent)
	svcFileName := strings.ToLower(name) + ".go"
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "service", svcFileName), svcContent); err != nil {
		return err
	}

	// test
	testContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.ServiceTestTmpl, data)
	if err != nil {
		return err
	}
	testContent = rewriteLegacyRefs(testContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "service", strings.ToLower(name)+"_test.go"), testContent); err != nil {
		return err
	}

	// repository
	domainContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.DomainTmpl, data)
	if err != nil {
		return err
	}
	domainContent = rewriteLegacyRefs(domainContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "domain", strings.ToLower(name)+".go"), domainContent); err != nil {
		return err
	}
	repoPortContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.RepositoryPortTmpl, data)
	if err != nil {
		return err
	}
	repoPortContent = rewriteLegacyRefs(repoPortContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "repository", strings.ToLower(name)+"_repository.go"), repoPortContent); err != nil {
		return err
	}
	daoContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.DAOTmpl, data)
	if err != nil {
		return err
	}
	daoContent = rewriteLegacyRefs(daoContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "repository", "dao", strings.ToLower(name)+".go"), daoContent); err != nil {
		return err
	}
	daoIFContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.DAOInterfaceTmpl, data)
	if err != nil {
		return err
	}
	daoIFContent = rewriteLegacyRefs(daoIFContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "repository", "dao", strings.ToLower(name)+"_dao.go"), daoIFContent); err != nil {
		return err
	}
	cacheContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.CacheTmpl, data)
	if err != nil {
		return err
	}
	cacheContent = rewriteLegacyRefs(cacheContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "repository", "cache", strings.ToLower(name)+".go"), cacheContent); err != nil {
		return err
	}
	repoImplContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.RepositoryImplTmpl, data)
	if err != nil {
		return err
	}
	repoImplContent = rewriteLegacyRefs(repoImplContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "repository", strings.ToLower(name)+".go"), repoImplContent); err != nil {
		return err
	}

	// event
	eventContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.EventTmpl, data)
	if err != nil {
		return err
	}
	eventContent = rewriteLegacyRefs(eventContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "event", strings.ToLower(name)+".go"), eventContent); err != nil {
		return err
	}
	producerContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.EventProducerTmpl, data)
	if err != nil {
		return err
	}
	producerContent = rewriteLegacyRefs(producerContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "event", "producer.go"), producerContent); err != nil {
		return err
	}
	consumerContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.EventConsumerTmpl, data)
	if err != nil {
		return err
	}
	consumerContent = rewriteLegacyRefs(consumerContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "event", "consumer.go"), consumerContent); err != nil {
		return err
	}

	// mocks
	mockSvcContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.MockServiceTmpl, data)
	if err != nil {
		return err
	}
	mockSvcContent = rewriteLegacyRefs(mockSvcContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "mocks", strings.ToLower(name)+"_mock.go"), mockSvcContent); err != nil {
		return err
	}
	mockRepoContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.MockRepositoryTmpl, data)
	if err != nil {
		return err
	}
	mockRepoContent = rewriteLegacyRefs(mockRepoContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "internal", "mocks", strings.ToLower(name)+"_repo_mock.go"), mockRepoContent); err != nil {
		return err
	}

	// go.mod
	goModContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.GoModTmpl, data)
	if err != nil {
		return err
	}
	goModContent = rewriteGoMod(goModContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "go.mod"), goModContent); err != nil {
		return err
	}

	// README
	readmeContent, err := scaffold.ExecuteTemplate(scaffoldtmpl.ReadmeTmpl, data)
	if err != nil {
		return err
	}
	readmeContent = rewriteLegacyRefs(readmeContent)
	if err := scaffold.GenerateFile(filepath.Join(name, "README.md"), readmeContent); err != nil {
		return err
	}

	return nil
}

func rewriteLegacyRefs(s string) string {
	// code imports
	s = strings.ReplaceAll(s, "github.com/ink-yht-code/gintx", "github.com/ink-yht-code/sprout/sproutx")
	s = strings.ReplaceAll(s, "github.com/ink-yht-code/gint", "github.com/ink-yht-code/sprout")
	return s
}

func rewriteGoMod(s string) string {
	// rewrite module deps to single-module sprout
	s = rewriteLegacyRefs(s)
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "github.com/ink-yht-code/sprout/sproutx") {
			// sproutx is a subpackage within sprout module, not a separate module
			continue
		}
		out = append(out, line)
	}
	res := strings.Join(out, "\n")
	// Make generated project build locally when generated inside this repo.
	// The generated go.mod still requires github.com/ink-yht-code/sprout, but we redirect it to local path.
	if !strings.Contains(res, "replace github.com/ink-yht-code/sprout =>") {
		res = res + "\nreplace github.com/ink-yht-code/sprout => ../\n"
	}
	return res
}
