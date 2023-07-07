package email

import (
	"testing"

	"github.com/binbomb/goapp/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := " A test Email sender "
	content := `
		<h1> Hello worlds </h1>
		<p>
			just test
		</p>
	`
	to := []string{"tonghoameo@gmail.com"}
	attachFiles := []string{"../simple_readme.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
