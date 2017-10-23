package dsl

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"

	lua "github.com/yuin/gopher-lua"
)

type tac struct {
	filename string
	fd       *os.File
	scanner  *TacScanner
}

func (t *tac) open() error {
	fd, err := os.Open(t.filename)
	if err != nil {
		return err
	}
	t.fd = fd
	t.scanner = newTacScanner(fd)
	return nil
}

func (d *dslConfig) dslTacOpen(L *lua.LState) int {
	t := &tac{filename: L.CheckString(1)}
	if err := t.open(); err != nil {
		L.RaiseError("open error: %s", err.Error())
		return 0
	}
	log.Printf("[INFO] Create new tac scanner for file `%s`\n", t.filename)
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable("tac"))
	L.Push(ud)
	return 1
}

// получение tac из lua-state
func checkTac(L *lua.LState) *tac {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*tac); ok {
		return v
	}
	L.ArgError(1, "tac expected")
	return nil
}

// получение строки
func (c *dslConfig) dslTacLine(L *lua.LState) int {
	t := checkTac(L)
	if t.scanner == nil {
		L.RaiseError("tac not initialized")
		return 0
	}
	if t.scanner.Scan() {
		text := t.scanner.Text()
		L.Push(lua.LString(text))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

// остановка
func (c *dslConfig) dslTacClose(L *lua.LState) int {
	t := checkTac(L)
	if t.fd == nil {
		L.RaiseError("tac not initialized")
		return 0
	}
	t.fd.Close()
	log.Printf("[INFO] Close tac scanner for file `%s`\n", t.filename)
	t = nil
	return 0
}

const maxBufSize = 8 * 1024

type TacScanner struct {
	r            io.ReadSeeker
	split        bufio.SplitFunc
	buf          []byte
	offset       int64
	atEOF        bool
	tokens       [][]byte
	partialToken int
	err          error
}

func newTacScanner(r io.ReadSeeker) *TacScanner {
	b := &TacScanner{
		r:     r,
		buf:   make([]byte, 4096),
		atEOF: true,
		split: bufio.ScanLines,
	}
	b.offset, b.err = r.Seek(0, 2)
	return b
}

func (b *TacScanner) fillbuf() error {
	b.tokens = b.tokens[:0]
	if b.offset == 0 {
		return io.EOF
	}
	space := len(b.buf) - b.partialToken
	if space == 0 {
		if len(b.buf) >= maxBufSize {
			return errors.New("token too long")
		}
		n := len(b.buf) * 2
		if n > maxBufSize {
			n = maxBufSize
		}
		newBuf := make([]byte, n)
		copy(newBuf, b.buf[0:b.partialToken])
		b.buf = newBuf
		space = len(b.buf) - b.partialToken
	}
	if int64(space) > b.offset {
		b.buf = b.buf[0 : b.partialToken+int(b.offset)]
		space = len(b.buf) - b.partialToken
	}
	newOffset := b.offset - int64(space)
	copy(b.buf[space:], b.buf[0:b.partialToken])
	_, err := b.r.Seek(newOffset, 0)
	if err != nil {
		return err
	}
	b.offset = newOffset
	if _, err := io.ReadFull(b.r, b.buf[0:space]); err != nil {
		return err
	}
	if b.offset > 0 {
		advance, _, err := b.split(b.buf, b.atEOF)
		if err != nil {
			return err
		}
		b.partialToken = advance
		if advance == 0 || advance == len(b.buf) {
			return b.fillbuf()
		}
	} else {
		b.partialToken = 0
	}
	for i := b.partialToken; i < len(b.buf); {
		advance, token, err := b.split(b.buf[i:], b.atEOF)
		if err != nil {
			b.tokens = b.tokens[:0]
			return err
		}
		if advance == 0 {
			break
		}
		b.tokens = append(b.tokens, token)
		i += advance
	}
	b.atEOF = false
	if len(b.tokens) == 0 {
		return b.fillbuf()
	}
	return nil
}

func (b *TacScanner) Scan() bool {
	if len(b.tokens) > 0 {
		b.tokens = b.tokens[0 : len(b.tokens)-1]
	}
	if len(b.tokens) > 0 {
		return true
	}
	if b.err != nil {
		return false
	}
	b.err = b.fillbuf()
	return len(b.tokens) > 0
}

func (b *TacScanner) Split(split bufio.SplitFunc) {
	b.split = split
}

func (b *TacScanner) Bytes() []byte {
	return b.tokens[len(b.tokens)-1]
}

func (b *TacScanner) Text() string {
	return string(b.Bytes())
}

func (b *TacScanner) Err() error {
	if len(b.tokens) > 0 {
		return nil
	}
	if b.err == io.EOF {
		return nil
	}
	return b.err
}
