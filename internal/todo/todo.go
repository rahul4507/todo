package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Item struct {
	Text string
	Done bool
}

func NewItem(text string) Item {
	return Item{
		Text: text,
		Done: false,
	}
}

type List struct {
	Items []Item
}

func NewList() *List {
	return &List{
		Items: []Item{},
	}
}

func (l *List) Add(text string) error {
	item := NewItem(text)
	// here check that this should not be in the list already
	for _, existing := range l.Items {
		if existing.Text == item.Text {
			return errors.New("Item already exists in the list")
		}
	}
	l.Items = append(l.Items, item)
	return nil
}

func (l *List) Complete(index int) error {
	// here not need to extract that particular element wrapped in a try exception
	// todo: lets just as a basic one for now
	// add a check
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	l.Items[index].Done = true
	return nil
}

func (l *List) String() string {
	if len(l.Items) == 0 {
		return "No items to return"
	}

	result := "TODO List:\n"

	for i, item := range l.Items {
		status := " "
		if item.Done {
			status = "✓"
		}
		result += fmt.Sprintf("%d. [%s] %s\n", i+1, status, item.Text)
	}
	return result
}

// Save writes the todo list to a file in JSON format
func (l *List) Save(filename string) error {
	data, err := json.Marshal(l)

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Load reads a todo list from a file
func (l *List) Load(filename string) error {
	data, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	return json.Unmarshal(data, l)
}
