create table if not exists rentals (
  id uuid primary key,
  user_id uuid not null,
  car_id uuid not null,
  start_date date not null,
  end_date date not null,
  insurance boolean not null default false,
  total_price numeric(10,2) not null,
  status text not null check (status in ('pending','confirmed','active','completed','cancelled')),
  late_fee numeric(10,2) not null default 0,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
