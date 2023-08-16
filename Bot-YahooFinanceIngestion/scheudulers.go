package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func scheduleDataFetcherJob() error {
	// Get the current executable file path
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Build the cron job command
	cmd := exec.Command("crontab", "-l")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	cronJobCmd := fmt.Sprintf("%s %s", "@daily", execPath)

	// Check if the cron job already exists
	if !strings.Contains(string(out), cronJobCmd) {
		// Add the cron job command to the crontab
		cmd := exec.Command("bash", "-c", fmt.Sprintf("echo '%s %s' | crontab -", "0 18 * * *", execPath))
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
func ScheduleChecker() {
	var c bool
	if strings.Contains(runtime.GOOS, "linux") {
		c = checkCronJob()
		if !c {
			scheduleDataFetcherJob()
		}
	} else if strings.Contains(runtime.GOOS, "windows") {
		c = isDataFetcherScheduled()
		if !c {
			scheduleDataFetcher()
		}
	}
}

func isDataFetcherScheduled() bool {
	cmd := exec.Command("schtasks", "/query", "/TN", "DataFetcher")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if strings.Contains(string(output), "DataFetcher") {
		fmt.Println("DataFetcher is scheduled.")
		return true
	}
	fmt.Println("DataFetcher is not scheduled. Scheduling now...")

	return false
}
func checkCronJob() bool {
	out, err := exec.Command("crontab", "-l").Output()
	if err != nil {
		// handle error
		return false
	}
	// split output into lines
	lines := strings.Split(string(out), "\n")
	// loop through lines to find "DataFetcher" entry
	for _, line := range lines {
		if strings.Contains(line, "DataFetcher") {
			return true
		}
	}
	return false
}
func scheduleDataFetcher() {
	// Get the path of the current executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Build the command to create a new scheduled task, starts at 6 PM every day
	cmd := exec.Command("schtasks", "/create", "/tn", "DataFetcher", "/tr", exePath, "/sc", "daily", "/st", "18:00")

	// Set the working directory to the directory of the executable
	cmd.Dir = filepath.Dir(exePath)

	// Run the command and print any errors
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("DataFetcher scheduled successfully.")
}
