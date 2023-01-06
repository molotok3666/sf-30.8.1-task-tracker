package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// Задача.
type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Tasks возвращает список задач из БД.
func (s *Storage) Tasks(taskID int, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
			INSERT INTO tasks (author_id, assigned_id, title, content)
			VALUES ($1, $2, $3, $4) RETURNING id;
		`,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

func (s *Storage) UpdateTask(t Task) error {
	var id int
	err := s.db.QueryRow(context.Background(), `
			UPDATE tasks
			SET opened = $1, closed = $2, author_id = $3, assigned_id = $4, title = $5, content = $6
			WHERE id = $7
			RETURNING id
			;
		`,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
		t.ID,
	).Scan(&id)
	return err
}

// Удаляет задачу по ID
func (s *Storage) DeleteTask(id int) error {
	err := s.db.QueryRow(context.Background(), `
		DELETE FROM tasks WHERE id = $1 RETURNING id
	`,
		id,
	).Scan(&id)
	return err
}

type User struct {
	ID   int
	Name string
}

// NewUser создает нового пользователя и возвращает его id
func (s *Storage) NewUser(u User) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
			INSERT INTO users (name)
			VALUES ($1) RETURNING id;
		`,
		u.Name,
	).Scan(&id)
	return id, err
}

type Label struct {
	ID   int
	Name string
}

// NewLabel создает новую метку и возвращает её id
func (s *Storage) NewLabel(label Label) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
			INSERT INTO labels (name)
			VALUES ($1) RETURNING id;
		`,
		label.Name,
	).Scan(&id)
	return id, err
}

type TaskLabel struct {
	TaskId  int
	LabelId int
}

// NewTaskLabel создает связь задачи и метки и возвращает её id
func (s *Storage) NewTaskLabel(taskId int, labelId int) error {
	err := s.db.QueryRow(context.Background(), `
			INSERT INTO tasks_labels (task_id, label_id)
			VALUES ($1, $2)
			RETURNING task_id;
		`,
		taskId,
		labelId,
	).Scan(&taskId)
	return err
}

func (s *Storage) TasksByLabel(label string) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
			SELECT tasks.*  FROM labels
			INNER JOIN tasks_labels ON tasks_labels.label_id = labels.id
			INNER JOIN tasks ON tasks.id = tasks_labels.task_id
			WHERE labels.name = $1
			ORDER BY tasks.id
		`,
		label,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)
	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()

}
