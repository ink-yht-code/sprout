package scaffoldtmpl

// DAOTmpl DAO 模型模板（数据模型）
// 使用数据库特定类型，映射数据库表结构
var DAOTmpl = `package dao

import (
	"time"

	"gorm.io/gorm"
)

// {{.NameUpper}} {{.Name}} 数据模型
// 使用数据库特定类型，映射数据库表结构
type {{.NameUpper}} struct {
	ID        string ` + "`" + `gorm:"primaryKey;type:varchar(64)"` + "`" + `
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt ` + "`" + `gorm:"index"` + "`" + `
	Name      string ` + "`" + `gorm:"type:varchar(256);not null"` + "`" + `
	// TODO: 添加其他字段
	// 注意：使用数据库特定类型（sql.NullString等），包含 gorm 标签
}

// TableName 指定表名
func ({{.NameUpper}}) TableName() string {
	return "{{.Name}}_table"
}
`
