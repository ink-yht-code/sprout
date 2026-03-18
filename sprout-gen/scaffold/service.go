package scaffold

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

// GenerateFile 在目标文件不存在时创建文件。
//
// 若文件已存在则直接返回 nil，不会覆盖用户内容。
func GenerateFile(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// GenerateFileForce 强制写入文件（会覆盖已有文件）。
func GenerateFileForce(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// ExecuteTemplate 执行文本模板并返回渲染后的内容。
func ExecuteTemplate(tmplStr string, data any) (string, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ServiceData 是 scaffold 模板渲染所需的数据。
type ServiceData struct {
	Name           string
	NameUpper      string
	ServiceID      int
	HasHTTP        bool
	HasRPC         bool
	DAOType        string
	CacheType      string
	SproutVersion  string
	SproutxVersion string

	GintVersion  string
	GintxVersion string
}
