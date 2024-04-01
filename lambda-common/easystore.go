//
//
//

package main

import (
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

	_, err := es.Create(obj)
	return err
}

func getEasystoreObjectByKey(es uvaeasystore.EasyStore, namespace string, identifier string, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObject, error) {
	return es.GetByKey(namespace, identifier, what)
}

func getEasystoreObjectsByFields(es uvaeasystore.EasyStore, namespace string, fields uvaeasystore.EasyStoreObjectFields, what uvaeasystore.EasyStoreComponents) (uvaeasystore.EasyStoreObjectSet, error) {
	return es.GetByFields(namespace, fields, what)
}

func putEasystoreObject(es uvaeasystore.EasyStore, obj uvaeasystore.EasyStoreObject, what uvaeasystore.EasyStoreComponents) error {
	_, err := es.Update(obj, what)
	return err
}

//
// end of file
//
