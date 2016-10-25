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


## store data in vault

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
  "nfs-server-address":"10.0.0.40",
}
$> ./vault write ${VAULT_HASH} @vault-ip.json
$> ./vault read ${VAULT_HASH}
```

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
--stemcell-name ubuntu-trusty \
--infer-from-cloud \
--nfs-share-path /var/vcap/store \
--doppler-drain-buffer-size 256 \
--nfs-allow-from-network-cidr 10.0.0.1/24
```

## everything else
*typically one would want to just use the cloudfoundry plugin flag to rotate
passwords and certs rather than populating those fields manually*

## more on vault integration

```
...
--vault-domain "http://127.0.0.1:8200" \
--vault-hash-password "secret/pcf-np-1-password" \
--vault-hash-keycert "secret/pcf-np-1-keycert" \
--vault-hash-host "secret/pcf-np-1-hostname" \
--vault-hash-ip "secret/pcf-np-1-ips" \
--vault-token "xxxx-xxxx-xxxx-xxxx-xxxxxxxx" \
--vault-rotate \
```
*you might notice the optional use of vault flags. this does a few things.*
- Flags:
  - `--vault-rotate` : when used will auto-populate the values in vault for the secret stores defined at flags `vault-hash-password` & `vault-hash-keycert`
    - when this flag is set, you will get a unique set of passwords, usernames, certs for all relevant cloud foundry values sent to the targetted vault
  - `--vault-domain` : this is where your vault lives
  - `--vault-token` : this is the token to access your targetted vault and secret stores
  - `--vault-hash-password` : a vault bucket the user defines which will store  passwords/usernames/secrets for pcf installation
    - the values in this store will be used to autopopulate required fields for the PCF plugin, so the user doesnt have to enter values manually
    - the values in this store will be automatically re/generated & overwritten if the `vault-rotate` flag is set
  - `--vault-hash-keycert` : a vault bucket the user defines which will store  self-signed internal keys and certificates for pcf installation
    - the values in this store will be used to autopopulate required fields for the PCF plugin, so the user doesnt have to enter values manually
    - the values in this store will be automatically re/generated & overwritten if the `vault-rotate` flag is set
  - `--vault-hash-ip` : a vault bucket the user defines which will store ip & other information for pcf installation
    - a user can leverage this store to populate any key/value pair they wish.
    - the keys in this store will be mapped to any matching flagname, and the value will be used to set that flags value (so the user does not have to enter it)
    - these values will not be overwritten or rotated regardless of if the `vault-rotate` flag is set or not
  - `--vault-hash-host` : a vault bucket the user defines which will store host and other information for pcf installation
    - a user can leverage this store to populate any key/value pair they wish.
    - the keys in this store will be mapped to any matching flagname, and the value will be used to set that flags value (so the user does not have to enter it)
    - these values will not be overwritten or rotated regardless of if the `vault-rotate` flag is set or not


## more on cloud config integration

```
...
--infer-from-cloud \
```
*you might notice the optional use of infer from cloud. this allows the user to let the plugin try to guess options from the cloud config*
- Flags: 
  - `--infer-from-cloud` : this will pull in the targetted bosh's cloud config and use the values defined there to set the plugin flag argument values (user doesnt need to set flags manually)
    - (disktype, vmtype, az, network) information will all be used to populate the following flags:
      - "mysql-disk-type"
      - "diego-db-disk-type"
      - "diego-cell-disk-type"
      - "diego-brain-disk-type"
      - "etcd-disk-type"
      - "nfs-disk-type"
      - "diego-brain-vm-type"
      - "errand-vm-type"
      - "clock-global-vm-type"
      - "doppler-vm-type"
      - "uaa-vm-type"
      - "diego-cell-vm-type"
      - "diego-db-vm-type"
      - "router-vm-type"
      - "haproxy-vm-type"
      - "nats-vm-type"
      - "consul-vm-type"
      - "etcd-vm-type"
      - "nfs-vm-type"
      - "mysql-vm-type"
      - "mysql-proxy-vm-type"
      - "cc-worker-vm-type"
      - "cc-vm-type"
      - "loggregator-traffic-controller-vmtype"
      - "bootstrap-vm-type"
      - "acceptance-tests-vm-type"
      - "az"
      - "network"
