package scaffoldtmpl

// EventTmpl Event 层模板（事件定义）
var EventTmpl = `package event

// {{.NameUpper}}EventType {{.Name}} 事件类型
type {{.NameUpper}}EventType string

const (
	// {{.NameUpper}}Created {{.Name}}创建事件
	{{.NameUpper}}Created {{.NameUpper}}EventType = "{{.Name}}.created"
	// {{.NameUpper}}Updated {{.Name}}更新事件
	{{.NameUpper}}Updated {{.NameUpper}}EventType = "{{.Name}}.updated"
	// {{.NameUpper}}Deleted {{.Name}}删除事件
	{{.NameUpper}}Deleted {{.NameUpper}}EventType = "{{.Name}}.deleted"
)

// {{.NameUpper}}Event {{.Name}} 事件
type {{.NameUpper}}Event struct {
	Type      {{.NameUpper}}EventType ` + "`" + `json:"type"` + "`" + `
	{{.NameUpper}}ID  string           ` + "`" + `json:"{{.Name}}_id"` + "`" + `
	Timestamp int64            ` + "`" + `json:"timestamp"` + "`" + `
	// TODO: 添加其他事件字段
}
`
