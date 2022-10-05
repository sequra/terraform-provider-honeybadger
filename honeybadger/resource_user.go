package honeybadger

import (
	"context"
	"log"
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
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"team": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_admin": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"user_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
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

	userEmail := d.Get("email").(string)
	teams := d.Get("team")
	for _, item := range teams.(*schema.Set).List() {
		team := item.(map[string]interface{})
		teamID := team["id"].(int)
		isAdmin := team["is_admin"].(bool)
		err := c.CreateUser(userEmail, isAdmin, teamID)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("User " + userEmail + " will be inserted into team  " + strconv.Itoa(teamID))
	}

	d.SetId(userEmail)
	d.Set("last_updated", time.Now().Format(time.RFC850))

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	userEmail := d.Id()
	if d.HasChange("team") {
		err := updateUserTeam(userEmail, d, m)
		if err != nil {
			return diag.FromErr(err)
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
	teams := d.Get("team")
	for _, item := range teams.(*schema.Set).List() {
		team := item.(map[string]interface{})
		teamID := team["id"].(int)
		userID := team["user_id"].(int)
		err := c.DeleteUser(userID, teamID)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("User " + userEmail + " will be deleted from team  " + strconv.Itoa(teamID))
	}

	d.SetId("")

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*hbc.HoneybadgerClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var unstructuredUserTeams []map[string]interface{}
	userEmail := d.Id()
	userTeams, err := c.GetUserFromTeams(userEmail)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("Reading user with email %s", userEmail)
	for _, user := range userTeams {
		log.Printf("Found user with email %s in team %d", userEmail, user.TeamID)
		unstructuredUserTeams = append(unstructuredUserTeams, map[string]interface{}{
			"is_admin": user.IsAdmin,
			"user_id":  user.ID,
			"id":       user.TeamID,
		})

	}

	d.Set("email", userEmail)
	if err := d.Set("team", unstructuredUserTeams); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateUserTeam(userEmail string, d *schema.ResourceData, m interface{}) error {
	c := m.(*hbc.HoneybadgerClient)
	oldState, newState := d.GetChange("team")

	if oldState == nil {
		oldState = new(schema.Set)
	}
	if newState == nil {
		newState = new(schema.Set)
	}

	os := oldState.(*schema.Set)
	ns := newState.(*schema.Set)
	updateOperation := calculateRealDifference(os, ns).List()
	removeOperation := os.Difference(ns).List()
	addOperation := ns.Difference(os).List()

	// Delete user from team
	for _, operation := range removeOperation {
		team := operation.(map[string]interface{})
		teamID := team["id"].(int)
		userID := team["user_id"].(int)
		log.Printf("User %s with id %d it will be deleted from team %d", userEmail, userID, teamID)
		err := c.DeleteUser(userID, teamID)
		if err != nil {
			return err
		}
	}
	// Add user to a team
	for _, operation := range addOperation {
		team := operation.(map[string]interface{})
		teamID := team["id"].(int)
		isAdmin := team["is_admin"].(bool)
		log.Printf("User %s it will be invited to team %d with admin to %t", userEmail, teamID, isAdmin)
		err := c.CreateUser(userEmail, isAdmin, teamID)
		if err != nil {
			return err
		}
	}

	// Update user in a team
	for _, operation := range updateOperation {
		team := operation.(map[string]interface{})
		teamID := team["id"].(int)
		isAdmin := team["is_admin"].(bool)
		user, _ := c.GetUserForTeam(userEmail, teamID) //userID is 0 because new state could not preserve the ID
		log.Printf("User %s with ID %d it will be updated in team %d with admin value %t", userEmail, user.ID, teamID, isAdmin)
		err := c.UpdateUser(user.ID, isAdmin, teamID)
		if err != nil {
			return err
		}
	}

	return nil
}

func calculateRealDifference(oldState *schema.Set, newState *schema.Set) *schema.Set {
	updateOperation := schema.NewSet(oldState.F, []interface{}{})

	for _, os := range oldState.List() {
		teamRem := os.(map[string]interface{})
		teamIDRem := teamRem["id"].(int)
		for _, ns := range newState.List() {
			teamAdd := ns.(map[string]interface{})
			teamIDAdd := teamAdd["id"].(int)
			if teamIDAdd == teamIDRem { // It is s an upate operation
				updateOperation.Add(ns)
				oldState.Remove(os)
				newState.Remove(ns)
			}
		}
	}
	return updateOperation
}
