package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"
	"math/rand"
	"strings"
	"strconv"
	"net/http"
	"fmt"
)

type Pronoun struct {
	userPronouns     *regexp.Regexp
	newPronouns []string
}

type Eliza struct {
	responses     []Pronoun
	wordSwitch []Pronoun
}

//Reads a .txt file and casts the string values into an array
func readPronouns(filePath string) []Pronoun {
	//Adapted from https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
	//Opens the given file, log error if it fails
	file, err := os.Open(filePath)
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

//Response function, generates elizas answer
func (me *Eliza) E_Responder(input string) string {
	//Loop to control and find a response
	for _, response := range me.responses {
		//If to check if the user input matches anything
		if matchFound := response.userPronouns.FindStringSubmatch(input); matchFound != nil {
			//Regexp
			boundaries := regexp.MustCompile(`[\s,.?!]+`)
			//Gets a random response from the botResponses.txt
			e_answer := response.newPronouns[rand.Intn(len(response.newPronouns))]
			//Loop to switch the pronouns 
			for m, match := range matchFound[1:] {
				tokens := boundaries.Split(match, -1)
				for t, token := range tokens {
					for _, substitution := range me.wordSwitch {
						if substitution.userPronouns.MatchString(token) {
							tokens[t] = substitution.newPronouns[rand.Intn(len(substitution.newPronouns))]
							break
						}
					}
				}
				e_answer = strings.Replace(e_answer, "$"+strconv.Itoa(m+1), strings.Join(tokens, " "), -1)
			
			}
			return e_answer
		
		}
	}
	return "I cannot answer you."
}

func ElizaFromFiles(responsesPath string, pronounPath string) Eliza {
	
	getResponse := Eliza{}
	getResponse.wordSwitch = readPronouns(pronounPath)
	getResponse.responses = readPronouns(responsesPath)

	return getResponse
}

func elizaHandler(w http.ResponseWriter, r *http.Request) {
	
		userInput := r.URL.Query().Get("input")
		eliza := ElizaFromFiles("botResponses.txt", "pronouns.txt")
		fmt.Fprintf(w, eliza.E_Responder(userInput))
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