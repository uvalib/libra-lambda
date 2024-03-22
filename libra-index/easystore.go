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

func getEasystoreObject(es uvaeasystore.EasyStore, namespace string, identifier string) (uvaeasystore.EasyStoreObject, error) {

	// we just want the fields and the metadata
	return es.GetByKey(namespace, identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
}

//
// end of file
//
