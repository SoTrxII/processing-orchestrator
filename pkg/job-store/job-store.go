package job_store

import (
	"context"
	"encoding/json"
	"processing-orchestrator/internal/utils"
	"sync"
)

type ExternalJobStore struct {
	client    utils.StateSaver
	cache     map[string]*JobState
	component string
	storeKey  string
	mu        sync.Mutex
}

func NewJobStore(client utils.StateSaver, component string) *ExternalJobStore {
	return &ExternalJobStore{client: client, component: component, storeKey: "jobs", cache: map[string]*JobState{}}
}

// Upsert a job to the store
func (m *ExternalJobStore) Upsert(job *JobState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[job.Id] = job

	// Save the entire map
	return m.updateBacking()
}

// Get a job by id
func (m *ExternalJobStore) Get(id string) (*JobState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if job, ok := m.cache[id]; ok {
		return job, nil
	}
	return nil, nil
}

// GetAll jobs, optionally reloading from the backing store
func (m *ExternalJobStore) GetAll(reload bool) ([]*JobState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if reload {
		err := m.fetchBacking()
		if err != nil {
			return nil, err
		}
	}
	jobs := []*JobState{}
	for _, job := range m.cache {
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Delete a job by id
func (m *ExternalJobStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cache, id)

	// Save the entire map
	return m.updateBacking()
}

// Sync the backing store with the in-memory cache
// the cache is the source of truth
func (m *ExternalJobStore) updateBacking() error {
	bytes, err := json.Marshal(m.cache)
	if err != nil {
		return err
	}
	return m.client.SaveState(context.Background(), m.component, m.storeKey, bytes, map[string]string{})
}

// Fetch the backing store into the in-memory cache
func (m *ExternalJobStore) fetchBacking() error {
	item, err := m.client.GetState(context.Background(), m.component, m.storeKey, map[string]string{})
	if err != nil {
		return err
	}
	if item.Value == nil {
		return nil
	}
	return json.Unmarshal(item.Value, &m.cache)
}
