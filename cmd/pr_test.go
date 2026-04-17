package cmd

import "testing"

func TestBranchToTitle(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"feature prefix", "feature/add-login", "Add login"},
		{"feat prefix", "feat/add-login", "Add login"},
		{"fix prefix", "fix/broken-auth", "Broken auth"},
		{"bugfix prefix", "bugfix/null-pointer", "Null pointer"},
		{"bug prefix", "bug/crash-on-start", "Crash on start"},
		{"hotfix prefix", "hotfix/prod-outage", "Prod outage"},
		{"chore prefix", "chore/update-deps", "Update deps"},
		{"docs prefix", "docs/readme-update", "Readme update"},
		{"refactor prefix", "refactor/split-auth", "Split auth"},
		{"test prefix", "test/add-coverage", "Add coverage"},
		{"ci prefix", "ci/fix-workflow", "Fix workflow"},
		{"no prefix", "my-feature-branch", "My feature branch"},
		{"underscores", "feature/my_cool_thing", "My cool thing"},
		{"mixed separators", "fix/foo-bar_baz", "Foo bar baz"},
		{"plain name", "main", "Main"},
		{"empty string", "", ""},
		{"only prefix stripped to empty", "feature/", ""},
		{"already title case", "feature/Add-Login", "Add Login"},
		{"unicode first char", "feat/ünit-test", "Ünit test"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := branchToTitle(tc.input)
			if got != tc.want {
				t.Errorf("branchToTitle(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
