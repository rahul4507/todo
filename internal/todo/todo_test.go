package todo

import (
	"os"
	"strings"
	"testing"
	"time"
)

// Helper functions for tests to avoid unchecked errors
func mustAdd(t *testing.T, list *List, text string) {
	t.Helper()
	if err := list.Add(text); err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}
}

func mustComplete(t *testing.T, list *List, index int) {
	t.Helper()
	if err := list.Complete(index); err != nil {
		t.Fatalf("Failed to complete item: %v", err)
	}
}

func mustSetPriority(t *testing.T, list *List, index int, priority Priority) {
	t.Helper()
	if err := list.SetPriority(index, priority); err != nil {
		t.Fatalf("Failed to set priority: %v", err)
	}
}

func mustSetDueDate(t *testing.T, list *List, index int, dueDate time.Time) {
	t.Helper()
	if err := list.SetDueDate(index, dueDate); err != nil {
		t.Fatalf("Failed to set due date: %v", err)
	}
}

func mustAddTag(t *testing.T, list *List, index int, tag string) {
	t.Helper()
	if err := list.AddTag(index, tag); err != nil {
		t.Fatalf("Failed to add tag: %v", err)
	}
}

func TestNewItem(t *testing.T) {

	item := NewItem("Learn go")

	if item.Text != "Learn go" {
		t.Errorf("Expected text to be 'Learn go', got '%s'", item.Text)
	}

	if item.Done {
		t.Error("New item should no be marked as Done")
	}

}

func TestAddItem(t *testing.T) {
	list := NewList()

	err := list.Add("Buy Milk")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(list.Items) != 1 {
		t.Errorf("Expected 1 item in list, got '%d'", len(list.Items))
	}

	if list.Items[0].Text != "Buy Milk" {
		t.Errorf("Expected text to be 'Buy Milk', but got '%s'", list.Items[0].Text)
	}

	// Test adding duplicate item
	err = list.Add("Buy Milk")
	if err == nil {
		t.Error("Expected error when adding duplicate item")
	}
}

func TestCompleteItem(t *testing.T) {
	list := NewList()

	if err := list.Add("Buy Milk"); err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	err := list.Complete(0)

	if err != nil {
		t.Errorf("Unexpected Error : %v", err)
	}

	if !list.Items[0].Done {
		t.Error("Expected item to be marked as done")
	}

	err = list.Complete(1)
	if err == nil {
		t.Error("Expected Error when completing a not existent Item")
	}
}

func TestSaveAndLoad(t *testing.T) {

	// need to create a temp file for testing this
	tmpfile, err := os.CreateTemp("", "todo-test")
	if err != nil {
		t.Fatalf("Could not create the temp file: %v", err)
	}

	// Tear Down
	defer os.Remove(tmpfile.Name())

	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")
	mustComplete(t, list, 0) // Task 1 gets completed and moves to bottom

	if err := list.Save(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to save the file: %v", err)
	}

	// load the list from the file
	loadedList := NewList()
	if err := loadedList.Load(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to load list: %v", err)
	}

	// verify the loaded list matches original (3 items, sorted with incomplete first)
	if len(loadedList.Items) != 3 {
		t.Errorf("Expected 3 items in loaded list, got %d", len(loadedList.Items))
	}

	// After sorting: incomplete tasks (Task 2, Task 3) should be first
	if loadedList.Items[0].Text != "Task 2" {
		t.Errorf("Expected first item text to be 'Task 2' but got %s", loadedList.Items[0].Text)
	}

	if loadedList.Items[0].Done {
		t.Error("First item should be incomplete")
	}

	// Completed task (Task 1) should be last
	if loadedList.Items[2].Text != "Task 1" {
		t.Errorf("Expected last item text to be 'Task 1' but got %s", loadedList.Items[2].Text)
	}

	if !loadedList.Items[2].Done {
		t.Error("Last item should be completed")
	}
}

