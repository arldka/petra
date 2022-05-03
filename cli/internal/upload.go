package module

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"gopkg.in/yaml.v2"
)

var gcsBucket *GCSBackend

type GCSBackend struct {
	client *storage.Client
	bucket string
}

type Metadata struct {
	Owner string
	Team  string
}

type PetraConfig struct {
	Namespace string
	Name      string
	Provider  string
	Version   string
	Metadata  Metadata
}

func InitGCSBackend(bckt string) error {
	ctx := context.Background()
	fmt.Println("bucket name :", bckt)
	client, err := storage.NewClient(ctx)
	fmt.Println("Client : ", client)
	if err != nil {
		return err
	}

	gcsBucket = &GCSBackend{
		client: client,
		bucket: bckt,
	}

	attrs, err := gcsBucket.client.Bucket(gcsBucket.bucket).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		fmt.Fprintln(os.Stderr, "The", gcsBucket.bucket, "bucket does not exist")
		return err
	}
	if err != nil {
		// Other error to handle
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println("The", gcsBucket.bucket, "bucket exists and has attributes:", attrs)
	return err
}

func Tar(moduleDirectory string) error {
	file, err := os.Create("module.zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(moduleDirectory, walker)
	if err != nil {
		panic(err)
	}
	return nil
}

func GetPetraConfig(modulePath string) (*PetraConfig, error) {
	config := PetraConfig{}
	configPath := modulePath + ".petra-config.yaml"

	fmt.Println(configPath)

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	fmt.Printf("%+v\n", config)

	// check required fields
	if config.Namespace == "" {
		return nil, fmt.Errorf("error: required field (namespace) is missing in the config file")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("error: required field (name) is missing in the config file")
	}
	if config.Provider == "" {
		return nil, fmt.Errorf("error: required field (provider) is missing in the config file")
	}
	if config.Version == "" {
		return nil, fmt.Errorf("error: required field (version) is missing in the config file")
	}

	return &config, nil
}

func UploadModule(zipFilePath string, petraConf *PetraConfig) error {
	ctx := context.Background()
	// Open local file.
	f, err := os.Open(zipFilePath)
	if err != nil {
		return fmt.Errorf("os.Open: %v", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// {namespace}/{module}/{provider}/{namespace}-{module}-{provider}-{version}.tar.gz
	// e.g.: main/rabbitmq/helm/0.0.1/main-rabbitmq-helm-0.0.1.tar.gz

	// {namespace}/{module}/{provider}/
	objectDirectory := petraConf.Namespace + "/" + petraConf.Name + "/" + petraConf.Provider + "/"
	// {namespace}-{module}-{provider}-{version}.tar.gz
	object := petraConf.Namespace + "-" + petraConf.Name + "-" + petraConf.Provider + "-" + petraConf.Version + ".tar.gz"

	o := gcsBucket.client.Bucket(gcsBucket.bucket).Object(objectDirectory + object)

	wc := o.NewWriter(ctx)

	// Update the object to set the metadata:
	// - owner
	// - team
	objectAttrs := map[string]string{
		"owner": petraConf.Metadata.Owner,
		"team":  petraConf.Metadata.Team,
	}
	wc.ObjectAttrs.Metadata = objectAttrs

	// Upload an object with storage.Writer.
	if _, err = io.Copy(wc, f); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}
	//fmt.Fprintf(w, "Blob %v uploaded.\n", filename)
	return nil
}
