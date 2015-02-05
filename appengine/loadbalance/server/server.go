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
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/mjibson/appstats"
)

func init() {
	http.Handle("/reliable", appstats.NewHandler(reliable))
	http.Handle("/flaky", appstats.NewHandler(flaky))
}

func reliable(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json")
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

type appError struct {
	err string
}

func (e appError) Error() string {
	return e.err
}
