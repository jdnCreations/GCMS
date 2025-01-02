-- name: GetCustomerById :one
SELECT * FROM customers where id = $1;

-- name: CreateCustomer :one
INSERT INTO customers (id, first_name, last_name, email)
VALUES (
  gen_random_uuid(),
  $1,
  $2,
  $3
)
RETURNING *;

-- name: DeleteCustomerById :exec
DELETE FROM customers where id = $1;

-- name: GetAllCustomers :many
SELECT * FROM customers ORDER BY first_name;