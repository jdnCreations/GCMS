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

-- name: UpdateCustomer :one
UPDATE customers 
SET first_name = COALESCE(NULLIF($1, ''), first_name),
    last_name = COALESCE(NULLIF($2, ''), last_name),
    email = COALESCE(NULLIF($3, ''), email)
WHERE id = $4
RETURNING id, first_name, last_name, email;