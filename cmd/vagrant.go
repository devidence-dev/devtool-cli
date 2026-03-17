package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pxblo-x1/devtool-cli/internal/ui"
	"github.com/spf13/cobra"
)

var vagrantCmd = &cobra.Command{
	Use:   "vagrant",
	Short: "Vagrant related commands",
}

var vagrantCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up Vagrant boxes and libvirt resources",
	RunE:  runVagrantCleanup,
}

func init() {
	vagrantCmd.AddCommand(vagrantCleanupCmd)
}

func runVagrantCleanup(cmd *cobra.Command, args []string) error {
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)).
		WithTextStyle(pterm.NewStyle(pterm.Bold)).
		WithFullWidth().
		Println("Vagrant & Libvirt Cleanup")

	if !commandExists("vagrant") {
		ui.Error("Vagrant is not installed on this system")
		return fmt.Errorf("vagrant not found")
	}

	libvirtInstalled := commandExists("virsh")
	if !libvirtInstalled {
		ui.Warning("libvirt/virsh is not installed. Libvirt cleanup will be skipped")
	}

	showCurrentState(libvirtInstalled)

	options := []string{
		"1. Destroy all Vagrant VMs",
		"2. Remove unused Vagrant boxes",
		"3. Remove all Vagrant boxes",
		"4. Clean orphaned libvirt volumes",
		"5. Clean inactive libvirt domains",
		"6. Full cleanup (all of the above)",
		"7. Exit",
	}

	selected, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Select an option")
	if err != nil {
		return err
	}
	fmt.Println()

	switch {
	case strings.HasPrefix(selected, "1."):
		ok, _ := pterm.DefaultInteractiveConfirm.Show("Are you sure you want to destroy all VMs?")
		if ok {
			withSpinner("Destroying all Vagrant VMs...", "VMs destroyed", destroyAllVMs)
		}

	case strings.HasPrefix(selected, "2."):
		ok, _ := pterm.DefaultInteractiveConfirm.Show("Remove unused boxes?")
		if ok {
			withSpinner("Removing unused boxes...", "Unused boxes removed", func() {
				runCmd("vagrant", "box", "prune", "-f")
			})
		}

	case strings.HasPrefix(selected, "3."):
		ok, _ := pterm.DefaultInteractiveConfirm.Show("WARNING! Remove ALL boxes?")
		if ok {
			withSpinner("Removing all boxes...", "All boxes removed", removeAllBoxes)
		}

	case strings.HasPrefix(selected, "4."):
		if !libvirtInstalled {
			ui.Error("libvirt is not available")
			return fmt.Errorf("libvirt not found")
		}
		ok, _ := pterm.DefaultInteractiveConfirm.Show("Remove orphaned volumes?")
		if ok {
			withSpinner("Cleaning orphaned libvirt volumes...", "Orphaned volumes cleaned", cleanLibvirtVolumes)
		}

	case strings.HasPrefix(selected, "5."):
		if !libvirtInstalled {
			ui.Error("libvirt is not available")
			return fmt.Errorf("libvirt not found")
		}
		ok, _ := pterm.DefaultInteractiveConfirm.Show("Remove inactive domains containing 'vagrant'?")
		if ok {
			withSpinner("Cleaning inactive libvirt domains...", "Inactive domains cleaned", cleanLibvirtDomains)
		}

	case strings.HasPrefix(selected, "6."):
		ui.Warning("FULL CLEANUP: All VMs, boxes and libvirt resources will be removed")
		ok, _ := pterm.DefaultInteractiveConfirm.Show("Are you COMPLETELY sure?")
		if !ok {
			ui.Info("Operation cancelled")
			return nil
		}

		steps := 3
		if libvirtInstalled {
			steps = 5
		}
		pb, _ := pterm.DefaultProgressbar.WithTotal(steps).WithTitle("Full cleanup").Start()

		pb.UpdateTitle("Destroying Vagrant VMs...")
		destroyAllVMs()
		pb.Increment()

		pb.UpdateTitle("Removing Vagrant boxes...")
		removeAllBoxes()
		pb.Increment()

		if libvirtInstalled {
			pb.UpdateTitle("Cleaning libvirt domains...")
			cleanLibvirtDomainsForce()
			pb.Increment()

			pb.UpdateTitle("Cleaning libvirt volumes...")
			cleanLibvirtVolumes()
			pb.Increment()

			pb.UpdateTitle("Cleaning libvirt networks...")
			cleanLibvirtNetworks()
			pb.Increment()
		}

		pb.UpdateTitle("Cleaning Vagrant temp directory...")
		os.RemoveAll(os.Getenv("HOME") + "/.vagrant.d/tmp")

		_, _ = pb.Stop()
		ui.Success("Full cleanup complete!")

	case strings.HasPrefix(selected, "7."):
		ui.Info("Exiting...")
		return nil
	}

	fmt.Println()
	ui.Success("Done")
	return nil
}

func withSpinner(start, done string, fn func()) {
	spinner, _ := pterm.DefaultSpinner.Start(start)
	fn()
	spinner.Success(done)
}

func showCurrentState(libvirtInstalled bool) {
	pterm.DefaultSection.Println("Current State")

	ui.Info("Vagrant virtual machines:")
	runCmd("vagrant", "global-status", "--prune")
	fmt.Println()

	ui.Info("Installed Vagrant boxes:")
	runCmd("vagrant", "box", "list")
	fmt.Println()

	if libvirtInstalled {
		ui.Info("Libvirt domains:")
		runCmd("virsh", "-c", "qemu:///system", "list", "--all")
		fmt.Println()

		ui.Info("Volumes in default pool:")
		runCmd("virsh", "-c", "qemu:///system", "vol-list", "default")
		fmt.Println()
	}
}

func destroyAllVMs() {
	out, err := exec.Command("vagrant", "global-status", "--prune").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "running") || strings.Contains(line, "poweroff") || strings.Contains(line, "saved") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				runCmd("vagrant", "destroy", fields[0], "-f")
			}
		}
	}
}

func removeAllBoxes() {
	out, err := exec.Command("vagrant", "box", "list").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			runCmd("vagrant", "box", "remove", "-f", "--all", fields[0])
		}
	}
}

func cleanLibvirtVolumes() {
	out, err := exec.Command("virsh", "-c", "qemu:///system", "vol-list", "default").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "vagrant") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			runCmd("virsh", "-c", "qemu:///system", "vol-delete", fields[0], "default")
		}
	}
}

func cleanLibvirtDomains() {
	out, err := exec.Command("virsh", "-c", "qemu:///system", "list", "--all").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "shut off") || !strings.Contains(line, "vagrant") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			runCmd("virsh", "-c", "qemu:///system", "undefine", fields[1], "--remove-all-storage")
		}
	}
}

func cleanLibvirtDomainsForce() {
	out, err := exec.Command("virsh", "-c", "qemu:///system", "list", "--all").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "vagrant") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			runCmd("virsh", "-c", "qemu:///system", "destroy", fields[1])
			runCmd("virsh", "-c", "qemu:///system", "undefine", fields[1], "--remove-all-storage")
		}
	}
}

func cleanLibvirtNetworks() {
	out, err := exec.Command("virsh", "-c", "qemu:///system", "net-list", "--all").Output()
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.Contains(line, "vagrant") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			runCmd("virsh", "-c", "qemu:///system", "net-destroy", fields[0])
			runCmd("virsh", "-c", "qemu:///system", "net-undefine", fields[0])
		}
	}
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
