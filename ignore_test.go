// Implement tests for the `ignore` library
package ignore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Validate "CompileIgnoreLines()"
func TestCompileIgnoreLines(test *testing.T) {
	gitIgnore := []string{"# This is a comment in a .gitignore file!", "/node_modules", "*.swp", "/nonexistent", "!/nonexistent/foo", "/baz", "/foo/*.wat"}
	object, err := CompileIgnoreLines(gitIgnore...)
	assert.Nil(test, err, "error from CompileIgnoreLines should be nil")

	// Paths which are targeted by the above "lines"
	assert.Equal(test, true, object.MatchesPath("node_modules/"), "node_modules should match")
	assert.Equal(test, true, object.MatchesPath("yo.swp"), "should ignore all swp files")
	assert.Equal(test, true, object.MatchesPath("foo/bar.wat"), "should ignore all wat files in foo")

	// Paths which are not targeted by the above "lines"
	assert.Equal(test, false, object.MatchesPath("nonexistent/foo/wat"), "should accept unignored files in ignored directories")
	assert.Equal(test, false, object.MatchesPath("othernonexistent//wat"), "should accept unignored files in ignored directories")
	assert.Equal(test, false, object.MatchesPath("nonexistent/foo"), "shouldn't match unignored files")
}
