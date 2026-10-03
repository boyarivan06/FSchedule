DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS children;
DROP TABLE IF EXISTS users;migrate -path ./migrations -database "postgres://go_user:machine_banana@localhost:5432/f_schedule?sslmode=disable" up