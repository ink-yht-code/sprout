package parser

import (
	"os"
	"testing"
)

func TestParseContent(t *testing.T) {
	content := `syntax = "v1"

type ExampleReq {
    Name string ` + "`" + `json:"name" validate:"required"` + "`" + `
}

type ExampleResp {
    Message string ` + "`" + `json:"message"` + "`" + `
}

server {
    prefix "/api"
}

service TestService {
    public {
        GET "/ping" Ping -> ExampleResp
        POST "/create" Create -> ExampleResp
    }
    
    private {
        GET "/info" GetInfo -> ExampleResp
    }
}
`

	file, err := parseContent(content)
	if err != nil {
		t.Fatalf("ParseContent failed: %v", err)
	}

	if len(file.Types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(file.Types))
	}

	if file.Types[0].Name != "ExampleReq" {
		t.Errorf("Expected type name 'ExampleReq', got '%s'", file.Types[0].Name)
	}

	if len(file.Types[0].Fields) != 1 {
		t.Errorf("Expected 1 field in ExampleReq, got %d", len(file.Types[0].Fields))
	}

	if file.Types[0].Fields[0].Name != "Name" {
		t.Errorf("Expected field name 'Name', got '%s'", file.Types[0].Fields[0].Name)
	}

	if file.Server.Prefix != "/api" {
		t.Errorf("Expected prefix '/api', got '%s'", file.Server.Prefix)
	}

	if len(file.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(file.Services))
	}

	if file.Services[0].Name != "TestService" {
		t.Errorf("Expected service name 'TestService', got '%s'", file.Services[0].Name)
	}

	if len(file.Services[0].Public) != 2 {
		t.Errorf("Expected 2 public endpoints, got %d", len(file.Services[0].Public))
	}

	if len(file.Services[0].Private) != 1 {
		t.Errorf("Expected 1 private endpoint, got %d", len(file.Services[0].Private))
	}

	if file.Services[0].Public[0].Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", file.Services[0].Public[0].Method)
	}

	if file.Services[0].Public[0].Path != "/ping" {
		t.Errorf("Expected path '/ping', got '%s'", file.Services[0].Public[0].Path)
	}

	if file.Services[0].Public[0].Handler != "Ping" {
		t.Errorf("Expected handler 'Ping', got '%s'", file.Services[0].Public[0].Handler)
	}
}

func TestParseFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.sprout")
	if err != nil {
		t.Fatalf("CreateTemp failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `syntax = "v1"

type TestReq {
    ID int ` + "`" + `json:"id"` + "`" + `
}

service Test {
    public {
        GET "/test" Test -> TestReq
    }
}
`

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}
	tmpFile.Close()

	file, err := ParseFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if len(file.Types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(file.Types))
	}
}

func TestParseField(t *testing.T) {
	tests := []struct {
		line     string
		name     string
		typeName string
		tags     string
	}{
		{
			line:     `Name string ` + "`" + `json:"name" validate:"required"` + "`" + ``,
			name:     "Name",
			typeName: "string",
			tags:     `json:"name" validate:"required"`,
		},
		{
			line:     `ID int`,
			name:     "ID",
			typeName: "int",
			tags:     "",
		},
		{
			line:     `Count int ` + "`" + `json:"count"` + "`" + ``,
			name:     "Count",
			typeName: "int",
			tags:     `json:"count"`,
		},
	}

	for _, tt := range tests {
		field := parseField(tt.line)
		if field == nil {
			t.Errorf("parseField(%q) returned nil", tt.line)
			continue
		}

		if field.Name != tt.name {
			t.Errorf("Expected name '%s', got '%s'", tt.name, field.Name)
		}

		if field.Type != tt.typeName {
			t.Errorf("Expected type '%s', got '%s'", tt.typeName, field.Type)
		}

		if field.Tags != tt.tags {
			t.Errorf("Expected tags '%s', got '%s'", tt.tags, field.Tags)
		}
	}
}

func TestParseEndpoint(t *testing.T) {
	tests := []struct {
		line     string
		method   string
		path     string
		handler  string
		request  string
		response string
	}{
		{
			line:     `GET "/ping" Ping -> ExampleResp`,
			method:   "GET",
			path:     "/ping",
			handler:  "Ping",
			request:  "",
			response: "ExampleResp",
		},
		{
			line:     `POST "/create" Create(CreateReq) -> CreateResp`,
			method:   "POST",
			path:     "/create",
			handler:  "Create",
			request:  "CreateReq",
			response: "CreateResp",
		},
		{
			line:     `PUT "/update" Update -> ExampleResp`,
			method:   "PUT",
			path:     "/update",
			handler:  "Update",
			request:  "",
			response: "ExampleResp",
		},
	}

	for _, tt := range tests {
		endpoint := parseEndpoint(tt.line)
		if endpoint == nil {
			t.Errorf("parseEndpoint(%q) returned nil", tt.line)
			continue
		}

		if endpoint.Method != tt.method {
			t.Errorf("Expected method '%s', got '%s'", tt.method, endpoint.Method)
		}

		if endpoint.Path != tt.path {
			t.Errorf("Expected path '%s', got '%s'", tt.path, endpoint.Path)
		}

		if endpoint.Handler != tt.handler {
			t.Errorf("Expected handler '%s', got '%s'", tt.handler, endpoint.Handler)
		}

		if endpoint.Request != tt.request {
			t.Errorf("Expected request '%s', got '%s'", tt.request, endpoint.Request)
		}

		if endpoint.Response != tt.response {
			t.Errorf("Expected response '%s', got '%s'", tt.response, endpoint.Response)
		}
	}
}
