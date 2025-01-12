package auth

import "testing"

func TestPassword(t *testing.T) {
	pass := "jbord"

	encrytped, err := HashPassword(pass)
	if err != nil {
		t.Errorf("could not hash password")
	}

	t.Logf("encrypted pw: %v", encrytped)

	err = CheckPasswordHash(pass, encrytped)
	if err != nil {
		t.Errorf("invalid password")
	}
}