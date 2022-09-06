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
				Required: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"admin": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneyBadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userEmail := d.Get("email").(string)
	userID, err := c.CreateUser(userEmail)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(userID))
	d.Set("last_updated", time.Now().Format(time.RFC850))

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneyBadgerClient)

	userID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("admin") {
		isAdmin := d.Get("admin").(bool)
		err := c.UpdateUser(userID, isAdmin)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneyBadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteUser(userID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneyBadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userId, err := strconv.Atoi(d.Id())
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get usr ID from Terraform State",
			Detail:   "Unable to get usr ID from Terraform State",
		})

		return diag.FromErr(err)
	}

	user, err := c.FindUserByID(userId)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not found",
			Detail:   "User " + d.Id() + " not found in HoneyBadger.",
		})
		return diag.FromErr(err)
	}

	d.Set("name", user.Name)
	d.Set("email", user.Email)
	d.Set("admin", user.IsAdmin)

	return diags
}
