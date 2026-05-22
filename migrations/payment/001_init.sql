create table if not exists payments (
  id uuid primary key,
  rental_id uuid not null,
  user_id uuid not null,
  amount numeric(10,2) not null,
  status text not null check (status in ('success','failed','refunded')),
  method text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
