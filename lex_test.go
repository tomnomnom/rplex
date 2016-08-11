package rplex

import (
	"testing"
	"unicode"
)

type testToken struct {
	TextToken
}

func TestAccept(t *testing.T) {
	l := New("abc")

	ts := l.Run(func(l *Lexer) LexFn {
		l.Accept("a")
		l.Accept("d")
		l.Accept("bc")
		l.Accept("bc")

		l.Emit(&testToken{})
		return nil
	})

	if len(ts) != 1 {
		t.Fatalf("have length %s; want 1", len(ts))
	}

	if ts[0].Text() != "abc" {
		t.Errorf("have text '%s'; want 'abc'", ts[0].Text())
	}
}

func TestAcceptRun(t *testing.T) {

	l := New("abc123")

	ts := l.Run(func(l *Lexer) LexFn {
		l.AcceptRun("abc")
		l.Emit(&testToken{})

		l.AcceptRun("1234")
		l.Emit(&testToken{})
		return nil
	})

	if len(ts) != 2 {
		t.Fatalf("have length %d; want 2", len(ts))
	}

	if ts[0].Text() != "abc" {
		t.Errorf("have text '%s'; want 'abc'", ts[0].Text())
	}

	if ts[1].Text() != "123" {
		t.Errorf("have text '%s'; want '123'", ts[1].Text())
	}
}

func TestPeek(t *testing.T) {
	l := New("abc")

	ts := l.Run(func(l *Lexer) LexFn {
		l.Accept("a")
		if l.Peek() == 'b' {
			l.Accept("b")
		}

		l.Emit(&testToken{})

		return nil
	})

	if len(ts) != 1 {
		t.Fatalf("have length %s; want 1", len(ts))
	}

	if ts[0].Text() != "ab" {
		t.Errorf("have text '%s'; want 'ab'", ts[0].Text())
	}
}

func TestIgnore(t *testing.T) {
	l := New("abc")

	ts := l.Run(func(l *Lexer) LexFn {
		l.Accept("a")
		l.Ignore()
		l.AcceptRun("bc")

		l.Emit(&testToken{})

		return nil
	})

	if len(ts) != 1 {
		t.Fatalf("have length %s; want 1", len(ts))
	}

	if ts[0].Text() != "bc" {
		t.Errorf("have text '%s'; want 'bc'", ts[0].Text())
	}
}

func TestAcceptFunc(t *testing.T) {
	l := New("abc")

	ts := l.Run(func(l *Lexer) LexFn {
		l.AcceptFunc(func(r rune) bool {
			return r == 'a'
		})

		l.AcceptFunc(func(r rune) bool {
			return r == 'a'
		})

		l.Emit(&testToken{})

		return nil
	})

	if len(ts) != 1 {
		t.Fatalf("have length %s; want 1", len(ts))
	}

	if ts[0].Text() != "a" {
		t.Errorf("have text '%s'; want 'a'", ts[0].Text())
	}
}

func TestAcceptRunFunc(t *testing.T) {
	l := New("123a")

	ts := l.Run(func(l *Lexer) LexFn {
		l.AcceptRunFunc(unicode.IsNumber)

		l.Emit(&testToken{})

		return nil
	})

	if len(ts) != 1 {
		t.Fatalf("have length %s; want 1", len(ts))
	}

	if ts[0].Text() != "123" {
		t.Errorf("have text '%s'; want '123'", ts[0].Text())
	}
}

func TestAcceptUntil(t *testing.T) {
	l := New("123abc")

	ts := l.Run(func(l *Lexer) LexFn {
		l.AcceptUntil("a")
		l.Emit(&testToken{})

		l.AcceptUntil("z")
		l.Emit(&testToken{})

		return nil
	})

	if len(ts) != 2 {
		t.Fatalf("have length %s; want 2", len(ts))
	}

	if ts[0].Text() != "123" {
		t.Errorf("have text '%s'; want '123'", ts[0].Text())
	}

	if ts[1].Text() != "abc" {
		t.Errorf("have text '%s'; want 'abc'", ts[0].Text())
	}
}

func TestAcceptUntilUnescaped(t *testing.T) {
	l := New(`123\"abc"def`)

	ts := l.Run(func(l *Lexer) LexFn {
		l.AcceptUntilUnescaped(`"`)
		l.Emit(&testToken{})

		l.AcceptUntilUnescaped("z")
		l.Emit(&testToken{})

		return nil
	})

	if len(ts) != 2 {
		t.Fatalf("have length %s; want 2", len(ts))
	}

	if ts[0].Text() != `123\"abc` {
		t.Errorf(`have text '%s'; want '123\"abc'`, ts[0].Text())
	}

	if ts[1].Text() != `"def` {
		t.Errorf(`have text '%s'; want '"def'`, ts[0].Text())
	}

}
