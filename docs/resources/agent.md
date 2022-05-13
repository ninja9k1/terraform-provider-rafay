---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rafay_agent Resource - terraform-provider-rafay"
subcategory: ""
description: |-
  
---

# rafay_agent (Resource)



## Example Usage

```terraform
resource "rafay_agent" "tfdemoagent1" {
  metadata {
    name    = "tfdemoagent1"
    project = "upgrade"
  }
  spec {
    type = "ClusterAgent"
    cluster {
      name = "dev-test"
    }
    active = true
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **metadata** (Block List, Max: 1) Metadata of the agent resource (see [below for nested schema](#nestedblock--metadata))
- **spec** (Block List, Max: 1) Specification of the agent resource (see [below for nested schema](#nestedblock--spec))
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Optional:

- **annotations** (Map of String) annotations of the resource
- **description** (String) description of the resource
- **labels** (Map of String) labels of the resource
- **name** (String) name of the resource
- **project** (String) Project of the resource


<a id="nestedblock--spec"></a>
### Nested Schema for `spec`

Optional:

- **active** (Boolean) flag to indicate if the agent is active
- **cluster** (Block List, Max: 1) metadata of cluster metadata (see [below for nested schema](#nestedblock--spec--cluster))
- **type** (String) type of agent

<a id="nestedblock--spec--cluster"></a>
### Nested Schema for `spec.cluster`

Optional:

- **name** (String) name of the cluster



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

