package pantech

import (
	"os"
	"regexp"
	"testing"
)

// TestVersionMatchesChangelog keeps Version (the User-Agent) in step with the
// newest released heading.
func TestVersionMatchesChangelog(t *testing.T) {
	data, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatal(err)
	}
	m := regexp.MustCompile(`(?m)^## v(\d+\.\d+\.\d+) `).FindStringSubmatch(string(data))
	if m == nil {
		t.Fatal("CHANGELOG.md has no ## vX.Y.Z heading")
	}
	if m[1] != Version {
		t.Fatalf("CHANGELOG.md's newest version is %s, Version is %s", m[1], Version)
	}
}
