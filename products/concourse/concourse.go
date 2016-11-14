package concourse

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/atc"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/baggageclaim"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/garden"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/groundcrew"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/postgresql"
	"github.com/enaml-ops/omg-product-bundle/products/concourse/enaml-gen/tsa"
	"github.com/xchapter7x/lo"
)

const (
	concourseReleaseName string = "concourse"
	gardenReleaseName    string = "garden-runc"
)

//Deployment -
type Deployment struct {
	//enaml.Deployment
	manifest            *enaml.DeploymentManifest
	ConcourseURL        string
	ConcourseUserName   string
	ConcoursePassword   string
	NetworkName         string
	WebIPs              []string
	AZs                 []string
	WorkerInstances     int
	DeploymentName      string
	PostgresPassword    string
	WebVMType           string `omg:"web-vm-type"`
	WorkerVMType        string
	DatabaseVMType      string
	DatabaseStorageType string
	ConcourseReleaseVer string
	ConcourseReleaseURL string
	ConcourseReleaseSHA string
	StemcellVersion     string
	StemcellAlias       string
	StemcellOS          string
	GardenReleaseVer    string
	GardenReleaseURL    string
	GardenReleaseSHA    string
	TLSCert             string
	TLSKey              string
}

//NewDeployment -
func NewDeployment() (d Deployment) {
	d = Deployment{}
	d.manifest = new(enaml.DeploymentManifest)
	return
}

func (d *Deployment) doCloudConfigValidation(data []byte) (err error) {
	lo.G.Debug("Cloud Config:", string(data))
	c := &enaml.CloudConfigManifest{}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}

	for _, azName := range d.AZs {
		if !c.ContainsAZName(azName) {
			err = fmt.Errorf("AZ [%s] is not defined as a AZ in cloud config", azName)
			return
		}
	}

	if !c.ContainsVMType(d.WebVMType) {
		err = fmt.Errorf("WebVMType[%s] is not defined as a VMType in cloud config", d.WebVMType)
		return
	}
	if !c.ContainsVMType(d.WorkerVMType) {
		err = fmt.Errorf("WorkerVMType[%s] is not defined as a VMType in cloud config", d.WorkerVMType)
		return
	}
	if !c.ContainsVMType(d.DatabaseVMType) {
		err = fmt.Errorf("DatabaseVMType[%s] is not defined as a VMType in cloud config", d.DatabaseVMType)
		return
	}
	if !c.ContainsDiskType(d.DatabaseStorageType) {
		err = fmt.Errorf("DatabaseStorageType[%s] is not defined as a DiskType in cloud config", d.DatabaseStorageType)
		return
	}
	return
}

//Initialize -
func (d *Deployment) Initialize(cloudConfig []byte) error {

	err := d.doCloudConfigValidation(cloudConfig)
	if err != nil {
		return err
	}

	if d.ConcoursePassword == "" {
		return fmt.Errorf("Must supply concourse password")
	}

	var web *enaml.InstanceGroup
	var db *enaml.InstanceGroup
	var worker *enaml.InstanceGroup
	d.manifest.SetName(d.DeploymentName)
	d.manifest.AddRelease(enaml.Release{
		Name:    concourseReleaseName,
		URL:     d.ConcourseReleaseURL,
		SHA1:    d.ConcourseReleaseSHA,
		Version: d.ConcourseReleaseVer,
	})
	d.manifest.AddRelease(enaml.Release{
		Name:    gardenReleaseName,
		URL:     d.GardenReleaseURL,
		SHA1:    d.GardenReleaseSHA,
		Version: d.GardenReleaseVer,
	})
	d.manifest.AddStemcell(enaml.Stemcell{
		Alias:   d.StemcellAlias,
		OS:      d.StemcellOS,
		Version: d.StemcellVersion,
	})

	update := d.CreateUpdate()
	d.manifest.SetUpdate(update)

	if web, err = d.CreateWebInstanceGroup(); err != nil {
		return err
	}
	d.manifest.AddInstanceGroup(web)

	if db, err = d.CreateDatabaseInstanceGroup(); err != nil {
		return err
	}
	d.manifest.AddInstanceGroup(db)

	if worker, err = d.CreateWorkerInstanceGroup(); err != nil {
		return err
	}

	d.manifest.AddInstanceGroup(worker)
	return nil
}

