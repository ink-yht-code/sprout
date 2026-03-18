package session

import (
	"context"
	"testing"
	"time"
)

func TestMemoryProvider_Create(t *testing.T) {
	provider := NewMemoryProvider()
	defer provider.Close()

	ctx := context.Background()
	data := map[string]interface{}{"user_id": "123"}
	expire := 1 * time.Hour

	err := provider.Create(ctx, "test-session", data, expire)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	sess, err := provider.Get(ctx, "test-session")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if sess.ID != "test-session" {
		t.Errorf("Expected ID 'test-session', got '%s'", sess.ID)
	}

	if sess.Data["user_id"] != "123" {
		t.Errorf("Expected user_id '123', got '%v'", sess.Data["user_id"])
	}
}

func TestMemoryProvider_Get_NotFound(t *testing.T) {
	provider := NewMemoryProvider()
	defer provider.Close()

	ctx := context.Background()
	_, err := provider.Get(ctx, "non-existent")
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound, got %v", err)
	}
}

func TestMemoryProvider_Update(t *testing.T) {
	provider := NewMemoryProvider()
	defer provider.Close()

	ctx := context.Background()
	data := map[string]interface{}{"user_id": "123"}
	provider.Create(ctx, "test-session", data, 1*time.Hour)

	newData := map[string]interface{}{"user_id": "456", "role": "admin"}
	err := provider.Update(ctx, "test-session", newData)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	sess, _ := provider.Get(ctx, "test-session")
	if sess.Data["user_id"] != "456" {
		t.Errorf("Expected user_id '456', got '%v'", sess.Data["user_id"])
	}

	if sess.Data["role"] != "admin" {
		t.Errorf("Expected role 'admin', got '%v'", sess.Data["role"])
	}
}

func TestMemoryProvider_Delete(t *testing.T) {
	provider := NewMemoryProvider()
	defer provider.Close()

	ctx := context.Background()
	data := map[string]interface{}{"user_id": "123"}
	provider.Create(ctx, "test-session", data, 1*time.Hour)

	err := provider.Delete(ctx, "test-session")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = provider.Get(ctx, "test-session")
	if err != ErrSessionNotFound {
		t.Errorf("Expected ErrSessionNotFound after delete, got %v", err)
	}
}

func TestMemoryProvider_Refresh(t *testing.T) {
	provider := NewMemoryProvider()
	defer provider.Close()

	ctx := context.Background()
	data := map[string]interface{}{"user_id": "123"}
	provider.Create(ctx, "test-session", data, 1*time.Hour)

	oldExpiresAt := time.Now().Add(1 * time.Hour)
	sess, _ := provider.Get(ctx, "test-session")
	if sess.ExpiresAt.Before(oldExpiresAt.Add(-time.Minute)) || sess.ExpiresAt.After(oldExpiresAt.Add(time.Minute)) {
		t.Errorf("ExpiresAt not as expected")
	}

	err := provider.Refresh(ctx, "test-session", 2*time.Hour)
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	sess, _ = provider.Get(ctx, "test-session")
	newExpiresAt := time.Now().Add(2 * time.Hour)
	if sess.ExpiresAt.Before(newExpiresAt.Add(-time.Minute)) || sess.ExpiresAt.After(newExpiresAt.Add(time.Minute)) {
		t.Errorf("ExpiresAt not refreshed correctly")
	}
}

func TestMemoryProvider_Expiration(t *testing.T) {
	provider := NewMemoryProvider()
	defer provider.Close()

	ctx := context.Background()
	data := map[string]interface{}{"user_id": "123"}
	provider.Create(ctx, "test-session", data, 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	_, err := provider.Get(ctx, "test-session")
	if err != ErrSessionExpired {
		t.Errorf("Expected ErrSessionExpired, got %v", err)
	}
}

func TestSession_GetSet(t *testing.T) {
	sess := &Session{
		ID:        "test",
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	sess.Set("key1", "value1")
	sess.Set("key2", 123)

	val, ok := sess.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Expected 'value1', got %v", val)
	}

	val, ok = sess.Get("key2")
	if !ok || val != 123 {
		t.Errorf("Expected 123, got %v", val)
	}

	_, ok = sess.Get("key3")
	if ok {
		t.Error("Expected false for non-existent key")
	}
}

func TestSession_Del(t *testing.T) {
	sess := &Session{
		ID:        "test",
		Data:      map[string]interface{}{"key1": "value1", "key2": "value2"},
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	sess.Del("key1")

	_, ok := sess.Get("key1")
	if ok {
		t.Error("Expected key1 to be deleted")
	}

	val, ok := sess.Get("key2")
	if !ok || val != "value2" {
		t.Errorf("Expected 'value2', got %v", val)
	}
}
