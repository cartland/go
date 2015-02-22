/*
 * Copyright 2015 Chris Cartland
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"encoding/json"
	"net/http"
	"time"

	"github.com/mjibson/appstats"
)

type listing struct {
	Name       string `json:"name"`       // Resource name
	Location   string `json:"location"`   // URL
	Expiration int64  `json:"expiration"` // UNIX time in seconds
}

func init() {
	// curl -H "Content-Type: application/json" -d '{"name":"directory","location":"http://localhost:8888"}' http://localhost:8888 -f -X PUT
	http.Handle("/", appstats.NewHandler(register))
	// curl -H "Content-Type: application/json" http://localhost:8888/heartbeat -f -X POST
	http.Handle("/heartbeat", appstats.NewHandler(heartbeat))
}

func register(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		getDirectory(c, w, r)
		return
	case "PUT":
		putListing(c, w, r)
	default:
		c.Errorf("Method %v not expected.", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
	}
}

func getDirectory(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	listings, err := getAllListings(c)
	if err != nil {
		c.Errorf("listings GetAll query failed %v", err)
	}
	b, err := json.Marshal(listings)
	if err != nil {
		c.Errorf("getDirectory Marshal failed %v", err)
	}
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		c.Errorf("getDirectory %v", err)
	}
}

func getAllListings(c appengine.Context) ([]listing, error) {
	q := datastore.NewQuery("listing")
	var listings []listing
	_, err := q.GetAll(c, &listings)
	return listings, err
}

func putListing(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	newListing := listing{}
	requestData := make([]byte, 100)
	n, err := r.Body.Read(requestData)
	if err != nil {
		c.Errorf("addListing read body %v", err)
	}
	err = json.Unmarshal(requestData[:n], &newListing)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.Errorf("addListing Unmarshal %v", err)
		return
	}

	key, err := updateListing(c, newListing)
	if err != nil {
		c.Errorf("putListing updateListing err %v", err)
	}
	c.Infof("addListing key %v", key)

	// Retrieve the recently stored listing.
	err = datastore.Get(c, key, &newListing)
	if err != nil {
		c.Errorf("Could not get recently updated listing %v", err)
	}

	b, err := json.Marshal(newListing)
	if err != nil {
		c.Errorf("addListing Marshal %v", err)
	}
	_, err = w.Write(b)
	if err != nil {
		c.Errorf("addListing Write %v", err)
	}
}

func updateListing(c appengine.Context, newListing listing) (*datastore.Key, error) {
	stringId := "name" + newListing.Name + "location" + newListing.Location
	key := datastore.NewKey(c, "listing", stringId, 0, nil)
	key, err := datastore.Put(c, key, &newListing)
	return key, err
}

func extendListingExpiration(c appengine.Context, newListing listing) (*datastore.Key, error) {
	now := time.Now()
	newListing.Expiration = now.Add(time.Minute).Unix()

	stringId := "name" + newListing.Name + "location" + newListing.Location
	key := datastore.NewKey(c, "listing", stringId, 0, nil)
	key, err := datastore.Put(c, key, &newListing)
	return key, err
}

func heartbeat(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		checkOldListings(c, w, r)
		return
	default:
		c.Errorf("Method %v not expected.", r.Method)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
	}
}

func checkOldListings(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	q := datastore.NewQuery("listing")
	var listings []listing
	keys, err := q.GetAll(c, &listings)
	if err != nil {
		c.Errorf("getDirectory GetAll query failed %v", err)
	}

	for index, element := range listings {
		// Only check heartbeat if the listing has expired.
		if time.Now().Unix() > element.Expiration {
			if checkHeartbeat(c, element) {
				_, _ = extendListingExpiration(c, element)
			} else {
				key := keys[index]
				datastore.Delete(c, key)
			}
		}
	}
	getDirectory(c, w, r)
}

func checkHeartbeat(c appengine.Context, element listing) bool {
	client := urlfetch.Client(c)
	resp, err := client.Get(element.Location)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}