//CreateWebInstanceGroup -
func (d *Deployment) CreateWebInstanceGroup() (web *enaml.InstanceGroup, err error) {

	web = &enaml.InstanceGroup{
		Name:      "web",
		Instances: len(d.WebIPs),
		VMType:    d.WebVMType,
		AZs:       d.AZs,
		Stemcell:  d.StemcellAlias,
	}
	web.AddNetwork(enaml.Network{
		Name:      d.NetworkName,
		StaticIPs: d.WebIPs,
	})
	web.AddJob(d.CreateAtcJob())
	web.AddJob(d.CreateTsaJob())

	return
}

//CreateAtcJob -
func (d *Deployment) CreateAtcJob() (job *enaml.InstanceJob) {
	props := atc.AtcJob{
		ExternalUrl:        d.ConcourseURL,
		BasicAuthUsername:  d.ConcourseUserName,
		BasicAuthPassword:  d.ConcoursePassword,
		PostgresqlDatabase: "atc",
	}
	if d.TLSCert != "" {
		props.TlsCert = d.TLSCert
	}
	if d.TLSKey != "" {
		props.TlsKey = d.TLSKey
	}
	job = enaml.NewInstanceJob("atc", concourseReleaseName, props)
	return
}

//CreateTsaJob -
func (d *Deployment) CreateTsaJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("tsa", concourseReleaseName, tsa.TsaJob{})
	return
}

//CreateDatabaseInstanceGroup -
func (d *Deployment) CreateDatabaseInstanceGroup() (db *enaml.InstanceGroup, err error) {

	db = &enaml.InstanceGroup{
		Name:               "db",
		Instances:          1,
		PersistentDiskType: d.DatabaseStorageType,
		VMType:             d.DatabaseVMType,
		AZs:                d.AZs,
		Stemcell:           d.StemcellAlias,
	}
	db.AddNetwork(d.CreateNetwork())
	db.AddJob(d.CreatePostgresqlJob())

	return
}

//CreatePostgresqlJob -
func (d *Deployment) CreatePostgresqlJob() (job *enaml.InstanceJob) {
	dbs := make([]DBName, 1)
	dbs[0] = DBName{
		Name:     "atc",
		Role:     "atc",
		Password: d.PostgresPassword,
	}
	job = enaml.NewInstanceJob("postgresql", concourseReleaseName, postgresql.PostgresqlJob{
		Databases: dbs,
	})
	return
}

//CreateWorkerInstanceGroup -
func (d *Deployment) CreateWorkerInstanceGroup() (worker *enaml.InstanceGroup, err error) {
	worker = &enaml.InstanceGroup{
		Name:      "worker",
		Instances: d.WorkerInstances,
		VMType:    d.WorkerVMType,
		AZs:       d.AZs,
		Stemcell:  d.StemcellAlias,
	}

	worker.AddNetwork(d.CreateNetwork())
	worker.AddJob(d.CreateGroundCrewJob())
	worker.AddJob(d.CreateBaggageClaimJob())
	worker.AddJob(d.CreateGardenJob())

	return
}

//CreateGardenJob -
func (d *Deployment) CreateGardenJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("garden", gardenReleaseName, Garden{
		garden.Garden{
			ListenAddress:   "0.0.0.0:7777",
			ListenNetwork:   "tcp",
			AllowHostAccess: true,
		},
	})
	return
}

//CreateBaggageClaimJob -
func (d *Deployment) CreateBaggageClaimJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("baggageclaim", concourseReleaseName, baggageclaim.BaggageclaimJob{})
	return
}

//CreateGroundCrewJob -
func (d *Deployment) CreateGroundCrewJob() (job *enaml.InstanceJob) {
	job = enaml.NewInstanceJob("groundcrew", concourseReleaseName, groundcrew.GroundcrewJob{})
	return
}

//CreateNetwork -
func (d *Deployment) CreateNetwork() (network enaml.Network) {
	network = enaml.Network{
		Name: d.NetworkName,
	}
	return
}

//CreateUpdate -
func (d *Deployment) CreateUpdate() (update enaml.Update) {
	update = enaml.Update{
		Canaries:        1,
		MaxInFlight:     3,
		Serial:          false,
		CanaryWatchTime: "1000-60000",
		UpdateWatchTime: "1000-60000",
	}

	return
}

func (d Deployment) isStrongPass(pass string) (ok bool) {
	ok = false
	if len(pass) > 8 {
		ok = true
	}
	return
}

func insureHAInstanceCount(instances int) int {
	if instances < 2 {
		instances = 2
	}
	return instances
}

//GetDeployment -
func (d Deployment) GetDeployment() enaml.DeploymentManifest {
	return *d.manifest
}
