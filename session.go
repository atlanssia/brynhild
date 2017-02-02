package brynhild

import (
	"bufio"
	"net"
)

type session struct {
	sessionId uint64
	connection net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	scanner *bufio.Scanner
}

// new session instance
func newSession(conn net.Conn, sessionId uint64) *session {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(conn)
	instance := &session{sessionId, conn, reader, writer, scanner}
	return instance
}

// send response to client
func (session *session) reply(msg string) error {
	session.writer.WriteString(msg)
	return nil
}