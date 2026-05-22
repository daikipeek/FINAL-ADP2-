create table if not exists users (
  id uuid primary key,
  name text not null,
  email text unique not null,
  password_hash text not null,
  role text not null check (role in ('customer','admin')),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
