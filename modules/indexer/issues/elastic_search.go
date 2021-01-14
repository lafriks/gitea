// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"code.gitea.io/gitea/modules/log"

	"github.com/olivere/elastic/v7"
)

var (
	_ Indexer = &ElasticSearchIndexer{}
)

// docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" -v /home/lafriks/tmp/elasticdata0:/usr/share/elasticsearch/data docker.elastic.co/elasticsearch/elasticsearch:7.10.1

// ElasticSearchIndexer implements Indexer interface
type ElasticSearchIndexer struct {
	client               *elastic.Client
	indexerName          string
	available            bool
	availabilityCallback func(bool)
	stopTimer            chan struct{}
}

type elasticLogger struct {
	log.LevelLogger
}

func (l elasticLogger) Printf(format string, args ...interface{}) {
	_ = l.Log(2, l.GetLevel(), format, args...)
}

// NewElasticSearchIndexer creates a new elasticsearch indexer
func NewElasticSearchIndexer(url, indexerName string) (*ElasticSearchIndexer, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetHealthcheckInterval(10 * time.Second),
		elastic.SetGzip(false),
	}

	logger := elasticLogger{log.GetLogger(log.DEFAULT)}

	if logger.GetLevel() == log.TRACE || logger.GetLevel() == log.DEBUG {
		opts = append(opts, elastic.SetTraceLog(logger))
	} else if logger.GetLevel() == log.ERROR || logger.GetLevel() == log.CRITICAL || logger.GetLevel() == log.FATAL {
		opts = append(opts, elastic.SetErrorLog(logger))
	} else if logger.GetLevel() == log.INFO || logger.GetLevel() == log.WARN {
		opts = append(opts, elastic.SetInfoLog(logger))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	indexer := &ElasticSearchIndexer{
		client:      client,
		indexerName: indexerName,
		available:   true,
		stopTimer:   make(chan struct{}),
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				indexer.checkAvailability()
			case <-indexer.stopTimer:
				ticker.Stop()
				return
			}
		}
	}()

	return indexer, nil
}

const (
	defaultMapping = `{
		"mappings": {
			"properties": {
				"id": {
					"type": "integer",
					"index": true
				},
				"repo_id": {
					"type": "integer",
					"index": true
				},
				"title": {
					"type": "text",
					"index": true
				},
				"content": {
					"type": "text",
					"index": true
				},
				"comments": {
					"type" : "text",
					"index": true
				}
			}
		}
	}`
)

// Init will initialize the indexer
func (b *ElasticSearchIndexer) Init() (bool, error) {
	ctx := context.Background()
	exists, err := b.client.IndexExists(b.indexerName).Do(ctx)
	if err != nil {
		return false, b.checkError(err)
	}

	if !exists {
		var mapping = defaultMapping

		createIndex, err := b.client.CreateIndex(b.indexerName).BodyString(mapping).Do(ctx)
		if err != nil {
			return false, b.checkError(err)
		}
		if !createIndex.Acknowledged {
			return false, errors.New("init failed")
		}

		return false, nil
	}
	return true, nil
}

// SetAvailabilityChangeCallback sets callback that will be triggered when availability changes
func (b *ElasticSearchIndexer) SetAvailabilityChangeCallback(callback func(bool)) {
	b.availabilityCallback = callback
}

// Ping checks if elastic is available
func (b *ElasticSearchIndexer) Ping() bool {
	return b.available
}

// Index will save the index data
func (b *ElasticSearchIndexer) Index(issues []*IndexerData) error {
	if len(issues) == 0 {
		return nil
	} else if len(issues) == 1 {
		issue := issues[0]
		_, err := b.client.Index().
			Index(b.indexerName).
			Id(fmt.Sprintf("%d", issue.ID)).
			BodyJson(map[string]interface{}{
				"id":       issue.ID,
				"repo_id":  issue.RepoID,
				"title":    issue.Title,
				"content":  issue.Content,
				"comments": issue.Comments,
			}).
			Do(context.Background())
		return b.checkError(err)
	}

	reqs := make([]elastic.BulkableRequest, 0)
	for _, issue := range issues {
		reqs = append(reqs,
			elastic.NewBulkIndexRequest().
				Index(b.indexerName).
				Id(fmt.Sprintf("%d", issue.ID)).
				Doc(map[string]interface{}{
					"id":       issue.ID,
					"repo_id":  issue.RepoID,
					"title":    issue.Title,
					"content":  issue.Content,
					"comments": issue.Comments,
				}),
		)
	}

	_, err := b.client.Bulk().
		Index(b.indexerName).
		Add(reqs...).
		Do(context.Background())
	return b.checkError(err)
}

// Delete deletes indexes by ids
func (b *ElasticSearchIndexer) Delete(ids ...int64) error {
	if len(ids) == 0 {
		return nil
	} else if len(ids) == 1 {
		_, err := b.client.Delete().
			Index(b.indexerName).
			Id(fmt.Sprintf("%d", ids[0])).
			Do(context.Background())
		return b.checkError(err)
	}

	reqs := make([]elastic.BulkableRequest, 0)
	for _, id := range ids {
		reqs = append(reqs,
			elastic.NewBulkDeleteRequest().
				Index(b.indexerName).
				Id(fmt.Sprintf("%d", id)),
		)
	}

	_, err := b.client.Bulk().
		Index(b.indexerName).
		Add(reqs...).
		Do(context.Background())
	return b.checkError(err)
}

// Search searches for issues by given conditions.
// Returns the matching issue IDs
func (b *ElasticSearchIndexer) Search(keyword string, repoIDs []int64, limit, start int) (*SearchResult, error) {
	kwQuery := elastic.NewMultiMatchQuery(keyword, "title", "content", "comments")
	query := elastic.NewBoolQuery()
	query = query.Must(kwQuery)
	if len(repoIDs) > 0 {
		var repoStrs = make([]interface{}, 0, len(repoIDs))
		for _, repoID := range repoIDs {
			repoStrs = append(repoStrs, repoID)
		}
		repoQuery := elastic.NewTermsQuery("repo_id", repoStrs...)
		query = query.Must(repoQuery)
	}
	searchResult, err := b.client.Search().
		Index(b.indexerName).
		Query(query).
		Sort("id", true).
		From(start).Size(limit).
		Do(context.Background())
	if err != nil {
		return nil, b.checkError(err)
	}

	hits := make([]Match, 0, limit)
	for _, hit := range searchResult.Hits.Hits {
		id, _ := strconv.ParseInt(hit.Id, 10, 64)
		hits = append(hits, Match{
			ID: id,
		})
	}

	return &SearchResult{
		Total: searchResult.TotalHits(),
		Hits:  hits,
	}, nil
}

// Close implements indexer
func (b *ElasticSearchIndexer) Close() {
	close(b.stopTimer)
}

func (b *ElasticSearchIndexer) checkError(err error) error {
	if errors.Is(err, elastic.ErrNoClient) && b.available {
		b.available = false
		if b.availabilityCallback != nil {
			b.availabilityCallback(b.available)
		}
	}
	return err
}

func (b *ElasticSearchIndexer) checkAvailability() {
	if b.available {
		return
	}

	// Request cluster state to check if elastic is available again
	_, err := b.client.ClusterState().Do(context.Background())
	if err != nil {
		return
	}
	b.available = true
	if b.availabilityCallback != nil {
		b.availabilityCallback(b.available)
	}
}
