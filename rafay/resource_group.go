package rafay

import (
	"context"
	"fmt"
	"time"

	"github.com/RafaySystems/rctl/pkg/config"
	"github.com/RafaySystems/rctl/pkg/group"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	err := group.CreateGroup(d.Get("name").(string), d.Get("description").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := group.GetGroupByName(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	p, err := group.NewGroupFromResponse([]byte(resp))
	if err != nil {
		return diag.FromErr(err)
	} else if p == nil {
		d.SetId("")
		return diags
	}

	d.SetId(p.ID)

	return diags
}

func getGroupById(id string) (string, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := "/auth/v1/groups/"
	uri = uri + fmt.Sprintf("%s/", id)
	return auth.AuthAndRequest(uri, "GET", nil)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	//resp, err := project.GetProjectByName(d.Get("name").(string))
	resp, err := getGroupById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	p, err := group.NewGroupFromResponse([]byte(resp))
	if err != nil {
		return diag.FromErr(err)
	} else if p == nil {
		d.SetId("")
		return diags
	}

	if err := d.Set("name", p.Name); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO implement update project
	var diags diag.Diagnostics
	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	err := group.DeleteGroupById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
