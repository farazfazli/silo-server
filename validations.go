package main

import "regexp"

// Matches a-z A-Z 0-9 @ . and > 0 characters
// I chose this in favor of a complex email
// regex because I actually understand it :P
// `` allows only having to escape the . once
func Validate(email string) string {
	regex := regexp.MustCompile(`[^a-zA-Z0-9@\.]+`)
	return regex.ReplaceAllString(email, "")
}

// TODO: Add valid emails to redis after confirmation link is clicked
