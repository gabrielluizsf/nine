package main

import (
	"fmt"
	"strings"

	"github.com/i9si-sistemas/command"
	"github.com/i9si-sistemas/safeos"
)

var Root = &safeos.Root{
	Dir: "./version",
}

func main() {
	b, err := Root.ReadFile("version")
	if err != nil {
		panic("??????????")
	}
	version := string(b)
	fmt.Println("🔄 creating a new tag...")
	if output, err := command.New().
		Execute("git", "tag", version).
		CombinedOutput(); err != nil {
		panic(fmt.Sprintf("❌ Cannot create tag %s: %s", version, strings.Join([]string{string(output), err.Error()}, "\n")))
	}
	fmt.Println("✅ tag created.")
	if output, err := command.New().
		Execute("git", "push", "origin", version).
		CombinedOutput(); err != nil {
		panic(fmt.Sprintf("❌ Cannot push tag %s: %s", version, strings.Join([]string{string(output), err.Error()}, "\n")))
	}
	fmt.Println("✅ tag pushed.")
}
