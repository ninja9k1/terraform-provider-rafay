package rafay

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/RafaySystems/rctl/pkg/agent"
	"github.com/RafaySystems/rctl/pkg/commands"
	"github.com/RafaySystems/rctl/pkg/config"
	"github.com/RafaySystems/rctl/pkg/models"
	"github.com/RafaySystems/rctl/pkg/project"
	"gopkg.in/yaml.v3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentCreate,
		ReadContext:   resourceAgentRead,
		UpdateContext: resourceAgentUpdate,
		DeleteContext: resourceAgentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"projectname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agent_filepath": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAgentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	filePath := d.Get("agent_filepath").(string)
	var agentYaml commands.AgentYamlConfig
	//make sure this is the correct file path
	log.Println("filepath: ", filePath)
	log.Printf("create gitops agent resource")
	//get project id with project name, p.id used to refer to project id -> need p.ID for calling createAgent
	resp, err := project.GetProjectByName(d.Get("projectname").(string))
	if err != nil {
		log.Printf("project does not exist, error %s", err.Error())
		return diag.FromErr(err)
	}
	p, err := project.NewProjectFromResponse([]byte(resp))
	if err != nil {
		return diag.FromErr(err)
	} else if p == nil {
		d.SetId("")
		return diags
	}
	projectId := p.ID
	//open file path and retirve config spec from yaml file (from run function in commands/create_agent.go)
	//read and capture file from file path
	if f, err := os.Open(filePath); err == nil {
		c, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println("error cpaturing file")
		}
		//implement createClusterOverride from commands/create_cluster_override.go -> then call clusteroverride.CreateClusterOverride
		agentDefinition := c
		err = yaml.Unmarshal(agentDefinition, &agentYaml)
		if err != nil {
			log.Printf("Failed Unmarshal correctly")
		}
		// check if project is provided from yaml file
		if agentYaml.Metadata.Project != "" {
			projectId, err = config.GetProjectIdByName(agentYaml.Metadata.Project)
			if err != nil {
				log.Println("error getting project ID from yaml file")
			}
		}
		// create the agent
		if strings.ToLower(agentYaml.Spec.Template.Type) == strings.ToLower(string(agent.ClusterAgent)) {
			spec := models.AgentSpec{
				AgentType:   models.AgentType(agent.ClusterAgent),
				ClusterName: agentYaml.Spec.Template.ClusterName,
			}
			if agentYaml.Spec.Template.ClusterName == "" {
				log.Println("you must provide a clusterName in the config")
			}
			err = agent.CreateAgent(agentYaml.Metadata.Name, projectId, spec)
		}

		if err != nil {
			log.Printf("Failed to create Agent: %s\n", agentYaml.Metadata.Name)
		} else {
			log.Printf("Successfully created Agent: %s\n", agentYaml.Metadata.Name)
		}
	}
	//set id to metadata.Name
	d.SetId(agentYaml.Metadata.Name)
	return diags
}

func resourceAgentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceAgentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Println("rctl doesn't have update functionality for agent")
	return diags
}

func resourceAgentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	log.Printf("resource agent delete id %s", d.Id())

	//get project id with project name, p.id used to refer to project id -> need p.ID for calling DeleteAgent
	resp, err := project.GetProjectByName(d.Get("projectname").(string))
	if err != nil {
		log.Printf("project does not exist, error %s", err.Error())
		return diag.FromErr(err)
	}
	p, err := project.NewProjectFromResponse([]byte(resp))
	if err != nil {
		return diag.FromErr(err)
	} else if p == nil {
		d.SetId("")
		return diags
	}
	projectId := p.ID
	//convert namesapce interface to passable list for function
	err = agent.DeleteAgent(d.Id(), projectId)
	if err != nil {
		log.Println("error deleting agent")
	} else {
		log.Println("Deleted agent: ", d.Id())
	}
	return diags
}
