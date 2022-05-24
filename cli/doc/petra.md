## petra

Private terraform registry cli

### Synopsis

CLI to manage terraform modules in our private registry in a Google Cloud Storage bucket.

```
petra [flags]
```

### Options

```
      --gcs-bucket string         Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket
  -h, --help                      help for petra
      --module-directory string   Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/
```

### SEE ALSO

* [petra remove](petra_remove.md)	 - Remove the module from a private registry
* [petra update](petra_update.md)	 - Update one or multiple config values of a module.
* [petra upload](petra_upload.md)	 - Upload a terraform module to a private registry
