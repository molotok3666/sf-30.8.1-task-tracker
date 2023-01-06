package storage

type Interface interface {
	Tasks(int, int) ([]Task, error)
	NewTask(Task) (int, error)
	DeleteTask(int)
	NewUser(string) (int, error)
	NewLabel(string) (int, error)
	NewTaskLabel(Task, Label)
}