func TestDeleteItem(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")

	err := list.Delete(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(list.Items) != 2 {
		t.Errorf("Expected 2 items after deletion, got %d", len(list.Items))
	}

	if list.Items[1].Text != "Task 3" {
		t.Errorf("Expected second item to be 'Task 3', got '%s'", list.Items[1].Text)
	}

	// Test deleting invalid index
	err = list.Delete(10)
	if err == nil {
		t.Error("Expected error when deleting invalid index")
	}
}

func TestEditItem(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Original Task")

	err := list.Edit(0, "Updated Task")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.Items[0].Text != "Updated Task" {
		t.Errorf("Expected text to be 'Updated Task', got '%s'", list.Items[0].Text)
	}

	// Test editing with empty text
	err = list.Edit(0, "")
	if err == nil {
		t.Error("Expected error when editing with empty text")
	}

	// Test editing invalid index
	err = list.Edit(10, "Test")
	if err == nil {
		t.Error("Expected error when editing invalid index")
	}
}

func TestUncompleteItem(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustComplete(t, list, 0) // Task 1 gets completed and moves to bottom (index 1)

	// After sorting, completed task is at the end
	if !list.Items[1].Done {
		t.Error("Task should be completed before uncompleting")
	}

	if list.Items[1].Text != "Task 1" {
		t.Errorf("Expected 'Task 1' at index 1, got '%s'", list.Items[1].Text)
	}

	// Uncomplete the completed task (now at index 1)
	err := list.Uncomplete(1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// After uncompleting and re-sorting, both tasks are incomplete
	// They maintain their relative order, so Task 2 is still first
	if list.Items[1].Done {
		t.Error("Task at index 1 should be marked as incomplete")
	}

	if list.Items[1].Text != "Task 1" {
		t.Errorf("Expected 'Task 1' at index 1, got '%s'", list.Items[1].Text)
	}

	// Verify both are incomplete now
	for i, item := range list.Items {
		if item.Done {
			t.Errorf("All tasks should be incomplete, but item %d is done", i)
		}
	}

	// Test uncompleting invalid index
	err = list.Uncomplete(10)
	if err == nil {
		t.Error("Expected error when uncompleting invalid index")
	}
}

func TestClearCompleted(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")
	mustAdd(t, list, "Task 4")

	mustComplete(t, list, 0)
	mustComplete(t, list, 1)

	count := list.ClearCompleted()

	if count != 2 {
		t.Errorf("Expected 2 items cleared, got %d", count)
	}

	if len(list.Items) != 2 {
		t.Errorf("Expected 2 items remaining, got %d", len(list.Items))
	}

	for _, item := range list.Items {
		if item.Done {
			t.Error("No completed items should remain after clearing")
		}
	}
}

func TestGetStats(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")
	mustComplete(t, list, 0)
	mustComplete(t, list, 1)

	stats := list.GetStats()

	if stats.Total != 3 {
		t.Errorf("Expected total 3, got %d", stats.Total)
	}

	if stats.Completed != 2 {
		t.Errorf("Expected completed 2, got %d", stats.Completed)
	}

	if stats.Pending != 1 {
		t.Errorf("Expected pending 1, got %d", stats.Pending)
	}
}

func TestSort(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")
	mustAdd(t, list, "Task 4")

	// Complete some tasks
	list.Items[0].Done = true
	list.Items[2].Done = true

	// Sort the list
	list.Sort()

	// First two items should be incomplete
	if list.Items[0].Done || list.Items[1].Done {
		t.Error("First items should be incomplete after sorting")
	}

	// Last two items should be completed
	if !list.Items[2].Done || !list.Items[3].Done {
		t.Error("Last items should be completed after sorting")
	}

	if list.Items[0].Text != "Task 2" || list.Items[1].Text != "Task 4" {
		t.Error("Incomplete tasks not in correct order")
	}

	if list.Items[2].Text != "Task 1" || list.Items[3].Text != "Task 3" {
		t.Error("Completed tasks not in correct order")
	}
}

func TestPriorityString(t *testing.T) {
	tests := []struct {
		priority Priority
		expected string
	}{
		{PriorityHigh, "HIGH"},
		{PriorityMedium, "MEDIUM"},
		{PriorityLow, "LOW"},
		{Priority(99), "NONE"}, // Invalid priority
	}

	for _, tt := range tests {
		result := tt.priority.String()
		if result != tt.expected {
			t.Errorf("Priority.String() = %s, expected %s", result, tt.expected)
		}
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		input    string
		expected Priority
	}{
		{"HIGH", PriorityHigh},
		{"high", PriorityHigh},
		{"H", PriorityHigh},
		{"h", PriorityHigh},
		{"MEDIUM", PriorityMedium},
		{"medium", PriorityMedium},
		{"MED", PriorityMedium},
		{"M", PriorityMedium},
		{"LOW", PriorityLow},
		{"low", PriorityLow},
		{"L", PriorityLow},
		{"l", PriorityLow},
		{"invalid", PriorityMedium}, // Default
		{"", PriorityMedium},         // Default
	}

	for _, tt := range tests {
		result := ParsePriority(tt.input)
		if result != tt.expected {
			t.Errorf("ParsePriority(%s) = %v, expected %v", tt.input, result, tt.expected)
		}
	}
}

func TestSetPriority(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")

	// Test setting valid priority
	err := list.SetPriority(0, PriorityHigh)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.Items[0].Priority != PriorityHigh {
		t.Errorf("Expected priority HIGH, got %v", list.Items[0].Priority)
	}

	// Test setting different priorities
	err = list.SetPriority(1, PriorityLow)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.Items[1].Priority != PriorityLow {
		t.Errorf("Expected priority LOW, got %v", list.Items[1].Priority)
	}

	// Test invalid index (negative)
	err = list.SetPriority(-1, PriorityHigh)
	if err == nil {
		t.Error("Expected error for negative index")
	}

	// Test invalid index (too large)
	err = list.SetPriority(10, PriorityHigh)
	if err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestSetDueDate(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")

	dueDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	// Test setting valid due date
	err := list.SetDueDate(0, dueDate)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if list.Items[0].DueDate == nil {
		t.Error("DueDate should not be nil")
	} else if !list.Items[0].DueDate.Equal(dueDate) {
		t.Errorf("Expected due date %v, got %v", dueDate, list.Items[0].DueDate)
	}

	// Test invalid index (negative)
	err = list.SetDueDate(-1, dueDate)
	if err == nil {
		t.Error("Expected error for negative index")
	}

	// Test invalid index (too large)
	err = list.SetDueDate(10, dueDate)
	if err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestAddTag(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")

	// Test adding first tag
	err := list.AddTag(0, "work")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(list.Items[0].Tags) != 1 || list.Items[0].Tags[0] != "work" {
		t.Error("Tag was not added correctly")
	}

	// Test adding second tag
	err = list.AddTag(0, "urgent")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(list.Items[0].Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(list.Items[0].Tags))
	}

	// Test adding duplicate tag
	err = list.AddTag(0, "work")
	if err == nil {
		t.Error("Expected error when adding duplicate tag")
	}

	// Test invalid index (negative)
	err = list.AddTag(-1, "test")
	if err == nil {
		t.Error("Expected error for negative index")
	}

	// Test invalid index (too large)
	err = list.AddTag(10, "test")
	if err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestRemoveTag(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAddTag(t, list, 0, "work")
	mustAddTag(t, list, 0, "urgent")
	mustAddTag(t, list, 0, "personal")

	// Test removing existing tag
	err := list.RemoveTag(0, "urgent")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(list.Items[0].Tags) != 2 {
		t.Errorf("Expected 2 tags after removal, got %d", len(list.Items[0].Tags))
	}

	// Verify the correct tag was removed
	for _, tag := range list.Items[0].Tags {
		if tag == "urgent" {
			t.Error("Tag 'urgent' should have been removed")
		}
	}

	// Test removing non-existent tag
	err = list.RemoveTag(0, "nonexistent")
	if err == nil {
		t.Error("Expected error when removing non-existent tag")
	}

	// Test invalid index (negative)
	err = list.RemoveTag(-1, "work")
	if err == nil {
		t.Error("Expected error for negative index")
	}

	// Test invalid index (too large)
	err = list.RemoveTag(10, "work")
	if err == nil {
		t.Error("Expected error for out of range index")
	}
}

func TestSearch(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Buy groceries")
	mustAdd(t, list, "Learn Go programming")
	mustAdd(t, list, "Write documentation")
	mustAddTag(t, list, 0, "shopping")
	mustAddTag(t, list, 1, "coding")
	mustAddTag(t, list, 2, "coding")

	// Test search by text
	results := list.Search("go")
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].Text != "Learn Go programming" {
		t.Error("Wrong item returned from search")
	}

	// Test search by tag
	results = list.Search("coding")
	if len(results) != 2 {
		t.Errorf("Expected 2 results for tag search, got %d", len(results))
	}

	// Test case-insensitive search
	results = list.Search("BUY")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for case-insensitive search, got %d", len(results))
	}

	// Test no results
	results = list.Search("nonexistent")
	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}

	// Test partial match
	results = list.Search("doc")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for partial match, got %d", len(results))
	}
}

