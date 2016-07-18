package director 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Http struct {

	/*Port - Descr: Port of the Registry to connect to Default: 25777
*/
	Port interface{} `yaml:"port,omitempty"`

	/*Password - Descr: Password to access the Registry Default: <nil>
*/
	Password interface{} `yaml:"password,omitempty"`

	/*User - Descr: User to access the Registry Default: <nil>
*/
	User interface{} `yaml:"user,omitempty"`

}