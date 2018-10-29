package jiraquery

import (
	"fmt"
	"strings"
	"time"
)

type Query interface {
	String() string
}

type Condition interface {
	String() string
}

func And(left, right Condition) Condition {
	return &binaryOp{
		left,
		"AND",
		right,
	}
}

func Or(left, right Condition) Condition {
	return &binaryOp{
		left,
		"OR",
		right,
	}
}

func Not(cond Condition) Condition {
	return &unaryOp{"NOT", cond}
}

type binaryOp struct {
	field Condition
	op    string
	value Condition
}

func (b *binaryOp) String() string {
	return fmt.Sprintf("%s %s %s", b.field.String(), b.op, b.value.String())
}

type unaryOp struct {
	op   string
	cond Condition
}

func (u *unaryOp) String() string {
	return fmt.Sprintf("%s %s", u.op, u.cond.String())
}

func Eq(field, value Condition) Condition {
	return &binaryOp{
		field,
		"=",
		value,
	}
}

func NotEq(field, value Condition) Condition {
	return &binaryOp{
		field,
		"!=",
		value,
	}
}

func GreaterThan(field, value Condition) Condition {
	return &binaryOp{
		field,
		">",
		value,
	}
}

func LessThan(field, value Condition) Condition {
	return &binaryOp{
		field,
		"<",
		value,
	}
}

type Word string

func (w Word) String() string { return string(w) }

type listOp struct {
	formatString string
	op           string
	conditions   []Condition
}

func strListToConditionList(strs []string) []Condition {
	conds := make([]Condition, len(strs))
	for i, str := range strs {
		conds[i] = Word(str)
	}
	return conds
}

func conditionListToStrList(conds []Condition) []string {
	strs := make([]string, len(conds))
	for i, cond := range conds {
		strs[i] = cond.String()
	}
	return strs
}

func (l *listOp) String() string {
	return fmt.Sprintf(l.formatString, strings.Join(conditionListToStrList(l.conditions), l.op))
}

func List(vals ...string) Condition {
	return &listOp{
		"( %s )",
		", ",
		strListToConditionList(vals),
	}
}

func In(field Word, values Condition) Condition {
	return &binaryOp{
		field,
		"IN",
		values,
	}
}

type wrapped struct {
	Condition
}

func (w *wrapped) String() string {
	return fmt.Sprintf("( %s )", w.Condition.String())
}

func Wrapped(cond Condition) Condition {
	return &wrapped{cond}
}

func MultiOr(conds ...Condition) Condition {
	return &listOp{
		"%s",
		" OR ",
		conds,
	}
}

func MultiAnd(conds ...Condition) Condition {
	return &listOp{
		"%s",
		" AND ",
		conds,
	}
}

func Before(field Word, when time.Time) Condition {
	return &binaryOp{
		field,
		"<",
		Word(when.Format("2006-1-2 15:04")),
	}
}

func After(field Word, when time.Time) Condition {
	return &binaryOp{
		field,
		">",
		Word(when.Format("2006-1-2 15:04")),
	}
}

// Vanity functions.

func Project(str string) Condition {
	return Eq(Word("project"), Word(fmt.Sprintf("%q", str)))
}

func IssueType(str string) Condition {
	return Eq(Word("issueType"), Word(fmt.Sprintf("%q", str)))
}

func CreatedBefore(t time.Time) Condition {
	return Before(Word("created"), t)
}

func CreatedAfter(t time.Time) Condition {
	return After(Word("created"), t)
}

// How should a builder for this look, it'd be nice to be able to have:
//
// Project("MX").IssueType("Bug").CreatedBefore(time.Now())
//
// We should have an and and an or builder.
//
// AndBuilder().Project("MX").etc...

type QueryBuilder interface {
	// Vanity functions

	Project(string) QueryBuilder
	IssueType(string) QueryBuilder
	CreatedAfter(time.Time) QueryBuilder
	CreatedBefore(time.Time) QueryBuilder

	// General functions
	Not(cond Condition) QueryBuilder
	Eq(field, value Condition) QueryBuilder
	NotEq(field, value Condition) QueryBuilder
	GreaterThan(field, value Condition) QueryBuilder
	LessThan(field, value Condition) QueryBuilder
	In(field Word, values Condition) QueryBuilder
	Wrapped(cond Condition) QueryBuilder

	Value() Query
}

type queryBuilder struct {
	connector  string
	conditions []Condition
}

func AndBuilder() QueryBuilder {
	return &queryBuilder{" AND ", nil}
}

func OrBuilder() QueryBuilder {
	return &queryBuilder{" OR ", nil}
}

func (q *queryBuilder) Project(str string) QueryBuilder {
	q.conditions = append(q.conditions, Project(str))
	return q
}

func (q *queryBuilder) IssueType(str string) QueryBuilder {
	q.conditions = append(q.conditions, IssueType(str))
	return q
}

func (q *queryBuilder) CreatedAfter(t time.Time) QueryBuilder {
	q.conditions = append(q.conditions, CreatedAfter(t))
	return q
}

func (q *queryBuilder) CreatedBefore(t time.Time) QueryBuilder {
	q.conditions = append(q.conditions, CreatedBefore(t))
	return q
}

func (q *queryBuilder) Not(cond Condition) QueryBuilder {
	q.conditions = append(q.conditions, Not(cond))
	return q
}

func (q *queryBuilder) Eq(field, value Condition) QueryBuilder {
	q.conditions = append(q.conditions, Eq(field, value))
	return q
}

func (q *queryBuilder) NotEq(field, value Condition) QueryBuilder {
	q.conditions = append(q.conditions, NotEq(field, value))
	return q
}

func (q *queryBuilder) GreaterThan(field, value Condition) QueryBuilder {
	q.conditions = append(q.conditions, GreaterThan(field, value))
	return q
}

func (q *queryBuilder) LessThan(field, value Condition) QueryBuilder {
	q.conditions = append(q.conditions, LessThan(field, value))
	return q
}

func (q *queryBuilder) In(field Word, values Condition) QueryBuilder {
	q.conditions = append(q.conditions, In(field, values))
	return q
}

func (q *queryBuilder) Wrapped(cond Condition) QueryBuilder {
	q.conditions = append(q.conditions, Wrapped(cond))
	return q
}

func (q *queryBuilder) Value() Query {
	return &listOp{
		"%s",
		q.connector,
		q.conditions,
	}
}
