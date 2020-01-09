package runcommands

import (
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// S3Target is the struct for the copypasta target destination
type S3Target struct {
	Name            string `yaml:"name"`
	AccessKey       string `yaml:"accesskey"`
	SecretAccessKey string `yaml:"secretaccesskey"`
	BucketName      string `yaml:"bucketname"`
	Endpoint        string `yaml:"endpoint"`
	Location        string `yaml:"location"`
}

// Update the configuration file
func Update(name, accessKey, secretAccessKey, bucketName, endpoint, location string) error {
	var target *S3Target
	var err error

	freshTarget := &S3Target{
		Name:            name,
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
		Endpoint:        endpoint,
		Location:        location,
	}

	target = freshTarget
	targetContents, err := yaml.Marshal(&target)
	if err != nil {
		log.Errorf("Error marshalling data: %s", err.Error())
		return err
	}

	err = ioutil.WriteFile(filepath.Join(os.Getenv("HOME"), ".copypastarc"), targetContents, 0666)
	if err != nil {
		log.Errorf("Error writing config to file: %s", err.Error())
		return err
	}

	return nil
}

// Load retrieves the S3 target configuration from a rc file
func Load() (*S3Target, error) {
	var target *S3Target

	byteContent, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".copypastarc"))
	if err != nil {
		log.Errorf("Unable to load the S3 target, make sure ~/.copypastarc exists %s", err.Error())
		return nil, err
	}
	err = yaml.Unmarshal(byteContent, &target)
	if err != nil {
		log.Errorf("Unable to parse data %s", err.Error())
		return nil, err
	}

	return target, nil
}
