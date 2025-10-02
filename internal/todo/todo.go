package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Priority int

const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

func (p Priority) String() string {
	switch p {
	case PriorityHigh:
		return "HIGH"
	case PriorityMedium:
		return "MEDIUM"
	case PriorityLow:
		return "LOW"
	default:
		return "NONE"
	}
}

func ParsePriority(s string) Priority {
	switch strings.ToUpper(s) {
	case "HIGH", "H":
		return PriorityHigh
	case "MEDIUM", "MED", "M":
		return PriorityMedium
	case "LOW", "L":
		return PriorityLow
	default:
		return PriorityMedium
	}
}

type Item struct {
	Text      string
	Done      bool
	Priority  Priority
	DueDate   *time.Time `json:"DueDate,omitempty"`
	Tags      []string   `json:"Tags,omitempty"`
	CreatedAt time.Time
}

func NewItem(text string) Item {
	return Item{
		Text:      text,
		Done:      false,
		Priority:  PriorityMedium,
		Tags:      []string{},
		CreatedAt: time.Now(),
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

	// Sort: move completed tasks to the bottom
	l.Sort()
	return nil
}

// Sort reorders the list so incomplete tasks come first, completed tasks go to the bottom
func (l *List) Sort() {
	var incomplete []Item
	var completed []Item

	for _, item := range l.Items {
		if item.Done {
			completed = append(completed, item)
		} else {
			incomplete = append(incomplete, item)
		}
	}

	l.Items = append(incomplete, completed...)
}

// Delete removes an item from the list by index
func (l *List) Delete(index int) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	l.Items = append(l.Items[:index], l.Items[index+1:]...)
	return nil
}

// Edit changes the text of an existing item
func (l *List) Edit(index int, newText string) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	if newText == "" {
		return errors.New("Task text cannot be empty")
	}
	l.Items[index].Text = newText
	return nil
}

// Uncomplete marks a task as incomplete
func (l *List) Uncomplete(index int) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	l.Items[index].Done = false
	l.Sort()
	return nil
}

// ClearCompleted removes all completed tasks from the list
func (l *List) ClearCompleted() int {
	var incomplete []Item
	count := 0

	for _, item := range l.Items {
		if !item.Done {
			incomplete = append(incomplete, item)
		} else {
			count++
		}
	}

	l.Items = incomplete
	return count
}

// Stats represents statistics about the todo list
type Stats struct {
	Total     int
	Completed int
	Pending   int
}

// GetStats returns statistics about the todo list
func (l *List) GetStats() Stats {
	stats := Stats{Total: len(l.Items)}

	for _, item := range l.Items {
		if item.Done {
			stats.Completed++
		} else {
			stats.Pending++
		}
	}

	return stats
}

// SetPriority sets the priority of a task
func (l *List) SetPriority(index int, priority Priority) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	l.Items[index].Priority = priority
	return nil
}

// SetDueDate sets the due date of a task
func (l *List) SetDueDate(index int, dueDate time.Time) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	l.Items[index].DueDate = &dueDate
	return nil
}

// AddTag adds a tag to a task
func (l *List) AddTag(index int, tag string) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	// Check if tag already exists
	for _, t := range l.Items[index].Tags {
		if t == tag {
			return errors.New("Tag already exists")
		}
	}
	l.Items[index].Tags = append(l.Items[index].Tags, tag)
	return nil
}

// RemoveTag removes a tag from a task
func (l *List) RemoveTag(index int, tag string) error {
	if index < 0 || index >= len(l.Items) {
		return errors.New("Item index out of Range")
	}
	tags := l.Items[index].Tags
	for i, t := range tags {
		if t == tag {
			l.Items[index].Tags = append(tags[:i], tags[i+1:]...)
			return nil
		}
	}
	return errors.New("Tag not found")
}

// Search returns items that match the query in text or tags
func (l *List) Search(query string) []Item {
	var results []Item
	query = strings.ToLower(query)

	for _, item := range l.Items {
		// Search in text
		if strings.Contains(strings.ToLower(item.Text), query) {
			results = append(results, item)
			continue
		}

		// Search in tags
		for _, tag := range item.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, item)
				break
			}
		}
	}

	return results
}

// FilterByPriority returns items with the specified priority
func (l *List) FilterByPriority(priority Priority) []Item {
	var results []Item
	for _, item := range l.Items {
		if item.Priority == priority {
			results = append(results, item)
		}
	}
	return results
}

// FilterByTag returns items with the specified tag
func (l *List) FilterByTag(tag string) []Item {
	var results []Item
	for _, item := range l.Items {
		for _, t := range item.Tags {
			if t == tag {
				results = append(results, item)
				break
			}
		}
	}
	return results
}

// GetOverdue returns items that are past their due date
func (l *List) GetOverdue() []Item {
	var results []Item
	now := time.Now()

	for _, item := range l.Items {
		if item.DueDate != nil && item.DueDate.Before(now) && !item.Done {
			results = append(results, item)
		}
	}

	return results
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

		// Priority indicator
		prioritySymbol := ""
		switch item.Priority {
		case PriorityHigh:
			prioritySymbol = "🔴 "
		case PriorityMedium:
			prioritySymbol = "🟡 "
		case PriorityLow:
			prioritySymbol = "🟢 "
		}

		result += fmt.Sprintf("%d. [%s] %s%s", i+1, status, prioritySymbol, item.Text)

		// Add due date if present
		if item.DueDate != nil {
			dueStr := item.DueDate.Format("2006-01-02")
			if item.DueDate.Before(time.Now()) && !item.Done {
				result += fmt.Sprintf(" 📅 %s (OVERDUE!)", dueStr)
			} else {
				result += fmt.Sprintf(" 📅 %s", dueStr)
			}
		}

		// Add tags if present
		if len(item.Tags) > 0 {
			result += fmt.Sprintf(" 🏷️  %s", strings.Join(item.Tags, ", "))
		}

		result += "\n"
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
