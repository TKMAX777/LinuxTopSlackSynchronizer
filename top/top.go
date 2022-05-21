package top

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

type Process struct {
	PID     int
	User    string
	Pr      int
	Ni      int
	Virt    string
	Res     string
	SHR     int
	S       string
	CPU     float32
	Memory  float32
	Time    string
	Command string
}

var commandRegexp = regexp.MustCompile(`\s*(\d+)\s+(\S+)\s+(\d+)\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\S+)\s+(\d+\.*\d*)\s+(\d+\.*\d*)\s+(\d+:\d+\.\d+)\s+(.+)\s+`)

func Get() ([]Process, error) {
	cmd := exec.Command("/usr/bin/top", "-bcn", "1", "-w")

	output, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "Output")
	}

	fmt.Printf("%s\n", output)

	var processes = []Process{}

	for _, sm := range commandRegexp.FindAllStringSubmatch(string(output), -1) {
		PID, _ := strconv.Atoi(sm[1])
		Pr, _ := strconv.Atoi(sm[3])
		Ni, _ := strconv.Atoi(sm[4])
		SHR, _ := strconv.Atoi(sm[7])

		CPU, _ := strconv.ParseFloat(sm[9], 32)
		Memory, _ := strconv.ParseFloat(sm[10], 32)

		var process = Process{
			PID:     PID,
			User:    sm[2],
			Pr:      Pr,
			Ni:      Ni,
			Virt:    sm[5],
			Res:     sm[6],
			SHR:     SHR,
			S:       sm[8],
			CPU:     float32(CPU),
			Memory:  float32(Memory),
			Time:    sm[11],
			Command: sm[12],
		}

		processes = append(processes, process)
	}

	return processes, nil
}
