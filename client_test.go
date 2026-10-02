package opendata

import "testing"

func TestNewClientRequiresKey(t *testing.T) {
	if _, err := NewClient(""); err == nil {
		t.Fatal("want error for empty key")
	}
}

func TestNewClientFromEnv(t *testing.T) {
	t.Setenv(EnvServiceKey, "")
	if _, err := NewClientFromEnv(); err == nil {
		t.Fatal("want error when env is empty")
	}
	t.Setenv(EnvServiceKey, "abc")
	c, err := NewClientFromEnv()
	if err != nil || c.key != "abc" || c.baseURL != DefaultBaseURL {
		t.Fatalf("got %+v, %v", c, err)
	}
}
