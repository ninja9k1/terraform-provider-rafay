package rafay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RafaySystems/rctl/pkg/cluster"
	"github.com/RafaySystems/rctl/pkg/clusterctl"
	"github.com/RafaySystems/rctl/pkg/config"
	glogger "github.com/RafaySystems/rctl/pkg/log"
	"github.com/RafaySystems/rctl/pkg/project"
	"github.com/RafaySystems/rctl/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/go-yaml/yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type configMetadata struct {
	Name    string `yaml:"name"`
	Project string `yaml:"project"`
}

type configResourceType struct {
	Meta *configMetadata `yaml:"metadata"`
}

func resourceEKSCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEKSClusterCreate,
		ReadContext:   resourceEKSClusterRead,
		UpdateContext: resourceEKSClusterUpdate,
		DeleteContext: resourceEKSClusterDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"yamlfilepath": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"projectname": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func findResourceNameFromConfig(configBytes []byte) (string, string, error) {
	var config configResourceType
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return "", "", nil
	} else if config.Meta == nil {
		return "", "", fmt.Errorf("%s", "Invalid resource: No metadata found")
	} else if config.Meta.Name == "" {
		return "", "", fmt.Errorf("%s", "Invalid resource: No name specified in metadata")
	}
	return config.Meta.Name, config.Meta.Project, nil
}

func collateConfigsByName(rafayConfigs, clusterConfigs [][]byte) (map[string][]byte, []error) {
	var errs []error
	configsMap := make(map[string][][]byte)
	// First find all rafay spec configurations
	for _, config := range rafayConfigs {
		name, _, err := findResourceNameFromConfig(config)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, ok := configsMap[name]; ok {
			errs = append(errs, fmt.Errorf(`duplicate "cluster" resource with name "%s" found`, name))
			continue
		}
		configsMap[name] = append(configsMap[name], config)
	}
	// Then append the cluster specific configurations
	for _, config := range clusterConfigs {
		name, _, err := findResourceNameFromConfig(config)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, ok := configsMap[name]; !ok {
			errs = append(errs, fmt.Errorf(`error finding "Cluster" configuration for name "%s"`, name))
			continue
		}
		configsMap[name] = append(configsMap[name], config)
	}
	// Remove any configs that don't have the tail end (cluster related configs)
	result := make(map[string][]byte)
	for name, configs := range configsMap {
		if len(configs) <= 1 {
			errs = append(errs, fmt.Errorf(`no "ClusterConfig" found for cluster "%s"`, name))
			continue
		}
		collatedConfigBytes, err := utils.JoinYAML(configs)
		if err != nil {
			errs = append(errs, fmt.Errorf(`error collating YAML files for cluster "%s": %s`, name, err))
			continue
		}
		result[name] = collatedConfigBytes
		log.Printf(`final Configuration for cluster "%s": %#v`, name, string(collatedConfigBytes))
	}
	return result, errs
}

func resourceEKSClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	log.Printf("create EKS cluster resource")
	c := config.GetConfig()
	logger := glogger.GetLogger()

	YamlConfigFilePath := d.Get("yamlfilepath").(string)

	log.Printf("yaml file path  %s", YamlConfigFilePath)

	fileBytes, err := utils.ReadYAMLFileContents(YamlConfigFilePath)
	if err != nil {
		return diag.FromErr(err)
	}

	// split the file and update individual resources
	y, uerr := utils.SplitYamlAndGetListByKind(fileBytes)
	if uerr != nil {
		return diag.FromErr(err)
	}

	var rafayConfigs, clusterConfigs [][]byte
	rafayConfigs = y["Cluster"]
	clusterConfigs = y["ClusterConfig"]
	if len(rafayConfigs) > 1 {
		return diag.FromErr(fmt.Errorf("%s", "only one cluster per config is supported"))
	}
	for _, yi := range rafayConfigs {
		log.Println("rafayConfig:", string(yi))
		name, project, err := findResourceNameFromConfig(yi)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s", "failed to get cluster name"))
		}
		log.Println("rafayConfig name:", name, "project:", project)
		if name != d.Get("name").(string) {
			return diag.FromErr(fmt.Errorf("%s", "cluster name does not match config file "))
		}
		if project != d.Get("projectname").(string) {
			return diag.FromErr(fmt.Errorf("%s", "project name does not match config file"))
		}
	}

	for _, yi := range clusterConfigs {
		log.Println("clusterConfig", string(yi))
		name, _, err := findResourceNameFromConfig(yi)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s", "failed to get cluster name"))
		}
		if name != d.Get("name").(string) {
			return diag.FromErr(fmt.Errorf("%s", "ClusterConfig name does not match config file"))
		}
	}

	configMap, errs := collateConfigsByName(rafayConfigs, clusterConfigs)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println("error in collateConfigsByName", err)
		}
		return diag.FromErr(fmt.Errorf("%s", "failed in collateConfigsByName"))
	}
	// Make request
	for clusterName, configBytes := range configMap {
		log.Println("create cluster:", clusterName, "config:", string(configBytes))
		if err := clusterctl.Apply(logger, c, clusterName, configBytes, true); err != nil {
			return diag.FromErr(fmt.Errorf("error performing apply on cluster %s: %s", clusterName, err))
		}
	}

	// get project details
	resp, err := project.GetProjectByName(d.Get("projectname").(string))
	if err != nil {
		fmt.Print("project does not exist")
		return diags
	}
	project, err := project.NewProjectFromResponse([]byte(resp))
	if err != nil {
		fmt.Printf("project does not exist")
		return diags
	}

	s, err := cluster.GetCluster(d.Get("name").(string), project.ID)
	if err != nil {
		log.Printf("error while getCluster %s", err.Error())
		return diag.FromErr(err)
	}

	log.Printf("resource eks cluster created %s", s.ID)
	d.SetId(s.ID)

	return diags
}

func resourceEKSClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	resp, err := project.GetProjectByName(d.Get("projectname").(string))
	if err != nil {
		fmt.Print("project name missing in the resource")
		return diags
	}

	project, err := project.NewProjectFromResponse([]byte(resp))
	if err != nil {
		fmt.Printf("project does not exist")
		return diags
	}
	_, err = cluster.GetCluster(d.Get("name").(string), project.ID)
	if err != nil {
		log.Printf("error in get cluster %s", err.Error())
		return diag.FromErr(err)
	}

	return diags
}

func resourceEKSClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("update EKS cluster resource")
	return diags
}

func resourceEKSClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("resource cluster delete id %s", d.Id())

	resp, err := project.GetProjectByName(d.Get("projectname").(string))
	if err != nil {
		fmt.Print("project  does not exist")
		return diags
	}

	project, err := project.NewProjectFromResponse([]byte(resp))
	if err != nil {
		fmt.Printf("project  does not exist")
		return diags
	}

	errDel := cluster.DeleteCluster(d.Get("name").(string), project.ID)
	if errDel != nil {
		log.Printf("delete cluster error %s", errDel.Error())
		return diag.FromErr(errDel)
	}

	return diags
}
