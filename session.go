package brynhild

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	// The client has connected, and is awaiting our first response
	clientGreeting = iota
	// We have responded to the client's connection and are awaiting a command
	clientCmd
	// We have received the sender and recipient information
	clientData
	// We have agreed with the client to secure the connection over TLS
	clientStartTLS
	// Server will shutdown, client to shutdown on next command turn
	clientShutdown
)

type session struct {
	sessionId  uint64
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
	scanner    *bufio.Scanner
}

type command struct {
	line   string
	action string
	fields []string
	params []string
}

// new session instance
func newSession(conn net.Conn, sessionId uint64) *session {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(conn)
	instance := &session{sessionId, conn, reader, writer, scanner}
	return instance
}

// TODO
func (session *session) upgradeToTLS(tlsConfig *tls.Config) error {
	return nil
}

// send response to client
func (session *session) sendResponse(code int, msg string) error {
	fmt.Fprintf(session.writer, "%d %s\r\n", code, msg)
	return nil
}

//TODO
func (session *session) close() {
	return
}

func (session *session) parseLine(line string) (cmd command) {

	cmd.line = line
	cmd.fields = strings.Fields(line)

	if len(cmd.fields) > 0 {
		cmd.action = strings.ToUpper(cmd.fields[0])
		if len(cmd.fields) > 1 {
			cmd.params = strings.Split(cmd.fields[1], ":")
		}
	}

	return

}

func (session *session) handle() {
	// the welcoming message
	greeting := fmt.Sprintf("%s - Session id: %d, Time: %s", s.option.Welcoming, sessionId, time.Now().Format(time.RFC3339))
	helo := fmt.Sprintf("250 %s", s.option.Hostname)
	ehlo := fmt.Sprintf("250-%s\r\n", s.option.Hostname)
	maxMsgSize := fmt.Sprintf(messageSize, s.option.MaxMessageSize)

	session.sendResponse(220, greeting)
	for {
		for session.scanner.Scan() {
			session.handleLine(session.scanner.Text())
		}

		err := session.scanner.Err()
		if err != nil {
			log.Error(err)
			continue
		}
		break
	}
}

func (session *session) handleLine(line string) {
	if session.server.ProtocolLogger != nil {
		session.server.ProtocolLogger.Printf("%s < %s", session.conn.RemoteAddr(), line)
	}
	cmd := parseLine(line)

	// Commands are dispatched to the appropriate handler functions.
	// If a network error occurs during handling, the handler should
	// just return and let the error be handled on the next read.

	switch cmd.action {

	case "HELO":
		session.handleHELO(cmd)
		return

	case "EHLO":
		session.handleEHLO(cmd)
		return

	case "MAIL":
		session.handleMAIL(cmd)
		return

	case "RCPT":
		session.handleRCPT(cmd)
		return

	case "STARTTLS":
		session.handleSTARTTLS(cmd)
		return

	case "DATA":
		session.handleDATA(cmd)
		return

	case "RSET":
		session.handleRSET(cmd)
		return

	case "NOOP":
		session.handleNOOP(cmd)
		return

	case "QUIT":
		session.handleQUIT(cmd)
		return

	case "AUTH":
		session.handleAUTH(cmd)
		return

	case "XCLIENT":
		session.handleXCLIENT(cmd)
		return

	}

	session.sendResponse("502 Unsupported command.")

}
