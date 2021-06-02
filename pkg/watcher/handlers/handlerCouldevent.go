package handlers

import (
	"context"
	"fmt"
	"log"

	// cloudevents "github.com/cloudevents/sdk-go"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
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
	// Pull metav1.Object out of the object
	if o, err := meta.Accessor(obj); err == nil {
		c.Log.Info("resource added...", "object", o)
	} else {
		c.Log.Error(err, "OnAdd missing Meta",
			"object", obj, "type", fmt.Sprintf("%T", obj))
		return
	}

	// Pull the runtime.Object out of the object
	if o, ok := obj.(runtime.Object); ok {
		c.Log.Info("resource added again...", "object", o)
	} else {
		c.Log.Error(nil, "OnAdd missing runtime.Object",
			"object", obj, "type", fmt.Sprintf("%T", obj))
		return
	}
	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"hello": "world"})

	ctx := cloudevents.ContextWithTarget(context.Background(), c.Sink)
	// Send that Event.
	if result := c.Client.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Fatalf("failed to send, %v", result)
	}

	c.Log.Info("resource added %v", obj)
}

func (c *CouldeventHandler) OnUpdate(oldObj, newObj interface{}) {
	c.Log.Info("resource update")

}

func (c *CouldeventHandler) OnDelete(obj interface{}) {
	c.Log.Info("resource delete")
}
