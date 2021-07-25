package rafay

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"path/filepath"

	rctlconfig "github.com/RafaySystems/rctl/pkg/config"
	rctlcontext "github.com/RafaySystems/rctl/pkg/context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func New(_ string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"provider_config_file": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RAFAY_PROVIDER_CONFIG", "~/.rafay/cli/config.json"),
				},
				"ignore_insecure_tls_error": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"rafay_project":          resourceProject(),
				"rafay_group":            resourceGroup(),
				"rafay_groupassociation": resourceGroupAssociation(),
				"rafay_cloud_credential": resourceCloudCredential(),
				"rafay_eks_cluster":      resourceEKSCluster(),
				"rafay_import_cluster":   resourceImportCluster(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				/*
					"rafay_apikey":            dataSourceUser(),
					"rafay_project":           dataSourceProject(),
					"rafay_group":             dataSourceGroup(),
					"rafay_cluster_blueprint": dataSourceClusterBlueprint(),

					"rafay_cloudaccount_aws":     dataSourceCloudAccountAws(),
					"rafay_cloudaccount_azure":   dataSourceCloudAccountAzure(),
					"rafay_cloudaccount_gcp":     dataSourceCloudAccountGcp(),
					"rafay_cloudaccount_vsphere": dataSourceCloudAccountVsphere(),
				*/
			},
			ConfigureContextFunc: providerConfigure,
		}

		return p
	}
}

func providerConfigure(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config_file := rd.Get("provider_config_file").(string)
	ignoreTlsError := rd.Get("ignore_insecure_tls_error").(bool)

	log.Printf("rafay provider config file %s", config_file)
	var diags diag.Diagnostics

	configPath := filepath.Dir(config_file)
	fileName := filepath.Base(config_file)
	cliCtx := rctlcontext.GetContext()
	cliCtx.ConfigFile = fileName
	cliCtx.ConfigDir = configPath
	err := rctlconfig.InitConfig(cliCtx)

	if err != nil {
		log.Printf("rafay provider config init error %s", err.Error())
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create rafay provider",
			Detail:   "Unable to init config for authenticated rafay provider",
		})
		return nil, diags
	}

	if ignoreTlsError {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return rctlconfig.GetConfig(), diags

}
