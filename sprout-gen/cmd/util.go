package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func title(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func parseFields(fieldsStr string) []EntityField {
	if strings.TrimSpace(fieldsStr) == "" {
		return []EntityField{}
	}
	var fields []EntityField
	for _, f := range strings.Split(fieldsStr, ",") {
		parts := strings.Split(strings.TrimSpace(f), ":")
		if len(parts) == 2 {
			fields = append(fields, EntityField{Name: parts[0], Type: parts[1]})
		}
	}
	return fields
}

func goTypeFromName(name string) string {
	switch name {
	case "string":
		return "string"
	case "int":
		return "int"
	case "int64":
		return "int64"
	case "int32":
		return "int32"
	case "float", "float64":
		return "float64"
	case "bool":
		return "bool"
	default:
		return name
	}
}

func getModuleNameFromGoMod(dir string) string {
	goModPath := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func generateFile(path string, content string) error {
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
