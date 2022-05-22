package email

import (
	"errors"
	"sync"

	"github.com/go-eagle/eagle/pkg/log"
)

// Client Email sending client
var Client Driver

// Lock Read-write lock
var Lock sync.RWMutex

var (
	// ErrChanNotOpen Mail queue is not open
	ErrChanNotOpen = errors.New("email queue does not open")
)

// Config email config
type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Name      string
	Address   string
	ReplyTo   string
	KeepAlive int
}

// Init Initialize the client
func Init(cfg Config) {
	log.Info("email init")
	Lock.Lock()
	defer Lock.Unlock()

	// make sure it's closed
	if Client != nil {
		Client.Close()
	}

	client := NewSMTPClient(SMTPConfig{
		Name:      cfg.Name,
		Address:   cfg.Address,
		ReplyTo:   cfg.ReplyTo,
		Host:      cfg.Host,
		Port:      cfg.Port,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Keepalive: cfg.KeepAlive,
	})

	Client = client
}
