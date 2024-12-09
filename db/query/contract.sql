-- name: CreateContract :one
INSERT INTO public.contracts
(id, provider_name, provider_email, document_url, status, created_at, is_signed, date_signed, envelop_id, "access_id", "tenant_id", contract_type, id_control_batch)
VALUES(nextval('contracts_id_seq'::regclass), $1, $2, $3,'pending'::character varying, now(), false, NULL, $4, $5, $6, $7, $8)
    RETURNING *;


-- name: AssignContract :one
UPDATE public.contracts
SET status='assign', is_signed=true, date_signed=now()
WHERE envelop_id=$1 AND
      is_signed=false AND
      status!='assign'
    RETURNING *;

-- -- name: AssignContractAssignment :one
-- UPDATE public.contracts
-- SET status='assign', is_signed=true, date_signed=now()
-- WHERE envelop_id=$1 AND
--     is_signed=false AND
--     status!='assign'AND
--       contract_type = 'term_of_assignment' AND
--       id_control_batch IS NULL
--     RETURNING *;

-- name: GetContractAll :one
SELECT *
FROM public.contracts
WHERE access_id=$1 AND tenant_id=$2;