func TestFilterByPriority(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")
	mustAdd(t, list, "Task 4")

	mustSetPriority(t, list, 0, PriorityHigh)
	mustSetPriority(t, list, 1, PriorityMedium)
	mustSetPriority(t, list, 2, PriorityHigh)
	mustSetPriority(t, list, 3, PriorityLow)

	// Filter by high priority
	results := list.FilterByPriority(PriorityHigh)
	if len(results) != 2 {
		t.Errorf("Expected 2 high priority items, got %d", len(results))
	}

	// Filter by medium priority
	results = list.FilterByPriority(PriorityMedium)
	if len(results) != 1 {
		t.Errorf("Expected 1 medium priority item, got %d", len(results))
	}

	// Filter by low priority
	results = list.FilterByPriority(PriorityLow)
	if len(results) != 1 {
		t.Errorf("Expected 1 low priority item, got %d", len(results))
	}
}

func TestFilterByTag(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1")
	mustAdd(t, list, "Task 2")
	mustAdd(t, list, "Task 3")
	mustAdd(t, list, "Task 4")

	mustAddTag(t, list, 0, "work")
	mustAddTag(t, list, 1, "personal")
	mustAddTag(t, list, 2, "work")
	mustAddTag(t, list, 2, "urgent")

	// Filter by "work" tag
	results := list.FilterByTag("work")
	if len(results) != 2 {
		t.Errorf("Expected 2 items with 'work' tag, got %d", len(results))
	}

	// Filter by "personal" tag
	results = list.FilterByTag("personal")
	if len(results) != 1 {
		t.Errorf("Expected 1 item with 'personal' tag, got %d", len(results))
	}

	// Filter by non-existent tag
	results = list.FilterByTag("nonexistent")
	if len(results) != 0 {
		t.Errorf("Expected 0 items with non-existent tag, got %d", len(results))
	}

	// Filter by "urgent" tag
	results = list.FilterByTag("urgent")
	if len(results) != 1 {
		t.Errorf("Expected 1 item with 'urgent' tag, got %d", len(results))
	}
}

