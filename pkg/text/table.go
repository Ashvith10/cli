package text

import (
	"github.com/jedib0t/go-pretty/table"

	"github.com/renderinc/cli/pkg/client"
	clientjob "github.com/renderinc/cli/pkg/client/jobs"
	"github.com/renderinc/cli/pkg/deploy"
	"github.com/renderinc/cli/pkg/resource"
)

func ResourceTable(v []resource.Resource) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Project", "Environment", "Type", "ID"})
	for _, r := range v {
		t.AppendRow(table.Row{r.Name(), r.ProjectName(), r.EnvironmentName(), r.Type(), r.ID()})
	}
	return FormatString(t.Render())
}

func JobTable(v []*clientjob.Job) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Command", "Started", "Finished", "Plan", "ID"})
	for _, r := range v {
		t.AppendRow(table.Row{r.StartCommand, r.StartedAt, r.FinishedAt, r.PlanId, r.Id})
	}
	return FormatString(t.Render())
}

func DeployTable(v []*client.Deploy) string {
	t := table.NewWriter()
	t.AppendHeader(toRow(deploy.Header()))
	for _, r := range v {
		t.AppendRow(toRow(deploy.Row(r)))
	}
	return FormatString(t.Render())
}

func ProjectTable(v []*client.Project) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "ID"})
	for _, r := range v {
		t.AppendRow(table.Row{r.Name, r.Id})
	}
	return FormatString(t.Render())
}

func EnvironmentTable(v []*client.Environment) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Protected", "ID"})
	for _, r := range v {
		t.AppendRow(table.Row{r.Name, r.ProtectedStatus, r.Id})
	}
	return FormatString(t.Render())
}

func toRow(r []string) table.Row {
	row := table.Row{}
	for _, r := range r {
		row = append(row, r)
	}

	return row
}
