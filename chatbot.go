package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"
	"math/rand"
	"net/http"
	"fmt"
)

type Pronoun struct {
	userPronouns     *regexp.Regexp
	newPronouns []string
}

//Reads a .txt file and casts the string values into an array
func readPronouns(path string) []Pronoun {
	//Adapted from https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
	//Opens the given file, log error if it fails
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	//close the file once done
	defer file.Close()

	//variable of type Pronoun with empty array
	var switchPros []Pronoun
	//for loop that reads in each line
	for scanner, readCheck := bufio.NewScanner(file), false; scanner.Scan(); {
		switch line := scanner.Text(); {
			//if the lenght of the line is 0 then we've reached the end of the line
			case len(line) == 0:
				//set check to false
				readCheck = false
			//if check is false then append the rest of the line to the array and reset check to true
			case readCheck == false:
				switchPros = append(switchPros, Pronoun{userPronouns: regexp.MustCompile(line)})
				readCheck = true
			default:
				//switch the users pronouns with the reflected pronoun from the read in file
				switchPros[len(switchPros)-1].newPronouns = append(switchPros[len(switchPros)-1].newPronouns, line)
		}
	}
	//return the switched pronouns
	return switchPros
}


func elizaHandler(w http.ResponseWriter, r *http.Request) {
	
		userInput := r.URL.Query().Get("input")
		fmt.Fprintf(w, userInput)
}

func main() {
	//To pick a random response from the ElizaResponse func
	rand.Seed(time.Now().UTC().UnixNano())

    // Webpage Directory
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/chatbot", elizaHandler)

	http.ListenAndServe(":8080", nil)
	
}