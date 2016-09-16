# p-rabbitmq plugin

This is a product plugin for Pivotal-RabbitMQ.

## Register Plugin

```bash
$ omg-linux register-plugin -type product -pluginpath ./p-rabbitmq-plugin-linux
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
  p-rabbitmq-plugin-linux \
    --system-domain sys.example.com \
    --rabbit-public-ip 1.0.0.3 \
    --infer-from-cloud \
    --rabbit-server-ip 1.0.0.4 --rabbit-server-ip 1.0.0.5
    --broker-ip 1.0.0.6 \
    --syslog-address 1.0.0.7 \
    --nats-machine-ip 1.0.0.8 --nats-machine-ip 1.0.0.9 \
    --nats-pass natspassword \
    --system-services-password systemservices \
    --doppler-zone zone \
    --doppler-shared-secret secret \
    --etcd-machine-ip 1.0.0.10 --etcd-machine-ip 1.0.0.11
```

Note that several of these flags must match the values used in the CF deployment.
The recommended approach is to pull these values from Vault instead of specifying them manually.

## Vault Integration

The deployment can be simplified by pulling values from a Vault instance.
If you've used `omg` to deploy elastic runtime, many of these values should already be present in Vault.

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
  p-rabbitmq-plugin-linux \
    --system-domain sys.example.com \
    --vault-domain http://10.0.0.200:8200 \
    --vault-token qwertyuiopasdfjkl \
    --vault-hash secret/pcf-1-passwords \
    --vault-hash secret/pcf-1-ips \
    --vault-hash secret/pcf-1-hosts \
    --infer-from-cloud \
```
