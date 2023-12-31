/*
 Licensed to the Apache Software Foundation (ASF) under one
 or more contributor license agreements.  See the NOTICE file
 distributed with this work for additional information
 regarding copyright ownership.  The ASF licenses this file
 to you under the Apache License, Version 2.0 (the
 "License"); you may not use this file except in compliance
 with the License.  You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package events

import (
	"sync"

	"github.com/apache/yunikorn-core/pkg/metrics"
	"github.com/apache/yunikorn-scheduler-interface/lib/go/si"
)

var maxEventStoreSize = 1000

// The EventStore operates under the following assumptions:
//  - for a given ObjectID only one (the latest) event is stored
//  - there is a cap for the number of events stored
//  - the CollectEvents() function clears the currently stored events in the EventStore
// Assuming the rate of events generated by the scheduler component in a given time period
// is high, calling CollectEvents() periodically should be fine.
// Objects are stored using their ObjectID, so different types of objects
// with the same id overwrite each other.
type EventStore interface {
	Store(event *si.EventRecord)
	CollectEvents() []*si.EventRecord
	CountStoredEvents() int
}

type defaultEventStore struct {
	eventMap     map[string]*si.EventRecord
	storedEvents int

	sync.RWMutex
}

func newEventStoreImpl() EventStore {
	return &defaultEventStore{
		eventMap: make(map[string]*si.EventRecord),
	}
}

func (es *defaultEventStore) Store(event *si.EventRecord) {
	es.Lock()
	defer es.Unlock()

	// limiting the size of the store
	if es.storedEvents >= maxEventStoreSize {
		metrics.GetEventMetrics().IncEventsNotStored()
		return
	}

	if _, ok := es.eventMap[event.ObjectID]; !ok {
		es.storedEvents++
	}

	es.eventMap[event.ObjectID] = event
	metrics.GetEventMetrics().IncEventsStored()
}

func (es *defaultEventStore) CollectEvents() []*si.EventRecord {
	es.Lock()
	defer es.Unlock()

	messages := make([]*si.EventRecord, 0, len(es.eventMap))
	for _, v := range es.eventMap {
		messages = append(messages, v)
	}

	es.eventMap = make(map[string]*si.EventRecord)
	es.storedEvents = 0

	metrics.GetEventMetrics().AddEventsCollected(len(messages))
	return messages
}

func (es *defaultEventStore) CountStoredEvents() int {
	es.RLock()
	defer es.RUnlock()

	return len(es.eventMap)
}
