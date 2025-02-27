package vultr

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vultr/govultr/v2"
)

func dataSourceVultrInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVultrInstancesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"main_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpu_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"date_created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allowed_bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"netmask_v4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_v4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"power_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"v6_network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"v6_main_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"v6_network_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kvm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backups": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"os_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"features": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"backups_schedule": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_network_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"vpc_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceVultrInstancesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client).govultrClient()

	filters, filtersOk := d.GetOk("filter")

	if !filtersOk {
		return diag.Errorf("issue with filter: %v", filtersOk)
	}

	var serverList []govultr.Instance
	f := buildVultrDataSourceFilter(filters.(*schema.Set))
	options := &govultr.ListOptions{}
	for {
		servers, meta, err := client.Instance.List(ctx, options)
		if err != nil {
			return diag.Errorf("error getting servers: %v", err)
		}

		for _, s := range servers {
			// we need convert the a struct INTO a map so we can easily manipulate the data here
			sm, err := structToMap(s)

			if err != nil {
				return diag.FromErr(err)
			}

			if filterLoop(f, sm) {
				serverList = append(serverList, s)
			}
		}

		if meta.Links.Next == "" {
			break
		} else {
			options.Cursor = meta.Links.Next
			continue
		}
	}

	if len(serverList) < 1 {
		return diag.Errorf("no results were found")
	}

	serverDetails := make([]interface{}, 0)
	for _, server := range serverList {
		schedule, err := client.Instance.GetBackupSchedule(ctx, server.ID)
		if err != nil {
			return diag.Errorf("error getting backup schedule: %v", err)
		}

		bsInfo := map[string]interface{}{
			"type": schedule.Type,
			"hour": strconv.Itoa(schedule.Hour),
			"dom":  strconv.Itoa(schedule.Dom),
			"dow":  strconv.Itoa(schedule.Dow),
		}

		vpcs, err := getVPCs(client, server.ID)
		if err != nil {
			return diag.Errorf(err.Error())
		}

		serverDetails = append(serverDetails, map[string]interface{}{
			"os":                  server.Os,
			"ram":                 server.RAM,
			"disk":                server.Disk,
			"main_ip":             server.MainIP,
			"vcpu_count":          server.VCPUCount,
			"region":              server.Region,
			"date_created":        server.DateCreated,
			"allowed_bandwidth":   server.AllowedBandwidth,
			"netmask_v4":          server.NetmaskV4,
			"gateway_v4":          server.GatewayV4,
			"status":              server.Status,
			"power_status":        server.PowerStatus,
			"server_status":       server.ServerStatus,
			"plan":                server.Plan,
			"label":               server.Label,
			"internal_ip":         server.InternalIP,
			"kvm":                 server.KVM,
			"tag":                 server.Tag,
			"tags":                server.Tags,
			"os_id":               server.OsID,
			"app_id":              server.AppID,
			"image_id":            server.ImageID,
			"firewall_group_id":   server.FirewallGroupID,
			"v6_network":          server.V6Network,
			"v6_main_ip":          server.V6MainIP,
			"v6_network_size":     server.V6NetworkSize,
			"features":            server.Features,
			"hostname":            server.Hostname,
			"backups":             backupStatus(schedule.Enabled),
			"backups_schedule":    bsInfo,
			"private_network_ids": vpcs,
			"vpc_ids":             vpcs,
		})
	}

	d.SetId("instances")
	d.Set("instances", serverDetails)

	return nil
}
