package provider

import (
	"reflect"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		remoteURL string
		wantName  string
		wantCLI   string
	}{
		{"https://github.com/owner/repo.git", "github", "gh"},
		{"git@github.com:owner/repo.git", "github", "gh"},
		{"https://gitlab.com/owner/repo.git", "gitlab", "glab"},
		{"git@gitlab.com:owner/repo.git", "gitlab", "glab"},
		{"https://gitlab.mycompany.com/owner/repo.git", "gitlab", "glab"},
		{"https://bitbucket.org/owner/repo.git", "github", "gh"}, // unknown → github fallback
		{"", "github", "gh"}, // empty → github fallback
	}

	for _, tt := range tests {
		p := Detect(tt.remoteURL)
		if p.Name != tt.wantName {
			t.Errorf("Detect(%q).Name = %q, want %q", tt.remoteURL, p.Name, tt.wantName)
		}
		if p.CLI != tt.wantCLI {
			t.Errorf("Detect(%q).CLI = %q, want %q", tt.remoteURL, p.CLI, tt.wantCLI)
		}
	}
}

func TestBrowseRepoCmd(t *testing.T) {
	tests := []struct {
		remoteURL string
		want      []string
	}{
		{"https://github.com/owner/repo.git", []string{"browse"}},
		{"git@github.com:owner/repo.git", []string{"browse"}},
		{"https://gitlab.com/owner/repo.git", []string{"repo", "view", "--web"}},
		{"git@gitlab.com:owner/repo.git", []string{"repo", "view", "--web"}},
		{"https://gitlab.mycompany.com/owner/repo.git", []string{"repo", "view", "--web"}},
	}

	for _, tt := range tests {
		p := Detect(tt.remoteURL)
		if len(p.BrowseRepoCmd) != len(tt.want) {
			t.Errorf("Detect(%q).BrowseRepoCmd = %v, want %v", tt.remoteURL, p.BrowseRepoCmd, tt.want)
			continue
		}
		for i := range tt.want {
			if p.BrowseRepoCmd[i] != tt.want[i] {
				t.Errorf("Detect(%q).BrowseRepoCmd[%d] = %q, want %q", tt.remoteURL, i, p.BrowseRepoCmd[i], tt.want[i])
			}
		}
	}
}

func TestViewPRCmd(t *testing.T) {
	tests := []struct {
		remoteURL string
		want      []string
	}{
		{"https://github.com/owner/repo.git", []string{"pr", "view", "--web"}},
		{"git@github.com:owner/repo.git", []string{"pr", "view", "--web"}},
		{"https://gitlab.com/owner/repo.git", []string{"mr", "view", "--web"}},
		{"git@gitlab.com:owner/repo.git", []string{"mr", "view", "--web"}},
		{"https://gitlab.mycompany.com/owner/repo.git", []string{"mr", "view", "--web"}},
	}

	for _, tt := range tests {
		p := Detect(tt.remoteURL)
		if len(p.ViewPRCmd) != len(tt.want) {
			t.Errorf("Detect(%q).ViewPRCmd = %v, want %v", tt.remoteURL, p.ViewPRCmd, tt.want)
			continue
		}
		for i := range tt.want {
			if p.ViewPRCmd[i] != tt.want[i] {
				t.Errorf("Detect(%q).ViewPRCmd[%d] = %q, want %q", tt.remoteURL, i, p.ViewPRCmd[i], tt.want[i])
			}
		}
	}
}

func TestProviderDoesNotExposeCreatePRFields(t *testing.T) {
	providerType := reflect.TypeOf(Provider{})
	for _, field := range []string{"SubCmd", "BaseBranchFlag", "BodyFlag", "SourceBranchFlag"} {
		if _, ok := providerType.FieldByName(field); ok {
			t.Errorf("Provider still exposes removed creation field %s", field)
		}
	}
}
