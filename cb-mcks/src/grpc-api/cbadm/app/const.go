package app

import (
	"fmt"
	"os"

	"github.com/cloud-barista/cb-mcks/src/utils/lang"
	"github.com/ghodss/yaml"
)

type IOptions interface {
	GetFilename() string
	GetData() string
}

type Options struct {
	OutStream  *os.File // output stream
	ConfigFile string   // config file
	Output     string   // output format (json/yaml)
	Filename   string   // file
	Data       string   // input data
	Namespace  string   // cloud-barista namespace
	Name       string   // object name
}

func (o *Options) GetFilename() string {
	return o.Filename
}
func (o *Options) GetData() string {
	return o.Data
}

func (o *Options) Println(format string, params ...interface{}) {
	msg := fmt.Sprintf(format+"\n", params...)
	if o.OutStream != nil {
		o.OutStream.WriteString(msg)
	} else {
		os.Stdout.WriteString(msg)
	}
}
func (o *Options) PrintlnError(err error) {
	o.Println("%+v\n", err)
}

func (o *Options) WriteBody(json []byte) {
	if o.Output == OUTPUT_JSON {
		o.OutStream.Write(lang.ToPrettyJSON(json))
	} else {
		if d, err := yaml.JSONToYAML(json); err == nil {
			o.OutStream.Write(d)
		} else {
			o.OutStream.Write(json)
		}
	}
}

type IOStreams struct {
	In     *os.File
	Out    *os.File
	ErrOut *os.File
}

const (
	OUTPUT_JSON = "json"
	OUTPUT_YAML = "yaml"
)

type OutputType string

type Output struct {
	Type   *string
	Stream *os.File
}

func (self *Output) Write(json []byte) {
	if *self.Type == OUTPUT_JSON {
		self.Stream.Write(lang.ToPrettyJSON(json))
	} else {
		if d, err := yaml.JSONToYAML(json); err == nil {
			self.Stream.Write(d)
		} else {
			self.Stream.Write(json)
		}
	}
}

func (self *Output) WriteString(format string, params ...interface{}) {
	line := fmt.Sprintf(format, params...)
	self.Stream.WriteString(line)
}
