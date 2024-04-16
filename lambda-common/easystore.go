//
//
//

package main

import (
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
)

func newEasystore(cfg *Config) (uvaeasystore.EasyStore, error) {

	// make better later
	config := uvaeasystore.DatastorePostgresConfig{
		DbHost:     cfg.EsDbHost,
		DbPort:     cfg.EsDbPort,
		DbName:     cfg.EsDbName,
		DbUser:     cfg.EsDbUser,
		DbPassword: cfg.EsDbPassword,
		DbTimeout:  30, // probably fix me later

		BusName:    cfg.BusName,
		SourceName: cfg.SourceName,

		//Log:        logger,
	}

	return uvaeasystore.NewEasyStore(config)
}

func createEasystoreObject(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject) error {

	obj, err := es.Create(obj)
	if err == nil {
		fmt.Printf("INFO: created new easystore object [%s/%s]\n", obj.Namespace(), obj.Id())
	}
	return err
}

func getEasystoreObjectByKey(es uvaeasystore.EasyStore, namespace string, identifier string, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObject, error) {
	return es.GetByKey(namespace, identifier, what)
}

func getEasystoreObjectsByFields(es uvaeasystore.EasyStore, namespace string, fields uvaeasystore.EasyStoreObjectFields, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObjectSet, error) {
	return es.GetByFields(namespace, fields, what)
}

func putEasystoreObject(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject, what uvaeasystore.EasyStoreComponents) error {
	obj, err := es.Update(obj, what)
	if err == nil {
		fmt.Printf("INFO: updated easystore object [%s/%s]\n", obj.Namespace(), obj.Id())
	}
	return err
}

//
// end of file
//
