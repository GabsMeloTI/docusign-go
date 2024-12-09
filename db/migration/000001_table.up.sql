CREATE TABLE contracts (
                           id BIGSERIAL PRIMARY KEY,
                           provider_name VARCHAR NOT NULL,
                           provider_email varchar NOT NULL,
                           document_url TEXT NOT NULL,
                           status VARCHAR(50) NOT NULL DEFAULT 'pending',
                           created_at TIMESTAMP DEFAULT NOW(),
                           is_signed BOOL,
                           date_signed TIMESTAMP,
                           envelop_id VARCHAR DEFAULT '' NOT NULL,
                           "access_id" BIGINT,
                           "tenant_id" UUID
);
alter table contracts add column contract_type varchar NOT NULL;
alter table contracts add column id_control_batch uuid;


CREATE TABLE public.anticipation (
                                     id bigserial NOT NULL,
                                     number_doc varchar NULL,
                                     value_original float8 NULL,
                                     date_issue varchar NULL,
                                     date_due varchar NULL,
                                     parcel int8 NULL,
                                     type_doc varchar NULL,
                                     motive_doc varchar NULL,
                                     id_client int8 NULL,
                                     id_control uuid NULL,
                                     requested bool NULL,
                                     status bool DEFAULT true NULL,
                                     create_at timestamp DEFAULT now() NOT NULL,
                                     update_at timestamp NULL,
                                     CONSTRAINT anticipation_pkey PRIMARY KEY (id),
                                     CONSTRAINT unique_number_doc UNIQUE (number_doc)
);
alter table anticipation add column access_id BIGINT;
alter table anticipation add column tenant_id UUID;

CREATE TABLE public.anticipation_solicit (
                                             id bigserial NOT NULL,
                                             id_anticipation int8 NULL,
                                             status bool DEFAULT true NULL,
                                             create_at timestamp DEFAULT now() NOT NULL,
                                             batch uuid NULL,
                                             id_solicit uuid NULL,
                                             paid bool DEFAULT false NULL,
                                             date_paid timestamp NULL,
                                             status_anticipation int8 NULL,
                                             CONSTRAINT anticipation_solicit_pkey PRIMARY KEY (id)
);