package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// EntityField 表示实体字段定义（name:type）。
type EntityField struct {
	Name string
	Type string
}

// NewEntityCommand 创建 entity 子命令。
//
// 用于生成 entity/port/dao/repository 四层代码。
func NewEntityCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "entity <name>",
		Short: "生成实体四层代码",
		Long: `生成实体的 entity、port、dao、repository 四层代码。

示例:
  sprout gen entity user
  sprout gen entity order --fields "name:string,price:int64"
  sprout gen entity product --dao gorm`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			fields, _ := cmd.Flags().GetString("fields")
			daoType, _ := cmd.Flags().GetString("dao")
			return generateEntity(name, fields, daoType)
		},
	}

	cmd.Flags().String("fields", "", "字段定义: name:string,age:int")
	cmd.Flags().String("dao", "memory", "DAO 类型: memory, gorm")
	return cmd
}

func generateEntity(name string, fieldsStr string, daoType string) error {
	fields := parseFields(fieldsStr)

	moduleName := getModuleNameFromGoMod(".")
	if moduleName == "" {
		moduleName = name
	}

	inServiceDir := false
	if _, err := os.Stat(filepath.Join("internal", "domain", "entity")); err == nil {
		inServiceDir = true
	}

	if err := genEntityFile(name, fields, inServiceDir); err != nil {
		return err
	}
	if err := genPortFile(name, moduleName, inServiceDir); err != nil {
		return err
	}
	if err := genDAOFile(name, moduleName, inServiceDir, daoType); err != nil {
		return err
	}
	if err := genRepositoryFile(name, moduleName, inServiceDir); err != nil {
		return err
	}

	return nil
}

func genEntityFile(name string, fields []EntityField, inServiceDir bool) error {
	nameTitle := title(name)
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("// %s %s 实体\n", nameTitle, name))
	buf.WriteString(fmt.Sprintf("type %s struct {\n", nameTitle))
	buf.WriteString("\tID        int64 `json:\"id\"`\n")
	for _, f := range fields {
		goType := goTypeFromName(f.Type)
		buf.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", title(f.Name), goType, f.Name))
	}
	buf.WriteString("\tCreatedAt int64 `json:\"created_at\"`\n")
	buf.WriteString("\tUpdatedAt int64 `json:\"updated_at\"`\n")
	buf.WriteString("}\n")

	entityPath := filepath.Join("internal", "domain", "entity", strings.ToLower(name)+".go")
	if !inServiceDir {
		entityPath = filepath.Join(name, "internal", "domain", "entity", strings.ToLower(name)+".go")
	}
	return generateFile(entityPath, "package entity\n\n"+buf.String())
}

func genPortFile(name, moduleName string, inServiceDir bool) error {
	nameTitle := title(name)
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("// %sRepository %s 仓储接口\n", nameTitle, name))
	buf.WriteString("type " + nameTitle + "Repository interface {\n")
	buf.WriteString(fmt.Sprintf("\tCreate(ctx context.Context, entity *entity.%s) error\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tFindByID(ctx context.Context, id int64) (*entity.%s, error)\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tUpdate(ctx context.Context, entity *entity.%s) error\n", nameTitle))
	buf.WriteString("\tDelete(ctx context.Context, id int64) error\n")
	buf.WriteString("}\n")

	portPath := filepath.Join("internal", "domain", "port", strings.ToLower(name)+"_repository.go")
	if !inServiceDir {
		portPath = filepath.Join(name, "internal", "domain", "port", strings.ToLower(name)+"_repository.go")
	}

	imports := fmt.Sprintf("import (\n\t\"context\"\n\t\"%s/internal/domain/entity\"\n)\n\n", moduleName)
	return generateFile(portPath, "package port\n\n"+imports+buf.String())
}

