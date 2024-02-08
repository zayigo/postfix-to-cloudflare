package tests

import (
	"os"
	"postfix_to_cf/mail"
	"reflect"
	"testing"
)

func TestParseEmailFromStdin(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		file     string
		wantErr  bool
		expected mail.EmailData
	}{
		{
			name:    "PlainTextEmail",
			file:    "./samples/plain_text.eml",
			wantErr: false,
			expected: mail.EmailData{
				To:      []mail.Contact{{Email: "jane.smith@example.com", Name: ptrToString("Jane Smith")}},
				From:    mail.Contact{Email: "john.doe@example.com", Name: ptrToString("John Doe")},
				Subject: "Test Email 1",
				Text:    ptrToString("This is a basic test email."),
			},
		},
		{
			name:    "HTMLEmail",
			file:    "./samples/html_content.eml",
			wantErr: false,
			expected: mail.EmailData{
				To:      []mail.Contact{{Email: "tester@example.com", Name: ptrToString("Tester")}},
				From:    mail.Contact{Email: "developer@example.com", Name: ptrToString("Developer")},
				Subject: "Test Email 3",
				Html:    ptrToString("<html>\n<body>\n<h1>This is an HTML Email</h1>\n<p>This is a test email with <b>HTML content</b>.</p>\n</body>\n</html>"),
			},
		},
		{
			name:    "MultipleRecipientsEmail",
			file:    "./samples/multiple_recipients.eml",
			wantErr: false,
			expected: mail.EmailData{
				To: []mail.Contact{
					{Email: "bob@example.com", Name: ptrToString("Bob Builder")},
					{Email: "charlie@example.com", Name: ptrToString("Charlie Day")},
				},
				Cc: []mail.Contact{
					{Email: "eve@example.com", Name: ptrToString("Eve")},
				},
				From:    mail.Contact{Email: "alice@example.com", Name: ptrToString("Alice Wonderland")},
				Subject: "Test Email 2",
				Text:    ptrToString("Testing email with multiple recipients and CC."),
			},
		},
		{
			name:    "ReplyToEmail",
			file:    "./samples/reply_to.eml",
			wantErr: false,
			expected: mail.EmailData{
				To: []mail.Contact{
					{Email: "service@example.com", Name: ptrToString("Service")},
				},
				ReplyTo: []mail.Contact{
					{Email: "feedback@example.com"},
				},
				From:    mail.Contact{Email: "client@example.com", Name: ptrToString("Client")},
				Subject: "Your Inquiry",
				Text:    ptrToString("Please send all replies to feedback@example.com."),
			},
		},
		{
			name:    "AttachmentEmail",
			file:    "./samples/attachment.eml",
			wantErr: false,
			expected: mail.EmailData{
				To: []mail.Contact{
					{Email: "user@example.com", Name: ptrToString("User")},
				},
				From:    mail.Contact{Email: "admin@example.com", Name: ptrToString("Admin")},
				Subject: "Test Email 4",
				Text:    ptrToString("This email contains an attachment.\n"),
			},
		},
		{
			name:    "MissingTo",
			file:    "./samples/missing_to.eml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Open the sample .eml file
			file, err := os.Open(tt.file)
			if err != nil {
				t.Fatalf("Failed to open file: %s, error: %v", tt.file, err)
			}
			defer file.Close()

			// Redirect Stdin to read from the .eml file
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = file

			// Call the function under test
			got, err := mail.ParseEmailFromStdin()

			// Check for expected error state
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEmailFromStdin() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Compare the expected and actual results
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseEmailFromStdin() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func ptrToString(s string) *string {
	return &s
}
