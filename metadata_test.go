package gbclog

import "testing"

func TestModuleMetadata(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "module path", got: ModulePath, want: "goark.dev/gbc-log"},
		{name: "repository", got: Repository, want: "goark-boot-contrib-log"},
		{name: "starter id", got: StarterID, want: "goark.boot.log"},
		{name: "system bean", got: BeanNameSystem, want: "goark.log.system"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("metadata mismatch: got %q, want %q", tt.got, tt.want)
			}
		})
	}
}
