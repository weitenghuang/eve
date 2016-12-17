package terraform

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/pborman/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	TMP_QUOIN_DIR          = "/tmp/quoin"
	TERRAFORM_PROCESS_NAME = "terraform"
	DEFAULT_STATE_FILE     = "terraform.tfstate"
	CUSTOM_VAR_FILE        = "varfile"
	DIR_CONFLICT           = "Directory conflict: "
	MAX_RETRY              = 15
	PERM                   = 0755
)

func PlanQuoin(name string, tarGz []byte) error {
	log.Println("Prepare work directory.")
	dir := setWorkDir(name, 0)
	log.Println(dir)
	defer os.RemoveAll(dir)

	if err := writeFileFromTarGz(dir, tarGz, PERM); err != nil {
		return err
	}
	if _, err := runTerraformPlan(dir, fmt.Sprintf("%s%s", name, ".tfplan")); err != nil {
		return err
	}
	return nil
}

func ApplyQuoin(name string, modules []byte, varfile []byte, remoteState string) error {
	log.Println("Prepare work directory.")
	dir := setWorkDir(name, 0)
	log.Println(dir)
	defer os.RemoveAll(dir)

	if err := writeFileFromTarGz(dir, modules, PERM); err != nil {
		return err
	}

	if varfile != nil {
		log.Println("Create var file.")
		if err := ioutil.WriteFile(filepath.Join(dir, CUSTOM_VAR_FILE), varfile, PERM); err != nil {
			return err
		}
	}

	log.Println("remote state:", remoteState)
	if err := runTerraformRemote(dir, remoteState); err != nil {
		return err
	}

	if err := runTerraform(dir, "apply"); err != nil {
		return err
	}
	return nil
}

func DeleteQuoin(name string, modules []byte, varfile []byte, remoteState string) error {
	log.Println("Prepare work directory.")
	dir := setWorkDir(name, 0)
	log.Println(dir)
	defer os.RemoveAll(dir)

	if err := writeFileFromTarGz(dir, modules, PERM); err != nil {
		return err
	}

	if varfile != nil {
		log.Println("Create var file.")
		if err := ioutil.WriteFile(filepath.Join(dir, CUSTOM_VAR_FILE), varfile, PERM); err != nil {
			return err
		}
	}

	log.Println("remote state:", remoteState)
	if err := runTerraformRemote(dir, remoteState); err != nil {
		return err
	}

	if err := runTerraform(dir, "destroy"); err != nil {
		return err
	}
	return nil
}

func checkDirConflictPanic(dir string) {
	if _, err := os.Stat(dir); err == nil {
		panic(fmt.Sprintf("%s%s", DIR_CONFLICT, dir))
	}
}

func setWorkDir(name string, retry int) (dir string) {
	defer func() {
		if r := recover(); r != nil && fmt.Sprintf("%s", r) == fmt.Sprintf("%s%s", DIR_CONFLICT, dir) {
			// Increase timestamp to create different path base
			time.Sleep(1 * time.Second)
			if retry > MAX_RETRY {
				panic(fmt.Sprintf("%s%s and failed retry over %s times", DIR_CONFLICT, dir, MAX_RETRY))
			}
			retry += 1
			dir = setWorkDir(fmt.Sprintf("%s-%s", name, retry), retry)
		}
	}()
	pathBase := fmt.Sprintf("%s%v", name, time.Now().Unix())
	pathUUID := uuid.NewMD5(uuid.NewRandom(), []byte(pathBase)).String()
	dir = filepath.Join(TMP_QUOIN_DIR, pathUUID)
	checkDirConflictPanic(dir)
	return dir
}

func writeFileFromTarGz(dir string, tarGz []byte, perm os.FileMode) error {
	r := bytes.NewReader(tarGz)
	gzf, err := gzip.NewReader(r)
	defer gzf.Close()
	if err != nil {
		return err
	}
	tr := tar.NewReader(gzf)

	if err = os.MkdirAll(dir, perm); err != nil {
		return err
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		info := hdr.FileInfo()
		path := filepath.Join(dir, hdr.Name)
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tr)
		if err != nil {
			return err
		}
		log.Printf("File %s is created at %s\n", hdr.Name, dir)
	}
	return nil
}

