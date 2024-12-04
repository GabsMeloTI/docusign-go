-- name: CreateContract :one
INSERT INTO public.contracts
(id, provider_name, document_url, status, created_at, is_signed, date_signed)
VALUES(nextval('contracts_id_seq'::regclass), $1, $2, 'pending'::character varying, now(), false, NULL)
    RETURNING *;