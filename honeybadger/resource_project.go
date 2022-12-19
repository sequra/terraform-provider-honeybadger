package honeybadger

import (
	"context"
	"strconv"
	"time"

	hbc "terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
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
			"language": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectName := d.Get("name").(string)
	hbProject, err := c.CreateProject(projectName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(hbProject.ID))
	d.Set("last_updated", time.Now().Format(time.RFC850))

	resourceProjectRead(ctx, d, m)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectID, _ := strconv.Atoi(d.Id())
	projectName := d.Get("name").(string)
	if projectID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Project not found",
			Detail:   "Project " + projectName + " with ID " + strconv.Itoa(projectID) + " not found in Honeybadger.",
		})
		return diags
	}

	if d.HasChange("name") {
		projectName := d.Get("name").(string)
		err := c.UpdateProject(projectName, projectID)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectID, _ := strconv.Atoi(d.Id())
	projectName := d.Get("name").(string)
	if projectID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Project not found",
			Detail:   "Project " + projectName + " with ID " + strconv.Itoa(projectID) + " not found in Honeybadger.",
		})
		return diags
	}

	err := c.DeleteProject(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectID, _ := strconv.Atoi(d.Id())
	project, err := c.FindProjectByID(projectID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Project not found",
			Detail:   "Project " + project.Name + " not found in Honeybadger.",
		})
		return diag.FromErr(err)
	}

	d.Set("name", project.Name)

	return diags
}
