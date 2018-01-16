/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import(
	"io"
	"k8s.io/kubernetes/pkg/kubectl/cmd/templates"
	"k8s.io/kubernetes/pkg/kubectl/util/i18n"
	"github.com/spf13/cobra"

	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"k8s.io/kubernetes/pkg/printers"
	"sort"
	"strings"
	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/registry/rbac/role"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubernetes/pkg/printers/internalversion"
	"text/tabwriter"
	"bytes"
)

// groupResource contains the APIGroup and APIResource
type groupResource struct {
	APIGroup    string
	APIResource metav1.APIResource
}

// ApiResourceOptions contains the input to the get command.
type ApiResourcesOptions struct {
	Out io.Writer

	Namespaced     bool
	ApiResource   string
	output     string
}

var (
	apiResourceLong = templates.LongDesc(`
		Print the supported API resource`)

	apiResourcesExample = templates.Examples(i18n.T(`
		# List all apisources with table.
		kubectl api-resources

		# List all apisources with more information (such as Verbs set).
		kubectl api-resources -o wide

		# List a single apiresource with specified NAME.
		kubectl apiresources apiresource ***

		# List a single apiresource in JSON/YAML output format.
		kubectl api-resources -o json apiresource ***

		# List all *** and *** together.
		kubectl api-resources apiresource ***,***

		# List all apiresources marked namespace.
		kubectl api-resources namespaced=true`))
)

// NewCmdApiResources creates a command object for the generic "apiresource" action, which
// print all apiresource in cluster.
func NewCmdApiResources(f cmdutil.Factory, out io.Writer) *cobra.Command {
	options := &ApiResourcesOptions{
		Out:    out,
	}

	cmd := &cobra.Command{
		Use:     "api-resources [(-o|--output=)json|yaml|wide] ([APIRESOURCETYPE NAME] | [NAMESPACED TRUE|FALSE] ...) [flags]",
		Short:   i18n.T("List all resources with different types"),
		Long:    apiResourceLong + "\n\n",
		Example: apiResourcesExample,
		Run: func(cmd *cobra.Command, args []string) {
			//cmdutil.CheckErr(options.ValidateArgs(cmd, args))
			cmdutil.CheckErr(options.RunApiResources(f, cmd, args))
		},
	}

	cmd.Flags().BoolVar(&options.Namespaced, "namespaced", options.Namespaced, "If present, list the resource type for the requested object(s).")
	cmd.Flags().StringVar(&options.ApiResource, "apiresource",  options.ApiResource, " from a server.")
	cmd.Flags().StringP("output", "o", "", "Output format. One of: json|yaml|wide.")
	cmd.Flags().Bool("no-headers", false, "When using the default or custom-column output format, don't print headers (default print headers).")

	return cmd
}

// Validate checks the set of flags provided by the user.
//func (options *ApiResourcesOptions) ValidateArgs(cmd *cobra.Command, args []string) error {
//
//	outputMode := cmd.Flags().Lookup("output").Value.String()
//	switch outputMode{
//	case "wide":
//		fallthrough
//	case "yaml":
//		fallthrough
//	case "json":
//		fallthrough
//	default:
//		return fmt.Errorf("unexpected -o output mode: %v. --output should be one of 'yaml'|'wide'|'json'", outputMode)
//	}
//
//	return nil
//}

// RunApiResources performs the get operation.
func (options *ApiResourcesOptions) RunApiResources(f cmdutil.Factory, cmd *cobra.Command, args []string) error {

	//r := f.NewBuilder().Unstructured().
	//	ResourceTypeOrNameArgs(true, args...).
	//	ContinueOnError().
	//	Latest().
	//	Flatten().
	//	Do()

	printOpts := ExtractApiResourcesPrintOptions(cmd)
	printer, err := f.PrinterForOptions(printOpts)
	if err != nil {
		return err
	}
	fmt.Println(printer)
	//filterOpts := cmdutil.ExtractCmdPrintOptions(cmd, options.AllNamespaces)
	//filterFuncs := f.DefaultResourceFilterFunc()

	//return options.printGeneric(printer, r, filterFuncs, filterOpts)
	return nil
}

