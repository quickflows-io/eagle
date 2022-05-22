package email

// Driver Email sending driver interface definition
type Driver interface {
	// Send email
	Send(to, subject, body string) error
	// Close link
	Close()
}

// Send email
func Send(to, subject, body string) error {
	Lock.RLock()
	defer Lock.RUnlock()

	if Client == nil {
		return nil
	}

	return Client.Send(to, subject, body)
}
