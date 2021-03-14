pragma foreign_keys = on;

create table if not exists main.plates (
  id integer primary key,
  code text not null unique,
  created_at text default current_timestamp,
  updated_at text default current_timestamp
);

create trigger if not exists plates_update_pre before update on plates
begin
  update plates
  set
    updated_at = current_timestamp,
    code = upper(code)
  where id = new.id;
end;

create table if not exists interpretations (
  id integer primary key,
  what text not null,
  why text not null,
  ip text not null,
  username text not null,
  plate_id integer not null,
  created_at text default current_timestamp,
  updated_at text default current_timestamp,

  unique (plate_id, what),

  foreign key (plate_id) references plates (id) on delete cascade
);

create trigger if not exists interpretations_update_pre before update on interpretations
begin
  update interpretations set updated_at = current_timestamp where id = new.id;
end;

create trigger if not exists interpretations_update_post after update on interpretations
begin
  update plates set updated_at = current_timestamp where id = new.plate_id;
end;
