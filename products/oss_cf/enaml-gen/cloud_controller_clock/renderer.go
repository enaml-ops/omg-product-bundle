package cloud_controller_clock 
/*
* File Generated by enaml generator
* !!! Please do not edit this file !!!
*/
type Renderer struct {

	/*DefaultResultsPerPage - Descr: Default number of results returned per page if user does not specify Default: 50
*/
	DefaultResultsPerPage interface{} `yaml:"default_results_per_page,omitempty"`

	/*MaxResultsPerPage - Descr: Maximum number of results returned per page Default: 100
*/
	MaxResultsPerPage interface{} `yaml:"max_results_per_page,omitempty"`

	/*MaxInlineRelationsDepth - Descr: Maximum depth of inlined relationships in the result Default: 2
*/
	MaxInlineRelationsDepth interface{} `yaml:"max_inline_relations_depth,omitempty"`

}