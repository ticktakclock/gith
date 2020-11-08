package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

type Branch []string

func main() {
	branches := GitBranch()
	if len(branches) <= 0 {
		fmt.Println("No Branches on this repository.")
		return
	}
	selected, err := branches.SelectOne()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Try to checkout branch: %q\n", selected)
	GitCheckout(selected)
}

func (this Branch) Trim() Branch {
	new := make(Branch, len(this))
	for index, branch := range this {
		new[index] = strings.Trim(strings.Trim(branch, "*"), " ")
	}
	return new
}

func (branches Branch) SelectOne() (string, error) {
	searcher := func(input string, index int) bool {
		item := branches[index]
		name := strings.Replace(strings.ToLower(item), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}
	prompt := promptui.Select{
		Label:    "Select Day",
		Items:    branches,
		Searcher: searcher,
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

func GitCheckout(branchName string) {
	out, err := exec.Command("git", "checkout", branchName).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}

func GitBranch() Branch {
	out, err := exec.Command("git", "branch").Output()

	if err != nil {
		fmt.Println(err)
		return Branch{}
	}
	s := string(out)

	var items Branch = strings.Split(s, "\n")
	return items.Trim()
}
