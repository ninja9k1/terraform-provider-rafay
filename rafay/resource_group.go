package rafay

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/RafaySystems/rctl/pkg/config"
	"github.com/RafaySystems/rctl/pkg/group"
	"github.com/RafaySystems/rctl/pkg/models"

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

	log.Printf("resource greoup create %s", d.Get("name").(string))
	err := group.CreateGroup(d.Get("name").(string), d.Get("description").(string))
	if err != nil {
		log.Printf("create group error %s", err.Error())
		return diag.FromErr(err)
	}

	resp, err := group.GetGroupByName(d.Get("name").(string))
	if err != nil {
		log.Printf("create group failed to get group, error %s", err.Error())
		return diag.FromErr(err)
	}

	g, err := group.NewGroupFromResponse([]byte(resp))
	if err != nil {
		log.Printf("create group failed to parse get response, error %s", err.Error())
		return diag.FromErr(err)
	} else if g == nil {
		log.Printf("create group failed to parse get response")
		d.SetId("")
		return diags
	}

	log.Printf("resource group created %s", g.ID)
	d.SetId(g.ID)

	return diags
}

func getGroupById(id string) (string, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := "/auth/v1/groups/"
	uri = uri + fmt.Sprintf("%s/", id)
	return auth.AuthAndRequest(uri, "GET", nil)
}

func getGroupFromResponse(json_data []byte) (*models.Group, error) {
	var gr models.Group
	if err := json.Unmarshal(json_data, &gr); err != nil {
		return nil, err
	}
	return &gr, nil
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	//resp, err := project.GetProjectByName(d.Get("name").(string))
	resp, err := getGroupById(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	p, err := getGroupFromResponse([]byte(resp))
	if err != nil {
		log.Printf("get group by id, error %s", err.Error())
		return diag.FromErr(err)
	} else if p == nil {
		log.Printf("get group response parse error")
		d.SetId("")
		return diags
	}

	if err := d.Set("name", p.Name); err != nil {
		log.Printf("get group set name error %s", err.Error())
		return diag.FromErr(err)
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	//TODO implement update project
	var diags diag.Diagnostics
	log.Printf("resource group update id %s", d.Id())
	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	log.Printf("resource group delete id %s", d.Id())
	err := group.DeleteGroupById(d.Id())
	if err != nil {
		log.Printf("delete group error %s", err.Error())
		return diag.FromErr(err)
	}

	return diags
}
