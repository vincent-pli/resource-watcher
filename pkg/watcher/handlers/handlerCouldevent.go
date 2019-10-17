package handlers

import (
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/cache"
)

type CouldeventHandler struct {
	Resource []string
	Client   cloudevents.Client
	Sink     string
	Log      logr.Logger
}

var _ cache.ResourceEventHandler = (*CouldeventHandler)(nil)

func (c *CouldeventHandler) OnAdd(obj interface{}) {
	c.Log.Info("resource added")

}

func (c *CouldeventHandler) OnUpdate(oldObj, newObj interface{}) {
	c.Log.Info("resource update")

}

func (c *CouldeventHandler) OnDelete(obj interface{}) {
	c.Log.Info("resource delete")
}