func TestGetOverdue(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task 1 - Past due")
	mustAdd(t, list, "Task 2 - Future due")
	mustAdd(t, list, "Task 3 - Past due but completed")
	mustAdd(t, list, "Task 4 - No due date")

	// Set past due date
	pastDate := time.Now().Add(-24 * time.Hour)
	mustSetDueDate(t, list, 0, pastDate)

	// Set future due date
	futureDate := time.Now().Add(24 * time.Hour)
	mustSetDueDate(t, list, 1, futureDate)

	// Set past due but mark as completed
	mustSetDueDate(t, list, 2, pastDate)
	mustComplete(t, list, 2)

	// Get overdue items
	results := list.GetOverdue()

	// Should only return Task 1 (past due and not completed)
	if len(results) != 1 {
		t.Errorf("Expected 1 overdue item, got %d", len(results))
	}

	if len(results) > 0 && results[0].Text != "Task 1 - Past due" {
		t.Errorf("Expected 'Task 1 - Past due', got '%s'", results[0].Text)
	}
}

func TestListString(t *testing.T) {
	list := NewList()

	// Test empty list
	result := list.String()
	if result != "No items to return" {
		t.Errorf("Expected 'No items to return' for empty list, got '%s'", result)
	}

	// Add items with various attributes
	mustAdd(t, list, "Simple task")
	mustAdd(t, list, "High priority task")
	mustAdd(t, list, "Task with due date")
	mustAdd(t, list, "Task with tags")
	mustAdd(t, list, "Completed task")

	// Set priority
	mustSetPriority(t, list, 1, PriorityHigh)

	// Set due date
	dueDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	mustSetDueDate(t, list, 2, dueDate)

	// Add tags
	mustAddTag(t, list, 3, "work")
	mustAddTag(t, list, 3, "urgent")

	// Complete a task
	mustComplete(t, list, 4)

	// Get string representation
	result = list.String()

	// Verify it contains expected elements
	if !strings.Contains(result, "TODO List:") {
		t.Error("String output should contain 'TODO List:'")
	}

	if !strings.Contains(result, "Simple task") {
		t.Error("String output should contain task text")
	}

	if !strings.Contains(result, "🔴") {
		t.Error("String output should contain high priority symbol")
	}

	if !strings.Contains(result, "2025-12-31") {
		t.Error("String output should contain due date")
	}

	if !strings.Contains(result, "work") {
		t.Error("String output should contain tags")
	}

	if !strings.Contains(result, "✓") {
		t.Error("String output should contain completed symbol")
	}
}

