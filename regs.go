// Regs is a package that builds rules, using complex regular expressions.
// It builds sets of regular expressions to match sequentially a string
// The main reason for the package, is to be able to build lexers and parsers
// for Natural Language Processing where rules can be numbered in the 100s.
//
// Using the most common techniques creates very unwieldy strings of
// regular expressions.
package regs

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	growthFactor = float32(2.0)  // growth by 100%
	shrinkFactor = float32(0.25) // shrink when size is 25% of capacity (0 means never shrink)

)

// These are options move to options struct
const (
	VERBOSE = iota
	LOG
)

// Rule holds a tuple of a regular expression and a rule name,
// as strings. It is used when logging to log which expression
// matched a string from a set of Rules.
type Rule struct {
	pattern string
	name    string
}

// Rules is a container for rules.
//type Rules []Rule

type Rules struct {
	elements []Rule
	size     int
}

// NewRules instantiates a new empty Rule list.
func NewRules() *Rules {
	return &Rules{}
}

// Add adds a new pattern rule
func (r *Rules) Add(value Rule) {
	r.growBy(1) // grow the array if necessary
	r.elements[r.size] = value
	r.size++ //added only one element
}

// AddSlice adds a slice of rules to Rules.
func (r *Rules) AddSlice(value []Rule) {
	r.growBy(len(value)) // grow the array if necessary
	for _, v := range value {
		r.elements[r.size] = v
		r.size++
	}

}

// Adds a string slice to the elements. Rules are then numbered
// Sequentially.
func (r *Rules) AddVector(value []string) {
	r.growBy(len(value)) // grow the array if necessary
	for i, v := range value {
		r.elements[r.size] = Rule{v, strconv.Itoa(i)}
		r.size++
	}

}

func (r *Rules) resize(cap int) {
	newElements := make([]Rule, cap, cap)
	copy(newElements, r.elements)
	r.elements = newElements
}

// Empty returns true if list does not contain any elements.
func (r *Rules) Empty() bool {
	return r.size == 0
}

// Clear removes all elements from the list.
func (r *Rules) Clear() {
	r.size = 0
	r.elements = []Rule{}
}

// Check that the index is within bounds of the list
func (r *Rules) withinRange(index int) bool {
	return index >= 0 && index < r.size
}

// Expand the array if necessary, i.e. capacity will be reached if we add n elements
func (list *Rules) growBy(n int) {
	// When capacity is reached, grow by a factor of growthFactor and add number of elements
	currentCapacity := cap(list.elements)
	if list.size+n >= currentCapacity {
		newCapacity := int(growthFactor * float32(currentCapacity+n))
		list.resize(newCapacity)
	}
}

// Values returns all elements in the Rule set.
func (r *Rules) Values() []Rule {
	newElements := make([]Rule, r.size, r.size)
	copy(newElements, r.elements[:r.size])
	return newElements
}

// String returns a string representation of the Rules
func (r *Rules) String() string {
	s := "Rule Set\n"
	values := []string{}
	for _, value := range r.elements[:r.size] {
		values = append(values, fmt.Sprintf("%v Rule: %v", value.name, value.pattern))
	}
	s += strings.Join(values, "\n")
	return s
}

// MustCompile compiles all the rules into a singular regular expression and returns
// the expression.
func (r *Rules) MustCompile() *regexp.Regexp {
	var ss []string
	elements := r.Values()
	for _, reg := range elements {
		ss = append(ss, reg.pattern)
	}
	s := strings.Join(ss, "|")
	re := regexp.MustCompile(s)
	return re
}

// MatchString matches the string against all the rules.
func (r *Rules) MatchString(s string) bool {
	re := r.MustCompile()
	if re.MatchString(s) {
		return true
	}

	return false
}

// Match the string against the slice of expressions
// and return the rule that matched
func (r *Rules) Verbose(s string) (string, bool) {
	elements := r.Values()
	for _, reg := range elements {
		re := regexp.MustCompile(reg.pattern)
		if re.MatchString(s) {
			return fmt.Sprintf("Rule %s", reg.name), true
		}

	}
	return "no match", false
}


// MakeRegexOr assembles a pattern string from a slice.
// Start regex symbols "^" or "$" must be included in the patterns
// slice
func MakeRegexOr(ss []string) string {
	var s string
	for _, v := range ss {
		s = s + v + `|`
	}

	//no need for an | a the end
	s = s[0 : len(s)-1]

	return s
}
