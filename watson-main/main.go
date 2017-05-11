package main

import (
	watson "github.com/SuccessRain/ibm-fpt-sdk-master"
	conversation "github.com/SuccessRain/ibm-fpt-sdk-master/conversation"
	"fmt"
	"os"
	"bufio"
	"strings"
	"flag"
)

type Watson struct {
	user string
	pass string
}

func New(user, pass string) *Watson {
	return &Watson{
		user: user,
		pass: pass,
	}
}

var target string
var inputFP string
var workspaces string
var username string
var password string

func AddSharedFlags(fs *flag.FlagSet) {
	fs.StringVar(&target, "t", "", "required, intent or entity")
	fs.StringVar(&inputFP, "i", "", "required, path to the input file")
	fs.StringVar(&workspaces, "w", "", "required, IBM workspaces id")
	fs.StringVar(&username, "u", "", "required, username conversation")
	fs.StringVar(&password, "p", "", "required, password conversation")
}

type ObjectFile struct{
	Intent string
	Name string
	Response string
}

func ReadIntentsFromFile(inputFP string) ([]ObjectFile, error) {
	input, err := os.Open(inputFP)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	var context []ObjectFile
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var con ObjectFile
		text := strings.Replace(scanner.Text(), `"`, ``, -1)
		tokens := strings.SplitN(text, ",", 2)
		//fmt.Print("TOKEN: \t"); fmt.Println(tokens)
		con.Intent, con.Name = strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		context = append(context, con)
		//fmt.Print("CONTEXT: \t"); fmt.Println(con.Intent); fmt.Println(con.Name); fmt.Println(con.Response)
	}

	return context, nil
}

func TestIntent(c conversation.Client, workspace_id string) {

	//fmt.Print("Intents Value:\t"); fmt.Println(reply.Intents[0].Intent)

	count := 0
	obj, err := ReadIntentsFromFile(inputFP)
	if err != nil{
		fmt.Println(err)
	}
	for i := 0; i < len(obj); i++ {
		reply, err := c.Message(workspace_id, obj[i].Name)
		if err != nil {
			fmt.Print("Message() failed:\t"); fmt.Println(err)
			return
		}
		if len(reply.Intents) > 0 {
			//fmt.Print("AAAAAAAAAAAAAAAA:\t"); fmt.Println(reply.Intents[0].Intent, obj[i].Intent, obj[i].Name)
			if strings.Compare(strings.TrimSpace(reply.Intents[0].Intent), strings.TrimSpace(obj[i].Intent)) == 0 {
				fmt.Print(i+1); fmt.Println(".\t Correct")
				count ++
			}else{
				fmt.Print(i+1); fmt.Println(".\t Incorrect")
			}
		}else{
			fmt.Print(i+1); fmt.Println(".\t Incorrect")
			//fmt.Print("BBBBBBBBBBBBBBBB:\t"); fmt.Println(obj[i].Intent, obj[i].Name)
		}
	}
	fmt.Printf("Success: %f \n", float64(count) / float64(len(obj)) * 100)
}

func test(c conversation.Client, workspace_id string){
	reply, err := c.Message(workspace_id, "hello")
	if err != nil {
		fmt.Print("Message() failed:\t"); fmt.Println(err)
		return
	}

	if len(reply.Intents) > 0 {
		fmt.Print("AAAAAAAAAAAAAAAA:\t"); fmt.Println(reply.Intents[0].Intent)
	}
	fmt.Print("BBBBBBBBBBBBBBBB:\t"); fmt.Println(reply)
}

func main() {

	trainCmd := flag.NewFlagSet("train", flag.ExitOnError)
	AddSharedFlags(trainCmd)

	testCmd := flag.NewFlagSet("test", flag.ExitOnError)
	AddSharedFlags(testCmd)

	if len(os.Args) < 2 {
		fmt.Println("Error: Input is not enough")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "train"{
		trainCmd.Parse(os.Args[2:])
	}else if command == "test" {
		testCmd.Parse(os.Args[2:])
	}else if command == "help" {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	if target != "intent" && target != "entity" {
		fmt.Println("Error: You must choose intent or entity")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if inputFP == "" {
		fmt.Println("Error: Input file is required but empty")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if workspaces == "" {
		fmt.Println("Error: Workspaces id is required")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if username == "" {
		fmt.Println("Error: username is required")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if password == "" {
		fmt.Println("Error: password id is required")
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	var creden watson.Credentials = watson.Credentials{
		Username: username,
		Password: password,
	}

	//"6fc4a9a2-aa9c-4e97-96a8-0b37a9528497" - CongVV
	//"a930e0a7-8624-4123-b682-fa97bb36ac1d" - Test
	//"2f8b859e-9df3-46af-a951-d2ad9d4125bb" - Test 2

	c, err := conversation.NewClient(watson.Config{Credentials: creden, })
	if err != nil {
		fmt.Print("NewClient() failed:\t"); fmt.Println(err)
		return
	}

	TestIntent(c, workspaces)

	/*
	var creden watson.Credentials = watson.Credentials{
		Username: "9a4a2bb2-12a0-469c-96bb-0a8387e0ae21",
		Password: "FdvIAuso3rNG",
		//Username: "647bf7bb-eb41-4884-b608-f6a4acd0f43f",
		//Password: "Zbj8CrIzQOvg",
	}
	c, err := conversation.NewClient(watson.Config{Credentials: creden, })
	if err != nil {
		fmt.Print("NewClient() failed:\t"); fmt.Println(err)
		return
	}
	test(c, "2f8b859e-9df3-46af-a951-d2ad9d4125bb")
	*/
}

const helpMessage string = `
api is CLI tool that helps you train and test IBM Watson in terminal

Usage: api <command> <option>
Available commands and corresponding options:
	train
	  -t string
	    	required, type of training (intent, entity)
	  -i string
	    	required, path to your input file
	  -workspaces id string
	  		required, IBM Watson workspaces id

	test
	  -t string
	    	required, type of training (intent, entity)
	  -i string
	    	required, path to your input file
	  -workspaces id string
	  		required, IBM Watson workspaces id

	help
`