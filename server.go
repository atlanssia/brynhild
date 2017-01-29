package brynhild

import (
	"crypto/tls"
	"net/smtp"
	"net"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"time"
)

type server struct {
	conf *Conf
	relay []string
	tlsConfig *tls.Config
}

const (
	// format greeting msg and sessionId , time.Now()
	welcomeMsg string = "220 %s SessionId:%d Time:%s"

	// format MaxMessageSize
	messageSize string = "250-SIZE %d\r\n"
	advStartTLS string = "250-STARTTLS"
	pipelining string = "250-PIPELINING\r\n"
	advEnhancedStatusCodes string = "250-ENHANCEDSTATUSCODES\r\n"
)

// new server instance
func NewServer(conf *Conf) (*server, error) {
	return &server{conf, nil, nil}, nil
}

// start a server
func (s *server) Start() error {
	listener, err := net.Listen("tcp", s.conf.ListenInterface)
	if err != nil {
		log.Panic(err)
	}

	err = s.configTLS()
	if err != nil {
		// TODO disable TLS support
		log.Error(err)
	}

	var sessionId uint64
	sessionId = 0
	for {
		conn, err := listener.Accept()
		sessionId++
		if err != nil {
			log.Error(err)
			continue
		}

		// handle a connection
		go s.handleSession(conn, sessionId)
	}
	return nil
}

// TODO signal to shutdown
func (s *server) Shutdown() error {
	return nil
}

// handle a connection
func (s *server) handleSession(conn net.Conn, sessionId uint64) {
	defer conn.Close()

	// the welcoming message
	greeting := fmt.Sprintf(welcomeMsg, s.conf.Welcoming, sessionId, time.Now().Format(time.RFC3339))
	helo := fmt.Sprintf("250 %s Hello", s.conf.Hostname)
	ehlo := fmt.Sprintf("250-%s Hello\r\n", s.conf.Hostname)
	maxMsgSize := fmt.Sprintf(messageSize, s.conf.MaxMessageSize)
}

func (s *server) configTLS() error {
	cert, err := tls.LoadX509KeyPair(s.conf.PublicKeyFile, s.conf.PrivateKeyFile)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ServerName:   s.conf.Hostname,
	}

	s.tlsConfig = tlsConfig
	return nil
}

func (s *server) sendMail(addr string, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello("xxxxxxxxxx.com"); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{InsecureSkipVerify: true}

		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}