func runTerraformGet(dir string) error {
	var outBuf, errorBuf bytes.Buffer
	getCommand := exec.Command(TERRAFORM_PROCESS_NAME, "get")
	getCommand.Dir = dir
	getCommand.Stdout = &outBuf
	getCommand.Stderr = &errorBuf
	if err := getCommand.Start(); err != nil {
		log.Println("Command starts with error:", err)
		return err
	}
	if err := getCommand.Wait(); err != nil {
		log.Println("Command exits with error:", err)
		errStr := errorBuf.String()
		log.Println(errStr)
		return fmt.Errorf(errStr)
	}
	if outBuf.Len() > 0 {
		log.Println(outBuf.String())
	}
	return nil
}

func runTerraformPlan(dir string, tfplan string) ([]byte, error) {
	if err := runTerraformGet(dir); err != nil {
		return nil, err
	}

	var outBuf, errorBuf bytes.Buffer
	varFile := filepath.Join(dir, "/terraform.tfvars")
	varFileArg := filepath.Join("-var-file=", varFile)
	stateFile := filepath.Join(dir, DEFAULT_STATE_FILE)
	stateFileArg := filepath.Join("-state=", stateFile)
	outFile := filepath.Join(dir, tfplan)
	outArg := filepath.Join("-out=", outFile)
	planCommand := exec.Command(TERRAFORM_PROCESS_NAME, "plan", varFileArg, stateFileArg, outArg)
	planCommand.Dir = dir
	planCommand.Stdout = &outBuf
	planCommand.Stderr = &errorBuf
	if err := planCommand.Start(); err != nil {
		log.Println("Command starts with error:", err)
		return nil, err
	}
	if err := planCommand.Wait(); err != nil {
		log.Println("Command exits with error:", err)
		errStr := errorBuf.String()
		log.Println(errStr)
		return nil, fmt.Errorf(errStr)
	}
	if outBuf.Len() > 0 {
		log.Printf("%s\n", outBuf.String())
		log.Println("Done!")
		tfplanBin, err := ioutil.ReadFile(outFile)
		if err != nil {
			return nil, err
		}
		return tfplanBin, nil
	}
	return nil, fmt.Errorf("Terraform plan executed with empty output: %s", errorBuf.String())
}

func runTerraformRemote(dir string, address string) error {
	var outBuf, errorBuf bytes.Buffer
	var argvs []string
	argvs = append(argvs, "remote")
	argvs = append(argvs, "config")
	argvs = append(argvs, "-backend=http")
	argvs = append(argvs, "-backend-config")
	// argvs = append(argvs, "address=http://localhost:8088/infrastructure/myvpc/state")
	argvs = append(argvs, fmt.Sprint("address=", address))
	log.Printf("argvs: %#v", argvs)
	planCommand := exec.Command(TERRAFORM_PROCESS_NAME, argvs...)
	planCommand.Dir = dir
	planCommand.Stdout = &outBuf
	planCommand.Stderr = &errorBuf
	if err := planCommand.Start(); err != nil {
		log.Println("Command starts with error:", err)
		return err
	}
	if err := planCommand.Wait(); err != nil {
		log.Println("Command exits with error:", err)
		errStr := errorBuf.String()
		log.Println(errStr)
		return fmt.Errorf(errStr)
	}
	log.Printf("%s\n", outBuf.String())
	return nil
}

func runTerraform(dir string, action string) error {
	if err := runTerraformGet(dir); err != nil {
		return err
	}

	var outBuf, errorBuf bytes.Buffer
	var argvs []string
	argvs = append(argvs, action)
	if action == "destroy" {
		argvs = append(argvs, "-force")
	}
	varFile := filepath.Join(dir, CUSTOM_VAR_FILE)
	varFileArg := filepath.Join("-var-file=", varFile)
	argvs = append(argvs, varFileArg)
	// stateFile := filepath.Join(dir, DEFAULT_STATE_FILE)
	// stateFileArg := filepath.Join("-state=", stateFile)
	// argvs = append(argvs, stateFileArg)
	planCommand := exec.Command(TERRAFORM_PROCESS_NAME, argvs...)
	planCommand.Dir = dir
	planCommand.Stdout = &outBuf
	planCommand.Stderr = &errorBuf
	if err := planCommand.Start(); err != nil {
		log.Println("Command starts with error:", err)
		return err
	}
	if err := planCommand.Wait(); err != nil {
		log.Println("Command exits with error:", err)
		errStr := errorBuf.String()
		log.Println(errStr)
		return fmt.Errorf(errStr)
	}
	if outBuf.Len() > 0 {
		log.Printf("%s\n", outBuf.String())
		log.Println("Done!")
		return nil
	}
	return fmt.Errorf("Terraform plan executed with empty output: %s", errorBuf.String())
}
