package cli

import (
	"testing"

	"github.com/edalferes/monogo/internal/config"
)

// mockModuleRunner implements ModuleRunner for testing.
type mockModuleRunner struct {
	runWithModulesCalled bool
	lastModules          []string
	lastConfig           *config.Config
	returnError          error
}

func (m *mockModuleRunner) RunWithModules(modules []string, cfg *config.Config) error {
	m.runWithModulesCalled = true
	m.lastModules = modules
	m.lastConfig = cfg
	return m.returnError
}

func TestCLI_validateAndCleanModules(t *testing.T) {
	cli := NewCLI()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "valid single module",
			input:    []string{"auth"},
			expected: []string{"auth"},
		},
		{
			name:     "valid multiple modules",
			input:    []string{"auth", "testmodule"},
			expected: []string{"auth", "testmodule"},
		},
		{
			name:     "mixed valid and invalid modules",
			input:    []string{"auth", "invalid", "testmodule"},
			expected: []string{"auth", "testmodule"},
		},
		{
			name:     "all invalid modules",
			input:    []string{"invalid1", "invalid2"},
			expected: []string{"auth"}, // default fallback
		},
		{
			name:     "empty modules",
			input:    []string{},
			expected: []string{"auth"}, // default fallback
		},
		{
			name:     "modules with whitespace",
			input:    []string{" auth ", " TESTMODULE "},
			expected: []string{"auth", "testmodule"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cli.validateAndCleanModules(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d modules, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected module[%d] = %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestCLI_isValidModule(t *testing.T) {
	cli := NewCLI()
	validModules := []string{"auth", "testmodule"}

	tests := []struct {
		module   string
		expected bool
	}{
		{"auth", true},
		{"testmodule", true},
		{"user", false},
		{"admin", false},
		{"invalid", false},
		{"", false},
		{"AUTH", false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.module, func(t *testing.T) {
			result := cli.isValidModule(tt.module, validModules)
			if result != tt.expected {
				t.Errorf("Expected isValidModule(%s) = %t, got %t",
					tt.module, tt.expected, result)
			}
		})
	}
}

func TestCLI_runApplication(t *testing.T) {
	mockRunner := &mockModuleRunner{}
	mockConfig := &config.Config{
		App: config.AppConfig{
			Name: "test-app",
		},
	}

	cli := NewCLIWithRunner(mockRunner, mockConfig)

	// Test with valid modules
	modules := []string{"auth", "testmodule"}
	err := cli.runApplication(modules)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !mockRunner.runWithModulesCalled {
		t.Error("Expected RunWithModules to be called")
	}

	if len(mockRunner.lastModules) != 2 {
		t.Errorf("Expected 2 modules, got %d", len(mockRunner.lastModules))
	}

	if mockRunner.lastModules[0] != "auth" || mockRunner.lastModules[1] != "testmodule" {
		t.Errorf("Expected [auth, testmodule], got %v", mockRunner.lastModules)
	}

	if mockRunner.lastConfig != mockConfig {
		t.Error("Expected config to be passed correctly")
	}
}

func TestNewCLI(t *testing.T) {
	cli := NewCLI()

	if cli == nil {
		t.Fatal("Expected CLI instance, got nil")
	}

	if cli.runner == nil {
		t.Error("Expected runner to be initialized")
	}

	if cli.config == nil {
		t.Error("Expected config to be initialized")
	}
}

func TestNewCLIWithRunner(t *testing.T) {
	mockRunner := &mockModuleRunner{}
	mockConfig := &config.Config{}

	cli := NewCLIWithRunner(mockRunner, mockConfig)

	if cli == nil {
		t.Fatal("Expected CLI instance, got nil")
	}

	if cli.runner != mockRunner {
		t.Error("Expected custom runner to be set")
	}

	if cli.config != mockConfig {
		t.Error("Expected custom config to be set")
	}
}
