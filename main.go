package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
)

var aliases = make(map[string]string)
var envar = make(map[string]string)
var envcheck, _ = regexp.Compile("[A-Z1-9]+=.+")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func runshuwu(wherethefuckisthefile string) {
	wtfitf := wherethefuckisthefile + "/.shuwurc" //i wanna be funny but "wherethefuckisthefile" is too long
	if _, err := os.Stat(wtfitf); err == nil {
		file, err := os.ReadFile(wtfitf)
		check(err)
		stfile := string(file)
		for _, script := range strings.Split(stfile, "\n") {
			if script == "" {
				print("")
			} else {
				runCommand(script)
			}

		}
	} else {
		myfile, e := os.Create(wtfitf)
		if e != nil {
			log.Fatal(e)
		}
		myfile.Close()
	}
}

func runCommand(commandStr string) error {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)

	if envcheck.MatchString(arrCommandStr[0]) {
		bruh := strings.Split(strings.Join(arrCommandStr, " "), "=")
		envar[bruh[0]] = bruh[1]
		return nil
	}

	switch arrCommandStr[0] {
	case "envs":
		for key, value := range envar {
			println(key, value)
		}
		return nil
	case "cd":
		if len(arrCommandStr) < 2 {
			os.Chdir(gethome())
		} else {
			os.Chdir(arrCommandStr[1])
		}
		return nil
	case "pwd":
		println(getcurrentdir())
		return nil
	case "exit":
		os.Exit(0)
	case "help":
		println("idfk man")
		return nil
	}
	if arrCommandStr[0] == "alias" {
		aliases[arrCommandStr[1]] = strings.Join(arrCommandStr[2:], " ")
		return nil
	}

	if command, exists := aliases[arrCommandStr[0]]; exists {
		runCommand(command)
		return nil
	}

	{
		cmd := exec.Command(arrCommandStr[0], arrCommandStr[1:]...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
}

func gethome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func getcurrentdir() string {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return mydir
}
func main() {

	runshuwu(gethome())
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(envar["PS1"] + " ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = runCommand(cmdString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
