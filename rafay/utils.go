package rafay

import (
	"log"
	"sort"

	commonpb "github.com/RafaySystems/rafay-common/proto/types/hub/commonpb"
	"k8s.io/apimachinery/pkg/api/resource"
)

type File struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

/*
type UploadedYAMLArtifact struct {
	Paths []*File `protobuf:"bytes,1,rep,name=paths,proto3" json:"paths,omitempty"`
}

type UploadedHelmArtifact struct {
	ChartPath   *File   `protobuf:"bytes,1,opt,name=chartPath,proto3" json:"chartPath,omitempty"`
	ValuesPaths []*File `protobuf:"bytes,2,rep,name=valuesPaths,proto3" json:"valuesPaths,omitempty"`
}

type YAMLInGitRepo struct {
	Repository string  `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Revision   string  `protobuf:"bytes,2,opt,name=revision,proto3" json:"revision,omitempty"`
	Paths      []*File `protobuf:"bytes,3,rep,name=paths,proto3" json:"paths,omitempty"`
}

type HelmInGitRepo struct {
	Repository  string  `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	Revision    string  `protobuf:"bytes,2,opt,name=revision,proto3" json:"revision,omitempty"`
	ChartPath   *File   `protobuf:"bytes,3,opt,name=chartPath,proto3" json:"chartPath,omitempty"`
	ValuesPaths []*File `protobuf:"bytes,4,rep,name=valuesPaths,proto3" json:"valuesPaths,omitempty"`
}

type HelmInHelmRepo struct {
	Repository   string  `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	ChartName    string  `protobuf:"bytes,2,opt,name=chartName,proto3" json:"chartName,omitempty"`
	ChartVersion string  `protobuf:"bytes,3,opt,name=chartVersion,proto3" json:"chartVersion,omitempty"`
	ValuesPaths  []*File `protobuf:"bytes,4,rep,name=valuesPaths,proto3" json:"valuesPaths,omitempty"`
}

type ManagedAlertManager struct {
	Configmap     *File `protobuf:"bytes,1,opt,name=configmap,proto3" json:"configmap,omitempty"`
	Secret        *File `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
	Configuration *File `protobuf:"bytes,3,opt,name=configuration,proto3" json:"configuration,omitempty"`
	Statefulset   *File `protobuf:"bytes,4,opt,name=statefulset,proto3" json:"statefulset,omitempty"`
}
*/

func toArrayString(in []interface{}) []string {
	out := make([]string, len(in))
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toArrayStringSorted(in []interface{}) []string {
	if in == nil {
		return nil
	}
	out := toArrayString(in)
	sort.Strings(out)
	return out
}

func toArrayInterface(in []string) []interface{} {
	out := make([]interface{}, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func toArrayInterfaceSorted(in []string) []interface{} {
	if in == nil {
		return nil
	}
	sort.Strings(in)
	out := toArrayInterface(in)
	return out
}

func toMapString(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for i, v := range in {
		if v == nil {
			out[i] = ""
			continue
		}
		out[i] = v.(string)
	}
	return out
}

func toMapByte(in map[string]interface{}) map[string][]byte {
	out := make(map[string][]byte)
	for i, v := range in {
		if v == nil {
			out[i] = []byte{}
			continue
		}
		value := v.(string)
		out[i] = []byte(value)
	}
	return out
}

func toMapInterface(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for i, v := range in {
		out[i] = v
	}
	return out
}

// Expanders

func expandMetaData(p []interface{}) *commonpb.Metadata {
	obj := &commonpb.Metadata{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})
	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}
	if v, ok := in["description"].(string); ok && len(v) > 0 {
		obj.Description = v
	}
	if v, ok := in["project"].(string); ok && len(v) > 0 {
		obj.Project = v
	}
	if v, ok := in["projectID"].(string); ok && len(v) > 0 {
		obj.ProjectID = v
	}
	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj.Id = v
	}

	if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}
	return obj
}

func expandPlacementLabels(p []interface{}) []*commonpb.PlacementLabel {
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := make([]*commonpb.PlacementLabel, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		label := commonpb.PlacementLabel{}

		if v, ok := in["key"].(string); ok {
			label.Key = v
		}
		if v, ok := in["value"].(string); ok {
			label.Value = v
		}
		obj[i] = &label
	}

	return obj
}

func expandPlacement(p []interface{}) *commonpb.PlacementSpec {
	obj := &commonpb.PlacementSpec{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return obj
	}

	in := p[0].(map[string]interface{})
	if v, ok := in["selector"].(string); ok && len(v) > 0 {
		obj.Selector = v
	}

	if v, ok := in["labels"].([]interface{}); ok {
		obj.Labels = expandPlacementLabels(v)
	}

	return obj
}

func expandFile(p []interface{}) *File {
	obj := File{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	in := p[0].(map[string]interface{})
	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	return &obj
}

func expandFiles(p []interface{}) []*File {
	if len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := make([]*File, len(p))
	for i := range p {
		of := File{}
		in := p[i].(map[string]interface{})
		if v, ok := in["name"].(string); ok && len(v) > 0 {
			of.Name = v
		}
		obj[i] = &of
	}
	return obj
}

func expandQuantity(p []interface{}) *resource.Quantity {
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})
	if v, ok := in["string"].(string); ok {
		log.Println("string v", v)
		ob, err := resource.ParseQuantity(v)
		if err == nil {
			log.Println("string v error", err, " ob ", ob)
			return &ob
		}
		log.Println("string v error", err)
	}

	return nil
}

