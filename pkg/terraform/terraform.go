package terraform

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/concur/eve/pkg/vault"
	"github.com/pborman/uuid"
	"io"
	"io/ioutil"
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
	PROVIDER_TEMPLATE      = `provider "aws" {
	access_key = "%s"
	secret_key = "%s"
	token = "%s"
  region = "%s"
  max_retries = 3
}
`
)

type Terraform struct {
	name        string
	dir         string
	remoteState string
	modules     []byte
	varfile     []byte
}

func NewTerraform(name string, remoteState string, modules []byte, varfile []byte) *Terraform {
	return &Terraform{
		name:        name,
		remoteState: remoteState,
		modules:     modules,
		varfile:     varfile,
	}
}

func (tf *Terraform) PlanQuoin() error {
	log.Println("Prepare work directory.")
	tf.dir = setWorkDir(tf.name, 0)
	log.Println(tf.dir)
	defer os.RemoveAll(tf.dir)

	if err := tf.writeFileFromTarGz(PERM); err != nil {
		return err
	}

	if _, err := tf.runTerraformPlan(); err != nil {
		return err
	}
	return nil
}

func (tf *Terraform) ApplyQuoin() error {
	log.Println("Prepare work directory.")
	tf.dir = setWorkDir(tf.name, 0)
	log.Println(tf.dir)
	defer os.RemoveAll(tf.dir)

	if err := tf.writeFileFromTarGz(PERM); err != nil {
		return err
	}

	if err := tf.writeVarFile(); err != nil {
		return err
	}

	log.Println("remote state:", tf.remoteState)
	if err := tf.runTerraformRemote(); err != nil {
		return err
	}

	if err := tf.runTerraform("apply"); err != nil {
		return err
	}
	return nil
}

func (tf *Terraform) DeleteQuoin() error {
	log.Println("Prepare work directory.")
	tf.dir = setWorkDir(tf.name, 0)
	log.Println(tf.dir)
	defer os.RemoveAll(tf.dir)

	if err := tf.writeFileFromTarGz(PERM); err != nil {
		return err
	}

	if err := tf.writeVarFile(); err != nil {
		return err
	}

	log.Println("remote state:", tf.remoteState)
	if err := tf.runTerraformRemote(); err != nil {
		return err
	}

	if err := tf.runTerraform("destroy"); err != nil {
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

func (tf *Terraform) writeFileFromTarGz(perm os.FileMode) error {
	r := bytes.NewReader(tf.modules)
	gzf, err := gzip.NewReader(r)
	defer gzf.Close()
	if err != nil {
		return err
	}
	tr := tar.NewReader(gzf)

	if err = os.MkdirAll(tf.dir, perm); err != nil {
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
		path := filepath.Join(tf.dir, hdr.Name)
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
		log.Printf("File %s is created at %s\n", hdr.Name, tf.dir)
	}
	return nil
}

func (tf *Terraform) runTerraformGet() error {
	var outBuf, errorBuf bytes.Buffer
	getCommand := exec.Command(TERRAFORM_PROCESS_NAME, "get")
	getCommand.Dir = tf.dir
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

func (tf *Terraform) runTerraformPlan() ([]byte, error) {
	tfplan := fmt.Sprintf("%s%s", tf.name, ".tfplan")
	if err := tf.runTerraformGet(); err != nil {
		return nil, err
	}

	var outBuf, errorBuf bytes.Buffer
	varFile := filepath.Join(tf.dir, "/terraform.tfvars")
	varFileArg := filepath.Join("-var-file=", varFile)
	stateFile := filepath.Join(tf.dir, DEFAULT_STATE_FILE)
	stateFileArg := filepath.Join("-state=", stateFile)
	outFile := filepath.Join(tf.dir, tfplan)
	outArg := filepath.Join("-out=", outFile)
	planCommand := exec.Command(TERRAFORM_PROCESS_NAME, "plan", varFileArg, stateFileArg, outArg)
	planCommand.Dir = tf.dir
	planCommand.Stdout = &outBuf
	planCommand.Stderr = &errorBuf
	awsEnv, err := addAWSCredEnv(planCommand.Env, "secret/quoin/providers/aws/credentials")
	if err != nil {
		log.Println("Command loads env with error:", err)
	}
	planCommand.Env = awsEnv
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

func (tf *Terraform) runTerraformRemote() error {
	var outBuf, errorBuf bytes.Buffer
	var argvs []string
	//default user name: terraform
	user, err := vault.GetLogicalData("secret/user/terraform")
	if err != nil {
		log.Println("Missing terraform user in vault:", err)
		return err
	}
	argvs = append(argvs, "remote")
	argvs = append(argvs, "config")
	argvs = append(argvs, "-backend=http")
	// argvs = append(argvs, "-backend-config")
	argvs = append(argvs, fmt.Sprint("-backend-config=address=", tf.remoteState))
	argvs = append(argvs, fmt.Sprint("-backend-config=skip_cert_verification=true"))
	argvs = append(argvs, fmt.Sprint("-backend-config=username=", user["name"]))
	argvs = append(argvs, fmt.Sprint("-backend-config=password=", user["password"]))
	log.Printf("argvs: %#v", argvs)
	planCommand := exec.Command(TERRAFORM_PROCESS_NAME, argvs...)
	planCommand.Dir = tf.dir
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

func (tf *Terraform) runTerraform(action string) error {
	if err := tf.runTerraformGet(); err != nil {
		return err
	}

	var outBuf, errorBuf bytes.Buffer
	var argvs []string
	argvs = append(argvs, action)
	if action == "destroy" {
		argvs = append(argvs, "-force")
	}
	varFile := filepath.Join(tf.dir, CUSTOM_VAR_FILE)
	varFileArg := filepath.Join("-var-file=", varFile)
	argvs = append(argvs, varFileArg)
	planCommand := exec.Command(TERRAFORM_PROCESS_NAME, argvs...)
	planCommand.Dir = tf.dir
	planCommand.Stdout = &outBuf
	planCommand.Stderr = &errorBuf
	awsEnv, err := addAWSCredEnv(planCommand.Env, "secret/quoin/providers/aws/credentials")
	if err != nil {
		log.Println("Command loads env with error:", err)
	}
	planCommand.Env = awsEnv
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

func (tf *Terraform) writeVarFile() error {
	if tf.varfile == nil {
		return nil
	}

	log.Println("Create var file.")
	if err := ioutil.WriteFile(filepath.Join(tf.dir, CUSTOM_VAR_FILE), tf.varfile, PERM); err != nil {
		return err
	}
	return nil
}

func addAWSCredEnv(env []string, secretPath string) ([]string, error) {
	// secretPath := "secret/quoin/providers/aws/credentials"
	awsCred, err := vault.GetLogicalData(secretPath)
	if err != nil {
		return nil, err
	}
	awsKey := "AWS_ACCESS_KEY_ID"
	awsSecret := "AWS_SECRET_ACCESS_KEY"
	awsToken := "AWS_SESSION_TOKEN"
	if len(env) == 0 {
		env = os.Environ()
	}
	awsEnv := append(env,
		fmt.Sprintf("%s=%s", awsKey, awsCred[awsKey]),
		fmt.Sprintf("%s=%s", awsSecret, awsCred[awsSecret]),
		fmt.Sprintf("%s=%s", awsToken, awsCred[awsToken]),
	)
	return awsEnv, nil
}