func genDAOFile(name, moduleName string, inServiceDir bool, daoType string) error {
	nameTitle := title(name)

	daoPath := filepath.Join("internal", "repository", "dao", strings.ToLower(name)+".go")
	if !inServiceDir {
		daoPath = filepath.Join(name, "internal", "repository", "dao", strings.ToLower(name)+".go")
	}

	switch daoType {
	case "gorm":
		return genGormDAO(daoPath, moduleName, name, nameTitle)
	default:
		return genMemoryDAO(daoPath, moduleName, name, nameTitle)
	}
}

func genMemoryDAO(daoPath, moduleName, name, nameTitle string) error {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("// %sDAO %s DAO 接口\n", nameTitle, name))
	buf.WriteString("type " + nameTitle + "DAO interface {\n")
	buf.WriteString(fmt.Sprintf("\tCreate(ctx context.Context, entity *entity.%s) error\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tFindByID(ctx context.Context, id int64) (*entity.%s, error)\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tUpdate(ctx context.Context, entity *entity.%s) error\n", nameTitle))
	buf.WriteString("\tDelete(ctx context.Context, id int64) error\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// %sDAOImpl 内存实现\n", nameTitle))
	buf.WriteString(fmt.Sprintf("type %sDAOImpl struct {\n", nameTitle))
	buf.WriteString("\tmu     sync.RWMutex\n")
	buf.WriteString(fmt.Sprintf("\tdata   map[int64]*entity.%s\n", nameTitle))
	buf.WriteString("\tnextID int64\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// New%sDAO 创建 DAO\n", nameTitle))
	buf.WriteString(fmt.Sprintf("func New%sDAO() *%sDAOImpl {\n", nameTitle, nameTitle))
	buf.WriteString("\treturn &" + nameTitle + "DAOImpl{\n")
	buf.WriteString("\t\tdata:   make(map[int64]*entity." + nameTitle + "),\n")
	buf.WriteString("\t\tnextID: 1,\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "DAOImpl) Create(ctx context.Context, entity *entity." + nameTitle + ") error {\n")
	buf.WriteString("\td.mu.Lock()\n\tdefer d.mu.Unlock()\n")
	buf.WriteString("\tentity.ID = d.nextID\n\td.nextID++\n")
	buf.WriteString("\td.data[entity.ID] = entity\n\treturn nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "DAOImpl) FindByID(ctx context.Context, id int64) (*entity." + nameTitle + ", error) {\n")
	buf.WriteString("\td.mu.RLock()\n\tdefer d.mu.RUnlock()\n")
	buf.WriteString("\treturn d.data[id], nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "DAOImpl) Update(ctx context.Context, entity *entity." + nameTitle + ") error {\n")
	buf.WriteString("\td.mu.Lock()\n\tdefer d.mu.Unlock()\n")
	buf.WriteString("\td.data[entity.ID] = entity\n\treturn nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "DAOImpl) Delete(ctx context.Context, id int64) error {\n")
	buf.WriteString("\td.mu.Lock()\n\tdefer d.mu.Unlock()\n")
	buf.WriteString("\tdelete(d.data, id)\n\treturn nil\n")
	buf.WriteString("}\n")

	imports := fmt.Sprintf("import (\n\t\"context\"\n\t\"sync\"\n\t\"%s/internal/domain/entity\"\n)\n\n", moduleName)
	return generateFile(daoPath, "package dao\n\n"+imports+buf.String())
}

func genGormDAO(daoPath, moduleName, name, nameTitle string) error {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("// %sDAO %s DAO 接口\n", nameTitle, name))
	buf.WriteString("type " + nameTitle + "DAO interface {\n")
	buf.WriteString(fmt.Sprintf("\tCreate(ctx context.Context, entity *entity.%s) error\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tFindByID(ctx context.Context, id int64) (*entity.%s, error)\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tUpdate(ctx context.Context, entity *entity.%s) error\n", nameTitle))
	buf.WriteString("\tDelete(ctx context.Context, id int64) error\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// %sGormDAO GORM 实现\n", nameTitle))
	buf.WriteString(fmt.Sprintf("type %sGormDAO struct {\n", nameTitle))
	buf.WriteString("\tdb *gorm.DB\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// New%sGormDAO 创建 DAO\n", nameTitle))
	buf.WriteString(fmt.Sprintf("func New%sGormDAO(db *gorm.DB) *%sGormDAO {\n", nameTitle, nameTitle))
	buf.WriteString(fmt.Sprintf("\treturn &%sGormDAO{db: db}\n", nameTitle))
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "GormDAO) Create(ctx context.Context, entity *entity." + nameTitle + ") error {\n")
	buf.WriteString("\treturn d.db.Create(entity).Error\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "GormDAO) FindByID(ctx context.Context, id int64) (*entity." + nameTitle + ", error) {\n")
	buf.WriteString("\tvar entity entity." + nameTitle + "\n")
	buf.WriteString("\terr := d.db.First(&entity, id).Error\n")
	buf.WriteString("\treturn &entity, err\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "GormDAO) Update(ctx context.Context, entity *entity." + nameTitle + ") error {\n")
	buf.WriteString("\treturn d.db.Save(entity).Error\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (d *" + nameTitle + "GormDAO) Delete(ctx context.Context, id int64) error {\n")
	buf.WriteString("\treturn d.db.Delete(&entity." + nameTitle + "{}, id).Error\n")
	buf.WriteString("}\n")

	imports := fmt.Sprintf("import (\n\t\"context\"\n\t\"gorm.io/gorm\"\n\t\"%s/internal/domain/entity\"\n)\n\n", moduleName)
	return generateFile(daoPath, "package dao\n\n"+imports+buf.String())
}

func genRepositoryFile(name, moduleName string, inServiceDir bool) error {
	nameTitle := title(name)
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("// %sRepository %s 仓储实现\n", nameTitle, name))
	buf.WriteString(fmt.Sprintf("type %sRepository struct {\n", nameTitle))
	buf.WriteString(fmt.Sprintf("\tdao dao.%sDAO\n", nameTitle))
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("// New%sRepository 创建仓储\n", nameTitle))
	buf.WriteString(fmt.Sprintf("func New%sRepository(dao dao.%sDAO) port.%sRepository {\n", nameTitle, nameTitle, nameTitle))
	buf.WriteString(fmt.Sprintf("\treturn &%sRepository{dao: dao}\n", nameTitle))
	buf.WriteString("}\n\n")

	buf.WriteString("func (r *" + nameTitle + "Repository) Create(ctx context.Context, entity *entity." + nameTitle + ") error {\n")
	buf.WriteString("\treturn r.dao.Create(ctx, entity)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (r *" + nameTitle + "Repository) FindByID(ctx context.Context, id int64) (*entity." + nameTitle + ", error) {\n")
	buf.WriteString("\treturn r.dao.FindByID(ctx, id)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (r *" + nameTitle + "Repository) Update(ctx context.Context, entity *entity." + nameTitle + ") error {\n")
	buf.WriteString("\treturn r.dao.Update(ctx, entity)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func (r *" + nameTitle + "Repository) Delete(ctx context.Context, id int64) error {\n")
	buf.WriteString("\treturn r.dao.Delete(ctx, id)\n")
	buf.WriteString("}\n")

	repoPath := filepath.Join("internal", "repository", strings.ToLower(name)+".go")
	if !inServiceDir {
		repoPath = filepath.Join(name, "internal", "repository", strings.ToLower(name)+".go")
	}

	imports := fmt.Sprintf("import (\n\t\"context\"\n\t\"%s/internal/domain/entity\"\n\t\"%s/internal/domain/port\"\n\t\"%s/internal/repository/dao\"\n)\n\n", moduleName, moduleName, moduleName)
	return generateFile(repoPath, "package repository\n\n"+imports+buf.String())
}
