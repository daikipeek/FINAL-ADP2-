create table if not exists cars (
  id uuid primary key,
  brand text not null,
  model text not null,
  type text not null,
  location text not null,
  status text not null default 'available',
  daily_rate numeric(10,2) not null,
  image_url text not null default '',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists favorites (
  user_id uuid not null,
  car_id uuid not null references cars(id) on delete cascade,
  created_at timestamptz not null default now(),
  primary key(user_id, car_id)
);

create table if not exists reviews (
  id uuid primary key,
  user_id uuid not null,
  car_id uuid not null references cars(id) on delete cascade,
  rating int not null check (rating between 1 and 5),
  comment text not null default '',
  created_at timestamptz not null default now()
);

insert into cars(id,brand,model,type,location,status,daily_rate,image_url) values
('11111111-1111-1111-1111-111111111111','Toyota','Camry','sedan','Almaty','available',55,'https://images.unsplash.com/photo-1621007947382-bb3c3994e3fb?auto=format&fit=crop&w=1200&q=80'),
('22222222-2222-2222-2222-222222222222','BMW','X5','suv','Astana','available',120,'https://images.unsplash.com/photo-1555215695-3004980ad54e?auto=format&fit=crop&w=1200&q=80'),
('33333333-3333-3333-3333-333333333333','Hyundai','Elantra','sedan','Shymkent','available',48,'https://images.unsplash.com/photo-1619767886558-efdc259cde1a?auto=format&fit=crop&w=1200&q=80')
on conflict do nothing;
