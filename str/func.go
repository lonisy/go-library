package str

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

func UriParseArgs(str string) string {
    re := regexp.MustCompile(`/(\d+)/`)
    for {
        _str := re.ReplaceAllString(str, "/:arg/")
        if _str != str {
            str = _str
        } else {
            break
        }
    }
    re2 := regexp.MustCompile(`/:arg`)
    count := 0
    return re2.ReplaceAllStringFunc(str, func(match string) string {
        arg := ":arg" + strconv.Itoa(count)
        count++
        return "/" + arg
    })
}

// ///api/checkBrief/get/198
func UriParseArgs2(str string) string {
    re := regexp.MustCompile(`/(\d+)$`)
    for {
        _str := re.ReplaceAllString(str, "/:arg")
        if _str != str {
            str = _str
        } else {
            break
        }
    }
    re2 := regexp.MustCompile(`/:arg`)
    count := 0
    return re2.ReplaceAllStringFunc(str, func(match string) string {
        arg := ":arg" + strconv.Itoa(count)
        count++
        return "/" + arg
    })
}
func UriParseArgsUUID(str string) string {
    re := regexp.MustCompile(`/([a-f0-9]{7}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})$`)
    for {
        _str := re.ReplaceAllString(str, "/:arg")
        if _str != str {
            str = _str
        } else {
            break
        }
    }
    re2 := regexp.MustCompile(`/:arg`)
    count := 0
    return re2.ReplaceAllStringFunc(str, func(match string) string {
        arg := ":arg" + strconv.Itoa(count)
        count++
        return "/" + arg
    })
}

// camel-case => Camel-Case
func CamelCase(in, sep string) string {
    tokens := strings.Split(in, sep)
    for i := range tokens {
        tokens[i] = strings.Title(strings.Trim(tokens[i], " "))
    }
    return strings.Join(tokens, sep)
}

func StringToCamel(s string) string {
    data := make([]byte, 0, len(s))
    j := false
    k := false
    num := len(s) - 1
    for i := 0; i <= num; i++ {
        d := s[i]
        if k == false && d >= 'A' && d <= 'Z' {
            k = true
        }
        if d >= 'a' && d <= 'z' && (j || k == false) {
            d = d - 32
            j = false
            k = true
        }
        if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
            j = true
            continue
        }
        data = append(data, d)
    }
    return string(data[:])
}

// GoCamelCase camel-cases a protobuf name for use as a Go identifier.
//
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
func GoCamelCase(s string) string {
    // Invariant: if the next letter is lower case, it must be converted
    // to upper case.
    // That is, we process a word at a time, where words are marked by _ or
    // upper case letter. Digits are treated as words.
    var b []byte
    for i := 0; i < len(s); i++ {
        c := s[i]
        switch {
        case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
            // Skip over '.' in ".{{lowercase}}".
        case c == '.':
            b = append(b, '_') // convert '.' to '_'
        case c == '_' && (i == 0 || s[i-1] == '.'):
            // Convert initial '_' to ensure we start with a capital letter.
            // Do the same for '_' after '.' to match historic behavior.
            b = append(b, 'X') // convert '_' to 'X'
        case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
            // Skip over '_' in "_{{lowercase}}".
        case isASCIIDigit(c):
            b = append(b, c)
        default:
            // Assume we have a letter now - if not, it's a bogus identifier.
            // The next word is a sequence of characters that must start upper case.
            if isASCIILower(c) {
                c -= 'a' - 'A' // convert lowercase to uppercase
            }
            b = append(b, c)

            // Accept lower case sequence that follows.
            for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ {
                b = append(b, s[i+1])
            }
        }
    }
    return string(b)
}
func isASCIILower(c byte) bool {
    return 'a' <= c && c <= 'z'
}
func isASCIIUpper(c byte) bool {
    return 'A' <= c && c <= 'Z'
}
func isASCIIDigit(c byte) bool {
    return '0' <= c && c <= '9'
}

func TruncateStr(str string, limit int) string {
    var i int
    var truncatedStr strings.Builder
    for _, r := range str {
        if i >= limit {
            break
        }
        truncatedStr.WriteRune(r)
        i++
    }
    return truncatedStr.String()
}

func ExtractLanguageNo(string2 string) string {
    re := regexp.MustCompile(`\((.*?)\)`)
    match := re.FindStringSubmatch(string2)
    if len(match) > 1 {
        language := match[1]
        return language
    }
    return ""
}

func StandardizedTranslationText(str string) string {
    if strings.Contains(str, `\n`) {
        str = strings.ReplaceAll(str, `\n`, "||1||")
    }
    if strings.Contains(str, "\n") {
        str = strings.ReplaceAll(str, "\n", "||2||")
    }
    //str = strings.ReplaceAll(str, `"`, `\"`)
    return str
}

func StandardizedTranslationTextToOriginal(str string) string {
    if strings.Contains(str, "||1||") {
        str = strings.ReplaceAll(str, "||1||", `\n`)
    }
    if strings.Contains(str, "||one||") {
        str = strings.ReplaceAll(str, "||one||", `\n`)
    }
    if strings.Contains(str, "||2||") {
        str = strings.ReplaceAll(str, "||2||", "\n")
    }
    if strings.Contains(str, "||two||") {
        str = strings.ReplaceAll(str, "||two||", "\n")
    }
    return str
}

func TranslationResultTrim(str string) string {
    str = strings.TrimPrefix(str, "```json")
    str = strings.TrimPrefix(str, "```JSON")
    str = strings.TrimPrefix(str, " ```JSON")
    str = strings.TrimPrefix(str, " ```JSON")
    str = strings.TrimSuffix(str, "```")
    str = strings.ReplaceAll(str, `{
  "`, `{"`)
    str = strings.ReplaceAll(str, `",
  "`, `","`)
    str = strings.ReplaceAll(str, `"
}`, `"}`)
    return str
}

func ExtractLanguageCode(input string) (string, error) {
    re := regexp.MustCompile(`\((\w+)\)`)
    matches := re.FindStringSubmatch(input)
    if len(matches) < 2 {
        return "", fmt.Errorf("no language code found in string")
    }
    return matches[1], nil
}
