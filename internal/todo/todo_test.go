package todo

import (
	"os"
	"testing"
)

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

	list.Add("Buy Milk")

	if len(list.Items) != 1 {
		t.Errorf("Expected 1 item in list, got '%d'", len(list.Items))
	}

	if list.Items[0].Text != "Buy Milk" {
		t.Errorf("Expected text to be 'Buy Milk', but got '%s'", list.Items[0].Text)
	}
}

func TestCompleteItem(t *testing.T) {
	list := NewList()

	list.Add("Buy Milk")

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
	list.Add("Task 1")
	list.Add("Task 2")
	list.Add("Task 3")
	list.Complete(0) // Task 1 gets completed and moves to bottom

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
	list.Add("Task 1")
	list.Add("Task 2")
	list.Add("Task 3")

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
	list.Add("Original Task")

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
	list.Add("Task 1")
	list.Add("Task 2")
	list.Complete(0) // Task 1 gets completed and moves to bottom (index 1)

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
	list.Add("Task 1")
	list.Add("Task 2")
	list.Add("Task 3")
	list.Add("Task 4")

	list.Complete(0)
	list.Complete(1)

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
	list.Add("Task 1")
	list.Add("Task 2")
	list.Add("Task 3")
	list.Complete(0)
	list.Complete(1)

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
	list.Add("Task 1")
	list.Add("Task 2")
	list.Add("Task 3")
	list.Add("Task 4")

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