func TestGetOverdueWithOverdueDate(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Overdue task")

	// Set a date in the past
	overdueDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	mustSetDueDate(t, list, 0, overdueDate)

	results := list.GetOverdue()

	if len(results) != 1 {
		t.Errorf("Expected 1 overdue task, got %d", len(results))
	}

	// Verify the String() method shows OVERDUE
	output := list.String()
	if !strings.Contains(output, "OVERDUE") {
		t.Error("String output should show OVERDUE for past due dates")
	}
}

func TestItemCreation(t *testing.T) {
	item := NewItem("Test task")

	// Verify default values
	if item.Text != "Test task" {
		t.Errorf("Expected text 'Test task', got '%s'", item.Text)
	}

	if item.Done {
		t.Error("New item should not be marked as done")
	}

	if item.Priority != PriorityMedium {
		t.Errorf("Expected default priority MEDIUM, got %v", item.Priority)
	}

	if len(item.Tags) != 0 {
		t.Errorf("Expected empty tags, got %d tags", len(item.Tags))
	}

	if item.DueDate != nil {
		t.Error("New item should not have a due date")
	}

	if item.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestSaveError(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Test task")

	// Try to save to an invalid path (directory doesn't exist)
	err := list.Save("/nonexistent/directory/todos.json")
	if err == nil {
		t.Error("Expected error when saving to invalid path")
	}
}

func TestLoadError(t *testing.T) {
	list := NewList()

	// Try to load from non-existent file
	err := list.Load("/nonexistent/file.json")
	if err == nil {
		t.Error("Expected error when loading non-existent file")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	// Create a temp file with invalid JSON
	tmpfile, err := os.CreateTemp("", "invalid-json")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write invalid JSON
	_, err = tmpfile.WriteString("invalid json content {{{")
	if err != nil {
		t.Fatalf("Could not write to temp file: %v", err)
	}
	tmpfile.Close()

	// Try to load invalid JSON
	list := NewList()
	err = list.Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestStringWithAllPriorities(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Low priority")
	mustAdd(t, list, "Medium priority")
	mustAdd(t, list, "High priority")

	mustSetPriority(t, list, 0, PriorityLow)
	mustSetPriority(t, list, 1, PriorityMedium)
	mustSetPriority(t, list, 2, PriorityHigh)

	output := list.String()

	// Verify all priority symbols are present
	if !strings.Contains(output, "🟢") {
		t.Error("String output should contain low priority symbol")
	}

	if !strings.Contains(output, "🟡") {
		t.Error("String output should contain medium priority symbol")
	}

	if !strings.Contains(output, "🔴") {
		t.Error("String output should contain high priority symbol")
	}
}

func TestStringWithFutureDueDate(t *testing.T) {
	list := NewList()
	mustAdd(t, list, "Task with future due date")

	// Set a future due date (not overdue)
	futureDate := time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
	mustSetDueDate(t, list, 0, futureDate)

	output := list.String()

	// Should contain date but NOT "OVERDUE"
	if !strings.Contains(output, "2099-12-31") {
		t.Error("String output should contain due date")
	}

	if strings.Contains(output, "OVERDUE") {
		t.Error("String output should not show OVERDUE for future dates")
	}
}
