package utils

import (
	"fmt"
	"strings"
)

// PrintCleanLog formats raw LLM output for readable server
func PrintCleanLog(source, response string) {
	clean := strings.ReplaceAll(response, "\\n", "\n")
	clean = strings.ReplaceAll(clean, "\\\"", "\"")

	// Print with a visible border
	fmt.Println("\n================ LLM RESPONSE START ================")
	fmt.Printf("SOURCE: %s\n", source)
	fmt.Println("----------------------------------------------------")
	fmt.Println(clean)
	fmt.Println("================= LLM RESPONSE END =================\n")
}
