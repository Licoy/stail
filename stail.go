package stail

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strconv"
)

type STail interface {
	// Tail 在{filepath}文件中从末尾的{tailLine}行开始获取数据
	Tail(filepath string, tailLine int, call func(content string)) (item STailItem, err error)
	// TailTotal 在{filepath}文件从开头获取数据
	TailTotal(filepath string, call func(content string)) (item STailItem, err error)
}

type STailItem interface {
	Watch()
	Close() (err error)
}

type sTail struct {
	os  string
	opt Options
}

type sTailItem struct {
	stdout      io.ReadCloser
	callback    func(string)
	ctxCancelFn context.CancelFunc
	closed      bool
}

func New(opt Options) (res STail, err error) {
	st := &sTail{os: runtime.GOOS, opt: opt}
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
	res = st
	return
}

func (s *sTail) Tail(filepath string, tailLine int, callback func(string)) (item STailItem, err error) {
	ctx, cancelFn := context.WithCancel(context.Background())
	cmd, err := s.getCommand(ctx, filepath, tailLine)
	if err != nil {
		cancelFn()
		return
	}
	var stdout io.ReadCloser
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		err = errors.New("get system pipe failed")
		cancelFn()
		return
	}
	err = cmd.Start()
	if err != nil {
		cancelFn()
		return
	}
	item = &sTailItem{stdout: stdout, callback: callback, ctxCancelFn: cancelFn}
	return
}

func (s *sTail) TailTotal(filepath string, callback func(string)) (item STailItem, err error) {
	return s.Tail(filepath, -1, callback)
}

func (s *sTail) getCommand(ctx context.Context, filepath string, tailLine int) (cmd *exec.Cmd, err error) {
	switch s.os {
	case OsWindows:
		cmd = s.windowsTail(ctx, filepath, tailLine)
	case OsLinux, OsDarwin:
		cmd = s.linuxTail(ctx, filepath, tailLine)
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

func (s *sTail) windowsTail(ctx context.Context, filepath string, tailLine int) *exec.Cmd {
	return exec.CommandContext(ctx, s.opt.PowerShellPath, "-Command", "Get-Content", "-Path", filepath,
		"-Tail", strconv.Itoa(tailLine), "-Wait")
}

func (s *sTail) linuxTail(ctx context.Context, filepath string, tailLine int) *exec.Cmd {
	params := make([]string, 0, 3)
	if tailLine >= 0 {
		params = append(params, fmt.Sprintf("-%df", tailLine))
	}
	params = append(params, filepath)
	return exec.CommandContext(ctx, s.opt.UnixTailPath, params...)
}

func (s *sTailItem) Watch() {
	reader := bufio.NewReader(s.stdout)
	for {
		line, errRs := reader.ReadString('\n')
		if errRs != nil || io.EOF == errRs {
			return
		}
		if s.callback != nil {
			s.callback(line)
		}
	}
}

func (s *sTailItem) Close() (err error) {
	if s.closed {
		return
	}
	if s.ctxCancelFn != nil {
		s.ctxCancelFn()
	}
	err = s.stdout.Close()
	if err != nil {
		return
	}
	s.closed = true
	return
}
