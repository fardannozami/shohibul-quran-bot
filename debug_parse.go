package main

import (
	"fmt"
	"github.com/fardannozami/shohibul-quran-bot/internal/parser"
)

func main() {
	p := parser.NewReportParser()
	msg := `#Alhamdulillah hari Ramadhan
1. Juz 21 ✅ Juz 22✅
2. Juz 23✅
3. Juz 24✅
4. Juz 25✅
5. Juz 26 ✅
6. Juz 27✅
7. Juz 28✅
8. Juz 29✅
9. Juz 30✅
10. Juz 1 ✅
11. Juz 2 ✅
12. Juz 3✅JUZ 4 ✅
13. Juz 5✅ Juz 6✅
14. Juz 7✅ Juz 8✅
15. Juz 9✅ Juz 10✅
16. Juz 11✅+ Al-Kahf✅ 
17. Juz 12✅
18. Juz 13✅
19. Juz 14✅
20. Juz 15,16,17,18✅
21. Juz 19, 20✅
22. Juz 21✅
23. Juz 22✅
24. Juz 23,24✅
25. Juz 25✅
26. Juz 26-30✅
27. Juz 1✅
28. Juz 2✅
29. Juz 3✅`

	results := p.Parse(msg)
	fmt.Printf("Number of results: %d\n", len(results))
	for i, r := range results {
		fmt.Printf("Result %d: Type=%s, Pages=%d, Surah=%s\n", i+1, r.ReportType, r.Pages, r.SurahName)
	}
}
