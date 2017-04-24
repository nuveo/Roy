package system

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Exec(name string, arg ...string) (ret string, err error) {

	cmd := exec.Command(name, arg...)
	//cmd.Stdin = strings.NewReader("some input")
	var outBuff bytes.Buffer
	var errBuff bytes.Buffer
	cmd.Stdout = &outBuff
	cmd.Stderr = &errBuff
	err = cmd.Run()
	if err != nil {
		return
	}
	fmt.Println(outBuff.String())
	fmt.Println(errBuff.String())
	ret = outBuff.String()
	return
}
