package regs

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var cases = []string{
	`(?i)War`,
	`POSTAGE DUE`,
}

func TestCases(t *testing.T) {
	s := "WAR TAX STAMPS"
	re := regexp.MustCompile(MakeRegexOr(cases))
	ok := re.MatchString(s)
	assert.Equal(t, ok, true, "Case Test Failed")
}

func TestRules(t *testing.T) {

	rules := NewRules()
	rules.Add(Rule{"^Testing_", "A testing rule"})
	rules.Add(Rule{"(?i)^War", "Matched War"})
	rules.Add(Rule{"POSTAGE", "Matched Postage"})

	elements := rules.Values()

	assert.Equal(t, rules.size, len(elements), "Elements must be 2")

	test := "WAR ON TEXT"

	assert.Equal(t, true, rules.MustCompile().MatchString(test), "")

	assert.Equal(t, true, rules.MatchString(test), "")

	fmt.Println(rules.Verbose(test))
}
