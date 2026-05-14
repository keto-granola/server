-- name: InsertProduct :one
INSERT INTO products (name) VALUES ($1) RETURNING id;