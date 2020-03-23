package znc

import (
	"testing"
	zncv1 "znc-operator/pkg/apis/znc/v1"
)

func TestRenderConfiguration(t *testing.T) {
	_, err := RenderConfiguration(&zncv1.ZNCSpec{
		Version: "1.7.5",
		Config: zncv1.ZNCSpecConfig{
			AnonIPLimit:  0,
			ConnectDelay: 0,
			HideVersion:  false,
			LoadModules: []string{
				"webadmin",
				"modperl",
				"modpython",
			},
			Users: []zncv1.ZNCSpecConfigUser{
				{
					Admin: false,
					LoadModules: []string{
						"controlpanel",
						"chansaver",
					},
				},
			},
		},
	})
	if err != nil {
		t.Error("rendering config caused an unexpected error", err)
	}
}
