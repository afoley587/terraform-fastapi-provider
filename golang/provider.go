// provider.go
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"garbage_des_encrypt": desEncrypt(),
			"garbage_des_decrypt": desDecrypt(),
		},
	}
}
