-- +goose Up
CREATE TABLE Todos (
  id INTEGER PRIMARY KEY,
  title TEXT NOT NULL,
  completed INTEGER NOT NULL
);

-- +goose Down
DROP TABLE Todos;
