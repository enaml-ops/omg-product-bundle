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

    omg deploy-product \
        --bosh-url <BOSH_URL> \
        --bosh-port <BOSH_PORT> \
        --bosh-user <BOSH_USER> \
        --bosh-pass <BOSH_PASS> \
        --ssl-ignore \
        concourse-plugin \
            --web-vm-type small \
            --worker-vm-type small \
            --database-vm-type small \
            --network-name private \
            --url <CONCOURSE_URL> \
            --username <CONCOURSE_USER> \
            --password <CONCOURSE_PASS> \
            --web-instances 1 \
            --web-azs z1 \
            --worker-azs z1 \
            --database-azs z1 \
            --postresql-db-pwd <PASSWORD> \
            --database-storage-type medium \
            --bosh-stemcell-alias trusty \
            --remote-stemcell-url https://d26ekeud912fhb.cloudfront.net/bosh-stemcell/aws/light-bosh-stemcell-3262.2-aws-xen-hvm-ubuntu-trusty-go_agent.tgz \
            --stemcell-ver 3262.2 \
            --remote-stemcell-sha 64234353e233be1630f6f033c85f0a9fea21b25e

## Tips and tricks

- Set `LOG_LEVEL=debug` for verbose output
- Add the `--print-manifest` flag to see the manifest you are about to deploy:

    `omg deploy-product --print-manifest ...`
