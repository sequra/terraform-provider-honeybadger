package honeybadger

import (
	"context"
	"strconv"
	"time"

	hbc "terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
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
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	teamName := d.Get("name").(string)
	hbTeam, err := c.CreateTeam(teamName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(hbTeam.ID))
	d.Set("last_updated", time.Now().Format(time.RFC850))

	resourceTeamRead(ctx, d, m)

	return diags
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	teamID, _ := strconv.Atoi(d.Id())
	teamName := d.Get("name").(string)
	if teamID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Team not found",
			Detail:   "Team " + teamName + " with ID " + strconv.Itoa(teamID) + " not found in Honeybadger.",
		})
		return diags
	}

	if d.HasChange("name") {
		teamName := d.Get("name").(string)
		err := c.UpdateTeam(teamName, teamID)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceTeamRead(ctx, d, m)
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	teamID, _ := strconv.Atoi(d.Id())
	teamName := d.Get("name").(string)
	if teamID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Team not found",
			Detail:   "Team " + teamName + " with ID " + strconv.Itoa(teamID) + " not found in Honeybadger.",
		})
		return diags
	}

	err := c.DeleteTeam(teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	teamID, _ := strconv.Atoi(d.Id())
	team, err := c.FindTeamByID(teamID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Team not found",
			Detail:   "Team " + team.Name + " not found in Honeybadger.",
		})
		return diag.FromErr(err)
	}

	d.Set("name", team.Name)

	return diags
}
