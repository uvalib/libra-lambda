//
//
//

package main

import (
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"log"
)

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
	return es.ObjectGetByKey(namespace, identifier, what)
}

func getEasystoreObjectsByFields(es uvaeasystore.EasyStoreReadonly, namespace string, fields uvaeasystore.EasyStoreObjectFields, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObjectSet, error) {
	return es.ObjectGetByFields(namespace, fields, what)
}

func putEasystoreObject(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject, what uvaeasystore.EasyStoreComponents) error {
	obj, err := es.ObjectUpdate(obj, what)
	if err == nil {
		fmt.Printf("INFO: updated easystore object [%s/%s]\n", obj.Namespace(), obj.Id())
	}
	return err
}

//
// end of file
//
