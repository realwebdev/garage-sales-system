package role_test

import (
	"encoding/json"
	"testing"

	"github.com/realwebdev/garage-sales-system/business/types/role"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    role.Role
		wantErr bool
	}{
		{"Valid Admin", "Admin", role.Admin, false},
		{"Valid User", "User", role.User, false},
		{"Invalid Casing", "ADMIN", role.Role{}, true},
		{"Invalid Role", "Ghost", role.Role{}, true},
		{"Empty String", "", role.Role{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := role.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	t.Run("Marshal", func(t *testing.T) {
		input := role.Admin
		want := `"Admin"`

		got, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		if string(got) != want {
			t.Errorf("got %s, want %s", string(got), want)
		}
	})

	t.Run("UnmarshalValid", func(t *testing.T) {
		input := `"User"`
		var got role.Role

		if err := json.Unmarshal([]byte(input), &got); err != nil {
			t.Fatalf("json.Unmarshall failed: %v", err)
		}

		if !got.Equal(role.User) {
			t.Errorf("got %v, want %v", got, role.User)
		}
	})

	t.Run("UnmarshalInvalide", func(t *testing.T) {
		input := `"SuperAdmin"`
		var got role.Role

		if err := json.Unmarshal([]byte(input), &got); err == nil {
			t.Error("expected error for invalid role in JSON, got nil")
		}
	})
}

func TestSet(t *testing.T) {
	items := role.Set()
	if len(items) != 2 {
		t.Errorf("expected 2 roles, got %d", len(items))
	}
}