// ExtractCmdPrintOptions parses printer specific commandline args and
// returns a PrintOptions object.
// Requires that printer flags have been added to cmd (see AddPrinterFlags)
func ExtractApiResourcesPrintOptions(cmd *cobra.Command) *printers.PrintOptions {
	flags := cmd.Flags()

	options := &printers.PrintOptions{
		NoHeaders:          cmdutil.GetFlagBool(cmd, "no-headers"),
		Wide:               cmdutil.GetWideFlag(cmd),
	}

	var outputFormat string
	if flags.Lookup("output") != nil {
		outputFormat = cmdutil.GetFlagString(cmd, "output")
	}

	if flags.Lookup("sort-by") != nil {
		options.SortBy = cmdutil.GetFlagString(cmd, "sort-by")
	}

	options.OutputFormatType = outputFormat

	return options
}

func (o *ApiResourcesOptions) GetApiResource(f cmdutil.Factory) ([]groupResource, error) {
	discoveryClient, err := f.DiscoveryClient()
	if err != nil {
		return nil, err
	}

	// Always request fresh data from the server
	discoveryClient.Invalidate()

	lists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, fmt.Errorf("get available api resources from server failed: %v", err)
	}

	resources := []groupResource{}

	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}
		gv, err := schema.ParseGroupVersion(list.GroupVersion)
		if err != nil {
			glog.V(1).Infof("Unable to parse groupversion %s:%s", list.GroupVersion, err.Error())
			continue
		}
		for _, resource := range list.APIResources {
			if len(resource.Verbs) == 0 {
				continue
			}
			//// filter apiGroup
			//if o.apiGroup != gv.Group {
			//	continue
			//}
			//// filter namespaced
			//if o.namespaced != resource.Namespaced {
			//	continue
			//}
			resources = append(resources, groupResource{
				APIGroup:    gv.Group,
				APIResource: resource,
			})
		}
	}

	return resources, nil
}

func (options *ApiResourcesOptions) DescribeApiResource(f cmdutil.Factory, namespace, name string, describerSettings printers.DescriberSettings) (string, error) {
	//role, err := d.Rbac().ClusterRoles().Get(name, metav1.GetOptions{})
	//if err != nil {
	//	return "", err
	//}
	//
	//breakdownRules := []rbac.PolicyRule{}
	//for _, rule := range role.Rules {
	//	breakdownRules = append(breakdownRules, validation.BreakdownRule(rule)...)
	//}
	//sort.Stable(rbac.SortableRuleSlice(compactRules))
	// 获取资源
	// 资源排序
	resources, err := options.GetApiResource(f)
	if err != nil {
		glog.V(1).Infof("Get available cluster resources failed: %v", err)
	}

	sort.Stable(sortableGroupResource(resources))

	return tabbedString(func(out io.Writer) error {
		w := internalversion.NewPrefixWriter(out)
		w.Write(1, "Resources\tNon-Resource URLs\tResource Names\tVerbs\n")
		w.Write(1, "---------\t-----------------\t--------------\t-----\n")
		for _, r := range resources {
			w.Write(1, "%s\t%v\t%v\t%v\n", combineResourceGroup(r.Resources, r.APIGroups), r.NonResourceURLs, r.ResourceNames, r.Verbs)
		}

		return nil
	})
}

func tabbedString(f func(io.Writer) error) (string, error) {
	out := new(tabwriter.Writer)
	buf := &bytes.Buffer{}
	out.Init(buf, 0, 8, 2, ' ', 0)

	err := f(out)
	if err != nil {
		return "", err
	}

	out.Flush()
	str := string(buf.String())
	return str, nil
}

type sortableGroupResource []groupResource

func (s sortableGroupResource) Len() int      { return len(s) }
func (s sortableGroupResource) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s sortableGroupResource) Less(i, j int) bool {
	ret := strings.Compare(s[i].APIGroup, s[j].APIGroup)
	if ret > 0 {
		return false
	} else if ret == 0 {
		return strings.Compare(s[i].APIResource.Name, s[j].APIResource.Name) < 0
	}
	return true
}
