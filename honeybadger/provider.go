package honeybadger

import (
	"context"
	"terraform-provider-honeybadger/cli"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HONEYBADGER_HOST", nil),
			},
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HONEYBADGER_API_KEY", nil),
			},
			"team_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HONEYBADGER_TEAM_ID", nil),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"honeybadger_users": dataSourceUsers(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"honeybadger_user": resourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	authToken := d.Get("api_key").(string)
	teamID := d.Get("team_id").(int)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (authToken != "") && (teamID != 0) {
		c := cli.NewClient(&host, &authToken, &teamID)

		return c, diags
	}

	if authToken == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Honeybadger client",
			Detail:   "Honeybadger Client cannot be created because 'api_key' provider parameter is not defined",
		})
	}

	if teamID == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Honeybadger client",
			Detail:   "Honeybadger Client cannot be created because 'team_id' provider parameter is not defined",
		})
	}

	return nil, diags
}
