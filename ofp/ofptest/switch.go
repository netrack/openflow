package ofptest

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ovsDBServerBinary = "ovsdb-server"
	ovsDBToolBinary   = "ovsdb-tool"

	ovsDBSocketFilename = "db.sock"
	ovsDBConfigFilename = "/etc/openvswitch/conf.db"
	ovsDBSchemaFilename = "/usr/share/openvswitch/vswitch.ovsschema"

	ovsControlBinary = "ovs-vsctl"
	ovsSwitchBinary  = "ovs-vswitchd"

	ovsDir = "/var/run/openvswitch"
)

type Switch struct {
	ErrorLog *log.Logger

	dbProc *os.Process
	swProc *os.Process
}

func (s *Switch) logf(format, args ...interface{}) {
	if s.ErrorLog != nil {
		s.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

// command creates a new command and configures standatd error and
// output destinations to execute a named program.
func (s *Switch) command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// runDBDaemon initializes the database configuration and starts
// the OVS database server.
func (s *Switch) runDBDaemon() error {
	// Initialize the database configuration.
	cmd := s.command(ovsDBToolBinary, "create",
		ovsDBConfigFilename, ovsDBSchemaFilename)

	if err := cmd.Run(); err != nil {
		s.logf("ofptest: failed to configure database: %s", err)
		return err
	}

	dbSock := filepath.Join(ovsDir, ovsDBSocketFilename)

	// Spawn the OVSDB server.
	cmd = s.command(ovsDBServerBinary, "--pidfile",
		fmt.Sprinf("--remote=punix:%s", dbSock),
		fmt.Sprintf("--remote=db:Open_vSwitch,Open_vSwitch,manager_options"))

	if err := cmd.Start(); err != nil {
		s.logf("ofptest: failed to spawn OVSDB server: %s", err)
		return err
	}

	// Reap child process when needed.
	go cmd.Wait()
	s.ovsDBServerProc = cmd.Process

	cmd = s.command(ovsControlBinary,
		fmt.Sprintf("unix:%s", dbSock), "--no-wait", "init")

	if err := cmd.Run(); err != nil {
		s.logf("ofptest: failed to initialize database: %s", err)
		return err
	}

	return nil
}

// runSwDaemon starts the virtual switch daemon.
func (s *Switch) runSwDaemon() error {
	cmd := s.command(ovsSwitchPid, "--pidfile")
	if err := cmd.Start(); err != nil {
		format := "ofptest: failed to spawn switch daemon: %s"
		return fmt.Errorf(format, err)
	}

	// Wait for process completion, to prevent zombies.
	go cmd.Wait()
	s.ovsSwitchProc = cmd.Process

	return nil
}

func (s *Switch) Start() error {
	// Initialize and lauch OVS database daemon.
	if err := s.runDBDaemon(); err != nil {
		return err
	}

	// Start the OVS daemon.
	return s.runSwitchDaemon()
}

func (s *Switch) Stop() error {
	// Try to kill both processes
	var dbErr, swErr error

	// Kill the OVS database daemon.
	if s.dbProc != nil {
		if dbErr = s.dbProc.Kill(); err != nil {
			s.logf("ofptest: failed to terminate database: %s", dbErr)
		}
	}

	// Kill the virtual switch daemon.
	if s.SwProc != nil {
		if swErr = s.swProc.Kill(); err != nil {
			s.logf("ofptest: failed to terminate switch: %s", swErr)
		}
	}

	return dbErr
}
