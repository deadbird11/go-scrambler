package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

const (
	url   = "http://app.aspell.net/create?max_size=60&spelling=US&max_variant=2&diacritic=strip&special=hacker&special=roman-numerals&download=wordlist&encoding=utf-8&format=inline"
	fname = "anagram_map.gob"
)

type (
	anagramMap map[string][]string
	sortRunes  []rune
)

func main() {
	if am, ok := loadMap(); ok {
		input := strings.ToLower(os.Args[1])
		if result, exists := am[calcKey(input)]; exists {
			for _, val := range result {
				if val != input {
					fmt.Println(val)
				}
			}
		} else {
			fmt.Println("Word not recognized.")
		}
	} else {
		return
	}
}

// loadMap - returns a map of strings to slices of english words
// where each lists contains a set of anagrams
func loadMap() (anagramMap, bool) {
	if result, ok := loadCached(); ok {
		return result, true
	}

	if result, ok := loadRemote(); ok {
		return result, true
	}

	return anagramMap{}, false
}

// loadCached - loads cached map data from previous uses
func loadCached() (anagramMap, bool) {
	var result anagramMap

	file, err := os.Open(fname)
	if err != nil {
		return anagramMap{}, false
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println(err)
		return anagramMap{}, false
	}

	return result, true
}

// loadRemote - downloads the data from the internet and
// saves constructed map to a cache file
func loadRemote() (anagramMap, bool) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return anagramMap{}, false
	}
	defer resp.Body.Close()

	words, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return anagramMap{}, false
	}

	r := bufio.NewReader(strings.NewReader(string(words)))
	result := anagramMap{}
	recording := false

	for {
		if l, _, err := r.ReadLine(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			if recording {
				line := strings.ToLower(string(l))
				key := calcKey(line)
				result[key] = append(result[key], line)
			} else if string(l) == "---" {
				recording = true
			}
		}
	}

	if !writeCache(result) {
		return anagramMap{}, false
	}

	return result, true
}

func writeCache(am anagramMap) bool {
	cacheFile, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer cacheFile.Close()

	encoder := gob.NewEncoder(cacheFile)

	if err := encoder.Encode(am); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func calcKey(s string) string {
	cpy := sortRunes(s)
	sort.Sort(cpy)
	return string(cpy)
}

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Len() int {
	return len(s)
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
