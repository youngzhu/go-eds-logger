package reportor

import "testing"

func Test_login(t *testing.T) {
	err := login()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("login success")
}
