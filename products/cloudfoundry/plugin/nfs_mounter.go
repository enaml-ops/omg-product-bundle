package cloudfoundry

import (
	"github.com/codegangsta/cli"
	"github.com/enaml-ops/enaml"
	nfsmounterlib "github.com/enaml-ops/omg-product-bundle/products/cloudfoundry/enaml-gen/nfs_mounter"
	"github.com/xchapter7x/lo"
)

//NewNFSMounter - Creates a new NFS Mounter
func NewNFSMounter(c *cli.Context) *NFSMounter {
	return &NFSMounter{
		NFSServerAddress: c.String("nfs-server-address"),
		SharePath:        c.String("nfs-share-path"),
	}
}

//CreateJob - Create the yaml job structure for NFSMounter
func (s *NFSMounter) CreateJob() enaml.InstanceJob {
	return enaml.InstanceJob{
		Name:    "nfs_mounter",
		Release: "cf",
		Properties: &nfsmounterlib.NfsMounterJob{
			NfsServer: &nfsmounterlib.NfsServer{
				Address: s.NFSServerAddress,
				Share:   s.SharePath,
			},
		},
	}
}

func (s *NFSMounter) hasValidValues() bool {

	lo.G.Debugf("checking '%s' for valid flags", "nfs mounter")

	if s.NFSServerAddress == "" {
		lo.G.Debugf("could not find a valid NFSServerAddress '%v'", s.NFSServerAddress)
	}

	if s.SharePath == "" {
		lo.G.Debugf("could not find a valid SharePath '%v'", s.SharePath)
	}
	return s.NFSServerAddress != "" &&
		s.SharePath != ""
}
