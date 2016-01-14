/*
The rules for parsing the input file are the same as the ones listed
in the Git docs here: http://git-scm.com/docs/gitignore
*/

package ignore

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// An IgnoreParser is an interface which exposes two methods:
//   MatchesPath() - Returns true if the path is targeted by the patterns compiled in the GitIgnore structure
// GitIgnore is a struct which contains a slice of regexp.Regexp
// patterns
type GitIgnore struct {
	patterns []*regexp.Regexp // List of regexp patterns which this ignore file applies
	negate   []bool           // List of booleans which determine if the pattern is negated
}

// This function pretty much attempts to mimic the parsing rules
// listed above at the start of this file
func getPatternFromLine(line string) (*regexp.Regexp, bool) {
	// Trim OS-specific carriage returns.
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, false
	}

	// Strip comments [Rule 2]
	if strings.HasPrefix(line, `#`) {
		return nil, false
	}

	// TODO: Handle [Rule 4] which negates the match for patterns leading with "!"
	negatePattern := false
	if string(line[0]) == "!" {
		negatePattern = true
		line = line[1:]
	}

	if line[0] == '/' {
		line = line[1:]
	}

	// If we encounter a foo/*.blah in a folder, prepend the / char
	if regexp.MustCompile(`([^\/+])/.*\*\.`).MatchString(line) {
		line = "/" + line
	}

	// Handle escaping the "." char
	line = regexp.MustCompile(`\.`).ReplaceAllString(line, `\.`)

	magicStar := "#$~"

	// Handle "/**/" usage
	if strings.HasPrefix(line, "/**/") {
		line = line[1:]
	}
	line = regexp.MustCompile(`/\*\*/`).ReplaceAllString(line, `(/|/.+/)`)
	line = regexp.MustCompile(`\*\*/`).ReplaceAllString(line, `(|.`+magicStar+`/)`)
	line = regexp.MustCompile(`/\*\*`).ReplaceAllString(line, `(|/.`+magicStar+`)`)

	// Handle escaping the "*" char
	line = regexp.MustCompile(`\\\*`).ReplaceAllString(line, `\`+magicStar)
	line = regexp.MustCompile(`\*`).ReplaceAllString(line, `([^/]*)`)

	// Handle escaping the "?" char
	line = strings.Replace(line, "?", `\?`, -1)

	line = strings.Replace(line, magicStar, "*", -1)

	// Temporary regex
	var expr = ""
	if strings.HasSuffix(line, "/") {
		expr = line + "(|.*)$"
	} else {
		expr = line + "(|/.*)$"
	}
	if strings.HasPrefix(expr, "/") {
		expr = "^(|/)" + expr[1:]
	} else {
		expr = "^(|.*/)" + expr
	}
	pattern, _ := regexp.Compile(expr)

	return pattern, negatePattern
}

// Accepts a variadic set of strings, and returns a GitIgnore object which
// converts and appends the lines in the input to regexp.Regexp patterns
// held within the GitIgnore objects "patterns" field
func CompileIgnoreLines(lines ...string) (*GitIgnore, error) {
	g := GitIgnore{}
	for _, line := range lines {
		pattern, negatePattern := getPatternFromLine(line)
		if pattern != nil {
			g.patterns = append(g.patterns, pattern)
			g.negate = append(g.negate, negatePattern)
		}
	}
	return &g, nil
}

// Accepts a ignore file as the input, parses the lines out of the file
// and invokes the CompileIgnoreLines method
func CompileIgnoreFile(fpath string) (*GitIgnore, error) {
	buffer, error := ioutil.ReadFile(fpath)
	if error == nil {
		s := strings.Split(string(buffer), "\n")
		return CompileIgnoreLines(s...)
	}
	return nil, error
}

// MatchesPath is a function for the IgnoreParser interface.
// It returns true if the given GitIgnore structure would target a given
// path string "f"
func (g GitIgnore) MatchesPath(f string) bool {
	// Replace OS-specific path separator.
	f = strings.Replace(f, string(os.PathSeparator), "/", -1)
	matchesPath := false
	for idx, pattern := range g.patterns {
		if pattern.MatchString(f) {
			matchesPath = true
			if g.negate[idx] {
				matchesPath = false
			}
		}
	}
	return matchesPath
}
