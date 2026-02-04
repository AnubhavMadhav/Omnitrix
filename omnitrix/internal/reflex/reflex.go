package reflex

import (
	"math/rand"
	"strings"
	"time"

	"github.com/cloudflare/ahocorasick"
)

type IntentType int

const (
	IntentNone     IntentType = iota
	IntentGreeting            // 1
	IntentBlocked             // 2
)

// Engine handles zero-latency pattern matching
type Engine struct {
	matcher        *ahocorasick.Matcher
	greetingOffset int // To know where greetings end and blocked words start
}

// Data Sources (Hardcoded for now)
var (
	// TODO: need to revamp the words in 'blocked'
	greetings = []string{"hello", "hi", "hey", "good morning", "good evening", "sup", "what's up", "greetings"}
	blocked   = []string{"kill", "murder", "hack", "ignore previous instructions", "exploit"}

	// Pre-written responses
	welcomeMsgs = []string{"Hello there!", "Greetings, human.", "System online. How can I help?", "Hi!"}
)

// New initializes the Aho-Corasick matcher
func New() *Engine {
	dictionary := append(greetings, blocked...)

	return &Engine{
		matcher:        ahocorasick.NewStringMatcher(dictionary),
		greetingOffset: len(greetings),
	}
}

// If IntentType is IntentNone, the Gateway should proceed to the LLM.
func (e *Engine) Search(input string) (IntentType, string) {

	normalizedInput := strings.ToLower(input)

	// Aho-Corasick - Match returns a slice of indices into the original dictionary
	matches := e.matcher.Match([]byte(normalizedInput))

	if len(matches) == 0 {
		return IntentNone, ""
	}

	// Safety First: Priority Check: Blocked words take precedence over Greetings.
	// E.g., "Hello, I want to kill" -> Should be Blocked, not Greeting.
	for _, index := range matches {
		// If index is greater than or equal to the greeting count, it's a Blocked word
		if index >= e.greetingOffset {
			return IntentBlocked, "I cannot answer that."
		}
	}

	// If the prompt is long (e.g., > 30 chars), the "Hi" is likely just a starting word in a long prompt with an ask
	// to a real question. We should skip Reflex and let the LLM handle it.
	if len(input) > 9 {
		return IntentNone, ""
	}

	return IntentGreeting, pickRandom(welcomeMsgs)
}

func pickRandom(opts []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return opts[r.Intn(len(opts))]
}
