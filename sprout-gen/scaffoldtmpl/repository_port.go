package scaffoldtmpl

// RepositoryPortTmpl repository 接口模板
var RepositoryPortTmpl = `//go:generate mockgen -source=./{{.Name}}_repository.go -package=mocks -destination=../../mocks/{{.Name}}_repo_mock.go {{.NameUpper}}Repository
package repository

import (
	"context"

	"{{.Name}}/internal/domain"
)

// {{.NameUpper}}Repository {{.Name}} 仓储接口
type {{.NameUpper}}Repository interface {
	// Create 创建实体，返回 ID
	Create(ctx context.Context, u *domain.{{.NameUpper}}) (string, error)
	// FindByID 根据 ID 查找
	FindByID(ctx context.Context, id string) (*domain.{{.NameUpper}}, error)
	// Update 更新实体
	Update(ctx context.Context, u *domain.{{.NameUpper}}) error
	// Delete 删除实体
	Delete(ctx context.Context, id string) error
}
`
