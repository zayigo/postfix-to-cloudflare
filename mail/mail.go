package mail

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strings"
)

type Contact struct {
	Email string  `json:"email"`
	Name  *string `json:"name,omitempty"`
}

type EmailData struct {
	To      []Contact `json:"to"`
	ReplyTo []Contact `json:"replyTo,omitempty"`
	Cc      []Contact `json:"cc,omitempty"`
	Bcc     []Contact `json:"bcc,omitempty"`
	From    Contact   `json:"from"`
	Subject string    `json:"subject"`
	Text    *string   `json:"text,omitempty"`
	Html    *string   `json:"html,omitempty"`
}

func ParseEmailFromStdin() (EmailData, error) {
	emailData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return EmailData{}, fmt.Errorf("error reading email from stdin: %v", err)
	}

	// Parse the email
	msg, err := mail.ReadMessage(bytes.NewReader(emailData))
	if err != nil {
		return EmailData{}, fmt.Errorf("error parsing email: %v", err)
	}

	// Process and extract email data
	emailObj, err := processEmail(msg)
	if err != nil {
		return EmailData{}, err
	}
	return emailObj, nil
}

func processEmail(msg *mail.Message) (EmailData, error) {
	header := msg.Header
	subject := header.Get("Subject")

	// Process the From field
	from := processAddressList(header.Get("From"))[0]

	// Process recipient fields
	to := processAddressList(header.Get("To"))
	if len(to) == 0 {
		return EmailData{}, fmt.Errorf("error: 'To' field is required and cannot be empty")
	}
	cc := processAddressList(header.Get("Cc"))
	bcc := processAddressList(header.Get("Bcc"))
	replyTo := processAddressList(header.Get("Reply-To"))

	// Initialize EmailData
	emailObj := EmailData{
		To:      to,
		ReplyTo: replyTo,
		Cc:      cc,
		Bcc:     bcc,
		From:    from,
		Subject: subject,
	}

	// Extract and set body (text and/or HTML)
	setEmailBody(msg, &emailObj)

	return emailObj, nil
}

func processAddressList(addressList string) []Contact {
	var contacts []Contact
	if addressList == "" {
		return contacts
	}
	addresses, err := mail.ParseAddressList(addressList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing address list: %v\n", err)
		return contacts
	}
	for _, addr := range addresses {
		contact := Contact{Email: addr.Address}
		if addr.Name != "" {
			contact.Name = &addr.Name
		}
		contacts = append(contacts, contact)
	}
	return contacts
}

func setEmailBody(msg *mail.Message, emailObj *EmailData) {
	mediaType, params, _ := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if mediaType == "" {
		mediaType = "text/plain"
	}

	// Check if the email content is multipart (e.g., both text and HTML parts)
	if strings.HasPrefix(mediaType, "multipart/") {
		// Create a new multipart reader
		mr := multipart.NewReader(msg.Body, params["boundary"])
		// Iterate through all parts of the multipart message
		for {
			p, err := mr.NextPart()
			// If we've reached the end of the parts, break out of the loop
			if err == io.EOF {
				break
			}
			// Handle any errors encountered while reading the next part
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading multipart section: %v\n", err)
				break
			}
			// Read the entire part body
			slurp, err := io.ReadAll(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading part body: %v\n", err)
				continue
			}

			// Determine the content type of the part
			partContentType := p.Header.Get("Content-Type")
			// If the part is plain text, set the email object's text field
			if strings.HasPrefix(partContentType, "text/plain") {
				text := string(slurp)
				emailObj.Text = &text
				// If the part is HTML, set the email object's HTML field
			} else if strings.HasPrefix(partContentType, "text/html") {
				html := string(slurp)
				emailObj.Html = &html
			}
		}
		// Handle non-multipart emails (either plain text or HTML)
	} else if mediaType == "text/plain" || mediaType == "text/html" {
		// Read the entire body of the email
		body, err := io.ReadAll(msg.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading body: %v\n", err)
			return
		}
		content := string(body)
		// Set the appropriate field in the email object based on the content type
		if mediaType == "text/plain" {
			emailObj.Text = &content
		} else {
			emailObj.Html = &content
		}
	}
}
