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
	Short: "Comandos relacionados a Vagrant",
}

var vagrantCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Limpia boxes de Vagrant y recursos de libvirt",
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
		ui.Error("Vagrant no está instalado en el sistema")
		return fmt.Errorf("vagrant no encontrado")
	}

	libvirtInstalled := commandExists("virsh")
	if !libvirtInstalled {
		ui.Warning("libvirt/virsh no está instalado. Se omitirá la limpieza de libvirt")
	}

	showCurrentState(libvirtInstalled)

	options := []string{
		"1. Destruir todas las VMs de Vagrant",
		"2. Eliminar boxes de Vagrant no utilizados",
		"3. Eliminar todos los boxes de Vagrant",
		"4. Limpiar volúmenes huérfanos de libvirt",
		"5. Limpiar dominios inactivos de libvirt",
		"6. Limpieza completa (todas las opciones anteriores)",
		"7. Salir",
	}

	selected, err := pterm.DefaultInteractiveSelect.WithOptions(options).Show("Selecciona una opción")
	if err != nil {
		return err
	}
	fmt.Println()

	switch {
	case strings.HasPrefix(selected, "1."):
		ok, _ := pterm.DefaultInteractiveConfirm.Show("¿Estás seguro de destruir todas las VMs?")
		if ok {
			withSpinner("Destruyendo todas las VMs de Vagrant...", "VMs destruidas", destroyAllVMs)
		}

	case strings.HasPrefix(selected, "2."):
		ok, _ := pterm.DefaultInteractiveConfirm.Show("¿Deseas eliminar los boxes no utilizados?")
		if ok {
			withSpinner("Eliminando boxes no utilizados...", "Boxes no utilizados eliminados", func() {
				runCmd("vagrant", "box", "prune", "-f")
			})
		}

	case strings.HasPrefix(selected, "3."):
		ok, _ := pterm.DefaultInteractiveConfirm.Show("¡ADVERTENCIA! ¿Deseas eliminar TODOS los boxes?")
		if ok {
			withSpinner("Eliminando todos los boxes...", "Todos los boxes eliminados", removeAllBoxes)
		}

	case strings.HasPrefix(selected, "4."):
		if !libvirtInstalled {
			ui.Error("libvirt no está disponible")
			return fmt.Errorf("libvirt no encontrado")
		}
		ok, _ := pterm.DefaultInteractiveConfirm.Show("¿Deseas eliminar volúmenes huérfanos?")
		if ok {
			withSpinner("Limpiando volúmenes huérfanos de libvirt...", "Volúmenes huérfanos limpiados", cleanLibvirtVolumes)
		}

	case strings.HasPrefix(selected, "5."):
		if !libvirtInstalled {
			ui.Error("libvirt no está disponible")
			return fmt.Errorf("libvirt no encontrado")
		}
		ok, _ := pterm.DefaultInteractiveConfirm.Show("¿Deseas eliminar dominios inactivos que contengan 'vagrant'?")
		if ok {
			withSpinner("Limpiando dominios inactivos de libvirt...", "Dominios inactivos limpiados", cleanLibvirtDomains)
		}

	case strings.HasPrefix(selected, "6."):
		ui.Warning("LIMPIEZA COMPLETA: Se eliminarán todas las VMs, boxes y recursos de libvirt")
		ok, _ := pterm.DefaultInteractiveConfirm.Show("¿Estás COMPLETAMENTE seguro?")
		if !ok {
			ui.Info("Operación cancelada")
			return nil
		}

		steps := 3
		if libvirtInstalled {
			steps = 5
		}
		pb, _ := pterm.DefaultProgressbar.WithTotal(steps).WithTitle("Limpieza completa").Start()

		pb.UpdateTitle("Destruyendo VMs de Vagrant...")
		destroyAllVMs()
		pb.Increment()

		pb.UpdateTitle("Eliminando boxes de Vagrant...")
		removeAllBoxes()
		pb.Increment()

		if libvirtInstalled {
			pb.UpdateTitle("Limpiando dominios de libvirt...")
			cleanLibvirtDomainsForce()
			pb.Increment()

			pb.UpdateTitle("Limpiando volúmenes de libvirt...")
			cleanLibvirtVolumes()
			pb.Increment()

			pb.UpdateTitle("Limpiando redes de libvirt...")
			cleanLibvirtNetworks()
			pb.Increment()
		}

		pb.UpdateTitle("Limpiando directorio temporal de Vagrant...")
		os.RemoveAll(os.Getenv("HOME") + "/.vagrant.d/tmp")

		pb.Stop()
		ui.Success("¡Limpieza completa finalizada!")

	case strings.HasPrefix(selected, "7."):
		ui.Info("Saliendo...")
		return nil
	}

	fmt.Println()
	ui.Success("Script finalizado")
	return nil
}

func withSpinner(start, done string, fn func()) {
	spinner, _ := pterm.DefaultSpinner.Start(start)
	fn()
	spinner.Success(done)
}

func showCurrentState(libvirtInstalled bool) {
	pterm.DefaultSection.Println("Estado Actual")

	ui.Info("Máquinas virtuales de Vagrant:")
	runCmd("vagrant", "global-status", "--prune")
	fmt.Println()

	ui.Info("Boxes de Vagrant instalados:")
	runCmd("vagrant", "box", "list")
	fmt.Println()

	if libvirtInstalled {
		ui.Info("Dominios de libvirt:")
		runCmd("virsh", "-c", "qemu:///system", "list", "--all")
		fmt.Println()

		ui.Info("Volúmenes en pool default:")
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
	cmd.Run()
}
