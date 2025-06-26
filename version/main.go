package main

import (
	"fmt"

	"github.com/i9si-sistemas/command"
	"github.com/i9si-sistemas/safeos"
	"github.com/i9si-sistemas/stringx"
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
	join := func(out []byte, err error) string {
		return stringx.ConvertStrings(string(out), err.Error()).Join("\n").String()
	}
	if output, err := command.New().
		Execute("git", "tag", version).
		CombinedOutput(); err != nil {
		panic(fmt.Sprintf("❌ Cannot create tag %s: %s", version, join(output, err)))
	}
	fmt.Println("✅ tag created.")
	if output, err := command.New().
		Execute("git", "push", "origin", version).
		CombinedOutput(); err != nil {
		panic(fmt.Sprintf("❌ Cannot push tag %s: %s", version, join(output, err)))
	}
	fmt.Println("✅ tag pushed.")
}
