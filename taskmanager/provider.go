package taskmanager

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager/datasources"
	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager/resources"
	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager_client"
)

type taskmanagerProvider struct {
	version string
	baseURL string
	token   string
	client  *taskmanager_client.Client
}

var tag string

func New() provider.Provider {
	const defaultVersion = "1.0"
	return &taskmanagerProvider{
		version: func() string {
			if tag == "" {
				return defaultVersion
			}
			return tag
		}(),
	}
}

func (tp *taskmanagerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "taskmanager"
	resp.Version = tp.version
}

func (tp *taskmanagerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Taskmanager Provider enables Terraform to interact eith Taskmanager REST API",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Required:    true,
				Description: "The tenant URL or domain",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`\S`), "must not be empty or whitespace"),
				},
			},
			"token": schema.StringAttribute{
				Required:    true,
				Description: "The authentication token used to connect with taskmanager API",
				Sensitive:   true,
			},
		},
	}
}

func (rp *taskmanagerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "COnfigure the taskmanager provider")

	var config taskmanager_client.TaskmanagerPlan
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := url.Parse(config.Tenant.ValueString())
	if err != nil || u.Scheme == "" || u.Host == "" {
		resp.Diagnostics.AddAttributeError(path.Root("tenant"), "Invalid Tenant URL", fmt.Sprintf("The provided tenant URL %q is not valid.", config.Tenant.ValueString()))
	}
	rp.baseURL = fmt.Sprintf("%s/api", config.Tenant.ValueString())

	if config.Token.IsNull() || config.Token.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing API Token",
			"You must provide an authentication token to connect to the Taskmanager API.",
		)
		return
	}
	rp.token = config.Token.ValueString()

	// Create API client
	client, err := taskmanager_client.NewClient(rp.baseURL, rp.token, rp.version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Taskmanager client",
			"Failed to authenticate or initialize the Taskmanager API client: "+err.Error(),
		)
		return
	}
	rp.client = client

	// Pass provider configuration to resources and data sources
	resp.ResourceData = client
	resp.DataSourceData = client
}

func (rp *taskmanagerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewDataSourceUser,
	}
}

func (rp *taskmanagerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewResourceUser,
	}
}
