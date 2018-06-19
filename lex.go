package rplex

import (
	"strings"
	"unicode/utf8"
)

// Lexer holds the state for lexing statements
type Lexer struct {
	Text       string  // The raw input text
	Pos        int     // The current byte offset in the text
	Width      int     // The width of the current rune in bytes
	Cur        rune    // The rune at the current position
	Prev       rune    // The rune at the previous position
	Tokens     []Token // The tokens that have been emitted
	TokenStart int     // The starting position of the current token
}

// A Token is a chunk of text
type Token interface {
	SetText(string)
	Text() string
}

// TextToken is a generic token that can be easily embedded into
// custom token types to meet the Token interface
type TextToken struct {
	text string
}

// SetText sets the text value of a TextToken
func (t *TextToken) SetText(text string) {
	t.text = text
}

// Text gets the text value of a TextToken
func (t *TextToken) Text() string {
	return t.text
}

// A LexFn does the meat of the work. It accepts a pointer
// to a Lexer, manipulates its state in some way, e.g. accepts
// runes and emits tokens, and then returns a new LexFn
// to deal with the next stage of the lexing - or nil if
// no lexing is left to be done.
type LexFn func(*Lexer) LexFn

// New returns a new Lexer for the provided input string
func New(text string) *Lexer {
	return &Lexer{
		Text:       text,
		Pos:        0,
		TokenStart: 0,
		Tokens:     make([]Token, 0),
	}
}

// Run runs the lexer and returns the lexed tokens
func (l *Lexer) Run(initial LexFn) []Token {

	for lexfn := initial; lexfn != nil; {
		lexfn = lexfn(l)
	}
	return l.Tokens
}

// Next gets the next rune in the input and updates the lexer state
func (l *Lexer) Next() rune {
	r, w := utf8.DecodeRuneInString(l.Text[l.Pos:])

	l.Pos += w
	l.Width = w

	l.Prev = l.Cur
	l.Cur = r

	return r
}

// Backup moves the lexer back one rune
// can only be used once per call of next()
func (l *Lexer) Backup() {
	l.Pos -= l.Width
}

// Peek returns the next rune in the input
// without moving the internal pointer
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// Ignore skips the current token
func (l *Lexer) Ignore() {
	l.TokenStart = l.Pos
}

// Emit adds the current token to the token slice and
// moves the tokenStart pointer to the current position
func (l *Lexer) Emit(t Token) {
	t.SetText(l.Text[l.TokenStart:l.Pos])
	l.TokenStart = l.Pos

	l.Tokens = append(l.Tokens, t)
}

// Accept moves the pointer if the next rune is in
// the set of valid runes
func (l *Lexer) Accept(valid string) bool {
	if strings.ContainsRune(valid, l.Next()) {
		return true
	}
	l.Backup()
	return false
}

// AcceptRun continually accepts runes from the
// set of valid runes
func (l *Lexer) AcceptRun(valid string) {
	for strings.ContainsRune(valid, l.Next()) {
	}
	l.Backup()
}

// RuneCheck is a function that determines if a rune is valid
// or not when using AcceptFunc or AcceptRunFunc. Some functions
// in the standard library, such as unicode.IsNumber() meet
// this interface already.
type RuneCheck func(rune) bool

// AcceptFunc accepts a rune if the provided runeCheck
// function returns true
func (l *Lexer) AcceptFunc(fn RuneCheck) {
	if fn(l.Next()) {
		return
	}
	l.Backup()
}

// AcceptRunFunc continually accepts runes for as long
// as the runeCheck function returns true
func (l *Lexer) AcceptRunFunc(fn RuneCheck) {
	for fn(l.Next()) {
	}
	l.Backup()
}

// AcceptUntil accepts runes until it hits a delimiter
// rune contained in the provided string
func (l *Lexer) AcceptUntil(delims string) {
	for !strings.ContainsRune(delims, l.Next()) {
		if l.Cur == utf8.RuneError {
			return
		}
	}
	l.Backup()
}

// AcceptUntilUnescaped accepts runes until it hits a delimiter
// rune contained in the provided string, unless that rune was
// escaped with a backslash
func (l *Lexer) AcceptUntilUnescaped(delims string) {

	// Read until we hit an unescaped rune or the end of the input
	inEscape := false
	for {
		r := l.Next()
		if r == '\\' && !inEscape {
			inEscape = true
			continue
		}
		if strings.ContainsRune(delims, r) && !inEscape {
			l.Backup()
			return
		}
		if l.Cur == utf8.RuneError {
			return
		}
		inEscape = false
	}
}
