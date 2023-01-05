package mailbox

import (
	"gomail/pkg/proto"
	"mime"
	"strings"
	"unicode/utf8"
)

func AddressStrings(as []*proto.Address) []string {
	res := make([]string, len(as))
	for i, a := range as {
		res[i] = AddressString(a)
	}
	return res
}
func AddressString(a *proto.Address) string {
	// Format address local@domain
	at := strings.LastIndex(a.Address, "@")
	var local, domain string
	if at < 0 {
		// This is a malformed address ("@" is required in addr-spec);
		// treat the whole address as local-part.
		local = a.Address
	} else {
		local, domain = a.Address[:at], a.Address[at+1:]
	}

	// Add quotes if needed
	quoteLocal := false
	for i, r := range local {
		if isAtext(r, false, false) {
			continue
		}
		if r == '.' {
			// Dots are okay if they are surrounded by atext.
			// We only need to check that the previous byte is
			// not a dot, and this isn't the end of the string.
			if i > 0 && local[i-1] != '.' && i < len(local)-1 {
				continue
			}
		}
		quoteLocal = true
		break
	}
	if quoteLocal {
		local = quoteString(local)

	}

	s := "<" + local + "@" + domain + ">"

	if a.Name == "" {
		return s
	}

	// If every character is printable ASCII, quoting is simple.
	allPrintable := true
	for _, r := range a.Name {
		// isWSP here should actually be isFWS,
		// but we don't support folding yet.
		if !isVchar(r) && !isWSP(r) || isMultibyte(r) {
			allPrintable = false
			break
		}
	}
	if allPrintable {
		return quoteString(a.Name) + " " + s
	}

	// Text in an encoded-word in a display-name must not contain certain
	// characters like quotes or parentheses (see RFC 2047 section 5.3).
	// When this is the case encode the name using base64 encoding.
	if strings.ContainsAny(a.Name, "\"#$%&'(),.:;<>@[]^`{|}~") {
		return mime.BEncoding.Encode("utf-8", a.Name) + " " + s
	}
	return mime.QEncoding.Encode("utf-8", a.Name) + " " + s
}

// isAtext reports whether r is an RFC 5322 atext character.
// If dot is true, period is included.
// If permissive is true, RFC 5322 3.2.3 specials is included,
// except '<', '>', ':' and '"'.
func isAtext(r rune, dot, permissive bool) bool {
	switch r {
	case '.':
		return dot

	// RFC 5322 3.2.3. specials
	case '(', ')', '[', ']', ';', '@', '\\', ',':
		return permissive

	case '<', '>', '"', ':':
		return false
	}
	return isVchar(r)
}

// isQtext reports whether r is an RFC 5322 qtext character.
func isQtext(r rune) bool {
	// Printable US-ASCII, excluding backslash or quote.
	if r == '\\' || r == '"' {
		return false
	}
	return isVchar(r)
}

// quoteString renders a string as an RFC 5322 quoted-string.
func quoteString(s string) string {
	var buf strings.Builder
	buf.WriteByte('"')
	for _, r := range s {
		if isQtext(r) || isWSP(r) {
			buf.WriteRune(r)
		} else if isVchar(r) {
			buf.WriteByte('\\')
			buf.WriteRune(r)
		}
	}
	buf.WriteByte('"')
	return buf.String()
}

// isVchar reports whether r is an RFC 5322 VCHAR character.
func isVchar(r rune) bool {
	// Visible (printing) characters.
	return '!' <= r && r <= '~' || isMultibyte(r)
}

// isMultibyte reports whether r is a multi-byte UTF-8 character
// as supported by RFC 6532
func isMultibyte(r rune) bool {
	return r >= utf8.RuneSelf
}

// isWSP reports whether r is a WSP (white space).
// WSP is a space or horizontal tab (RFC 5234 Appendix B).
func isWSP(r rune) bool {
	return r == ' ' || r == '\t'
}
