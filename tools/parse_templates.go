package main

import (
	"fmt"
	"html/template"
)

func main() {
	files := []string{"cmd/web/templates/base.layout.gohtml", "cmd/web/templates/terminal.page.gohtml"}
	t, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Println("Parse error:", err)
		return
	}
	fmt.Println("Parsed templates OK:", t.Name())
}
