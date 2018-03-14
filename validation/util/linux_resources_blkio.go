package util

import (
	"fmt"

	"github.com/mndrix/tap-go"
	rspec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/opencontainers/runtime-tools/cgroups"
)

// ValidateLinuxResourcesBlockIO validates linux.resources.blockIO.
func ValidateLinuxResourcesBlockIO(config *rspec.Spec, state *rspec.State) error {
	t := tap.New()
	t.Header(0)

	cg, err := cgroups.FindCgroup()
	t.Ok((err == nil), "find blkio cgroup")
	if err != nil {
		t.Diagnostic(err.Error())
		t.AutoPlan()
		return nil
	}

	lbd, err := cg.GetBlockIOData(state.Pid, config.Linux.CgroupsPath)
	t.Ok((err == nil), "get blkio cgroup data")
	if err != nil {
		t.Diagnostic(err.Error())
		t.AutoPlan()
		return nil
	}

	t.Ok(*lbd.Weight == *config.Linux.Resources.BlockIO.Weight, "blkio weight is set correctly")
	t.Diagnosticf("expect: %d, actual: %d", *config.Linux.Resources.BlockIO.Weight, *lbd.Weight)

	t.Ok(*lbd.LeafWeight == *config.Linux.Resources.BlockIO.LeafWeight, "blkio leafWeight is set correctly")
	t.Diagnosticf("expect: %d, actual: %d", *config.Linux.Resources.BlockIO.LeafWeight, *lbd.LeafWeight)

	for _, device := range config.Linux.Resources.BlockIO.WeightDevice {
		found := false
		for _, wd := range lbd.WeightDevice {
			if wd.Major == device.Major && wd.Minor == device.Minor {
				found = true
				t.Ok(*wd.Weight == *device.Weight, fmt.Sprintf("blkio weight for %d:%d is set correctly", device.Major, device.Minor))
				t.Diagnosticf("expect: %d, actual: %d", *device.Weight, *wd.Weight)

				t.Ok(*wd.LeafWeight == *device.LeafWeight, fmt.Sprintf("blkio leafWeight for %d:%d is set correctly", device.Major, device.Minor))
				t.Diagnosticf("expect: %d, actual: %d", *device.LeafWeight, *wd.LeafWeight)
			}
		}
		t.Ok(found, fmt.Sprintf("blkio weightDevice for %d:%d found", device.Major, device.Minor))
	}

	for _, device := range config.Linux.Resources.BlockIO.ThrottleReadBpsDevice {
		found := false
		for _, trbd := range lbd.ThrottleReadBpsDevice {
			if trbd.Major == device.Major && trbd.Minor == device.Minor {
				found = true
				t.Ok(trbd.Rate == device.Rate, fmt.Sprintf("blkio read bps for %d:%d is set correctly", device.Major, device.Minor))
				t.Diagnosticf("expect: %d, actual: %d", device.Rate, trbd.Rate)
			}
		}
		t.Ok(found, fmt.Sprintf("blkio read bps for %d:%d found", device.Major, device.Minor))
	}

	for _, device := range config.Linux.Resources.BlockIO.ThrottleWriteBpsDevice {
		found := false
		for _, twbd := range lbd.ThrottleWriteBpsDevice {
			if twbd.Major == device.Major && twbd.Minor == device.Minor {
				found = true
				t.Ok(twbd.Rate == device.Rate, fmt.Sprintf("blkio write bps for %d:%d is set correctly", device.Major, device.Minor))
				t.Diagnosticf("expect: %d, actual: %d", device.Rate, twbd.Rate)
			}
		}
		t.Ok(found, fmt.Sprintf("blkio write bps for %d:%d found", device.Major, device.Minor))
	}

	for _, device := range config.Linux.Resources.BlockIO.ThrottleReadIOPSDevice {
		found := false
		for _, trid := range lbd.ThrottleReadIOPSDevice {
			if trid.Major == device.Major && trid.Minor == device.Minor {
				found = true
				t.Ok(trid.Rate == device.Rate, fmt.Sprintf("blkio read iops for %d:%d is set correctly", device.Major, device.Minor))
				t.Diagnosticf("expect: %d, actual: %d", device.Rate, trid.Rate)
			}
		}
		t.Ok(found, fmt.Sprintf("blkio read iops for %d:%d found", device.Major, device.Minor))
	}

	for _, device := range config.Linux.Resources.BlockIO.ThrottleWriteIOPSDevice {
		found := false
		for _, twid := range lbd.ThrottleWriteIOPSDevice {
			if twid.Major == device.Major && twid.Minor == device.Minor {
				found = true
				t.Ok(twid.Rate == device.Rate, fmt.Sprintf("blkio write iops for %d:%d is set correctly", device.Major, device.Minor))
				t.Diagnosticf("expect: %d, actual: %d", device.Rate, twid.Rate)
			}
		}
		t.Ok(found, fmt.Sprintf("blkio write iops for %d:%d found", device.Major, device.Minor))
	}

	t.AutoPlan()
	return nil
}
