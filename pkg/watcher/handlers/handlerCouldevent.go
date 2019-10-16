package handlers

import (
	"k8s.io/client-go/tools/cache"
)

type couldeventHandler struct {
	resource []string
}

var _ cache.ResourceEventHandler = (*couldeventHandler)(nil)

func (c *couldeventHandler) OnAdd(obj interface{}) {

}

func (c *couldeventHandler) OnUpdate(oldObj, newObj interface{}) {

}

func (c *couldeventHandler) OnDelete(obj interface{}) {

}
