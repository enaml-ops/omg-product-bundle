# p-scs plugin

This is a product plugin for Spring Cloud Services for PCF.

## Register Plugin

```bash
$ omg-linux register-plugin -type product -pluginpath ./p-scs-plugin-linux
```

## Deployment

A general deployment looks like this:

```bash
$ omg-linux deploy-product \
  --print-manifest \
  --bosh-url https://1.0.0.1 \
  --bosh-port 25555 \
  --bosh-user admin \
  --bosh-pass pass \
  --ssl-ignore \
  p-scs-plugin-linux \
    --system-domain sys.example.com \
    --app-domain apps.example.com
    --infer-from-cloud \
    --uaa-admin-secret asdfghjkl \
    --admin-password qwertyyuiop
```

Note that several of these flags must match the values used in the CF deployment.
The recommended approach is to pull these values from Vault instead of specifying them manually.

## Vault Integration

The deployment can be simplified by pulling values from a Vault instance.
If you've used `omg` to deploy elastic runtime, these values should already be present in Vault.

Simply specify your Vault address, token, and a list of hashes to read from.
The hashes will be enumerated in the order specified.
Any value that matches that of a CLI flag will be used as a default.

```bash
$ omg-linux deploy-product \
  --print-manifest \
  --bosh-url https://1.0.0.1 \
  --bosh-port 25555 \
  --bosh-user admin \
  --bosh-pass pass \
  --ssl-ignore \
  p-scs-plugin-linux \
    --system-domain sys.example.com \
    --app-domain apps.example.com \
    --infer-from-cloud
    --vault-domain http://10.1.2.3:8200 \
    --vault-token qwertyuiopasdfjkl \
    --vault-hash secret/pcf-1-passwords
```
