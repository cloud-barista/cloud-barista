package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/cloud-barista/cb-mcks/src/grpc-api/cbadm/app"
)

type ConfigOptions struct {
	*app.Options
	Mcks_Server_addr    string
	Mcks_Timeout        string
	Mcks_Endpoint       string
	Mcks_Service_name   string
	Mcks_Sample_rate    string
	Spider_Server_addr  string
	Spider_Timeout      string
	Spider_Endpoint     string
	Spider_Service_name string
	Spider_Sample_rate  string
}

func (o *ConfigOptions) writeYaml(in interface{}) {
	if b, err := yaml.Marshal(in); err != nil {
		o.PrintlnError(err)
	} else {
		o.WriteBody(b)
	}
}

// returns a cobra command
func NewCommandConfig(options *app.Options) *cobra.Command {
	o := &ConfigOptions{
		Options: options,
	}

	// root
	cmds := &cobra.Command{
		Use:   "config",
		Short: "Configuration command",
		Long:  "This is a configuration command",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// add-context
	cmdC := &cobra.Command{
		Use:                   "add-context (NAME | --name NAME) [options]",
		Short:                 "Add a context",
		DisableFlagsInUseLine: true,
		Args:                  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if len(o.Name) == 0 {
					return fmt.Errorf("Name is required.")
				}
				if _, ok := app.Config.Contexts[o.Name]; ok {
					return fmt.Errorf("The context '%s' is alreaday exist", o.Name)
				} else {
					var sConf *app.CliConfig = new(app.CliConfig)

					sConf.ServerAddr = o.Spider_Server_addr
					sConf.Timeout = o.Spider_Timeout
					sConf.Interceptors.Opentracing.Jaeger.Endpoint = o.Spider_Endpoint
					sConf.Interceptors.Opentracing.Jaeger.ServiceName = o.Spider_Service_name
					sConf.Interceptors.Opentracing.Jaeger.SampleRate = o.Spider_Sample_rate

					var gConf *app.CliConfig = new(app.CliConfig)

					gConf.ServerAddr = o.Mcks_Server_addr
					gConf.Timeout = o.Mcks_Timeout
					gConf.Interceptors.Opentracing.Jaeger.Endpoint = o.Mcks_Endpoint
					gConf.Interceptors.Opentracing.Jaeger.ServiceName = o.Mcks_Service_name
					gConf.Interceptors.Opentracing.Jaeger.SampleRate = o.Mcks_Sample_rate

					app.Config.Contexts[o.Name] = &app.ConfigContext{
						Name:      o.Name,
						Namespace: o.Namespace,
						Mckscli:   gConf,
						Spidercli: sConf,
					}
				}
				app.Config.WriteConfig()
				o.writeYaml(app.Config)
				return nil
			}())
		},
	}
	cmdC.Flags().StringVarP(&o.Mcks_Server_addr, "mcks_server_addr", "", "127.0.0.1:50254", "Server Addr URL")
	cmdC.Flags().StringVarP(&o.Mcks_Timeout, "mcks_timeout", "", "1000s", "Timeout")
	cmdC.Flags().StringVarP(&o.Mcks_Endpoint, "mcks_endpoint", "", "localhost:6834", "endpoint URL")
	cmdC.Flags().StringVarP(&o.Mcks_Service_name, "mcks_service_name", "", "mcks grpc client", "Service Name")
	cmdC.Flags().StringVarP(&o.Mcks_Sample_rate, "mcks_sample_rate", "", "1", "sample rate")
	cmdC.Flags().StringVarP(&o.Spider_Server_addr, "spider_server_addr", "", "127.0.0.1:2048", "Server Addr URL")
	cmdC.Flags().StringVarP(&o.Spider_Timeout, "spider_timeout", "", "1000s", "Timeout")
	cmdC.Flags().StringVarP(&o.Spider_Endpoint, "spider_endpoint", "", "localhost:6832", "endpoint URL")
	cmdC.Flags().StringVarP(&o.Spider_Service_name, "spider_service_name", "", "spider grpc client", "Service Name")
	cmdC.Flags().StringVarP(&o.Spider_Sample_rate, "spider_sample_rate", "", "1", "sample rate")
	cmds.AddCommand(cmdC)

	// view
	cmds.AddCommand(&cobra.Command{
		Use:   "view",
		Short: "Get contexts",
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				o.writeYaml(app.Config)
				return nil
			}())
		},
	})

	// get context
	cmds.AddCommand(&cobra.Command{
		Use:   "get-context (NAME | --name NAME) [options]",
		Short: "Get a context",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					for k := range app.Config.Contexts {
						o.Println(k)
					}
				} else {
					if app.Config.Contexts[o.Name] != nil {
						o.writeYaml(app.Config.Contexts[o.Name])
					}
				}
				return nil
			}())
		},
	})

	// set context
	cmdS := &cobra.Command{
		Use:                   "set-context (NAME | --name NAME) [options]",
		Short:                 "Set a context",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					c.Help()
				} else if app.Config.Contexts[o.Name] != nil {
					app.Config.Contexts[o.Name].Name = o.Name
					if o.Mcks_Server_addr != "" {
						app.Config.Contexts[o.Name].Mckscli.ServerAddr = o.Mcks_Server_addr
					}
					if o.Mcks_Timeout != "" {
						app.Config.Contexts[o.Name].Mckscli.Timeout = o.Mcks_Timeout
					}
					if o.Mcks_Endpoint != "" {
						app.Config.Contexts[o.Name].Mckscli.Interceptors.Opentracing.Jaeger.Endpoint = o.Mcks_Endpoint
					}
					if o.Mcks_Service_name != "" {
						app.Config.Contexts[o.Name].Mckscli.Interceptors.Opentracing.Jaeger.ServiceName = o.Mcks_Service_name
					}
					if o.Mcks_Sample_rate != "" {
						app.Config.Contexts[o.Name].Mckscli.Interceptors.Opentracing.Jaeger.SampleRate = o.Mcks_Sample_rate
					}
					if o.Spider_Server_addr != "" {
						app.Config.Contexts[o.Name].Spidercli.ServerAddr = o.Spider_Server_addr
					}
					if o.Spider_Timeout != "" {
						app.Config.Contexts[o.Name].Spidercli.Timeout = o.Spider_Timeout
					}
					if o.Spider_Endpoint != "" {
						app.Config.Contexts[o.Name].Spidercli.Interceptors.Opentracing.Jaeger.Endpoint = o.Spider_Endpoint
					}
					if o.Spider_Service_name != "" {
						app.Config.Contexts[o.Name].Spidercli.Interceptors.Opentracing.Jaeger.ServiceName = o.Spider_Service_name
					}
					if o.Spider_Sample_rate != "" {
						app.Config.Contexts[o.Name].Spidercli.Interceptors.Opentracing.Jaeger.SampleRate = o.Spider_Sample_rate
					}
					o.writeYaml(app.Config.Contexts[o.Name])
				} else {
					o.Println("Not found a context (name=%s)", o.Name)
				}
				return nil
			}())
		},
	}
	cmdS.Flags().StringVarP(&o.Mcks_Server_addr, "mcks_server_addr", "", "127.0.0.1:50254", "Server Addr URL")
	cmdS.Flags().StringVarP(&o.Mcks_Timeout, "mcks_timeout", "", "1000s", "Timeout")
	cmdS.Flags().StringVarP(&o.Mcks_Endpoint, "mcks_endpoint", "", "localhost:6834", "endpoint URL")
	cmdS.Flags().StringVarP(&o.Mcks_Service_name, "mcks_service_name", "", "mcks grpc client", "Service Name")
	cmdS.Flags().StringVarP(&o.Mcks_Sample_rate, "mcks_sample_rate", "", "1", "sample rate")
	cmdS.Flags().StringVarP(&o.Spider_Server_addr, "spider_server_addr", "", "127.0.0.1:2048", "Server Addr URL")
	cmdS.Flags().StringVarP(&o.Spider_Timeout, "spider_timeout", "", "1000s", "Timeout")
	cmdS.Flags().StringVarP(&o.Spider_Endpoint, "spider_endpoint", "", "localhost:6832", "endpoint URL")
	cmdS.Flags().StringVarP(&o.Spider_Service_name, "spider_service_name", "", "spider grpc client", "Service Name")
	cmdS.Flags().StringVarP(&o.Spider_Sample_rate, "spider_sample_rate", "", "1", "sample rate")
	cmds.AddCommand(cmdS)

	// current-context (get/set)
	cmds.AddCommand(&cobra.Command{
		Use:                   "current-context (NAME | --name NAME) [options]",
		Short:                 "Get/Set a current context",
		DisableFlagsInUseLine: true,
		Args:                  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if len(o.Name) > 0 {
					_, ok := app.Config.Contexts[o.Name]
					if ok {
						app.Config.CurrentContext = o.Name
						app.Config.WriteConfig()
					} else {
						o.Println("context '%s' is not exist\n", o.Name)
					}
				}
				o.writeYaml(app.Config.GetCurrentContext().Name)
				return nil
			}())
		},
	})

	// delete-context
	cmds.AddCommand(&cobra.Command{
		Use:   "delete-context (NAME | --name NAME) [options]",
		Short: "Delete a context",
		Args:  app.BindCommandArgs(&o.Name),
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if o.Name == "" {
					return fmt.Errorf("Name Required.")
				}
				conf := app.Config
				if len(conf.Contexts) > 1 {
					delete(conf.Contexts, o.Name)
					if o.Name == conf.CurrentContext {
						conf.CurrentContext = func() string {
							if len(conf.Contexts) > 0 {
								for k := range conf.Contexts {
									return k
								}
							}
							return ""
						}()
					}
					conf.WriteConfig()
				}
				o.writeYaml(conf)
				return nil
			}())
		},
	})

	// set-namespace
	cmds.AddCommand(&cobra.Command{
		Use:                   "set-namespace (NAME | --name NAME) [options]",
		Short:                 "Set a namespace to context",
		Args:                  app.BindCommandArgs(&o.Name),
		DisableFlagsInUseLine: true,
		Run: func(c *cobra.Command, args []string) {
			app.ValidateError(c, func() error {
				if len(app.Config.GetCurrentContext().Name) == 0 {
					c.Help()
				} else {
					app.Config.GetCurrentContext().Namespace = args[0]
					app.Config.WriteConfig()
					o.writeYaml(app.Config.GetCurrentContext())
				}
				return nil
			}())
		},
	})

	return cmds
}
