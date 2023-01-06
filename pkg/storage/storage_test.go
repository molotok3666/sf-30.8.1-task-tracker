package storage

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var st *Storage

func TestNew(t *testing.T) {
	_, err := connect()
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_Tasks(t *testing.T) {
	st, _ := connect()
	data, err := st.Tasks(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
	data, err = st.Tasks(1, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestStorage_NewUser(t *testing.T) {
	st, _ := connect()
	user1 := User{
		Name: "user1",
	}

	_, err := st.NewUser(user1)
	if err != nil {
		log.Fatal(err)
	}

	user2 := User{
		Name: "user2",
	}
	_, err = st.NewUser(user2)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_NewTask(t *testing.T) {
	task := Task{
		AuthorID:   1,
		AssignedID: 2,
		Title:      "task",
		Content:    "complete task",
	}
	st, _ := connect()

	_, err := st.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_UpdateTask(t *testing.T) {
	st, _ := connect()
	task := Task{
		AuthorID:   1,
		AssignedID: 2,
		Title:      "task",
		Content:    "complete task",
	}

	taskId, err := st.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}
	task = Task{
		ID:         taskId,
		AuthorID:   2,
		AssignedID: 1,
		Title:      "new-task",
		Content:    "update task",
	}

	err = st.UpdateTask(task)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_DeleteTask(t *testing.T) {
	st, _ := connect()
	task := Task{
		AuthorID:   1,
		AssignedID: 2,
		Title:      "task",
		Content:    "complete task",
	}

	taskId, err := st.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}

	err = st.DeleteTask(taskId)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_NewLabel(t *testing.T) {
	st, _ := connect()
	label := Label{
		Name: "task-label",
	}
	_, err := st.NewLabel(label)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_NewTaskLabel(t *testing.T) {
	st, _ := connect()
	task := Task{
		AuthorID:   1,
		AssignedID: 2,
		Title:      "task",
		Content:    "complete task",
	}

	taskId, err := st.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}

	label := Label{
		Name: "task-label",
	}
	labelId, err := st.NewLabel(label)
	if err != nil {
		log.Fatal(err)
	}

	err = st.NewTaskLabel(taskId, labelId)
	if err != nil {
		log.Fatal(err)
	}
}

func TestStorage_TasksByLabel(t *testing.T) {
	st, _ := connect()
	task := Task{
		AuthorID:   1,
		AssignedID: 2,
		Title:      "task",
		Content:    "complete task",
	}

	taskId, err := st.NewTask(task)
	if err != nil {
		log.Fatal(err)
	}

	label := Label{
		Name: "task-label",
	}
	labelId, err := st.NewLabel(label)
	if err != nil {
		log.Fatal(err)
	}

	err = st.NewTaskLabel(taskId, labelId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = st.TasksByLabel(label.Name)
	if err != nil {
		log.Fatal(err)
	}
}

func connect() (*Storage, error) {
	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal(err)
	}

	user := os.Getenv("POSTGRES_USER")
	pwd := os.Getenv("POSTGRES_PASSWORD")
	dbService := os.Getenv("POSTGRES_DB_SERVICE")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || pwd == "" || dbService == "" || dbPort == "" || dbName == "" {
		return nil, errors.New("Empty environment variables")
	}

	connstr := "postgres://" + user + ":" + pwd +
		"@" + dbService + ":" + dbPort + "/" + dbName
	st, err := New(connstr)

	return st, err
}
