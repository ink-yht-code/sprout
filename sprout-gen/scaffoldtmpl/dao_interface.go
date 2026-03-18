package scaffoldtmpl

// DAOInterfaceTmpl DAO 接口模板
var DAOInterfaceTmpl = `//go:generate mockgen -source=./{{.Name}}_dao.go -package=mocks -destination=../../mocks/{{.Name}}_dao_mock.go {{.NameUpper}}DAO
package dao

import (
	"context"
	"gorm.io/gorm"
)

// {{.NameUpper}}DAO {{.Name}} DAO 接口
type {{.NameUpper}}DAO interface {
	Insert(ctx context.Context, u *{{.NameUpper}}) (string, error)
	FindByID(ctx context.Context, id string) (*{{.NameUpper}}, error)
	Update(ctx context.Context, u *{{.NameUpper}}) error
	UpdateNonZeroFields(ctx context.Context, u *{{.NameUpper}}) error
	Delete(ctx context.Context, id string) error
}

// gorm{{.NameUpper}}DAO GORM 实现
type gorm{{.NameUpper}}DAO struct {
	db *gorm.DB
}

// New{{.NameUpper}}DAO 创建 DAO
func New{{.NameUpper}}DAO(db *gorm.DB) {{.NameUpper}}DAO {
	return &gorm{{.NameUpper}}DAO{db: db}
}

func (d *gorm{{.NameUpper}}DAO) Insert(ctx context.Context, u *{{.NameUpper}}) (string, error) {
	err := d.db.WithContext(ctx).Create(u).Error
	return u.ID, err
}

func (d *gorm{{.NameUpper}}DAO) FindByID(ctx context.Context, id string) (*{{.NameUpper}}, error) {
	var u {{.NameUpper}}
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return &u, err
}

func (d *gorm{{.NameUpper}}DAO) Update(ctx context.Context, u *{{.NameUpper}}) error {
	return d.db.WithContext(ctx).Save(u).Error
}

// UpdateNonZeroFields 只更新非零字段
func (d *gorm{{.NameUpper}}DAO) UpdateNonZeroFields(ctx context.Context, u *{{.NameUpper}}) error {
	return d.db.WithContext(ctx).Model(u).Updates(u).Error
}

func (d *gorm{{.NameUpper}}DAO) Delete(ctx context.Context, id string) error {
	return d.db.WithContext(ctx).Delete(&{{.NameUpper}}{}, id).Error
}
`
