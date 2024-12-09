-- name: GetAnticipationByIdClient :many
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
    a.id_client = $3;


-- name: GetAllAnticipationsSolicit :many
select *
from public.anticipation_solicit
where status = true and paid = false;