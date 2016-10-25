# remove existing enaml-gen, minus any custom marshalling code we've written
find ./enaml-gen -type f | grep -v "_marshal" | grep -v "_test.go" | xargs rm -f

# remove any leftover empty directories
find ./enaml-gen -type d -empty -delete

# regenerate enaml structs
enaml generate https://bosh.io/d/github.com/cloudfoundry/cf-release?v=245
enaml generate https://bosh.io/d/github.com/cloudfoundry/diego-release?v=0.1487.0
enaml generate https://bosh.io/d/github.com/cloudfoundry-incubator/garden-linux-release?v=0.342.0
enaml generate https://bosh.io/d/github.com/cloudfoundry/uaa-release?v=20
