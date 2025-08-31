package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"strings"
	"sync"
)

// func init() {
// 	// Relax GC for speed
// 	debug.SetGCPercent(400)
// }

func processFile(path string, counts map[string]int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
	buf := make([]byte, 0, 4*1024*1024) // 4MB buffer
	scanner.Buffer(buf, 4*1024*1024)

	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())
		l := len(code)
		if l >= 8 && l <= 10 {
			mu.Lock()
			counts[code]++
			mu.Unlock()
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func readArtifacts(files []string) {
	fmt.Println("I'm here")

	counts := make(map[string]int, 10_000_000) // preallocate
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process each file simulataneously
	for _, f := range files {
		wg.Add(1)
		go processFile(f, counts, &mu, &wg)
	}
	wg.Wait()

	// Write valid codes (present in at least 2 files)
	outputPath := "../token/valid_codes.txt"
	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for code, cnt := range counts {
		if cnt >= 2 {
			writer.WriteString(code + "\n")
		}
	}
	writer.Flush()
}
