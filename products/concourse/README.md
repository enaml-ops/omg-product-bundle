# Concourse `omg` plugin

This plugin extends `omg` to be able to deploy the Concourse product.

## Download OMG

https://github.com/enaml-ops/omg-cli/releases/latest

## Download Concourse Plugin

https://github.com/enaml-ops/omg-product-bundle/releases/latest

Don't forget to make the plugin executable before installing it.

## Install Plugin

    omg register-plugin -type product -pluginpath ./concourse-plugin

Verify installation with:

    omg list-products

## Deploy Concourse

```
./omg-osx deploy-product \
--bosh-url <bosh-url> \
--bosh-port 25555 \
--bosh-user <bosh-user> \
--bosh-pass <bosh-pwd> \
--ssl-ignore \
--print-manifest \
concourse-plugin-osx \
--web-ip 192.168.10.31 \
--web-ip 192.168.10.32 \
--web-vm-type medium \
--worker-vm-type large.cpu \
--database-vm-type large.cpu \
--network-name concourse \
--concourse-ip 10.0.100.1 \
--concourse-username concourse \
--az z1 \
--database-storage-type medium \
--worker-instance-count 2
```
## Tips and tricks

- Set `LOG_LEVEL=debug` for verbose output
- Add the `--print-manifest` flag to see the manifest you are about to deploy:

    `omg deploy-product --print-manifest ...`
