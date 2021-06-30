package lang

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	// Random string generation
	letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var (
	// Random source by nano time
	randSrc = rand.NewSource(time.Now().UnixNano())
)

// NVL is null value logic
func NVL(str string, def string) string {
	if len(str) == 0 {
		return def
	}
	return str
}

// get store cluster key
func GetStoreClusterKey(namespace string, clusterName string) string {
	if clusterName == "" {
		return fmt.Sprintf("/ns/%s/clusters", namespace)
	} else {
		return fmt.Sprintf("/ns/%s/clusters/%s", namespace, clusterName)
	}
}

// get store node key
func GetStoreNodeKey(namespace string, clusterName string, nodeName string) string {
	if nodeName == "" {
		return fmt.Sprintf("/ns/%s/clusters/%s/nodes", namespace, clusterName)
	} else {
		return fmt.Sprintf("/ns/%s/clusters/%s/nodes/%s", namespace, clusterName, nodeName)
	}
}

// get uuid
func GetUid() string {
	return uuid.New().String()
}

// Random string generation
func GetRandomString(n int) string {
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

// get node name
func GetNodeName(clusterName string, role string, idx int) string {
	return fmt.Sprintf("%s-%s-%d-%s", clusterName, role[:1], idx, GetRandomString(5))
}

func GetIdxToInt(idx string) int {
	i, err := strconv.Atoi(idx)
	if err != nil {
		i = 0
	}
	return i
}

func GetMaxNumber(arr []int) int {
	max := 0
	for _, val := range arr {
		if val > max {
			max = val
		}
	}
	return max
}

func CheckName(name string) error {
	reg, _ := regexp.Compile("[a-z]([-a-z0-9]*[a-z0-9])?")
	filtered := reg.FindString(name)

	if filtered != name {
		return errors.New(name + ": The first character of name must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.")
	}
	return nil
}

func CheckIpCidr(name string, val string) error {
	reg, _ := regexp.Compile("^((?:\\d{1,3}.){3}\\d{1,3})\\/(\\d{1,2})$")
	filtered := reg.FindString(val)

	if filtered != val {
		return errors.New(fmt.Sprintf("%s %s : Type mismatch ex)10.244.0.0/16", name, val))
	}
	return nil
}
