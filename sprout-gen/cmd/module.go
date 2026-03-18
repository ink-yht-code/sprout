package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// NewModuleCommand 创建 module 子命令。
//
// 生成完整模块（基于 entity 命令生成四层代码，并补齐 service 层）。
func NewModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "module <name>",
		Short: "生成完整模块",
		Long: `生成完整的业务模块，包含 entity、port、dao、repository、service。

示例:
  sprout gen module user
  sprout gen module order --fields "name:string,price:int64"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			fields, _ := cmd.Flags().GetString("fields")
			return generateModule(name, fields)
		},
	}
	cmd.Flags().String("fields", "", "字段定义: name:string,age:int")
	return cmd
}

func generateModule(name string, fieldsStr string) error {
	moduleName := getModuleNameFromGoMod(".")
	if moduleName == "" {
		moduleName = name
	}

	inServiceDir := false
	if _, err := os.Stat(filepath.Join("internal", "domain", "entity")); err == nil {
		inServiceDir = true
	}

	if err := generateEntity(name, fieldsStr, "memory"); err != nil {
		return err
	}

	if err := genModuleService(name, moduleName, inServiceDir); err != nil {
		return err
	}

	return nil
}

func genModuleService(name, moduleName string, inServiceDir bool) error {
	nameTitle := title(name)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("// %sService %s 服务\n", nameTitle, name))
	buf.WriteString(fmt.Sprintf("type %sService struct {\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\trepo port.%sRepository\n", nameTitle))
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// New%sService 创建服务\n", nameTitle))
	buf.WriteString(fmt.Sprintf("func New%sService(repo port.%sRepository) *%sService {\n", nameTitle, nameTitle, nameTitle))
	buf.WriteString(fmt.Sprintf("\treturn &%sService{repo: repo}\n", nameTitle))
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// Create 创建%s\n", name))
	buf.WriteString(fmt.Sprintf("func (s *%sService) Create(ctx context.Context, entity *entity.%s) error {\n", nameTitle, nameTitle))
	buf.WriteString("\treturn s.repo.Create(ctx, entity)\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// GetByID 根据ID获取%s\n", name))
	buf.WriteString(fmt.Sprintf("func (s *%sService) GetByID(ctx context.Context, id int64) (*entity.%s, error) {\n", nameTitle, nameTitle))
	buf.WriteString("\treturn s.repo.FindByID(ctx, id)\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// Update 更新%s\n", name))
	buf.WriteString(fmt.Sprintf("func (s *%sService) Update(ctx context.Context, entity *entity.%s) error {\n", nameTitle, nameTitle))
	buf.WriteString("\treturn s.repo.Update(ctx, entity)\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// Delete 删除%s\n", name))
	buf.WriteString(fmt.Sprintf("func (s *%sService) Delete(ctx context.Context, id int64) error {\n", nameTitle))
	buf.WriteString("\treturn s.repo.Delete(ctx, id)\n")
	buf.WriteString("}\n")

	svcPath := filepath.Join("internal", "service", strings.ToLower(name)+".go")
	if !inServiceDir {
		svcPath = filepath.Join(name, "internal", "service", strings.ToLower(name)+".go")
	}

	imports := fmt.Sprintf("import (\n\t\"context\"\n\t\"%s/internal/domain/entity\"\n\t\"%s/internal/domain/port\"\n)\n\n", moduleName, moduleName)
	return generateFile(svcPath, "package service\n\n"+imports+buf.String())
}
