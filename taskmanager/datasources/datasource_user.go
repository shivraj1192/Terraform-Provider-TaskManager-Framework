package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager_client"
)

var (
	_ datasource.DataSource              = &DataSourceUser{}
	_ datasource.DataSourceWithConfigure = &DataSourceUser{}
)

type DataSourceUser struct {
	client *taskmanager_client.Client
	helper *ResourceUserHelper
}

type ResourceUserHelper struct{}

func NewDataSourceUser() datasource.DataSource {
	return &DataSourceUser{}
}

func NewResourceUserHelper() *ResourceUserHelper {
	return &ResourceUserHelper{}
}

func (du *DataSourceUser) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "taskmanager_user"
}

func (du *DataSourceUser) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*taskmanager_client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Provider client error",
			"REST API client is not configured",
		)
		tflog.Error(ctx, "Provider client is nil after Configure", map[string]interface{}{
			"method": "Configure",
		})
		return
	}

	du.client = client
	tflog.Info(ctx, "Configured DataSourceUser with Taskmanager client")
}

func (du *DataSourceUser) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Datasource for retrieving Taskmanager user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required:    true,
				Description: "Unique identifier for the resource",
			},
			"uname": schema.StringAttribute{
				Computed:    true,
				Description: "Username of user",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Name of user",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "Email of user",
			},
			"role": schema.StringAttribute{
				Computed:    true,
				Description: "Role of user",
			},
		},
	}
}

func (du *DataSourceUser) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading taskmanager_user resource")

	var state taskmanager_client.DataSourceUserPlan
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Read failed to get datasource user state", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	userID := state.ID.ValueInt64()
	user, err := du.client.GetUser(userID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user resource", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to read user resource, error:%#v", err))
		return
	}

	state.ID = types.Int64Value(int64(user.ID))

	planPtr, err := du.helper.getAndMapModelToPlan(ctx, state, du.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to set state after read",
			fmt.Sprintf("Error: %v", err),
		)
		tflog.Error(ctx, "Failed get and map user model to plan", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, planPtr)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to set state after read", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

}

func (ruh *ResourceUserHelper) getAndMapModelToPlan(ctx context.Context, plan taskmanager_client.DataSourceUserPlan, c *taskmanager_client.Client) (*taskmanager_client.DataSourceUserPlan, error) {
	userID := plan.ID.ValueInt64()

	tflog.Info(ctx, fmt.Sprintf("Reading role %s", userID))

	user, err := c.GetUser(userID)
	if err != nil {
		return nil, err
	}

	tflog.Info(ctx, fmt.Sprintf("Received user %#v", user))

	plan.Uname = types.StringValue(user.Uname)
	plan.Name = types.StringValue(user.Name)
	plan.Email = types.StringValue(user.Email)
	if (plan.Role.IsNull() || plan.Role.IsUnknown()) && user.Role == "" {
		plan.Role = types.StringNull()
	} else {
		plan.Role = types.StringValue(user.Role)
	}

	return &plan, nil
}
