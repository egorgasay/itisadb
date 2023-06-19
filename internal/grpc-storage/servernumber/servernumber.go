package servernumber

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

var ErrNotFound = errors.New("server_number was not found")

func Get(dir string) int32 {
	f, err := os.Open(dir + "/server_number")
	if err != nil {
		log.Println(err)
		return 0
	}
	defer f.Close()

	var b = make([]byte, 100)
	r, err := f.Read(b)
	if err != nil {
		log.Println(err)
		return 0
	}

	b = b[:r]

	num, err := strconv.Atoi(string(b))
	if err != nil {
		log.Println(err)
		return 0
	}

	return int32(num)
}

func Set(server int32) error {
	f, err := os.OpenFile("server_number", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprint(server))
	if err != nil {
		return fmt.Errorf("can't wrtite to file server_number: %w", err)
	}

	return f.Sync()
}
