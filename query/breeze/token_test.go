package breeze_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/utagai/look/query/breeze"
)

// This is the expected number of 'custom' Breeze tokens (aka, tokens that are
// not mapped to the ones found in the scanner package).
// Note that this should always match the length of the below map.
const expectedNumBreezeTokenTypes = 22

// This should always have a number of elements equal to the constant above.
var tokenToExampleStr = map[breeze.Token]string{
	breeze.TokenStageSeparator: "|",
	breeze.TokenFilter:         "filter",
	breeze.TokenSort:           "sort",
	breeze.TokenReduce:         "group",
	breeze.TokenMap:            "map",
	breeze.TokenLParen:         "(",
	breeze.TokenRParen:         ")",
	breeze.TokenLSqBracket:     "[",
	breeze.TokenRSqBracket:     "]",
	breeze.TokenComma:          ",",
	breeze.TokenContains:       "contains",
	breeze.TokenEquals:         "=",
	breeze.TokenGEQ:            ">",
	breeze.TokenExists:         "exists",
	breeze.TokenExistsNot:      "!exists",
	breeze.TokenPlus:           "+",
	breeze.TokenMinus:          "-",
	breeze.TokenMultiply:       "*",
	breeze.TokenDivide:         "/",
	breeze.TokenFalse:          "false",
	breeze.TokenTrue:           "true",
	breeze.TokenNull:           "null",
}

func TestTokenToExampleMapAndNumTokensIsInAgreement(t *testing.T) {
	require.Equal(t, expectedNumBreezeTokenTypes, len(tokenToExampleStr))
}

// This is a fairly simple test that takes each token mapped to a string that
// matches the token, constructs a string with all those tokens concatenated,
// and then runs the tokenizer through the string, calling Next() repeatedly. It
// then expects to get back the tokens that it concatenated in the same order.
func TestTokenizerRecognizesCustomTokens(t *testing.T) {
	allExampleStrings := make([]string, expectedNumBreezeTokenTypes)
	for i := 0; i < expectedNumBreezeTokenTypes; i++ {
		allExampleStrings[i] = tokenToExampleStr[breeze.Token(i)]
	}
	// We separate by space so that tokens don't join together and become
	// something different. The tokenizer does not return whitespace tokens, so
	// this is a good way to ensure we get back exactly the tokens we put in.
	allTokensStr := strings.Join(allExampleStrings, " ")

	tokenizer := breeze.NewTokenizer(allTokensStr)

	i := 0
	for tok := tokenizer.Next(); tok != breeze.TokenEOF; tok = tokenizer.Next() {
		require.Equal(t, i, int(tok), "note: this test may fail when tokens are removed/added")
		i++
	}
}

// Same as TestTokenizerRecognizesCustomTokens but it tests the scanner package
// tokens. This one isn't as 'general' as its sibling since it doesn't deal with
// custom tokens that could change.
func TestTokenizerRecognizesScannerTokens(t *testing.T) {
	input := "9.8 'c' c 2 \"world\""
	expectedTokens := []breeze.Token{
		breeze.TokenFloat,
		breeze.TokenChar,
		breeze.TokenIdent,
		breeze.TokenInt,
		breeze.TokenString,
	}
	tokenizer := breeze.NewTokenizer(input)
	actualTokens := make([]breeze.Token, 0, len(expectedTokens))
	for i := 0; i < len(expectedTokens); i++ {
		tok := tokenizer.Next()
		if tok == breeze.TokenEOF {
			t.Fatal("expected there to be more tokens, but there was none")
		}
		actualTokens = append(actualTokens, tok)
	}

	require.Equal(t, expectedTokens, actualTokens)
}

func TestTokenizerText(t *testing.T) {
	input := "hello = world"
	tokenizer := breeze.NewTokenizer(input)
	tok := tokenizer.Next()
	require.Equal(t, breeze.TokenIdent, tok)
	require.Equal(t, "hello", tokenizer.Text())
}

func TestTokenizerPeek(t *testing.T) {
	input := "hello 9.8 321"
	tokenizer := breeze.NewTokenizer(input)
	t.Run("peek at start of string is correct", func(t *testing.T) {
		tok, tokText := tokenizer.Peek()
		require.Equal(t, breeze.TokenIdent, tok)
		require.Equal(t, "hello", tokText)
	})

	t.Run("next after peek is correct", func(t *testing.T) {
		tok := tokenizer.Next()
		require.Equal(t, breeze.TokenIdent, tok)
	})

	t.Run("peeking once more after next is correct", func(t *testing.T) {
		tok, tokText := tokenizer.Peek()
		require.Equal(t, breeze.TokenFloat, tok)
		require.Equal(t, "9.8", tokText)
	})

	t.Run("peeking is idempotent", func(t *testing.T) {
		tok, tokText := tokenizer.Peek()
		require.Equal(t, breeze.TokenFloat, tok)
		require.Equal(t, "9.8", tokText)
		tok, tokText = tokenizer.Peek()
		require.Equal(t, breeze.TokenFloat, tok)
		require.Equal(t, "9.8", tokText)
	})

	t.Run("after idempotent peeks a next is still correct", func(t *testing.T) {
		tok := tokenizer.Next()
		require.Equal(t, breeze.TokenFloat, tok)
	})

	t.Run("peeking at last token is correct", func(t *testing.T) {
		tok, tokText := tokenizer.Peek()
		require.Equal(t, breeze.TokenInt, tok)
		require.Equal(t, "321", tokText)
	})

	// Advance to the end of the string:
	tok := tokenizer.Next()
	require.Equal(t, breeze.TokenInt, tok)

	// Peek after EOF is still EOF:
	t.Run("peeking after end of string is reached gives EOF returns", func(t *testing.T) {
		tok, tokText := tokenizer.Peek()
		require.Equal(t, breeze.TokenEOF.String(), tok.String())
		require.Equal(t, "", tokText)
	})
}

