// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: antecipations.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getAllAnticipationsSolicit = `-- name: GetAllAnticipationsSolicit :many
select id, id_anticipation, status, create_at, batch, id_solicit, paid, date_paid, status_anticipation
from public.anticipation_solicit
where status = true and paid = false
`

func (q *Queries) GetAllAnticipationsSolicit(ctx context.Context) ([]AnticipationSolicit, error) {
	rows, err := q.db.QueryContext(ctx, getAllAnticipationsSolicit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AnticipationSolicit
	for rows.Next() {
		var i AnticipationSolicit
		if err := rows.Scan(
			&i.ID,
			&i.IDAnticipation,
			&i.Status,
			&i.CreateAt,
			&i.Batch,
			&i.IDSolicit,
			&i.Paid,
			&i.DatePaid,
			&i.StatusAnticipation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAnticipationByIdClient = `-- name: GetAnticipationByIdClient :many
SELECT
    a.id,
    a.number_doc,
    a.value_original,
    a.date_issue,
    a.date_due,
    a.parcel,
    a.type_doc,
    a.motive_doc,
    a.requested,
    s.id as solicit_id,
    s.paid
FROM
    public.anticipation a
        left join anticipation_solicit s on s.id_anticipation = a.id
        and s.status = true
where
    a.status=true and
    a.access_id = $1 and
    a.tenant_id = $2 and
    a.id_client = $3
`

type GetAnticipationByIdClientParams struct {
	AccessID sql.NullInt64 `json:"access_id"`
	TenantID uuid.NullUUID `json:"tenant_id"`
	IDClient sql.NullInt64 `json:"id_client"`
}

type GetAnticipationByIdClientRow struct {
	ID            int64           `json:"id"`
	NumberDoc     sql.NullString  `json:"number_doc"`
	ValueOriginal sql.NullFloat64 `json:"value_original"`
	DateIssue     sql.NullString  `json:"date_issue"`
	DateDue       sql.NullString  `json:"date_due"`
	Parcel        sql.NullInt64   `json:"parcel"`
	TypeDoc       sql.NullString  `json:"type_doc"`
	MotiveDoc     sql.NullString  `json:"motive_doc"`
	Requested     sql.NullBool    `json:"requested"`
	SolicitID     sql.NullInt64   `json:"solicit_id"`
	Paid          sql.NullBool    `json:"paid"`
}

func (q *Queries) GetAnticipationByIdClient(ctx context.Context, arg GetAnticipationByIdClientParams) ([]GetAnticipationByIdClientRow, error) {
	rows, err := q.db.QueryContext(ctx, getAnticipationByIdClient, arg.AccessID, arg.TenantID, arg.IDClient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAnticipationByIdClientRow
	for rows.Next() {
		var i GetAnticipationByIdClientRow
		if err := rows.Scan(
			&i.ID,
			&i.NumberDoc,
			&i.ValueOriginal,
			&i.DateIssue,
			&i.DateDue,
			&i.Parcel,
			&i.TypeDoc,
			&i.MotiveDoc,
			&i.Requested,
			&i.SolicitID,
			&i.Paid,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
