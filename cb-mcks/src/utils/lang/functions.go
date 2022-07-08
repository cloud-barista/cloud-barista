package lang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	// Random string generation
	letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}

/* generate to a random string */
func GenerateNewRandomString(n int) string {
	randSrc := rand.NewSource(time.Now().UnixNano()) //Random source by nano time
	b := make([]byte, n)
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

/* generate to a new node name */
func GenerateNewNodeName(role string, idx int) string {
	return fmt.Sprintf("%s-%d-%s", role[:1], idx, GenerateNewRandomString(5))
}

/* get a idex of node name */
func GetNodeNameIndex(nodeName string) int {
	a := strings.Split(nodeName, "-")
	if len(a) >= 2 {
		if idx, err := strconv.ParseInt(a[1], 0, 64); err != nil {
			return 0
		} else {
			return int(idx)
		}
	}
	return 0
}

/* replace all */
func ReplaceAll(source string, olds []string, new string) string {

	for _, s := range olds {
		source = strings.ReplaceAll(source, s, new)
	}
	return source
}

/* verify cluster name */
func VerifyClusterName(name string) error {
	reg, _ := regexp.Compile("[a-z]([-a-z0-9]*[a-z0-9])?")
	filtered := reg.FindString(name)

	if filtered != name {
		return errors.New(name + ": The first character of name must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.")
	}
	return nil
}

/* verify CIDR */
func VerifyCIDR(name string, val string) error {
	reg, _ := regexp.Compile("^((?:\\d{1,3}.){3}\\d{1,3})\\/(\\d{1,2})$")
	filtered := reg.FindString(val)

	if filtered != val {
		return errors.New(fmt.Sprintf("%s %s : Type mismatch ex)10.244.0.0/16", name, val))
	}
	return nil
}

/* if it's not alpabet & number then replace to "" */
func GetOnlyLettersAndNumbers(name string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	val := reg.ReplaceAllString(name, "")

	return val
}

/* get a now string  */
func GetNowUTC() string {
	t := time.Now().UTC()
	return t.Format(time.RFC3339)
}

func ToPrettyJSON(data []byte) []byte {

	if len(data) > 0 {
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", "  "); err == nil {
			return buf.Bytes()
		}
	}
	return data
}

func ToTemplateBytes(tpl string, todo interface{}) ([]byte, error) {

	t, err := template.New("tpl").Funcs(
		template.FuncMap{
			"ToUpper": strings.ToUpper,
		}).Parse(tpl)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = t.Execute(&out, todo)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil

}
