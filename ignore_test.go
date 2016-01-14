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
	assert.Equal(test, true, object.MatchesPath("nonexistent/ignore"), "should ignore files inside of here")
	assert.Equal(test, true, object.MatchesPath("nonexistent"), "should also just ignore this whole thing")
	assert.Equal(test, true, object.MatchesPath("baz/yo.txt"), "should ignore files inside of here")
	// Paths which are not targeted by the above "lines"
	assert.Equal(test, false, object.MatchesPath("nonexistent/foo/wat"), "should accept unignored files in ignored directories")
	assert.Equal(test, false, object.MatchesPath("othernonexistent/"), "should allow things not in gitignore at all")
	assert.Equal(test, false, object.MatchesPath("nonexistent/foo"), "shouldn't match unignored files")
}

//tests for the ** and * rules in gitignore
func TestWithWildCards(test *testing.T) {
	gitIgnore := []string{"*.swp", "**/foo", "abc/**", "a/**b"}
	object, err := CompileIgnoreLines(gitIgnore...)
	assert.Nil(test, err, "error from CompileIgnoreLines should be nil")

	// Paths which are targeted by the above "lines"
	assert.Equal(test, true, object.MatchesPath("ayy.swp"), "should ignore swp files")
	assert.Equal(test, true, object.MatchesPath("lmao/ayy.swp"), "should ignore swp files in subdirectories")
	assert.Equal(test, true, object.MatchesPath("yo/lmao/ayy.swp"), "should ignore swp files in subdirectories")

	assert.Equal(test, true, object.MatchesPath("/ayy/lmao/alien/foo"), "foo should be ignored, according to the second rule here")
	assert.Equal(test, true, object.MatchesPath("/foo/boo"), "things called foo should be ignored everywhere, according to the second rule here")
	assert.Equal(test, true, object.MatchesPath("koo/yo/foo/boo"), "things called foo should be ignored everywhere, according to the second rule here")

	assert.Equal(test, true, object.MatchesPath("/abc/secret.txt"), "things inside abc should be ignored")
	assert.Equal(test, true, object.MatchesPath("/a/ab/b/hi.txt"), "things inside a/*b/ should be ignored")
}
