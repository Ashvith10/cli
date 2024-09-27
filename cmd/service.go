package cmd

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/renderinc/render-cli/pkg/client"
	"github.com/renderinc/render-cli/pkg/command"
	"github.com/renderinc/render-cli/pkg/environment"
	"github.com/renderinc/render-cli/pkg/project"
	"github.com/renderinc/render-cli/pkg/resource"
	"github.com/renderinc/render-cli/pkg/service"
	"github.com/renderinc/render-cli/pkg/tui"
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "List and manage services",
	RunE: func(cmd *cobra.Command, args []string) error {
		command.Wrap(cmd, loadResourceData, renderResources)(cmd.Context(), ListResourceInput{})
		return nil
	},
}

func loadResourceData(ctx context.Context, _ ListResourceInput) ([]resource.Resource, error) {
	resourceService, err := newResourceService()
	if err != nil {
		return nil, err
	}
	return resourceService.ListResources(ctx)
}

type ListResourceInput struct{}

func (l ListResourceInput) String() []string {
	return []string{}
}

func renderResources(ctx context.Context, loadData func(input ListResourceInput) ([]resource.Resource, error), in ListResourceInput) (tea.Model, error) {
	columns := []table.Column{
		{Title: "Type", Width: 12},
		{Title: "Project", Width: 15},
		{Title: "Environment", Width: 20},
		{Title: "Name", Width: 40},
		{Title: "ID", Width: 25},
	}

	return tui.NewTableModel[resource.Resource](
		"resources",
		func() ([]resource.Resource, error) {
			return loadData(in)
		},
		formatResourceRow,
		selectResource(ctx),
		columns,
		filterResource,
		[]tui.CustomOption[resource.Resource]{
			{
				Key:      "w",
				Title:    "Change Workspace",
				Function: resourceOptionSelectWorkspace(ctx),
			},
		},
	), nil
}

func formatResourceRow(r resource.Resource) table.Row {
	return []string{r.Type(), r.ProjectName(), r.EnvironmentName(), r.Name(), r.ID()}
}

func selectResource(ctx context.Context) func(resource.Resource) tea.Cmd {
	return func(r resource.Resource) tea.Cmd {
		return InteractiveCommandPalette(ctx, PaletteCommandInput{
			Commands: []PaletteCommand{
				{
					Name:        "logs",
					Description: "View resource logs",
					Action: func(ctx context.Context, args []string) tea.Cmd {
						return InteractiveLogs(ctx, LogInput{ResourceIDs: []string{r.ID()}})
					},
				},
				{
					Name:        "restart",
					Description: "Restart the service",
					Action: func(ctx context.Context, args []string) tea.Cmd {
						return InteractiveRestart(ctx, RestartInput{ServiceID: r.ID()})
					},
				},
			},
		})
	}
}

func filterResource(r resource.Resource, filter string) bool {
	searchFields := []string{r.ID(), r.Name(), r.ProjectName(), r.EnvironmentName(), r.Type()}
	for _, field := range searchFields {
		if strings.Contains(strings.ToLower(field), filter) {
			return true
		}
	}
	return false
}

func newResourceService() (*resource.Service, error) {
	httpClient := http.DefaultClient
	host := os.Getenv("RENDER_HOST")
	apiKey := os.Getenv("RENDER_API_KEY")

	c, err := client.ClientWithAuth(httpClient, host, apiKey)
	if err != nil {
		return nil, err
	}

	serviceRepo := service.NewRepo(c)
	environmentRepo := environment.NewRepo(c)
	projectRepo := project.NewRepo(c)

	serviceService := service.NewService(serviceRepo, environmentRepo, projectRepo)
	resourceService := resource.NewResourceService(serviceService, environmentRepo, projectRepo)

	return resourceService, nil
}

func resourceOptionSelectWorkspace(ctx context.Context) func(resource.Resource) tea.Cmd {
	return func(r resource.Resource) tea.Cmd {
		return InteractiveWorkspace(ctx, ListWorkspaceInput{})
	}
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}
