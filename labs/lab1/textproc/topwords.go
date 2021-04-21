// Find the top K most common words in a text document.
// Input path: location of the document, K top words
// Output: Slice of top K words
// For this excercise, word is defined as characters separated by a whitespace

// Note: You should use `checkError` to handle potential errors.

//Damian Hupka
//Cloud Native Application Architecture - ECGR 4090-Y04 Sp.21
//02/03/21

package textproc

import (
	"fmt"
	"log"
	"sort"
	"os"
	"bufio"

)

func topWords(path string, K int) []WordCount {
	// Your code here.....
	file, err := os.Open(path)
	checkError(err)

	defer file.Close()

	//Initializing scanner
	scanner := bufio.NewScanner(file)

	//Returns string with whitespace split and trailing spaces removed
	//Split will default to scan lines, this could similarly be used but adding an additional space argument 
	scanner.Split(bufio.ScanWords)

	//Initializing map to store key valued pairs of words and counts
	wordMap := make(map[string]int)

	//Initializing slice to return using WordCount struct 
	wordSlice := make([]WordCount, 0 )

	//Reading file line by line and adding word occurences to map
	for scanner.Scan(){
		//Most recent token generated as a string
		word := scanner.Text()
		//Stores initial key word and will increment value to one, etc.
		wordMap[word]++
	}

	//Adding wordMap entries to WordCount slice 
	for key, value := range wordMap {
		wordSlice = append(wordSlice, WordCount{key, value})
	}
	//Sorting output slice
	sortWordCounts(wordSlice)

	//Returning top K word frequencies (3 in the test case)
	return wordSlice[:K]
}

//--------------- DO NOT MODIFY----------------!

// A struct that represents how many times a word is observed in a document
type WordCount struct {
	Word  string
	Count int
}

// Method to convert struct to string format
func (wc WordCount) String() string {
	return fmt.Sprintf("%v: %v", wc.Word, wc.Count)
}

// Helper function to sort a list of word counts in place.
// This sorts by the count in decreasing order, breaking ties using the word.

func sortWordCounts(wordCounts []WordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		wc1 := wordCounts[i]
		wc2 := wordCounts[j]
		if wc1.Count == wc2.Count {
			return wc1.Word < wc2.Word
		}
		return wc1.Count > wc2.Count
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
