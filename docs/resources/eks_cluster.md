---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "rafay_eks_cluster Resource - terraform-provider-rafay"
subcategory: ""
description: |-
  
---

# rafay_eks_cluster (Resource)



## Example Usage

```terraform
resource "rafay_eks_cluster" "cluster" {
  name         = "demo-terraform2"
  projectname  = "dev3"
  yamlfilepath = "<file-path/eks-cluster.yaml>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)
- **projectname** (String)
- **yamlfilepath** (String)

### Optional

- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
- **update** (String)

