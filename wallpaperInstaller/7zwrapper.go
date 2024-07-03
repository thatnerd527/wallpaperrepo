package main

import (
	"os"
	"os/exec"
	"path"
	"syscall"
)

func extractZipToFolderWith7z(zipPath string, folderPath string) error {
	cmd := exec.Command("tools\\7z\\7za", "x", zipPath, "-o"+folderPath,"-aoa")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func smartExtract(zipPath string) (folderPath string, err error) {
	tempFolder := path.Join("temp", GenerateGUID())
	err = os.MkdirAll(tempFolder, 0755)
	if err != nil {
		return;
	}
	err = extractZipToFolderWith7z(zipPath, tempFolder)
	if err != nil {
		return;
	}
	// Determine if the author zipped the addon root folder or just the contents
	// If the zip file contains a single folder, then use that.
	// If the zip file has the manifest.json at the root, then use that.

	// Check if the zip file contains a single folder
	files, err := os.ReadDir(tempFolder)
	if err != nil {
		return
	}
	if len(files) == 1 && files[0].IsDir() {
		folderPath = path.Join(tempFolder, files[0].Name())
	} else {
		// Check if the zip file has the manifest.json at the root
		manifestPath := path.Join(tempFolder, "manifest.json")
		_, err = os.Stat(manifestPath)
		if err != nil {
			return
		}
		folderPath = tempFolder
	}

	return;
}