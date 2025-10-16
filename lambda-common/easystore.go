//
//
//

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/uvalib/easystore/uvaeasystore"
)

var maxEsRetries = 3
var esRetrySleepTime = 100 * time.Millisecond

func newEasystoreProxy(cfg *Config) (uvaeasystore.EasyStore, error) {

	config := uvaeasystore.ProxyConfigImpl{
		ServiceEndpoint: cfg.EsProxyUrl,
		Log:             log.Default(),
	}
	return uvaeasystore.NewEasyStoreProxy(config)
}

func newEasystoreReadonlyProxy(cfg *Config) (uvaeasystore.EasyStoreReadonly, error) {

	config := uvaeasystore.ProxyConfigImpl{
		ServiceEndpoint: cfg.EsProxyUrl,
		Log:             log.Default(),
	}
	return uvaeasystore.NewEasyStoreProxyReadonly(config)
}

func createEasystoreObject(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject) error {

	obj, err := es.ObjectCreate(obj)
	if err == nil {
		fmt.Printf("INFO: created new easystore object [%s/%s]\n", obj.Namespace(), obj.Id())
	}
	return err
}

func getEasystoreObjectByKey(es uvaeasystore.EasyStoreReadonly, namespace string, identifier string, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObject, error) {
	obj, err := es.ObjectGetByKey(namespace, identifier, what)
	if err == nil {
		fmt.Printf("INFO: got easystore object [%s/%s]\n", obj.Namespace(), obj.Id())
	}
	return obj, err
}

func getEasystoreObjectsByFields(es uvaeasystore.EasyStoreReadonly, namespace string, fields uvaeasystore.EasyStoreObjectFields, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObjectSet, error) {
	objSet, err := es.ObjectGetByFields(namespace, fields, what)
	if err == nil {
		fmt.Printf("INFO: got %d easystore objects\n", objSet.Count())
	}
	return objSet, err
}

func putEasystoreObject(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject, what uvaeasystore.EasyStoreComponents) error {
	obj, err := es.ObjectUpdate(obj, what)
	if err == nil {
		fmt.Printf("INFO: updated easystore object [%s/%s]\n", obj.Namespace(), obj.Id())
	}
	return err
}

func putEasystoreFieldWithRetry(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject, what uvaeasystore.EasyStoreComponents, field string, value string) (uvaeasystore.EasyStoreObject, error) {
	err := putEasystoreObject(es, obj, uvaeasystore.Fields)
	// happy day, return...
	if err == nil {
		return obj, err
	}

	// our retry loop
	for retry := 0; retry < maxEsRetries; retry++ {
		// if our object is stale
		if err == uvaeasystore.ErrStaleObject {

			fmt.Printf("WARNING: easystore object is stale [%s/%s], retry #%d\n", obj.Namespace(), obj.Id(), retry+1)

			// sleep for a bit before retrying
			time.Sleep(esRetrySleepTime)

			// try and get it again
			var newObj uvaeasystore.EasyStoreObject
			newObj, err = getEasystoreObjectByKey(es, obj.Namespace(), obj.Id(), what)
			// it's all over, return error
			if err != nil {
				return obj, err
			}
			obj = newObj
			fields := obj.Fields()
			fields[field] = value
			obj.SetFields(fields)
			err = putEasystoreObject(es, obj, uvaeasystore.Fields)
			// happy day, return...
			if err == nil {
				return obj, err
			}

			// otherwise, an error... if its stale, retry, otherwise abandon loop and return the error
		} else {
			// it's all over, return error
			fmt.Printf("WARNING: abandoning retry for [%s/%s] (%s)\n", obj.Namespace(), obj.Id(), err.Error())
			return obj, err
		}
	}

	// we have retried and are giving up
	fmt.Printf("WARNING: easystore object was stale [%s/%s], gave up\n", obj.Namespace(), obj.Id())
	return obj, err
}

//
// end of file
//