func expandResourceQuantity(p []interface{}) *commonpb.ResourceQuantity {
	obj := commonpb.ResourceQuantity{}
	if len(p) == 0 || p[0] == nil {
		return &obj
	}

	in := p[0].(map[string]interface{})
	if v, ok := in["memory"].([]interface{}); ok {
		obj.Memory = expandQuantity(v)
	}

	if v, ok := in["cpu"].([]interface{}); ok {
		obj.Cpu = expandQuantity(v)
	}

	log.Println("expandResourceQuantity obj", obj)
	return &obj
}

func expandProjectMeta(p []interface{}) []*commonpb.ProjectMeta {
	if len(p) == 0 {
		return []*commonpb.ProjectMeta{}
	}
	out := make([]*commonpb.ProjectMeta, len(p))
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := commonpb.ProjectMeta{}

		if v, ok := in["name"].(string); ok {
			obj.Name = v
		}
		if v, ok := in["id"].(string); ok {
			obj.Id = v
		}

		out[i] = &obj
	}

	log.Println("expandProjectMeta out", out)
	return out
}

func expandSharingSpec(p []interface{}) *commonpb.SharingSpec {
	obj := commonpb.SharingSpec{}
	if len(p) == 0 || p[0] == nil {
		return &obj
	}

	in := p[0].(map[string]interface{})
	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["projects"].([]interface{}); ok {
		obj.Projects = expandProjectMeta(v)
	}

	log.Println("expandSharingSpec obj", obj)
	return &obj
}

// Flatteners

func flattenMetaData(in *commonpb.Metadata) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.Name) > 0 {
		obj["name"] = in.Name
	}

	if len(in.Description) > 0 {
		obj["description"] = in.Description
	}

	if len(in.Project) > 0 {
		obj["project"] = in.Project
	}

	if len(in.ProjectID) > 0 {
		obj["projectID"] = in.ProjectID
	}

	if len(in.Id) > 0 {
		obj["id"] = in.Id
	}

	if len(in.Labels) > 0 {
		obj["labels"] = toMapInterface(in.Labels)
	}

	if len(in.Annotations) > 0 {
		obj["annotations"] = toMapInterface(in.Annotations)
	}

	return []interface{}{obj}
}

func flattenPlacement(in *commonpb.PlacementSpec) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})
	if len(in.Labels) > 0 {
		obj["labels"] = in.Labels
	}

	if len(in.Selector) > 0 {
		obj["selector"] = in.Selector
	}

	return []interface{}{obj}
}

func flattenFile(in *File) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})
	if len(in.Name) > 0 {
		obj["name"] = in.Name
	}
	return []interface{}{obj}
}

func flattenFiles(input []*File) []interface{} {
	if input == nil {
		return nil
	}

	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := map[string]interface{}{}
		if len(in.Name) > 0 {
			obj["name"] = in.Name
		}
		out[i] = obj
	}

	return out
}

func flattenResourceQuantity(in *commonpb.ResourceQuantity) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})
	if in.Memory != nil {
		obj1 := make([]interface{}, 1)
		obj2 := make(map[string]interface{})
		obj2["string"] = in.GetMemory().String()
		obj1[0] = obj2
		obj["memory"] = obj1
	}

	if in.Cpu != nil {
		obj1 := make([]interface{}, 1)
		obj2 := make(map[string]interface{})
		obj2["string"] = in.GetCpu().String()
		obj1[0] = obj2
		obj["cpu"] = obj1
	}

	log.Println("flattenResourceQuantity obj", obj)
	return []interface{}{obj}
}

func flattenResourceQuantities(in *commonpb.ResourceQuantity) []interface{} {
	if in == nil {
		return nil
	}
	objRoot := make([]interface{}, 1)

	obj := make(map[string]interface{})
	if in.Memory != nil {
		obj1 := make([]interface{}, 1)
		obj2 := make(map[string]interface{})
		obj2["string"] = in.GetMemory()
		obj1[0] = obj2
		obj["memory"] = obj1
	}

	if in.Cpu != nil {
		obj1 := make([]interface{}, 1)
		obj2 := make(map[string]interface{})
		obj2["string"] = in.GetCpu()
		obj1[0] = obj2
		obj["cpu"] = obj1
	}

	objRoot[0] = obj
	log.Println("flattenResourceQuantity obj", obj)
	return []interface{}{objRoot}
}

func flattenRatio(in *commonpb.ResourceRatio) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})
	obj["memory"] = in.Memory
	obj["cpu"] = in.Cpu

	return []interface{}{obj}
}

func flattenProjectMeta(input []*commonpb.ProjectMeta) []interface{} {
	if input == nil {
		return nil
	}

	out := make([]interface{}, len(input))
	for i, in := range input {
		obj := map[string]interface{}{}
		if len(in.Name) > 0 {
			obj["name"] = in.Name
		}
		if len(in.Id) > 0 {
			obj["id"] = in.Id
		}
		out[i] = obj
	}

	return out
}

func flattenSharingSpec(in *commonpb.SharingSpec) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})
	obj["enabled"] = in.Enabled
	if len(in.Projects) > 0 {
		obj["projects"] = flattenProjectMeta(in.Projects)
	}

	return []interface{}{obj}
}