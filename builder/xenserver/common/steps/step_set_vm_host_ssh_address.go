package steps

import (
	"context"
	"fmt"
	"github.com/xenserver/packer-builder-xenserver/builder/xenserver/common/xen"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepSetVmHostSshAddress struct{}

func (self *StepSetVmHostSshAddress) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {

	c := state.Get("client").(*xen.Connection)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Step: Set SSH address to VM host IP")

	uuid := state.Get("instance_uuid").(string)
	instance, err := c.GetClient().VM.GetByUUID(c.GetSessionRef(), uuid)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to get VM from UUID '%s': %s", uuid, err.Error()))
		return multistep.ActionHalt
	}

	host, err := c.GetClient().VM.GetResidentOn(c.GetSessionRef(), instance)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to get VM Host for VM '%s': %s", uuid, err.Error()))
	}

	address, err := c.GetClient().Host.GetAddress(c.GetSessionRef(), host)
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to get address from VM Host: %s", err.Error()))
	}

	state.Put("vm_host_address", address)
	ui.Say(fmt.Sprintf("Set host SSH address to '%s'.", address))

	return multistep.ActionContinue
}

func (self *StepSetVmHostSshAddress) Cleanup(state multistep.StateBag) {}