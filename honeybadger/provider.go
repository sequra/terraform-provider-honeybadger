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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"honeybadger_teams": dataSourceTeams(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"honeybadger_user":    resourceUser(),
			"honeybadger_team":    resourceTeam(),
			"honeybadger_project": resourceProject(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	authToken := d.Get("api_key").(string)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if authToken != "" {
		c := cli.NewClient(&host, &authToken)

		return c, diags
	}

	if authToken == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Honeybadger client",
			Detail:   "Honeybadger Client cannot be created because 'api_key' provider parameter is not defined",
		})
	}

	return nil, diags
}
