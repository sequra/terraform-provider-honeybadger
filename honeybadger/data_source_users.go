package honeybadger

import (
	"context"
	"strconv"
	"time"

	hbc "terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	users, err := c.GetUsers()
	if err != nil {
		return diag.FromErr(err)
	}

	var unstructuredUsers []map[string]interface{}

	for _, v := range users {
		user := make(map[string]interface{})

		user["user_id"] = v.ID
		user["name"] = v.Name
		user["admin"] = v.IsAdmin
		user["email"] = v.Email
		user["id"] = v.Email
		user["created_at"] = v.CreatedAt

		unstructuredUsers = append(unstructuredUsers, user)
	}

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
						"user_id": &schema.Schema{
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
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
