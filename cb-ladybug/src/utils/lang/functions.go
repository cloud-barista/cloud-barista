package lang

import (
	"fmt"
	"math/rand"
	"regexp"
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

// for worker node join command
func GetWorkerJoinCmd(cpInitResult string) string {
	var join1, join2 string
	joinRegex, _ := regexp.Compile("kubeadm\\sjoin\\s(.*?)\\s--token\\s(.*?)\\s")
	joinRegex2, _ := regexp.Compile("--discovery-token-ca-cert-hash\\ssha256:(.*?)\\n")

	if joinRegex.MatchString(cpInitResult) {
		res := joinRegex.FindStringSubmatch(cpInitResult)
		join1 = res[0]
	}
	if joinRegex2.MatchString(cpInitResult) {
		res := joinRegex2.FindStringSubmatch(cpInitResult)
		join2 = res[0]
	}

	return fmt.Sprintf("sudo %s %s", join1, join2)
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
func GetNodeName(clusterName string, role string) string {
	return fmt.Sprintf("%s-%s-%s", clusterName, role[:1], GetRandomString(5))
}
