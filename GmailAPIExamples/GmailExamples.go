package main

import (
	"./GmailCredentialManager"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"strings"
)

func main() {

	client := GmailCredentialManager.GetService()
	_ = GetAllMessagesFromUser(client, "", 4)

}

func GetLastNMessages(srv *gmail.Service, n int) []*gmail.Message {
	s := make([]*gmail.Message, n)
	r, e := srv.Users.Messages.List("me").Do()
	c(e)
	for i, e := range r.Messages {
		if i == n {
			break
		}
		t, x := srv.Users.Messages.Get("me", e.Id).Do()
		c(x)
		s[i] = t
	}

	return s
}

func GetLastNMessagesAsString(srv *gmail.Service, n int) []string {
	msgs := GetLastNMessages(srv, n)
	s := make([]string, n)

	for i, e := range msgs {
		parts := e.Payload.Parts

		t := make([]string, n)

		for _, k := range parts {
			b, err := base64.URLEncoding.DecodeString(k.Body.Data)
			c(err)

			t = append(t, string(b))
		}

		s[i] = strings.Join(t, "\n")
	}
	return s
}

func GetMessageBodyAsString(msg *gmail.Message) string {

	parts := msg.Payload.Parts
	t := make([]string, len(msg.Payload.Parts))

	for _, k := range parts {
		b, err := base64.URLEncoding.DecodeString(k.Body.Data)
		c(err)

		t = append(t, string(b))
	}

	return strings.Join(t, "\n")
}

func GetAllMessagesFromUser(srv *gmail.Service, sender string, depth int) []*gmail.Message {
	msgs := GetLastNMessages(srv, depth)
	t := make([]*gmail.Message, depth)
	tmp := 0
	for i, e := range msgs {
		if e.Payload.Headers[20].Value == sender {
			t = append(t, e)
		}
		fmt.Print(i)
		tmp = i
	}

	if len(t) != depth {
		x := make([]*gmail.Message, tmp)
		for _, e := range t {
			if e == nil {
				return x
			}
			x = append(x, e)
		}
		return x
	}
	return t
}

func GetAllMessagesFromUserAsStrings(client *gmail.Service, sender string, depth int) []string {
	msgs := GetAllMessagesFromUser(client, sender, depth)
	s := make([]string, len(msgs))
	for _, e := range msgs {
		s = append(s, GetMessageBodyAsString(e))
	}
	return s
}

//c checks that an error has not been thrown.
func c(e error) {
	if e != nil {
		fmt.Println(e)
		panic(9)
	}
}
