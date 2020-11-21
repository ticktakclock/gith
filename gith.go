package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
)

type Branch []string

func isRemoteBranch(branch string) bool {
	return strings.HasPrefix(string(branch), "remotes/")
}

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

	if selected == "" {
		fmt.Println("Nothing to do.")
		return
	}

	fmt.Printf("Try to checkout branch: %q\n", selected)
	if isRemoteBranch(selected) {
		var branch Branch = []string{strings.TrimPrefix(selected, "remotes/origin/")}
		checkoutBranchName, err := branch.SelectOneWithAdd()
		if err != nil {
			fmt.Println(err)
			return
		}
		GitCheckoutWithRemote(checkoutBranchName, selected)
	} else {
		GitCheckout(selected)
	}
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
		Label:    "Which branch do you want to checkout?",
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

func (branches Branch) SelectOneWithAdd() (string, error) {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    "What is a name of branch?",
			Items:    branches,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			branches = append(branches, result)
		}
	}
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	fmt.Printf("You choose %s\n", result)
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

func GitCheckoutWithRemote(branchName string, remoteBranchName string) {
	out, err := exec.Command("git", "checkout", "-b", branchName, remoteBranchName).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(out))
}

func GitBranch() Branch {
	out, err := exec.Command("git", "branch", "-a").Output()

	if err != nil {
		fmt.Println(err)
		return Branch{}
	}
	s := string(out)

	var items Branch = strings.Split(s, "\n")
	return items.Trim()
}
