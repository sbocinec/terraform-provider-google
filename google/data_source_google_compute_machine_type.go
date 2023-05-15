package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/compute/v1"
)

func dataSourceGoogleComputeMachineType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeMachineTypeRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeMachineTypeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return fmt.Errorf("Please specify zone to get appropriate machine types for zone. Unable to get zone: %s", err)
	}

	m, ok := d.GetOk("machine_type")
	if !ok {
		return fmt.Errorf("Please specify machine_type to get machine type details")
	}

	machineType, err := config.NewComputeClient(userAgent).MachineTypes.Get(project, zone, m.(string)).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Machine Type %s", m.(string)))
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("name", machineType.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", machineType.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("guest_cpus", machineType.GuestCpus); err != nil {
		return fmt.Errorf("Error setting guest_cpus: %s", err)
	}
	if err := d.Set("memory_mb", machineType.MemoryMb); err != nil {
		return fmt.Errorf("Error setting memory_mb: %s", err)
	}
	if err := d.Set("image_space_gb", machineType.ImageSpaceGb); err != nil {
		return fmt.Errorf("Error setting image_space_gb: %s", err)
	}
	if err := d.Set("scratch_disks", flattenMachineTypeScratchDisks(machineType.ScratchDisks)); err != nil {
		return fmt.Errorf("Error setting scratch_disks: %s", err)
	}
	if err := d.Set("maximum_persistent_disks", machineType.MaximumPersistentDisks); err != nil {
		return fmt.Errorf("Error setting maxium_persistent_disks: %s", err)
	}
	if err := d.Set("maximum_persistent_disks_size_gb", machineType.MaximumPersistentDisksSizeGb); err != nil {
		return fmt.Errorf("Error setting maxium_persistent_disks_size_gb: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("self_link", tpgresource.ConvertSelfLinkToV1(machineType.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("is_shared_cpu", machineType.IsSharedCpu); err != nil {
		return fmt.Errorf("Error setting is_shared_cpu: %s", err)
	}
	if err := d.Set("accelerators", flattenMachineTypeAccelerators(machineType.Accelerators)); err != nil {
		return fmt.Errorf("Error setting accelerators: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/machineTypes/%s", project, zone, machineType.Name))

	return nil
}

func flattenMachineTypeScratchDisks(disks []*compute.MachineTypeScratchDisks) []map[string]interface{} {
	disksSchema := make([]map[string]interface{}, len(disks))
	for i, disk := range disks {
		disksSchema[i] = map[string]interface{}{
			"disk_gb": disk.DiskGb,
		}
	}
	return disksSchema
}

func flattenMachineTypeAccelerators(accelerators []*compute.MachineTypeAccelerators) []map[string]interface{} {
	acceleratorsSchema := make([]map[string]interface{}, len(accelerators))
	for i, accelerator := range accelerators {
		acceleratorsSchema[i] = map[string]interface{}{
			"guest_accelerator_count": accelerator.GuestAcceleratorCount,
			"guest_accelerator_type":  accelerator.GuestAcceleratorType,
		}
	}
	return acceleratorsSchema
}
