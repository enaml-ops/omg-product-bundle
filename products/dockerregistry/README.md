# Docker Registry `omg` plugin

This plugin extends `omg` to be able to deploy the a HA docker-registry with NFS server shared by 1-N registry nodes and a Proxy to route https traffic to registry nodes running on port 5000.

## Download OMG

https://github.com/enaml-ops/omg-cli/releases/latest

## Download Docker Registry Plugin

https://github.com/enaml-ops/omg-product-bundle/releases/latest

Don't forget to make the plugin executable before installing it.

## Install Plugin

    omg-osx register-plugin -type product -pluginpath ./docker-registry-plugin-osx

Verify installation with:

    omg-osx list-products

## Deploy Docker Registry

```
./omg-osx deploy-product \
--bosh-url <bosh-url> \
--bosh-port 25555 \
--bosh-user <bosh-user> \
--bosh-pass <bosh-pwd> \
--ssl-ignore \
docker-registry-plugin-osx \
--proxy-ip 192.168.0.10 \
--registry-ip 192.168.0.11 \
--registry-ip 192.168.0.12 \
--nfs-server-ip 192.168.0.13 \
--registry-vm-type medium \
--proxy-vm-type medium \
--nfs-server-vm-type medium \
--nfs-server-disk-type large \
--network-name private \
--az z1 \
--public-host-name <optional public host/ip if nat'ed>
```
## Tips and tricks

- Set `LOG_LEVEL=debug` for verbose output
- Add the `--print-manifest` flag to see the manifest you are about to deploy:
