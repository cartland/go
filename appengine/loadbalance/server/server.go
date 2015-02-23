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
	"appengine/urlfetch"
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/mjibson/appstats"
)

type listing struct {
	Name       string `json:"name"`       // Resource name
	Location   string `json:"location"`   // URL
	Expiration int64  `json:"expiration"` // UNIX time in seconds
}

var serviceListings []listing

func init() {
	serviceListings = []listing{
		listing{
			Name:     "service",
			Location: "https://loadbalance-golang.appspot.com/reliable",
		},
		listing{
			Name:     "service",
			Location: "https://loadbalance-golang.appspot.com/flaky",
		},
	}

	http.Handle("/reliable", appstats.NewHandler(reliable))
	http.Handle("/flaky", appstats.NewHandler(flaky))
	http.Handle("/register", appstats.NewHandler(register))
}

func reliable(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(nil)
	if err == nil {
		_, e := w.Write(b)
		err = e
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		c.Errorf("%v", err)
	}
}

func flaky(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UTC().UnixNano())
	if rand.Intn(2) == 0 {
		reliable(c, w, r)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		c.Errorf("Flaky failed")
	}
}

func register(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	client := urlfetch.Client(c)
	for _, l := range serviceListings {
		registerListing(c, w, r, client, l)
	}
	reliable(c, w, r)
}

func registerListing(c appengine.Context, w http.ResponseWriter, r *http.Request, client *http.Client, l listing) {
	b, err := json.Marshal(l)
	body := bytes.NewReader(b)
	req, err := http.NewRequest("PUT", "https://directory-golang.appspot.com", body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		c.Errorf("%v", err)
	}
	_, err = client.Do(req)
}

type appError struct {
	err string
}

func (e appError) Error() string {
	return e.err
}
