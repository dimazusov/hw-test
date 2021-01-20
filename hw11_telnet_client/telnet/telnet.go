package telnet

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

type Client interface {
	Connect(ctx context.Context) error
	Close() error
	Send() error
	Receive() error
}

type tClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) Client {
	return &tClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (m *tClient) Connect(ctx context.Context) error {
	var err error

	m.conn, err = net.DialTimeout("tcp", m.address, m.timeout)
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}

	go func() {
		<-ctx.Done()
		m.conn.Close()
	}()

	return nil
}

func (m *tClient) Close() error {
	return m.conn.Close()
}

func (m *tClient) Send() error {
	r := bufio.NewReader(m.in)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return errors.Wrap(err, "cannot read")
		}

		b = append(b, '\n')
		if _, err = m.conn.Write(b); err != nil {
			return errors.Wrap(err, "cannot write")
		}
	}
}

func (m *tClient) Receive() error {
	r := bufio.NewReader(m.conn)
	for {
		b, _, err := r.ReadLine()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "cannot read")
		}

		b = append(b, '\n')
		_, err = m.out.Write(b)
		if err != nil {
			return errors.Wrap(err, "cannot write")
		}
	}
}
