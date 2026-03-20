package email

import (
	"github.com/emersion/go-imap/client"
)

func ConnectIMAP(host, port, user, pass string) (*client.Client, error) {
	c, err := client.DialTLS(host+":"+port, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Login(user, pass); err != nil {
		c.Logout()
		return nil, err
	}
	return c, nil
}
