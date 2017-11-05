// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	//"fmt"
	//"fmt"
	"time"

	"github.com/go-xorm/builder"
	"github.com/go-xorm/xorm"
)

// Reaction represents a reactions on issues and comments.
type Reaction struct {
	ID          int64     `xorm:"pk autoincr"`
	Type        string    `xorm:"INDEX UNIQUE(s) NOT NULL"`
	IssueID     int64     `xorm:"INDEX UNIQUE(s) NOT NULL"`
	CommentID   int64     `xorm:"INDEX UNIQUE(s)"`
	UserID      int64     `xorm:"INDEX UNIQUE(s) NOT NULL"`
	User        *User     `xorm:"-"`
	Created     time.Time `xorm:"-"`
	CreatedUnix int64     `xorm:"INDEX created"`
}

// AfterLoad is invoked from XORM after setting the values of all fields of this object.
func (s *Reaction) AfterLoad() {
	s.Created = time.Unix(s.CreatedUnix, 0).Local()
}

// FindReactionsOptions describes the conditions to Find reactions
type FindReactionsOptions struct {
	IssueID   int64
	CommentID int64
}

func (opts *FindReactionsOptions) toConds() builder.Cond {
	var cond = builder.NewCond()
	if opts.IssueID > 0 {
		cond = cond.And(builder.Eq{"reaction.issue_id": opts.IssueID})
	}
	if opts.CommentID > 0 {
		cond = cond.And(builder.Eq{"reaction.comment_id": opts.CommentID})
	}
	return cond
}

func findReactions(e Engine, opts FindReactionsOptions) ([]*Reaction, error) {
	reactions := make([]*Reaction, 0, 10)
	sess := e.Where(opts.toConds())
	return reactions, sess.
		Asc("reaction.issue_id", "reaction.comment_id", "reaction.created_unix", "reaction.id").
		Find(&reactions)
}

func createReaction(e *xorm.Session, opts *CreateReactionOptions) (*Reaction, error) {
	reaction := &Reaction{
		Type:    opts.Type,
		UserID:  opts.Doer.ID,
		IssueID: opts.Issue.ID,
	}
	if opts.Comment != nil {
		reaction.CommentID = opts.Comment.ID
	}
	if _, err := e.Insert(reaction); err != nil {
		return nil, err
	}

	return reaction, nil
}

// CreateReactionOptions defines options for creating reactions
type CreateReactionOptions struct {
	Type    string
	Doer    *User
	Issue   *Issue
	Comment *Comment
}

// CreateReaction creates reaction for issue or comment.
func CreateReaction(opts *CreateReactionOptions) (reaction *Reaction, err error) {
	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return nil, err
	}

	reaction, err = createReaction(sess, opts)
	if err != nil {
		return nil, err
	}

	if err = sess.Commit(); err != nil {
		return nil, err
	}
	return reaction, nil
}

// CreateIssueReaction creates a reaction on issue.
func CreateIssueReaction(doer *User, issue *Issue, content string) (*Reaction, error) {
	return CreateReaction(&CreateReactionOptions{
		Type:  content,
		Doer:  doer,
		Issue: issue,
	})
}

// CreateCommentReaction creates a reaction on comment.
func CreateCommentReaction(doer *User, issue *Issue, comment *Comment, content string) (*Reaction, error) {
	return CreateReaction(&CreateReactionOptions{
		Type:    content,
		Doer:    doer,
		Issue:   issue,
		Comment: comment,
	})
}

type ReactionList []*Reaction

func (list ReactionList) HasUser(userID int64) bool {
	for _, reaction := range list {
		if reaction.UserID == userID {
			return true
		}
	}
	return false
}
