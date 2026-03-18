package scaffoldtmpl

// ServiceTmpl service 模板
var ServiceTmpl = `//go:generate mockgen -source=./{{.Name}}.go -package=mocks -destination=../../mocks/{{.Name}}_mock.go {{.NameUpper}}Service
package service

import (
	"context"

	"github.com/ink-yht-code/sprout/jwt"

	"{{.Name}}/internal/domain"
	"{{.Name}}/internal/repository"
)

// {{.NameUpper}}Service {{.Name}} 服务接口
type {{.NameUpper}}Service interface {
	// Create 创建{{.Name}}，返回 ID
	Create(ctx context.Context, u *domain.{{.NameUpper}}) (string, error)
	// FindByID 根据 ID 查找
	FindByID(ctx context.Context, id string) (*domain.{{.NameUpper}}, error)
	// Update 更新{{.Name}}
	Update(ctx context.Context, u *domain.{{.NameUpper}}) error
	// Delete 删除{{.Name}}
	Delete(ctx context.Context, id string) error
}

// {{.NameUpper}}Service {{.Name}} 服务实现
type {{.Name}}Service struct {
	repo repository.{{.NameUpper}}Repository
	jwt  jwt.Manager
}

// New{{.NameUpper}}Service 创建服务
func New{{.NameUpper}}Service(repo repository.{{.NameUpper}}Repository, jwt jwt.Manager) {{.NameUpper}}Service {
	return &{{.Name}}Service{repo: repo, jwt: jwt}
}

// Create 创建{{.Name}}
func (s *{{.Name}}Service) Create(ctx context.Context, u *domain.{{.NameUpper}}) (string, error) {
	// TODO: 添加业务逻辑校验
	// TODO: 发布创建事件
	return s.repo.Create(ctx, u)
}

// FindByID 根据 ID 查找
func (s *{{.Name}}Service) FindByID(ctx context.Context, id string) (*domain.{{.NameUpper}}, error) {
	return s.repo.FindByID(ctx, id)
}

// Update 更新{{.Name}}
func (s *{{.Name}}Service) Update(ctx context.Context, u *domain.{{.NameUpper}}) error {
	// TODO: 添加业务逻辑校验
	// TODO: 发布更新事件
	return s.repo.Update(ctx, u)
}

// Delete 删除{{.Name}}
func (s *{{.Name}}Service) Delete(ctx context.Context, id string) error {
	// TODO: 添加业务逻辑校验
	// TODO: 发布删除事件
	return s.repo.Delete(ctx, id)
}

// TODO: 添加其他业务方法
// 示例：
// func (s *{{.NameUpper}}Service) FindByUserID(ctx context.Context, userID int64) ([]*domain.{{.NameUpper}}, error)
// func (s *{{.NameUpper}}Service) Search(ctx context.Context, keyword string, page, pageSize int) ([]*domain.{{.NameUpper}}, int64, error)
`
