package scaffoldtmpl

// DomainTmpl domain 模板（领域模型）
var DomainTmpl = `package domain

import "time"

// {{.NameUpper}} {{.Name}} 领域模型
// 使用简单类型，符合业务语义，不包含数据库特定字段
type {{.NameUpper}} struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	// TODO: 添加其他业务字段
	// 注意：使用简单类型（string、int64等），不使用 sql.NullString 等数据库类型
}
`
