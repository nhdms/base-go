package utils

import (
	"bytes"
	"compress/flate"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (

	// Define common regex patterns for API keys and tokens.
	// These are examples and can be extended.
	// We'll focus on patterns that capture the *value* after a key or prefix.

	// Regex for values following common API key names in query parameters or headers.
	// Use non-greedy matching `.*?` and lookbehind/lookahead if supported and needed,
	// but simple capture groups for values are often sufficient for common cases.
	//
	// Important: We want to match the VALUE, not the key name.
	// This regex targets common key names followed by an equals sign and then captures the value.
	// It assumes URL query parameters or similar key=value structures.
	apiKeyValRegex = regexp.MustCompile(`(?i)(?:api_?key|auth_?token|secret|access_?token|client_?secret)=([a-zA-Z0-9\-_=&\.]{30,})`)
	// We make sure the value is at least 30 characters long to avoid false positives for short strings.
	// The character set `[a-zA-Z0-9\-_=&\.]` is broad to cover various token characters.
	// We use `&` as a potential terminator for URL parameters, or `\.` for JWT-like parts.

	// Regex for Bearer tokens: "Bearer " followed by a long alphanumeric string
	bearerTokenRegex = regexp.MustCompile(`(?i)Bearer\s+([a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+|[a-zA-Z0-9\-_=]{30,})`)
	// Capture group 1 is the actual token value.

	// Regex for JWTs (three parts separated by dots, base64 encoded characters)
	// This should be applied carefully as JWTs can appear in various contexts.
	jwtRegex = regexp.MustCompile(`([a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+)`)
	// Capture group 1 is the full JWT.

	// Order of operations: More specific patterns first.
)

func SubStringArrays(src []string, sub []string) []string {
	if sub == nil || len(sub) <= 0 {
		return src
	}
	mx := map[string]bool{}
	for _, x := range sub {
		mx[x] = true
	}
	arrayResult := []string{}
	for _, x := range src {
		if _, ok := mx[x]; !ok {
			arrayResult = append(arrayResult, x)
		}
	}
	return arrayResult
}

func SubString(str string, start, end int) string {
	if len(str) < start+end {
		return str
	}

	str = str[start:end]

	// Remove invalid character
	if !utf8.ValidString(str) {
		v := make([]rune, 0, len(str))
		for i, r := range str {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(str[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		str = string(v)
	}

	return str
}

func SplitStringOnNewLine(str string) []string {
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "\n", -1)

	return strings.Split(str, "\n")
}

func RemoveUnderscoreAndPercentCharacter(str string) string {
	str = strings.Replace(str, "%", "\\%", -1)
	str = strings.Replace(str, "_", "\\_", -1)
	return str
}

func HasPrefixCaseInsensitive(s, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

func HasSuffixCaseInsensitive(s, prefix string) bool {
	return strings.HasSuffix(strings.ToLower(s), strings.ToLower(prefix))
}

func ContainsCaseInsensitive(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func NormalizeSpecializedSpace(input string) string {
	var output []string

	for _, i := range input {
		switch {
		case unicode.IsSpace(i):
			output = append(output, " ")
		default:
			output = append(output, string(i))
		}
	}
	return strings.Join(output, "")
}

func SubNString(str string, n int) string {
	if len(str) <= n {
		n = len(str)
	}
	return str[:n]
}

func IsLatin(str string) bool {
	reg := regexp.MustCompile(`^[A-Za-z1-9\s\x21-\x7f\xC0-\xD6\xD8-\xf6\xf8-\xff]+$`)
	return reg.MatchString(str)
}

// Remove any char is not number character
func ParseNumericFromString(value string) string {
	reg, _ := regexp.Compile("[^0-9]+")
	processedString := reg.ReplaceAllString(value, "")
	return processedString
}

// Check a string is valid CPF number format
func IsValidCPFNumberFormat(value string) bool {
	number := ParseNumericFromString(value)
	if len(number) == 11 {
		return true
	}
	return false
}

func EncodeToString(max int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func SortSliceStringOptionSize(sizeSlice []string) (size []string) {
	indexSortSliceSize := map[string]int64{}
	configIndexSortSliceSize := viper.GetStringMap("sort_option.index_sort_slice_size")

	for key, value := range configIndexSortSliceSize {
		indexSortSliceSize[strings.ToUpper(key)] = cast.ToInt64(value)
	}
	other := make([]string, 0)
	sizeSlice = UniqueStringSlice(sizeSlice)
	for _, entry := range sizeSlice {
		if _, ok := indexSortSliceSize[strings.ToUpper(entry)]; ok {
			size = append(size, entry)
		} else {
			other = append(other, entry)
		}
	}
	sort.Slice(size, func(i, j int) bool {
		return indexSortSliceSize[strings.ToUpper(size[i])] < indexSortSliceSize[strings.ToUpper(size[j])]
	})
	sort.Strings(other)
	size = append(size, other...)
	return size
}

func UniqueStringSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func CompareString(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

// Encode/Compress string by flate
func Deflate(rawContent string) (encodedContent string) {
	if len(rawContent) > 0 {
		// flate compress
		var b bytes.Buffer
		w, err := flate.NewWriter(&b, flate.BestSpeed)
		if err != nil {
			return encodedContent
		}

		w.Write([]byte(rawContent))
		w.Close()

		// base64 encode
		encodedContent = base64.StdEncoding.EncodeToString(b.Bytes())
	}

	return encodedContent
}

// Decode/Uncompress string by flate
func Inflate(encodedContent string) (rawContent string) {
	if len(encodedContent) > 0 {
		// base64 decode
		str64, err := base64.StdEncoding.DecodeString(encodedContent)
		if err != nil {
			return encodedContent
		}

		// flate uncompressed
		r := flate.NewReader(bytes.NewReader(str64))
		buf := new(bytes.Buffer)
		buf.ReadFrom(r)
		defer r.Close()

		rawContent = buf.String()
	}

	return rawContent
}

func IsInflateSettingData(data string) bool {
	mapData := map[string]interface{}{}
	if err := json.Unmarshal([]byte(data), &mapData); err != nil {
		return true
	}
	return false
}

func Join[T comparable](slice []T, sep string) string {
	resp := make([]string, len(slice))
	for i, v := range slice {
		resp[i] = cast.ToString(v)
	}
	return strings.Join(resp, sep)
}

// censorString replaces a portion of a string with asterisks
// while preserving the original length.
func censorString(s string) string {
	if len(s) == 0 {
		return ""
	}
	// Censor most of the string, leave a few characters at the beginning and end if it's long enough
	if len(s) > 8 { // For longer strings, show first 4 and last 4 chars
		return s[:4] + strings.Repeat("*", len(s)-8) + s[len(s)-4:]
	}
	return strings.Repeat("*", len(s)) // For shorter strings, censor entirely
}

// DetectAndCensorAPIKeysTokens finds and censors common API keys and tokens in a given string.
// It returns the censored string and a boolean indicating if any sensitive information was found.
func DetectAndCensorAPIKeysTokens(input string) (censoredString string) {
	censoredString = input
	foundSensitive := false

	// 1. Censor Bearer Tokens (most specific context)
	censoredString = bearerTokenRegex.ReplaceAllStringFunc(censoredString, func(s string) string {
		foundSensitive = true
		matches := bearerTokenRegex.FindStringSubmatch(s)
		if len(matches) > 1 {
			// matches[0] is the full match ("Bearer ABCDEF"), matches[1] is the captured group ("ABCDEF")
			return strings.Replace(s, matches[1], censorString(matches[1]), 1)
		}
		return censorString(s) // Fallback
	})

	// 2. Censor API key values in query parameters (e.g., access_token=VALUE)
	censoredString = apiKeyValRegex.ReplaceAllStringFunc(censoredString, func(s string) string {
		foundSensitive = true
		matches := apiKeyValRegex.FindStringSubmatch(s)
		if len(matches) > 1 {
			// matches[0] is the full match ("access_token=VALUE"), matches[1] is the captured group ("VALUE")
			return strings.Replace(s, matches[1], censorString(matches[1]), 1)
		}
		return censorString(s) // Fallback
	})

	// 3. Censor standalone JWTs (less specific, might catch non-tokens if not careful)
	// This might be redundant if the bearer token regex already catches JWTs.
	// You might want to skip this if your primary concern is `Bearer` tokens or `access_token=` values.
	censoredString = jwtRegex.ReplaceAllStringFunc(censoredString, func(s string) string {
		foundSensitive = true
		matches := jwtRegex.FindStringSubmatch(s)
		if len(matches) > 1 {
			return strings.Replace(s, matches[1], censorString(matches[1]), 1)
		}
		return censorString(s)
	})

	_ = foundSensitive
	return censoredString
}
