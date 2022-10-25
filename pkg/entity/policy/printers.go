package policy

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	v1 "github.com/strmprivacy/api-definitions-go/v2/api/entities/v1"
	"github.com/strmprivacy/api-definitions-go/v2/api/policies/v1"
	"strmprivacy/strm/pkg/common"
	"strmprivacy/strm/pkg/util"
)

var printer util.Printer

func configurePrinter(command *cobra.Command) util.Printer {
	outputFormat := util.GetStringAndErr(command.Flags(), common.OutputFormatFlag)

	p := availablePrinters()[outputFormat+command.Parent().Name()]

	if p == nil {
		common.CliExit(errors.New(fmt.Sprintf("Output format '%v' is not supported. Allowed values: %v", outputFormat, common.OutputFormatFlagAllowedValuesText)))
	}

	return p
}

func availablePrinters() map[string]util.Printer {
	return util.MergePrinterMaps(
		util.DefaultPrinters,
		map[string]util.Printer{
			common.OutputFormatPlain + common.ListCommandName:   listPlainPrinter{},
			common.OutputFormatTable + common.ListCommandName:   listTablePrinter{},
			common.OutputFormatPlain + common.GetCommandName:    getPlainPrinter{},
			common.OutputFormatPlain + common.CreateCommandName: createPlainPrinter{},
		},
	)
}

type listPlainPrinter struct{}
type listTablePrinter struct{}
type getPlainPrinter struct{}
type createPlainPrinter struct{}

func (p listPlainPrinter) Print(data interface{}) {
	listResponse, _ := (data).(*policies.ListPoliciesResponse)
	for _, policy := range listResponse.Policies {
		print1plain(policy)
	}
}

func (p listTablePrinter) Print(data interface{}) {
	listResponse, _ := (data).(*policies.ListPoliciesResponse)
	printTable(listResponse.Policies)
}

func (p getPlainPrinter) Print(data interface{}) {
	policy, _ := (data).(*v1.Policy)
	print1plain(policy)
}

func (p createPlainPrinter) Print(data interface{}) {
	policy, _ := (data).(*v1.Policy)
	print1plain(policy)
}

func print1plain(policy *v1.Policy) {
	fmt.Println("Name:", policy.Name)
	fmt.Println("Description:", policy.Description)
	fmt.Println("Retention(days):", policy.RetentionDays)
	fmt.Println("Legal Grounds:", policy.LegalGrounds)
	fmt.Println("State:", policy.State)
}

func printTable(policies []*v1.Policy) {
	rows := make([]table.Row, 0, len(policies))

	for _, policy := range policies {
		rows = append(rows, table.Row{
			policy.Name,
			policy.Description,
			policy.RetentionDays,
			policy.LegalGrounds,
			policy.State,
		})
	}

	util.RenderTable(
		table.Row{
			"Name",
			"Description",
			"Retention(days)",
			"Legal Grounds",
			"State",
		},
		rows,
	)
}
