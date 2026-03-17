package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/pxblo-x1/devtool-cli/internal/ui"
	"github.com/spf13/cobra"
)

var vscodeCmd = &cobra.Command{
	Use:   "vscode",
	Short: "Comandos relacionados a VS Code",
}

var vscodeKillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Termina procesos huérfanos de VS Code Remote SSH",
	RunE:  runVSCodeKill,
}

func init() {
	vscodeCmd.AddCommand(vscodeKillCmd)
}

func runVSCodeKill(cmd *cobra.Command, args []string) error {
	pterm.DefaultHeader.
		WithBackgroundStyle(pterm.NewStyle(pterm.BgCyan)).
		WithTextStyle(pterm.NewStyle(pterm.Bold)).
		WithFullWidth().
		Println("VS Code Server Process Killer")

	user := os.Getenv("USER")

	spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Buscando procesos de VS Code Server para %s...", user))
	pids, err := findVSCodePIDs(user)
	if err != nil {
		spinner.Fail("Error al buscar procesos")
		return err
	}

	if len(pids) == 0 {
		spinner.Success("No se encontraron procesos de VS Code Server")
		return nil
	}
	spinner.Success(fmt.Sprintf("%d proceso(s) encontrado(s)", len(pids)))
	fmt.Println()

	printProcessTable(pids)
	fmt.Println()

	ok, _ := pterm.DefaultInteractiveConfirm.Show("¿Terminar estos procesos?")
	if !ok {
		ui.Info("Operación cancelada")
		return nil
	}

	fmt.Println()
	killSpinner, _ := pterm.DefaultSpinner.Start("Terminando procesos con SIGTERM...")
	for _, pid := range pids {
		if p, err := os.FindProcess(pid); err == nil {
			_ = p.Signal(syscall.SIGTERM)
		}
	}
	time.Sleep(time.Second)

	var remaining []int
	for _, pid := range pids {
		if processExists(pid) {
			remaining = append(remaining, pid)
		}
	}

	if len(remaining) > 0 {
		killSpinner.Warning("Algunos procesos resistieron SIGTERM, usando SIGKILL...")
		for _, pid := range remaining {
			if p, err := os.FindProcess(pid); err == nil {
				_ = p.Signal(syscall.SIGKILL)
			}
		}
	}

	killSpinner.Success("Todos los procesos de VS Code Server han sido terminados")
	return nil
}

func findVSCodePIDs(user string) ([]int, error) {
	out, err := exec.Command("ps", "-u", user, "-o", "pid,command").Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(\.vscode-server/|vscode-server/bin)`)
	pidSet := map[int]bool{}

	for _, line := range strings.Split(string(out), "\n") {
		if re.MatchString(line) && !strings.Contains(line, "grep") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				if pid, err := strconv.Atoi(fields[0]); err == nil {
					pidSet[pid] = true
				}
			}
		}
	}

	pids := make([]int, 0, len(pidSet))
	for pid := range pidSet {
		pids = append(pids, pid)
	}
	sort.Ints(pids)
	return pids, nil
}

func printProcessTable(pids []int) {
	tableData := pterm.TableData{{"PID", "%CPU", "%MEM", "COMANDO"}}

	for _, pid := range pids {
		out, err := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid=,%cpu=,%mem=,command=").Output()
		if err != nil {
			continue
		}
		fields := strings.Fields(strings.TrimSpace(string(out)))
		if len(fields) >= 4 {
			command := strings.Join(fields[3:], " ")
			if len(command) > 60 {
				command = command[:60] + "..."
			}
			tableData = append(tableData, []string{fields[0], fields[1], fields[2], command})
		}
	}

	_ = pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
}

func processExists(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}
