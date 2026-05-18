package ecs

import (
	"fmt"
	"reflect"
)

const maxResourceTypes = 256

// Resources 管理 World 上的全局资源；通过 [World.Resources] 访问。
// 推荐使用 [Resource]、[AddResource] 与 [GetResource]。
type Resources struct {
	registry  resourceRegistry
	resources []any
}

type resourceRegistry struct {
	types    []reflect.Type
	ids      []uint8
	typeToID map[reflect.Type]uint8
}

func newResources() Resources {
	return Resources{
		registry:  newResourceRegistry(),
		resources: make([]any, maxResourceTypes),
	}
}

func newResourceRegistry() resourceRegistry {
	return resourceRegistry{
		types:    make([]reflect.Type, maxResourceTypes),
		typeToID: make(map[reflect.Type]uint8),
		ids:      make([]uint8, 0, 16),
	}
}

func (r *resourceRegistry) resourceID(tp reflect.Type) (uint8, bool) {
	if id, ok := r.typeToID[tp]; ok {
		return id, false
	}
	if len(r.typeToID) >= maxResourceTypes {
		panic(fmt.Sprintf("ecs: exceeded the maximum of %d resource types", maxResourceTypes))
	}
	id := uint8(len(r.typeToID))
	r.typeToID[tp] = id
	r.types[id] = tp
	r.ids = append(r.ids, id)
	return id, true
}

func (r *resourceRegistry) resourceType(id uint8) (reflect.Type, bool) {
	if int(id) >= len(r.types) {
		return nil, false
	}
	tp := r.types[id]
	if tp == nil {
		return nil, false
	}
	return tp, true
}

// Add 向 World 添加资源；res 应为指针。
//
// 若该类型资源已存在则 panic。
//
// 参见 [Resource.Add]。
func (r *Resources) Add(id ResID, res any) {
	if r.resources[id.id] != nil {
		panic(fmt.Sprintf("ecs: resource of ID %d was already added (type %v)", id.id, reflect.TypeOf(res)))
	}
	r.resources[id.id] = res
}

// Remove 从 World 移除资源。
//
// 若该类型资源不存在则 panic。
func (r *Resources) Remove(id ResID) {
	if r.resources[id.id] == nil {
		panic(fmt.Sprintf("ecs: resource of ID %d is not present", id.id))
	}
	r.resources[id.id] = nil
}

// Get 返回资源指针；不存在时返回 nil。
func (r *Resources) Get(id ResID) any {
	return r.resources[id.id]
}

// Has 是否已添加该类型资源。
func (r *Resources) Has(id ResID) bool {
	return r.resources[id.id] != nil
}

func (r *Resources) reset() {
	for i := range r.resources {
		r.resources[i] = nil
	}
}
