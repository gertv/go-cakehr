// Copyright 2019 - anova r&d bvba. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package cake

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Cake struct {
	key    string
	debug  bool
	client *http.Client
	base   string
}

type TimeOffPolicy struct {
	Name string `json:"name"`
}

type Employee struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type OutOfOfficeEntry struct {
	Id                int           `json:"id"`
	PolicyId          int           `json:"policy_id"`
	Policy            TimeOffPolicy `json:"policy"`
	EmployeeId        int           `json:"employee_id"`
	Employee          Employee      `json:"employee"`
	Details           string        `json:"string"`
	IsMultiDay        bool          `json:"is_multi_day"`
	IsSingleDay       bool          `json:"is_single_day"`
	IsPartOfDay       bool          `json:"is_part_of_day"`
	IsFirstPartOfDay  bool          `json:"first_part_of_day"`
	IsSecondPartOfDay bool          `json:"second_part_of_day"`
	Start             string        `json:"start_date"`
	End               string        `json:"end_date"`
	Hours             float32       `json:"hours"`
}

type OutOfOffice struct {
	Data []OutOfOfficeEntry `json:"data"`
}

// Start using the API here -
func CakeHR(subdomain, key string) Cake {
	return Cake{key, false, &http.Client{}, "https://" + subdomain + ".cake.hr/api"}
}

// Enable/disable debug mode. When debug mode is enabled,
// you will get additional logging showing the HTTP requests
// and responses
func (b *Cake) Debug(d bool) {
	b.debug = d
}

// Configure a custom HTTP client (e.g. to configure a proxy server)
func (b *Cake) Client(client *http.Client) {
	b.client = client
}

// Get a calendar that shows who's out
func (b Cake) OutOfOfficeOnDate(date string) (ooo OutOfOffice, err error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/leave-management/out-of-office-today?date=%s", b.base, date), nil)
	req.Header.Add("X-Auth-Token", b.key)
	log.Printf("%v", req)
	if err != nil {
		return
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, err := body(resp)
	if err != nil {
		return
	}
	if b.debug {
		log.Printf("Got response %s: %s", resp.Status, data)
	}

	err = json.Unmarshal(data, &ooo)

	return
}

// Extract body from the HTTP response
func body(resp *http.Response) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i OutOfOfficeEntry) StartTime() time.Time {
	time, err := time.Parse("2006-01-02", i.Start)
	if err != nil {
		panic(err)
	}
	return time
}

func (i OutOfOfficeEntry) EndTime() time.Time {
	time, err := time.Parse("2006-01-02", i.End)
	if err != nil {
		panic(err)
	}
	return time
}
