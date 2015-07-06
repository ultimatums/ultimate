package config

import "testing"

func TestLoadConfig(t *testing.T) {
	if _, err := LoadConfig("testdata/test.yml"); err != nil {
		t.Errorf("Error parsing %s: %s", "testdata/test.yml", err)
	}
}
