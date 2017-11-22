package main

import (
	"regexp"
	"time"
	"math/rand"
    "strings"
	"net/http"
	"html/template"
)

type text struct{
	input string
}

func ElizaResponse(input string, w http.ResponseWriter, r *http.Request) string {
    
    //Parse the input that the user inputs in the webpage
	r.ParseForm();
	guess := r.FormValue("input")
    
    //array of responses the eliza program will respond with
	responses := []string{
		"I’m not sure what you’re trying to say. Could you explain it to me?",
		"How does that make you feel?",
		"Why do you say that?",
	}

	iam := regexp.MustCompile("(?i)i(?:'| a|)?m(.*)")
	if iam.MatchString(guess) {
		return iam.ReplaceAllString(input, "How do I know you are $1")
	}

	//Adapted from https://golang.org/pkg/regexp/
	//Searchs the input to look for the word "Father" on its own and assigns it to this variable
	father, _ := regexp.MatchString(`(?i)\\bfather\\b`, input)
	
	//if the user input contains the word "father" it will return this string
	if (father) {
		return("Why don’t you tell me more about your father?")
	}

	//Puts the data into the guess.tmpl file
	m := text{input: responses[rand.Intn(len(responses))]}
	t, _ := template.ParseFiles("chatbot.tmpl")
	t.Execute(w, m)

	//returns the responses to the main function
	return responses[rand.Intn(len(responses))]
}

func reflection(input string) string{

	// List the pronouns to switch
	pronouns := [][]string{
		{`am`, `are`},
		{`I`, `you`},
		{`you`, `I`},
		{`me`, `you`},
		{`your`, `my`},
		{`my`, `your`},
	}

	// Split input into values
	boundaries := regexp.MustCompile(`\b`)

	values := boundaries.Split(input, -1)

	//Loop through the range of values and reflect the pronoun if it finds a match
	for i, token := range values {
		for _, reflection := range pronouns {
			if matched, _ := regexp.MatchString(reflection[0], token); matched {
				
				values[i] = reflection[1]
				break
			}
		}
	}
	
	//Join the string of values back together
	answer := strings.Join(values, ``)

	counterResp := []string{
		"Why do ",
		"How do you know that ",
		"I find it fasinating that ",
    }

	return (counterResp[rand.Intn(len(counterResp))] + answer)
}

func handler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "chatbot.html")
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}