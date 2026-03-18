package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// SproutFile 表示解析后的 .sprout 文件结构。
//
// Types 为类型定义列表；Server 为服务器配置；Services 为服务定义列表。
type SproutFile struct {
	Types    []TypeDefinition
	Server   ServerConfig
	Services []ServiceDefinition
}

// TypeDefinition 表示一个类型定义（type X { ... }）。
type TypeDefinition struct {
	Name   string
	Fields []Field
}

// Field 表示一个字段定义。
type Field struct {
	Name string
	Type string
	Tags string
}

// ServerConfig 表示 server 块配置。
type ServerConfig struct {
	Prefix string
}

// ServiceDefinition 表示 service 块定义。
type ServiceDefinition struct {
	Name    string
	Public  []APIEndpoint
	Private []APIEndpoint
}

// APIEndpoint 表示一个 API 端点定义。
type APIEndpoint struct {
	Method   string
	Path     string
	Handler  string // handler 方法名
	Request  string
	Response string
	Private  bool
}

// ParseFile 读取并解析 .sprout 文件。
func ParseFile(filePath string) (*SproutFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	content := string(data)
	return parseContent(content)
}

func parseContent(content string) (*SproutFile, error) {
	file := &SproutFile{
		Types:    []TypeDefinition{},
		Services: []ServiceDefinition{},
	}

	lines := strings.Split(content, "\n")
	var currentType *TypeDefinition
	var currentService *ServiceDefinition
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.HasPrefix(line, "type ") {
			if currentType != nil {
				file.Types = append(file.Types, *currentType)
			}
			currentType = &TypeDefinition{}
			re := regexp.MustCompile(`^type\s+(\w+)\s*\{`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				currentType.Name = matches[1]
			}
			continue
		}

		if currentType != nil {
			if strings.Contains(line, "{") || strings.Contains(line, "}") {
				if strings.Contains(line, "}") {
					file.Types = append(file.Types, *currentType)
					currentType = nil
				}
				continue
			}
			field := parseField(line)
			if field != nil {
				currentType.Fields = append(currentType.Fields, *field)
			}
		}

		if strings.HasPrefix(line, "server {") {
			currentSection = "server"
			continue
		}

		if strings.HasPrefix(line, "service ") {
			if currentService != nil {
				file.Services = append(file.Services, *currentService)
			}
			currentService = &ServiceDefinition{}
			re := regexp.MustCompile(`^service\s+(\w+)\s*\{`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				currentService.Name = matches[1]
			}
			currentSection = "service"
			continue
		}

		if strings.HasPrefix(line, "public {") {
			currentSection = "public"
			continue
		}

		if strings.HasPrefix(line, "private {") {
			currentSection = "private"
			continue
		}

		if strings.Contains(line, "}") {
			if currentSection == "service" && currentService != nil {
				file.Services = append(file.Services, *currentService)
				currentService = nil
			}
			if currentSection == "public" || currentSection == "private" {
				currentSection = "service"
			} else {
				currentSection = ""
			}
			continue
		}

		if currentSection == "server" && strings.HasPrefix(line, "prefix") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				file.Server.Prefix = strings.Trim(parts[1], `"`)
			}
		}

		if (currentSection == "public" || currentSection == "private") && currentService != nil {
			endpoint := parseEndpoint(line)
			if endpoint != nil {
				if currentSection == "public" {
					currentService.Public = append(currentService.Public, *endpoint)
				} else {
					currentService.Private = append(currentService.Private, *endpoint)
				}
			}
		}
	}

	return file, nil
}

func parseField(line string) *Field {
	// 支持: Name string `json:"name" validate:"required"`
	// 支持: Items []Item `json:"items"`
	// 使用 \x60 表示反引号
	re := regexp.MustCompile(`^(\w+)\s+(\S+)(?:\s+\x60([^\x60]*)\x60)?$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil
	}

	return &Field{
		Name: matches[1],
		Type: matches[2],
		Tags: matches[3],
	}
}

// GoType 将 .sprout 类型映射为 Go 类型
func GoType(sproutType string) string {
	// 处理数组类型: []Type
	if strings.HasPrefix(sproutType, "[]") {
		elemType := strings.TrimPrefix(sproutType, "[]")
		return "[]" + GoType(elemType)
	}
	switch sproutType {
	case "int", "int64", "int32", "string", "bool", "float64":
		return sproutType
	default:
		return sproutType // 自定义类型保持原样
	}
}

// Validate 验证 API 定义。
func (f *SproutFile) Validate() error {
	if len(f.Types) == 0 {
		return fmt.Errorf("no types defined")
	}
	if len(f.Services) == 0 {
		return fmt.Errorf("no services defined")
	}
	return nil
}

func parseEndpoint(line string) *APIEndpoint {
	// 格式1: GET "/ping" Ping -> ExampleResp
	// 格式2: POST "/create" Create(CreateReq) -> CreateResp
	re := regexp.MustCompile(`^(\w+)\s+"([^"]+)"\s+(\w+)(?:\((\w*)\))?\s*->\s*(\w+)$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 5 {
		return nil
	}

	return &APIEndpoint{
		Method:   matches[1],
		Path:     matches[2],
		Handler:  matches[3],
		Request:  matches[4], // 可能为空
		Response: matches[5],
	}
}
