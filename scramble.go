package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

const url = "https://svnweb.freebsd.org/csrg/share/dict/words?revision=61569&view=co"

type (
	anagramMap map[string][]string
	sortRunes  []rune
)

func main() {
	if am, err := loadMap(); err == nil {
		fmt.Println(am)
	} else {
		panic(err)
	}
}

func loadMap() (anagramMap, error) {
	return setupFromInternet()
}

func setupFromInternet() (anagramMap, error) {
	resp, err := http.Get(url)
	if err != nil {
		return anagramMap{}, err
	}
	defer resp.Body.Close()

	words, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return anagramMap{}, err
	}

	r := bufio.NewReader(strings.NewReader(string(words)))
	result := anagramMap{}

	for {
		if l, _, err := r.ReadLine(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			line := strings.ToLower(string(l))
			key := calcKey(line)
			result[key] = append(result[key], line)
		}
	}

	return result, nil
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
