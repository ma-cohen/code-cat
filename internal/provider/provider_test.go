package provider

import "testing"

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

func TestProviderFlags(t *testing.T) {
	gh := Detect("https://github.com/owner/repo.git")
	if gh.BaseBranchFlag != "--base" {
		t.Errorf("github BaseBranchFlag = %q, want --base", gh.BaseBranchFlag)
	}
	if gh.BodyFlag != "--body" {
		t.Errorf("github BodyFlag = %q, want --body", gh.BodyFlag)
	}
	if gh.SourceBranchFlag != "" {
		t.Errorf("github SourceBranchFlag = %q, want empty", gh.SourceBranchFlag)
	}

	gl := Detect("https://gitlab.com/owner/repo.git")
	if gl.BaseBranchFlag != "--target-branch" {
		t.Errorf("gitlab BaseBranchFlag = %q, want --target-branch", gl.BaseBranchFlag)
	}
	if gl.BodyFlag != "--description" {
		t.Errorf("gitlab BodyFlag = %q, want --description", gl.BodyFlag)
	}
	if gl.SourceBranchFlag != "--source-branch" {
		t.Errorf("gitlab SourceBranchFlag = %q, want --source-branch", gl.SourceBranchFlag)
	}
}
