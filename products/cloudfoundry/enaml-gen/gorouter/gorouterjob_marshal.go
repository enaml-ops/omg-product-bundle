package gorouter

// MarshalYAML implements the yaml.Marshaler interface.
func (j *GorouterJob) MarshalYAML() (interface{}, error) {
	result := make(map[string]interface{})
	if j.Nats != nil {
		result["nats"] = j.Nats
	}
	if j.MetronEndpoint != nil {
		result["metron_endpoint"] = j.MetronEndpoint
	}
	if j.RequestTimeoutInSeconds != nil {
		result["request_timeout_in_seconds"] = j.RequestTimeoutInSeconds
	}
	if j.Uaa != nil {
		result["uaa"] = j.Uaa
	}
	if j.Dropsonde != nil {
		result["dropsonde"] = j.Dropsonde
	}
	if j.Router != nil {
		result["router"] = j.Router
	}

	// The routing API is the tricky part that enaml can't solve alone.
	// Some of the fields are under "routing_api", and others are under
	// "routing-api".
	if j.RoutingApi != nil {
		result["routing_api"] = map[string]interface{}{
			"enabled": j.RoutingApi.Enabled,
		}
		result["routing-api"] = map[string]interface{}{
			"port":          j.RoutingApi.Port,
			"auth_disabled": j.RoutingApi.AuthDisabled,
		}
	}
	return result, nil
}
