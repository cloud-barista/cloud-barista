package app

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"github.com/cloud-barista/cb-mcks/src/utils/lang"
)

func ValidateError(c *cobra.Command, err error) {

	if err != nil {
		c.Help()
		msg := "\n" + err.Error()
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
		os.Exit(1)
	}

}

func BindCommandArgs(values ...*string) func(c *cobra.Command, args []string) error {

	return func(c *cobra.Command, args []string) error {

		for i, v := range args {
			if len(values) > i {
				*values[i] = v
			}
		}
		return nil
	}

}

func GetBody(o IOptions, tpl string) (buf []byte, err error) {
	// -f 옵션
	fileName := o.GetFilename()
	if len(fileName) > 0 {
		switch {
		case fileName == "-": // standard-in
			buf, err = io.ReadAll(os.Stdin)
		case strings.Index(fileName, "http://") == 0 || strings.Index(fileName, "https://") == 0: // http
			if _, err = url.Parse(fileName); err == nil {
				if resp, err := http.Get(fileName); err == nil {
					defer resp.Body.Close()
					buf, err = io.ReadAll(resp.Body)
				}
			}
		default:
			buf, err = os.ReadFile(fileName) // local file
		}

		if err == nil {
			buf, err = yaml.YAMLToJSON(buf)
		}
	} else if len(o.GetData()) > 0 {
		buf, err = lang.ToTemplateBytes(o.GetData(), o)
	} else {
		buf, err = lang.ToTemplateBytes(tpl, o)
	}

	return
}
