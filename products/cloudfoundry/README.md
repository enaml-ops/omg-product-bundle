# CF Product Plugin 
```
This is a product plugin for Elastic Runtime (CF)

The compiled version of this will
be runnable (like a tile) from the omg-cli
as a plugin.

Once executed it will spin up a cloud foundry
on whatever iaas your target bosh is
configured for.

```


## store ips in vault

```bash
$> # example to demonstrate how we would load up the vault with data
$> #set where your vault lives
$> VAULT_ADDR=http://192.168.0.1:8200
$> #set which hash your writing to
$> VAULT_HASH=secret/pcf-np-1-ips
$> cat vault-ip.json
{
  "nfs-ip": "10.0.0.31",
  "mysql-proxy-ip": "10.0.0.30",
  "mysql-ip": "10.0.0.29", 
  "haproxy-ip": "10.0.0.2", 
  "consul-ip": "10.0.0.23",
  "nats-machine-ip": "10.0.0.24", 
  "etcd-machine-ip": "10.0.0.25",
  "diego-brain-ip": "10.0.0.26",
  "doppler-ip": "10.0.0.27",
  "loggregator-traffic-controller-ip": "10.0.0.28",
  "router-ip": "10.0.0.20,10.0.0.21,10.0.0.22", 
  "diego-cell-ip": "10.0.0.1,10.0.0.2,10.0.0.3",
  "diego-db-ip": "10.0.1.1",
  "mysql-proxy-external-host": "10.0.0.1",
  "nfs-server-addres":"10.0.0.40",
  "bbs-api":"10.0.0.41"
}
$> ./vault write ${VAULT_HASH} @vault-ip.json
$> ./vault read ${VAULT_HASH}
```

## everything else
*typically one would want to just use the cloudfoundry plugin flag to rotate
passwords and certs rather than populating those fields manually*

##deploy cloud foundry
*remove the `--print-manifest` flag to actually deploy it to your target bosh*
```bash
$> omg deploy-product \
--print-manifest \
--bosh-url https://1.1.1.1 --bosh-port 25555 --bosh-user admin --bosh-pass pass --ssl-ignore \
cloudfoundry-plugin-osx \
--system-domain "sys.demo.io" \
--app-domain "app.demo.io" \
--vault-domain "http://127.0.0.1:8200" \
--vault-hash-password "secret/pcf-np-1-password" \
--vault-hash-keycert "secret/pcf-np-1-keycert" \
--vault-hash-host "secret/pcf-np-1-hostname" \
--vault-hash-ip "secret/pcf-np-1-ips" \
--vault-token "xxxx-xxxx-xxxx-xxxx-xxxxxxxx" \
--vault-rotate \
--az us-east-1c \
--stemcell-name ubuntu-trusty \
--network private \
--infer-from-cloud \
--nfs-share-path /var/vcap/store \
--doppler-drain-buffer-size 256 \
--nfs-allow-from-network-cidr 10.0.0.1/24
```
