package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager/helpers/imports"
	"github.com/shivraj1192/terraform-provider-taskmanager-framework/taskmanager_client"
)

var (
	_ resource.Resource                = &ResourceUser{}
	_ resource.ResourceWithConfigure   = &ResourceUser{}
	_ resource.ResourceWithImportState = &ResourceUser{}
)

type ResourceUser struct {
	client       *taskmanager_client.Client
	helper       *ResourceUserHelper
	importHelper *imports.ImportHelper
}

type ResourceUserHelper struct{}

func NewResourceUser() resource.Resource {
	return &ResourceUser{}
}

func NewResourceUserHelper() *ResourceUserHelper {
	return &ResourceUserHelper{}
}

func (ru *ResourceUser) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "taskmanager_user"
}

func (ru *ResourceUser) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	tflog.Info(ctx, "Configuring taskmanager_user resource")

	if req.ProviderData == nil {
		return
	}

	ru.client = req.ProviderData.(*taskmanager_client.Client)
	if ru.client == nil {
		resp.Diagnostics.AddError(
			"Provider client error",
			"REST API client is not configured",
		)
		tflog.Error(ctx, "Provider client is nil after Configure", map[string]interface{}{
			"method": "Configure",
		})
		return
	}

	ru.helper = NewResourceUserHelper()
	tflog.Info(ctx, "taskmanager_user resource configured")
}

func (ru *ResourceUser) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "taskmanager_user resource schema",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Unique identifier for the resource",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"uname": schema.StringAttribute{
				Required:    true,
				Description: "Username of user",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of user",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Email of user",
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"role": schema.StringAttribute{
				Optional:    true,
				Description: "Role of user",
			},
		},
	}
}

func (ru *ResourceUser) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "Creating taskmanager_user resource")

	var plan taskmanager_client.ResourceUserPlan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to read plan during resource user creation", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	user := taskmanager_client.User{
		Uname:    plan.Uname.ValueString(),
		Name:     plan.Name.ValueString(),
		Email:    plan.Email.ValueString(),
		Password: plan.Password.ValueString(),
	}

	if !plan.Role.IsNull() && !plan.Role.IsUnknown() {
		user.Role = plan.Role.ValueString()
	}

	userModel, err := ru.client.CreateUser(user)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create user", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to create user, error:%#v", err))
		return
	}

	user = userModel.User
	plan.ID = types.Int64Value(int64(user.ID))

	planPtr, err := ru.helper.getAndMapModelToPlan(ctx, plan, ru.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to set state after create",
			fmt.Sprintf("Error: %v", err),
		)
		tflog.Error(ctx, "Failed get and map user model to plan", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, planPtr)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to set state after create", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}
}

func (ru *ResourceUser) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading taskmanager_user resource")

	var state taskmanager_client.ResourceUserPlan
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Read failed to get resource user state", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	userID := state.ID.ValueInt64()
	user, err := ru.client.GetUser(userID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read user resource", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to read user resource, error:%#v", err))
		return
	}

	state.ID = types.Int64Value(int64(user.ID))

	planPtr, err := ru.helper.getAndMapModelToPlan(ctx, state, ru.client)
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

func (ru *ResourceUser) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "Update called for taskmanager_user")

	var plan, state taskmanager_client.ResourceUserPlan
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Update failed to get plan/state", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	user := taskmanager_client.User{
		Uname:    plan.Uname.ValueString(),
		Name:     plan.Name.ValueString(),
		Email:    plan.Email.ValueString(),
		Password: plan.Password.ValueString(),
	}

	userModel, err := ru.client.UpdateUser(user)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update user", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to update user, error:%#v", err))
		return
	}

	user = userModel.User
	plan.ID = types.Int64Value(int64(user.ID))

	planPtr, err := ru.helper.getAndMapModelToPlan(ctx, plan, ru.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to set state after update",
			fmt.Sprintf("Error: %v", err),
		)
		tflog.Error(ctx, "Failed get and map user model to plan", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, planPtr)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to set state after update", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}
}

func (ru *ResourceUser) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "delete called for taskmanager_user")

	var state taskmanager_client.ResourceUserPlan
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Delete failed to get state", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	err := ru.client.DeleteUser(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete user", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to delete user, error:%#v", err))
		return
	}
}

func (ru *ResourceUser) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importId := req.ID

	importData := &imports.ImportHelperData{
		ID: importId,
	}

	if err := ru.importHelper.ParseImportID([]string{"(?P<name>[^/]+)"}, importData); err != nil {
		resp.Diagnostics.AddError("Failed to parse import ID", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to parse import ID, error:%#v", err))
		return
	}
	userIDStr := importData.Fields["userID"]
	if strings.TrimSpace(userIDStr) == "" {
		resp.Diagnostics.AddError("Failed to import user", "Invalid userID")
		tflog.Error(ctx, "Failed to import user, Invalid userID")
		return
	}
	id, err := strconv.Atoi(userIDStr)
	userID := int64(id)

	tflog.Info(ctx, fmt.Sprintf("Importing user: %d", userID))

	user, err := ru.client.GetUser(userID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to import user", err.Error())
		tflog.Error(ctx, fmt.Sprintf("Failed to import user, error:%#v", err))
		return
	}

	plan := taskmanager_client.ResourceUserPlan{
		ID: types.Int64Value(int64(user.ID)),
	}

	planPtr, err := ru.helper.getAndMapModelToPlan(ctx, plan, ru.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to set state after import",
			fmt.Sprintf("Error: %v", err),
		)
		tflog.Error(ctx, "Failed get and map user model to plan", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, planPtr)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to set state after import", map[string]interface{}{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Imported user: %#v", planPtr))
}

func (ruh *ResourceUserHelper) getAndMapModelToPlan(ctx context.Context, plan taskmanager_client.ResourceUserPlan, c *taskmanager_client.Client) (*taskmanager_client.ResourceUserPlan, error) {
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
