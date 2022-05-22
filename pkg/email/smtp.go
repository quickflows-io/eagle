package email

import (
	"time"

	"github.com/go-mail/mail"

	"github.com/go-eagle/eagle/pkg/log"
)

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Name      string // sender name
	Address   string // sender address
	ReplyTo   string // Reply address
	Host      string // Server hostname
	Port      int    // server port
	Username  string // username
	Password  string // password
	Keepalive int    // Connection keep-alive duration, in seconds
}

// SMTP protocol end
type SMTP struct {
	Config SMTPConfig
	ch     chan *mail.Message
	chOpen bool
}

// NewSMTPClient Instantiate an SMTP client
func NewSMTPClient(config SMTPConfig) *SMTP {
	client := &SMTP{
		Config: config,
		ch:     make(chan *mail.Message, 30),
		chOpen: false,
	}

	return client
}

// Init Initialize the send queue
func (c *SMTP) Init() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.chOpen = false
				log.Error("send email queue err: %+v, retry after 10 second", err.(error))
				time.Sleep(time.Duration(10) * time.Second)
				c.Init()
			}
		}()

		d := mail.NewDialer(c.Config.Host, c.Config.Port, c.Config.Username, c.Config.Password)
		d.Timeout = time.Duration(c.Config.Keepalive+5) * time.Second
		c.chOpen = true

		var s mail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-c.ch:
				if !ok {
					log.Info("mail queue is close")
					c.chOpen = false
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err := mail.Send(s, m); err != nil {
					log.Warnf("email send failed, %v", err)
				} else {
					log.Info("email has send")
				}
			case <-time.After(time.Duration(c.Config.Keepalive) * time.Second):
				if open {
					if err := s.Close(); err != nil {
						log.Warnf("can not close smtp conn, %v", err)
					}
				}
				open = false
			}
		}
	}()
}

// Send email
func (c *SMTP) Send(to, subject, body string) error {
	if !c.chOpen {
		return ErrChanNotOpen
	}

	msg := mail.NewMessage()
	msg.SetAddressHeader("From", c.Config.Address, c.Config.Name)
	msg.SetAddressHeader("Reply-To", c.Config.ReplyTo, c.Config.Name)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("To", to)
	msg.SetBody("text/html", body)

	c.ch <- msg
	return nil
}

// Close queue
func (c *SMTP) Close() {
	if c.ch != nil {
		close(c.ch)
	}
}
