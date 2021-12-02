package stail

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
)

type sTail struct {
	os  string
	opt Options
}

func New(opt Options) (st *sTail, err error) {
	st = &sTail{os: runtime.GOOS, opt: opt}
	var cp string
	switch st.os {
	case OsWindows:
		if st.opt.PowerShellPath == "" {
			if cp, err = st.lookPath("powershell.exe", "pwsh.exe"); err != nil {
				return
			} else {
				st.opt.PowerShellPath = cp
			}
		}
	case OsLinux, OsDarwin:
		if st.opt.UnixTailPath == "" {
			if cp, err = st.lookPath("tail"); err != nil {
				return
			} else {
				st.opt.UnixTailPath = cp
			}
		}
	}
	return
}

func (s *sTail) Tail(filepath string, tailLine int, call func(content string)) (err error) {
	cmd, err := s.getCommand(filepath, tailLine)
	if err != nil {
		return
	}
	var stdout io.ReadCloser
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		err = errors.New("get system pipe failed")
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}
	reader := bufio.NewReader(stdout)
	for {
		line, errRs := reader.ReadString('\n')
		if errRs != nil || io.EOF == errRs {
			break
		}
		if call != nil {
			call(line)
		}
	}
	return
}

func (s *sTail) Total(filepath string, call func(content string)) (err error) {
	return s.Tail(filepath, -1, call)
}

func (s *sTail) getCommand(filepath string, tailLine int) (cmd *exec.Cmd, err error) {
	switch s.os {
	case OsWindows:
		cmd = s.windowsTail(filepath, tailLine)
	case OsLinux, OsDarwin:
		cmd = s.linuxTail(filepath, tailLine)
	default:
		err = errors.New("not supported on the current platform")
	}
	return
}

func (s *sTail) lookPath(filenames ...string) (findPath string, err error) {
	for _, n := range filenames {
		findPath, _ = exec.LookPath(n)
		if findPath != "" {
			return
		}
	}
	err = errors.New(fmt.Sprintf("not find files: %v", filenames))
	return
}

func (s *sTail) windowsTail(filepath string, tailLine int) *exec.Cmd {
	return exec.Command(s.opt.PowerShellPath, "-Command", "Get-Content", "-Path", filepath,
		"-Tail", strconv.Itoa(tailLine), "-Wait")
}

func (s *sTail) linuxTail(filepath string, tailLine int) *exec.Cmd {
	params := make([]string, 0, 3)
	if tailLine >= 0 {
		params = append(params, fmt.Sprintf("-%df", tailLine))
	}
	params = append(params, filepath)
	return exec.Command(s.opt.UnixTailPath, params...)
}
