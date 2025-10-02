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
	list.Complete(0)

	if err := list.Save(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to save the file: %v", err)
	}

	// load the list from the file
	loadedList := NewList()
	if err := loadedList.Load(tmpfile.Name()); err != nil {
		t.Fatalf("Failed to load list: %v", err)
	}

	// verify the loaded list matches original
	if len(loadedList.Items) != 2 {
		t.Errorf("Expected 2 items in loaded list, got %d", len(loadedList.Items))
	}

	if !loadedList.Items[0].Done {
		t.Errorf("Expected First Item to be complted")

	}

	if loadedList.Items[0].Text != "Task 1" {
		t.Errorf("Expected Text should be 'Task 1' but got %s", loadedList.Items[0].Text)
	}
}
