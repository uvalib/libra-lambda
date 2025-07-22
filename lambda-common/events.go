package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

func NewEventBus(eventBus string, eventSource string) (uvalibrabus.UvaBus, error) {
	// we will accept bad config and return nil quietly
	if len(eventBus) == 0 {
		fmt.Printf("INFO: Event bus is not configured, no telemetry emitted\n")
		return nil, uvalibrabus.ErrConfig
	}

	cfg := uvalibrabus.UvaBusConfig{BusName: eventBus, Source: eventSource, Log: nil}
	return uvalibrabus.NewUvaBus(cfg)
}

func pubAuditEvent(bus uvalibrabus.UvaBus, obj uvaeasystore.EasyStoreObject, who string, fname string, before string, after string) error {
	if bus == nil {
		return uvalibrabus.ErrConfig
	}
	detail, _ := auditPayload(who, fname, before, after)
	ev := uvalibrabus.UvaBusEvent{
		EventName:  uvalibrabus.EventFieldUpdate,
		Namespace:  obj.Namespace(),
		Identifier: obj.Id(),
		Detail:     detail,
	}
	return bus.PublishEvent(&ev)
}

func auditPayload(who string, fname string, before string, after string) (json.RawMessage, error) {
	pl := uvalibrabus.UvaAuditEvent{Who: who, FieldName: fname, Before: before, After: after}
	return pl.Serialize()
}

//
// end of file
//
