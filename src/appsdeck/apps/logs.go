package apps

import (
	"appsdeck/api"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

func Logs(app string, n int) error {
	nStr := strconv.FormatUint(uint64(n), 10)
	res, err := api.Logs(app, nStr)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(buffer))
	return nil
}

func LogsStream(app string) error {
	res, err := api.LogsStream(app)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	reader := bufio.NewReader(res.Body)
	for line, _, err := reader.ReadLine(); err == nil; line, _, err = reader.ReadLine() {
		if len(line) != 0 {
			fmt.Println(string(line))
		}
	}
	if err != io.EOF {
		return err
	}
	return nil
}
