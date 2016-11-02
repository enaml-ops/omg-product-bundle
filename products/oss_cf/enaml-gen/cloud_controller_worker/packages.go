package cloud_controller_worker 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Packages struct {

	/*WebdavConfig - Descr: The ca cert to use when communicating with webdav Default: 
*/
	WebdavConfig *PackagesWebdavConfig `yaml:"webdav_config,omitempty"`

	/*AppPackageDirectoryKey - Descr: Directory (bucket) used store app packages.  It does not have be pre-created. Default: cc-packages
*/
	AppPackageDirectoryKey interface{} `yaml:"app_package_directory_key,omitempty"`

	/*FogConnection - Descr: Fog connection hash Default: <nil>
*/
	FogConnection interface{} `yaml:"fog_connection,omitempty"`

	/*BlobstoreType - Descr: The type of blobstore backing to use. Valid values: ['fog', 'webdav'] Default: fog
*/
	BlobstoreType interface{} `yaml:"blobstore_type,omitempty"`

	/*MaxPackageSize - Descr: Maximum size of application package Default: 1073741824
*/
	MaxPackageSize interface{} `yaml:"max_package_size,omitempty"`

	/*FogAwsStorageOptions - Descr: Storage options passed to fog for aws blobstores. Valid keys: ['encryption']. Default: <nil>
*/
	FogAwsStorageOptions interface{} `yaml:"fog_aws_storage_options,omitempty"`

	/*Cdn - Descr: Key pair name for signed download URIs Default: 
*/
	Cdn *PackagesCdn `yaml:"cdn,omitempty"`

}