// We want to make sure that Next() at the end of a input string does not give
// the user a false boolean flag while still giving a meaningful (non-EOF)
// token.
func TestTokenizerGivesFalseAndEOFAtEnd(t *testing.T) {
	input := "hello 9.8"
	tokenizer := breeze.NewTokenizer(input)

	tok := tokenizer.Next()
	require.Equal(t, tok, breeze.TokenIdent)
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenFloat)
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenEOF)
}

func TestTokenizerWithDotPrefixedTokens(t *testing.T) {
	input := "hello .foo"
	tokenizer := breeze.NewTokenizer(input)

	tok := tokenizer.Next()
	require.Equal(t, tok, breeze.TokenIdent)
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenIdent)
	// We should check that the . is included here. This could be important for
	// the parser layer to distinguish the ident.
	require.Equal(t, tokenizer.Text(), ".foo")
}

func TestTokenizerWithFunctionStyle(t *testing.T) {
	input := "hello foo(2, \"hi\")"
	tokenizer := breeze.NewTokenizer(input)

	tok := tokenizer.Next()
	require.Equal(t, tok, breeze.TokenIdent, "expected TokenIdent")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenIdent, "expected TokenIdent")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenLParen, "expected TokenLParen")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenInt, "expected TokenInt")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenComma, "expected TokenComma")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenString, "expected TokenString")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenRParen, "expected TokenRParen")
}

func TestDetectsBinaryOpTokens(t *testing.T) {
	input := "+ - / *"
	tokenizer := breeze.NewTokenizer(input)

	tok := tokenizer.Next()
	// TODO: We can factor out these require.Equal() token checks to a helper.
	require.Equal(t, tok, breeze.TokenPlus, "expected TokenPlus")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenMinus, "expected TokenMinus")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenDivide, "expected TokenPlus")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenMultiply, "expected TokenMinus")
}

func TestDetectsArrayBrackets(t *testing.T) {
	input := "[1,2]"
	tokenizer := breeze.NewTokenizer(input)

	tok := tokenizer.Next()
	require.Equal(t, tok, breeze.TokenLSqBracket, "expected TokenLSqBracket")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenInt, "expected TokenInt")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenComma, "expected TokenComma")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenInt, "expected TokenInt")
	tok = tokenizer.Next()
	require.Equal(t, tok, breeze.TokenRSqBracket, "expected TokenRSqBracket")

}

// Basically, tokens found inside of a string should be detected as indivual
// tokens. The whole string should be one TokenString.
func TestKnownTokensInStringAreNotDetected(t *testing.T) {
	allExampleStrings := make([]string, expectedNumBreezeTokenTypes)
	for i := 0; i < expectedNumBreezeTokenTypes; i++ {
		allExampleStrings[i] = tokenToExampleStr[breeze.Token(i)]
	}
	input := fmt.Sprintf("%q", strings.Join(allExampleStrings, " "))
	tokenizer := breeze.NewTokenizer(input)
	tok := tokenizer.Next()
	require.Equal(t, tok, breeze.TokenString, "expected TokenString")
}

func TestTokenizerPosition(t *testing.T) {
	input := "hello 9.8"
	tokenizer := breeze.NewTokenizer(input)

	t.Run("tokenizer position before any tokenizing is zero", func(t *testing.T) {
		require.Equal(t, 0, tokenizer.Position())
	})

	tok := tokenizer.Next()
	require.Equal(t, tok, breeze.TokenIdent)

	t.Run("tokenizer position after single tokenize is correct", func(t *testing.T) {
		// NOTE: Position() follows the semantics of the scanner package's
		// Position.Column(), which is 1-indexed, NOT 0-indexed.
		require.Equal(t, 1, tokenizer.Position())
	})

	// Get to second token ('9.8').
	tok = tokenizer.Next()

	// Get to end
	tok = tokenizer.Next()

	t.Run("tokenizer position after EOF is correct", func(t *testing.T) {
		// NOTE: Position() follows the semantics of the scanner package's
		// Position.Column(), which is 1-indexed, NOT 0-indexed.
		require.Equal(t, len(input)+1, tokenizer.Position())
	})
}

// Note that this is not a perfect test. If we ever remove N tokens and add M
// tokens in a single change, and N = M, this test will pass. However, there is
// no real good way to test Go code for exhaustive switch statements.
// I have wondered if it would be worth it to literally read the Go file, parse
// it, and then check that for all Token* variables, there exists a matching
// case in Token#String() for it, but that seems overkill for something I
// imagine we'd catch in other tests like parser tests and query execution
// tests. Furthermore, note that if we do remove any pre-existing tokens, some
// of the tests above should fail (and this file should fail to compile), which
// should hopefully ensure we update things accordingly.
func TestTokenStringerIsExhaustive(t *testing.T) {
	for i := 0; i < expectedNumBreezeTokenTypes; i++ {
		// This will panic and the test will fail if we have removed a token.
		require.NotPanics(t, func() {
			_ = breeze.Token(i).String()
		}, "note: this test may fail when tokens are removed/added")
	}

	require.Panics(t, func() {
		// Now, since this is actually outside of the number of custom tokens we
		// have (tokens are zero-indexed), this _should_ panic. Unless of course,
		// we've added a new token and forgotten to update String().
		tokenStr := breeze.Token(expectedNumBreezeTokenTypes).String()
		t.Fatalf("got a value for token number %d, but did not expect to: %q", expectedNumBreezeTokenTypes, tokenStr)
	})
}
