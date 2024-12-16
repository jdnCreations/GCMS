-- name: GetCustomerById :one
SELECT * FROM customers where id = $1;