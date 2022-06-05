package logger

import (
	"testing"
)

func TestRetrieveLogContent(t *testing.T) {
	var c1 logContent
	var c2 *logContent

	c2, err := RetrieveLogContent()
	if err != nil {
		t.Fatal(err)
	}

	c1 = *c2

	t.Logf("c1:%v", c1)
	t.Log("c2:", c2)
}
