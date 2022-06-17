package yamlrawmessage

type YAMLRawMessage struct {
	unmarshal func(interface{}) error
}

func (rmsg *YAMLRawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rmsg.unmarshal = unmarshal
	return nil
}

func (rmsg *YAMLRawMessage) Unmarshal(v interface{}) {
	_ = rmsg.unmarshal(v)
}
