package uaa 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type LdapGroups struct {

	/*MaxSearchDepth - Descr: Set to number of levels a nested group search should go. Set to 1 to disable nested groups (default) Default: 1
*/
	MaxSearchDepth interface{} `yaml:"maxSearchDepth,omitempty"`

	/*GroupRoleAttribute - Descr: Used with groups-as-scopes, defines the attribute that holds the scope name(s). Default: <nil>
*/
	GroupRoleAttribute interface{} `yaml:"groupRoleAttribute,omitempty"`

	/*GroupSearchFilter - Descr: Search query filter to find the groups a user belongs to, or for a nested search, groups that a group belongs to Default: member={0}
*/
	GroupSearchFilter interface{} `yaml:"groupSearchFilter,omitempty"`

	/*SearchSubtree - Descr: Boolean value, set to true to search below the search base Default: true
*/
	SearchSubtree interface{} `yaml:"searchSubtree,omitempty"`

	/*SearchBase - Descr: Search start point for a user group membership search Default: 
*/
	SearchBase interface{} `yaml:"searchBase,omitempty"`

	/*ProfileType - Descr: What type of group integration should be used. Values are: 'no-groups', 'groups-as-scopes', 'groups-map-to-scopes' Default: no-groups
*/
	ProfileType interface{} `yaml:"profile_type,omitempty"`

	/*AutoAdd - Descr: Set to true when profile_type=groups_as_scopes to auto create scopes for a user. Ignored for other profiles. Default: true
*/
	AutoAdd interface{} `yaml:"autoAdd,omitempty"`

}