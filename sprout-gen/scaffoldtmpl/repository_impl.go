package scaffoldtmpl

// RepositoryImplTmpl repository 实现模板（带缓存和双模型转换）
var RepositoryImplTmpl = `package repository

import (
	"context"

	"{{.Name}}/internal/domain"
	"{{.Name}}/internal/repository/cache"
	"{{.Name}}/internal/repository/dao"
)

// Cached{{.NameUpper}}Repository {{.Name}} 仓储实现（带缓存）
type Cached{{.NameUpper}}Repository struct {
	dao   dao.{{.NameUpper}}DAO
	cache cache.{{.NameUpper}}Cache
}

// New{{.NameUpper}}Repository 创建仓储
func New{{.NameUpper}}Repository(d dao.{{.NameUpper}}DAO, c cache.{{.NameUpper}}Cache) {{.NameUpper}}Repository {
	return &Cached{{.NameUpper}}Repository{
		dao:   d,
		cache: c,
	}
}

// Create 创建实体
func (r *Cached{{.NameUpper}}Repository) Create(ctx context.Context, u *domain.{{.NameUpper}}) (string, error) {
	// Domain → DAO 转换
	ud := r.domainToEntity(u)
	id, err := r.dao.Insert(ctx, ud)
	if err != nil {
		return "", err
	}
	u.ID = id
	return id, nil
}

// FindByID 根据 ID 查找
func (r *Cached{{.NameUpper}}Repository) FindByID(ctx context.Context, id string) (*domain.{{.NameUpper}}, error) {
	// 先查缓存
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}

	// 缓存未命中，查数据库
	ud, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// DAO → Domain 转换
	u = r.entityToDomain(ud)
	
	// 写入缓存
	_ = r.cache.Set(ctx, u)
	
	return u, nil
}

// Update 更新实体
func (r *Cached{{.NameUpper}}Repository) Update(ctx context.Context, u *domain.{{.NameUpper}}) error {
	// Domain → DAO 转换
	ud := r.domainToEntity(u)
	// 使用 UpdateNonZeroFields 只更新非零字段
	err := r.dao.UpdateNonZeroFields(ctx, ud)
	if err != nil {
		return err
	}
	
	// 删除缓存
	_ = r.cache.Delete(ctx, u.ID)
	
	return nil
}

// Delete 删除实体
func (r *Cached{{.NameUpper}}Repository) Delete(ctx context.Context, id string) error {
	err := r.dao.Delete(ctx, id)
	if err != nil {
		return err
	}
	
	// 删除缓存
	_ = r.cache.Delete(ctx, id)
	
	return nil
}

// domainToEntity Domain 模型 → DAO 模型转换
func (r *Cached{{.NameUpper}}Repository) domainToEntity(u *domain.{{.NameUpper}}) *dao.{{.NameUpper}} {
	return &dao.{{.NameUpper}}{
		ID:   u.ID,
		Name: u.Name,
		// TODO: 添加其他字段转换
		// 注意：使用简单类型（string、int64等），不使用 sql.NullString 等数据库类型
	}
}

// entityToDomain DAO 模型 → Domain 模型转换
func (r *Cached{{.NameUpper}}Repository) entityToDomain(ud *dao.{{.NameUpper}}) *domain.{{.NameUpper}} {
	return &domain.{{.NameUpper}}{
		ID:   ud.ID,
		Name: ud.Name,
		// TODO: 添加其他字段转换
		// 注意：从 sql.NullString 转换为 string
		// 示例：Phone: ud.Phone.String,
	}
}
`
