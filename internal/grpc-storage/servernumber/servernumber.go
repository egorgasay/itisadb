package servernumber

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func Get(dir string) (int32, error) {
	f, err := os.Open(dir + "/server_number")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var b = make([]byte, 100)
	r, err := f.Read(b)
	if err != nil {
		return 0, err
	}

	b = b[:r]

	num, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}

	return int32(num), nil
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
