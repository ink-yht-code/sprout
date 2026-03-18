package scaffoldtmpl

// ServiceTestTmpl service 测试模板
var ServiceTestTmpl = `package service

import (
	"context"
	"testing"

	"github.com/ink-yht-code/sprout/jwt"

	"{{.Name}}/internal/repository"
	"{{.Name}}/internal/domain"
)

// mock{{.NameUpper}}Repository mock repository
type mock{{.NameUpper}}Repository struct {
	data map[string]*domain.{{.NameUpper}}
}

func newMock{{.NameUpper}}Repository() repository.{{.NameUpper}}Repository {
	return &mock{{.NameUpper}}Repository{data: make(map[string]*domain.{{.NameUpper}})}
}

func (m *mock{{.NameUpper}}Repository) Create(ctx context.Context, u *domain.{{.NameUpper}}) (string, error) {
	m.data[u.ID] = u
	return u.ID, nil
}

func (m *mock{{.NameUpper}}Repository) FindByID(ctx context.Context, id string) (*domain.{{.NameUpper}}, error) {
	return m.data[id], nil
}

func (m *mock{{.NameUpper}}Repository) Update(ctx context.Context, u *domain.{{.NameUpper}}) error {
	m.data[u.ID] = u
	return nil
}

func (m *mock{{.NameUpper}}Repository) Delete(ctx context.Context, id string) error {
	delete(m.data, id)
	return nil
}

func Test{{.NameUpper}}Service_Create(t *testing.T) {
	repo := newMock{{.NameUpper}}Repository()
	jwtManager := jwt.NewManager(jwt.Options{
		SignKey:        "test-secret",
		AccessExpire:  3600,
		RefreshExpire: 86400,
		Issuer:        "test",
	})
	svc := New{{.NameUpper}}Service(repo, jwtManager)
	id, err := svc.Create(context.Background(), &domain.{{.NameUpper}}{ID: "id-1", Name: "n"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != "id-1" {
		t.Fatalf("expected id 'id-1', got %s", id)
	}
}
`
