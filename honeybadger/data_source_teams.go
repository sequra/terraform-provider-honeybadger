package honeybadger

import (
	"context"
	"strconv"
	"time"

	hbc "terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	teams, err := c.GetTeams()
	if err != nil {
		return diag.FromErr(err)
	}

	var unstructuredTeams []map[string]interface{}

	for _, team := range teams {
		var unstructuredUsers []map[string]interface{}
		unstructuredTeam := make(map[string]interface{})

		unstructuredTeam["name"] = team.Name
		unstructuredTeam["id"] = team.ID
		unstructuredTeam["created_at"] = team.CreatedAt
		for _, user := range team.Users {
			unstructuredUser := make(map[string]interface{})
			unstructuredUser["id"] = user.ID
			unstructuredUser["name"] = user.Name
			unstructuredUser["email"] = user.Email
			unstructuredUser["admin"] = user.IsAdmin

			unstructuredUsers = append(unstructuredUsers, unstructuredUser)
		}
		unstructuredTeam["users"] = unstructuredUsers

		unstructuredTeams = append(unstructuredTeams, unstructuredTeam)
	}

	if err := d.Set("teams", unstructuredTeams); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTeamsRead,
		Schema: map[string]*schema.Schema{
			"teams": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"users": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"email": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"admin": &schema.Schema{
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
