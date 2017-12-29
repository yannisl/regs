package regs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"regexp"
	"strings"
	"fmt"
)

var cases = []string{
	`(?i)War`,
	`POSTAGE DUE`,
} 

func TestCases(t *testing.T) {
	s:="WAR TAX STAMPS"
	re:=regexp.MustCompile(MakeRegexOr(cases))
	ok:=re.MatchString(s)
	assert.Equal(t, ok, true, "Case Test Failed")
}

func TestRules(t *testing.T) {
	var str []string
	rules := NewRules()
	rules.Add(Rule{"(?i)^War","Matched War"})
	rules.Add(Rule{"POSTAGE","Matched Postage"})
	assert.Equal(t, len(cases),rules.size, "Length must be equal")

	elements := rules.Values()

	assert.Equal(t,2,len(elements), "Elements must be 2")

	test := "WAR ON TEXT"

	for _,reg:=range elements{
		str=append(str,reg.pattern)
		s:=strings.Join(str,"|")
		fmt.Println(reg.pattern, s)
		re:=regexp.MustCompile(s)
		assert.Equal(t,true,re.MatchString(test),"")
	}	

	assert.Equal(t,true,rules.MustCompile().MatchString(test),"")

	assert.Equal(t,true,rules.MatchString(test),"")
	
	fmt.Println(rules.Verbose(test))
}	