create table embeddings
(
    id          serial
        primary key,
    surah       integer not null,
    ayah        integer not null,
    data        jsonb   not null,
    embedding   vector(3072),
    token_count integer,
    created_at  timestamp with time zone default now(),
    updated_at  timestamp with time zone
);

alter table embeddings
    owner to postgres;

create index ix_embeddings_id
    on embeddings (id);

create index ix_embeddings_surah
    on embeddings (surah);

create index ix_embeddings_ayah
    on embeddings (ayah);

