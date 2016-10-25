package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Scim struct {

	/*Users - Descr: A list of users to be bootstrapped with authorities.
Each entry supports the following format:
  Short OpenStruct:
    - name: username
      password: password
      groups:
        - group1
        - group2
  Long OpenStruct:
    - name: username
      password: password
      groups:
        - group1
        - group2
      firstName: first name
      lastName: lastName
      email: email
      origin: origin-value - most commonly uaa
 Default: <nil>
*/
	Users interface{} `yaml:"users,omitempty"`

	/*UseridsEnabled - Descr: Enables the endpoint `/ids/Users` that allows consumers to translate user ids to name Default: true
*/
	UseridsEnabled interface{} `yaml:"userids_enabled,omitempty"`

	/*ExternalGroups - Descr: External group mappings. Either formatted as an OpenStruct.
As an OpenStruct, the mapping additionally specifies an origin to which the mapping is applied:
  origin1:
    external_group1:
      - internal_group1
      - internal_group2
      - internal_group3
    external_group2:
      - internal_group2
      - internal_group4
  origin2:
    external_group3:
      - internal_group3
      - internal_group4
      - internal_group5
 Default: <nil>
*/
	ExternalGroups interface{} `yaml:"external_groups,omitempty"`

	/*Groups - Descr: Contains a hash of group names and their descriptions. These groups will be added to the UAA database for the default zone but not associated with any user.
Example:
  uaa:
    scim:
      groups:
        my-test-group: 'My test group description'
        another-group: 'Another group description'
 Default: <nil>
*/
	Groups interface{} `yaml:"groups,omitempty"`

	/*User - Descr: If true override users defined in uaa.scim.users found in the database. Default: true
*/
	User *ScimUser `yaml:"user,omitempty"`

}