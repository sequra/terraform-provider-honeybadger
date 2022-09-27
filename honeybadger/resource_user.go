package honeybadger

import (
	"context"
	"strconv"
	"time"

	hbc "terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"admin": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	teams := d.Get("team_id").([]interface{})
	userEmail := d.Get("email").(string)
	for _, team := range teams {
		err := c.CreateUser(userEmail, team.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(userEmail)
	d.Set("last_updated", time.Now().Format(time.RFC850))

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userEmail := d.Id()
	userID := d.Get("user_id").(int)
	if userID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not found",
			Detail:   "User " + userEmail + " with ID " + strconv.Itoa(userID) + " not found in Honeybadger.",
		})
		return diags
	}

	if d.HasChange("admin") || d.HasChange("team_id") {
		isAdmin := d.Get("admin").(bool)
		teams := d.Get("team_id").([]interface{})
		for _, team := range teams {
			err := c.UpdateUser(userID, isAdmin, team.(int))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userEmail := d.Id()
	teams := d.Get("team_id").([]interface{})
	userID := d.Get("user_id").(int)
	if userID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not found",
			Detail:   "User " + userEmail + " with ID " + strconv.Itoa(userID) + " not found in Honeybadger.",
		})
		return diags
	}

	for _, team := range teams {
		err := c.DeleteUser(userID, team.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userEmail := d.Id()
	teams := d.Get("team_id").([]interface{})
	for _, team := range teams {
		user, err := c.FindUserByEmail(userEmail, team.(int))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "User not found",
				Detail:   "User " + userEmail + " not found in Honeybadger.",
			})
			return diag.FromErr(err)
		}

		d.Set("name", user.Name)
		d.Set("email", user.Email)
		d.Set("admin", user.IsAdmin)
		d.Set("user_id", user.ID)

		return diags
	}
	return diags
}
