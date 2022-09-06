package honeybadger

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	hbc "terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneyBadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	users, err := c.GetUsers()
	if err != nil {
		return diag.FromErr(err)
	}

	var unstructuredUsers []map[string]interface{}
	inrec, _ := json.Marshal(users)
	json.Unmarshal(inrec, &unstructuredUsers)

	if err := d.Set("users", unstructuredUsers); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUsersRead,
		Schema: map[string]*schema.Schema{
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": &schema.Schema{
							Type:     schema.TypeBool,
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
					},
				},
			},
		},
	}
}
