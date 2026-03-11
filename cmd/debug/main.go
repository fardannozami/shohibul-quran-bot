package main

import (
	"fmt"
	"regexp"
	"github.com/fardannozami/shohibul-quran-bot/internal/parser"
)

func main() {
	p := parser.NewReportParser()
	
	tests := []string{
		"Alhamdulillah asy-syura",
		"alhamdulillah asy syura",
		"Alhamdulillah Asy-Syura",
	}
	
	for _, msg := range tests {
		// Show what extractSurahRange sees
		workMsg := regexp.MustCompile(`(?i).*`).FindString(msg)
		_ = workMsg
		
		results := p.Parse(msg)
		fmt.Printf("Input: %q => %d result(s)\n", msg, len(results))
		for i, r := range results {
			fmt.Printf("  [%d] SurahName=%q\n", i, r.SurahName)
		}
	}
	
	// Also test FindSurahNumber directly  
	fmt.Println("\n--- FindSurahNumber tests ---")
	tests2 := []string{"asy-syura", "asy syura", "asy", "syura", "asysyura"}
	for _, s := range tests2 {
		num := parser.FindSurahNumber(s)
		fmt.Printf("FindSurahNumber(%q) = %d (%s)\n", s, num, parser.GetSurahName(num))
	}
}
