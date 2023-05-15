package config

import "testing"

func TestLoadTemplateConfig(t *testing.T) {
	name := "../A3.json"
	c, err := LoadTemplateConfig(name)
	if err != nil {
		t.Fatal("err: ", err)
	}
	t.Logf("%#v", c)
}
