package rplex

import (
	"strings"
	"unicode/utf8"
)

// A lexer holds the state for lexing statements
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

type TextToken struct {
	text string
}

func (t *TextToken) SetText(text string) {
	t.text = text
}

func (t *TextToken) Text() string {
	return t.text
}

// A LexFn accepts a pointer to a lexer,
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

// next gets the next rune in the input and updates the lexer state
func (l *Lexer) Next() rune {
	r, w := utf8.DecodeRuneInString(l.Text[l.Pos:])

	l.Pos += w
	l.Width = w

	l.Prev = l.Cur
	l.Cur = r

	return r
}

// backup moves the lexer back one rune
// can only be used once per call of next()
func (l *Lexer) Backup() {
	l.Pos -= l.Width
}

// peek returns the next rune in the input
// without moving the internal pointer
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// ignore skips the current token
func (l *Lexer) Ignore() {
	l.TokenStart = l.Pos
}

// emit adds the current token to the token slice and
// moves the tokenStart pointer to the current position
func (l *Lexer) Emit(t Token) {
	t.SetText(l.Text[l.TokenStart:l.Pos])
	l.TokenStart = l.Pos

	l.Tokens = append(l.Tokens, t)
}

// accept moves the pointer if the next rune is in
// the set of valid runes
func (l *Lexer) Accept(valid string) bool {
	if strings.ContainsRune(valid, l.Next()) {
		return true
	}
	l.Backup()
	return false
}

// acceptRun continually accepts runes from the
// set of valid runes
func (l *Lexer) AcceptRun(valid string) {
	for strings.ContainsRune(valid, l.Next()) {
	}
	l.Backup()
}

// a runeCheck is a function that determines if a rune is valid
// or not so that we can do complex checks against runes
type RuneCheck func(rune) bool

// acceptFunc accepts a rune if the provided runeCheck
// function returns true
func (l *Lexer) AcceptFunc(fn RuneCheck) {
	if fn(l.Next()) {
		return
	}
	l.Backup()
}

// acceptRunFunc continually accepts runes for as long
// as the runeCheck function returns true
func (l *Lexer) AcceptRunFunc(fn RuneCheck) {
	for fn(l.Next()) {
	}
	l.Backup()
}

// acceptUntil accepts runes until it hits a delimiter
// rune contained in the provided string
func (l *Lexer) AcceptUntil(delims string) {
	for !strings.ContainsRune(delims, l.Next()) {
		if l.Cur == utf8.RuneError {
			return
		}
	}
	l.Backup()
}

// acceptUntilUnescaped accepts runes until it hits a delimiter
// rune contained in the provided string, unless that rune was
// escaped with a backslash
func (l *Lexer) AcceptUntilUnescaped(delims string) {
	// Read until we hit an unescaped rune or the end of the input
	for {
		if strings.ContainsRune(delims, l.Next()) && l.Prev != '\\' {
			l.Backup()
			return
		}
		if l.Cur == utf8.RuneError {
			return
		}
	}
}
