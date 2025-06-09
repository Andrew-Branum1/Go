package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	// add new packages as necessary
)

// add new functions and declerations as necessary
type WordCount struct {
	Word  string
	Count int
}

// detemines intermediate file
func hashKey(key string, R int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % R
}

// allows reruns
func cleanUpFiles() {
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "int") && strings.HasSuffix(f.Name(), ".txt") {
			os.Remove(f.Name())
		}
		if strings.HasPrefix(f.Name(), "red-") && strings.HasSuffix(f.Name(), ".out") {
			os.Remove(f.Name())
		}
	}
}

func WordCount_MR_DMP(inFile string, outFile string, nMap int, nReduce int) {

}

func WordCount_MR_SMP(inFile string, outFile string, nMap int, nReduce int) {

}

func WordCount_MR_S(inFile string, outFile string, nMap int, nReduce int) {

	cleanUpFiles()

	data, _ := os.ReadFile(inFile)

	tokens := strings.Fields(string(data))
	tokenCount := len(tokens)
	chunk := tokenCount / nMap

	for m := 0; m < nMap; m++ {

		start := m * chunk
		end := start + chunk
		//fmt.Printf("Mapper %d: tokens %d → %d\n", m, start, end)

		if m == nMap-1 {

			end = tokenCount

		}

		chunk := tokens[start:end]
		reducers := make(map[int]*os.File)

		for r := 0; r < nReduce; r++ {

			fname := fmt.Sprintf("int%d-%d.txt", m, r)

			f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

			if err != nil {
				fmt.Printf("open file error")
				return
			}

			reducers[r] = f
		}

		for _, word := range chunk {

			rid := hashKey(word, nReduce)

			fmt.Fprintf(reducers[rid], "%s 1\n", word)
		}

		for _, f := range reducers {

			f.Close()

		}
	}

	for r := 0; r < nReduce; r++ {

		wordCounts := make(map[string]int)

		for m := 0; m < nMap; m++ {

			file := fmt.Sprintf("int%d-%d.txt", m, r)

			data, err := os.ReadFile(file)

			if err != nil {
				fmt.Printf("read file error, 1")
				return
			}

			lines := strings.Split(string(data), "\n")

			for _, line := range lines {

				parts := strings.SplitN(line, " ", 2)

				word := parts[0]
				count, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
				wordCounts[word] += count
			}
		}

		out := fmt.Sprintf("red-%d.out", r)
		f, err := os.Create(out)

		if err != nil {
			fmt.Printf("create file error, 1")
			return
		}

		for word, count := range wordCounts {

			fmt.Fprintf(f, "%s %d\n", word, count)

		}
		f.Close()
	}

	totalCounts := make(map[string]int)

	for r := 0; r < nReduce; r++ {

		file := fmt.Sprintf("red-%d.out", r)

		content, err := os.ReadFile(file)

		if err != nil {
			fmt.Printf("read file error, 2")
			return
		}

		lines := strings.Split(string(content), "\n")

		for _, line := range lines {

			parts := strings.SplitN(line, " ", 2)
			word := parts[0]
			count, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			totalCounts[word] += count
		}
	}

	final := []WordCount{}

	for word, count := range totalCounts {

		final = append(final, WordCount{word, count})

	}
	sort.Slice(final, func(i, j int) bool {

		if final[i].Count == final[j].Count {

			return final[i].Word < final[j].Word

		}

		return final[i].Count > final[j].Count
	})

	out, err := os.Create(outFile)

	if err != nil {
		fmt.Printf("create file error, 2")
		return
	}

	defer out.Close()

	for _, wc := range final {

		fmt.Fprintf(out, "%s %d\n", wc.Word, wc.Count)

	}
}

func WordCount_GO(inFile string, outFile string) {

	file, _ := os.Open(inFile)

	defer file.Close()

	hashTable := make(map[string]int)
	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {

		word := scanner.Text()
		hashTable[word]++

	}

	counts := []WordCount{}

	for word, count := range hashTable {

		counts = append(counts, WordCount{word, count})

	}

	sort.Slice(counts, func(i, j int) bool {

		if counts[i].Count == counts[j].Count {

			return counts[i].Word < counts[j].Word

		}
		return counts[i].Count > counts[j].Count

	})

	out, _ := os.Create(outFile)

	defer out.Close()

	for _, wc := range counts {

		fmt.Fprintf(out, "%s %d\n", wc.Word, wc.Count)

	}
}

func WordCount_UNIX(inFile string, outFile string) {

	cmdStr := fmt.Sprintf(`cat "%s" |
	tr -s '[:space:]' '\n' |
	grep -v '^$' |
	sort |
	uniq -c |
	awk '{print $2, $1}' |
	sort -k2,2nr -k1,1 > "%s"`, inFile, outFile)

	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Run()
}

func main() {
	if len(os.Args) < 4 || len(os.Args) > 6 {
		fmt.Printf("%s:\n\tUsage: ./p3 run-mode inFile outFile <NMap> <NReduce>\n", os.Args[0])
		fmt.Printf("\tRun Modes:\n\t\t1 Unix Pipeline\n\t\t2 Simple wordcount in GO\n\t\t3 MapReduce Sequential\n\t\t4 MapReduce SMP\n\t\t5 MapReduce DMP\n")
		return
	}

	switch os.Args[1] {
	case "1":
		WordCount_UNIX(os.Args[2], os.Args[3])
	case "2":
		WordCount_GO(os.Args[2], os.Args[3])
	case "3":
		nMap, _ := strconv.Atoi(os.Args[4])
		nReduce, _ := strconv.Atoi(os.Args[5])
		WordCount_MR_S(os.Args[2], os.Args[3], nMap, nReduce)
	case "4":
		nMap, _ := strconv.Atoi(os.Args[4])
		nReduce, _ := strconv.Atoi(os.Args[5])
		WordCount_MR_SMP(os.Args[2], os.Args[3], nMap, nReduce)
	case "5":
		nMap, _ := strconv.Atoi(os.Args[4])
		nReduce, _ := strconv.Atoi(os.Args[5])
		WordCount_MR_DMP(os.Args[2], os.Args[3], nMap, nReduce)
	default:
		fmt.Printf("Unknown run-mode\n")
	}
}
