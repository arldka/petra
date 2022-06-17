## petractl

Private terraform registry cli

### Synopsis

CLI to manage terraform modules in our private registry in a Google Cloud Storage bucket.

```
petractl [flags]
```

### Options

```
      --gcs-bucket string         Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket
  -h, --help                      help for petractl
      --module-directory string   Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/
```

### SEE ALSO

* [petractl remove](petractl_remove.md)	 - Remove the module from a private registry
* [petractl update](petractl_update.md)	 - Update one or multiple config values of a module.
* [petractl upload](petractl_upload.md)	 - Upload a terraform module to a private registry
