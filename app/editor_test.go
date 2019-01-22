package app

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createFile(data []byte) (*os.File, error) {
	f, err := os.OpenFile(
		fmt.Sprintf("/tmp/%s", randString(32)),
		os.O_CREATE|os.O_RDWR,
		777,
	)
	if err != nil {
		return nil, err
	}
	_, err = f.Write(data)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestEditorStarts(t *testing.T) {

	f, err := createFile([]byte{0x7F, 0x45, 0x4C, 0x46})
	require.Nil(t, err)

	_, err = NewEditor(f)
	require.Nil(t, err)
}
