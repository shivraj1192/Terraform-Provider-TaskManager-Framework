package taskmanager_client

import "github.com/hashicorp/terraform-plugin-framework/types"

type TaskmanagerPlan struct {
	Tenant types.String `tfsdk:"base_url"`
	Token  types.String `tfsdk:"token"`
}

type DataSourceUserPlan struct {
	ID    types.Int64  `tfsdk:"id"`
	Uname types.String `tfsdk:"uname"`
	Name  types.String `tfsdk:"name"`
	Email types.String `tfsdk:"email"`
	Role  types.String `tfsdk:"role"`
}

type ResourceUserPlan struct {
	ID       types.Int64  `tfsdk:"id"`
	Uname    types.String `tfsdk:"uname"`
	Name     types.String `tfsdk:"name"`
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`
	Role     types.String `tfsdk:"role"`
}
