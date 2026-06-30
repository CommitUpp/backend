alter table public.groups enable row level security;
alter table public.group_members enable row level security;

drop policy if exists "authenticated users can create groups" on public.groups;
create policy "authenticated users can create groups"
on public.groups
for insert
to public
with check (auth.uid() is not null);

drop policy if exists "authenticated users can update groups" on public.groups;
create policy "authenticated users can update groups"
on public.groups
for update
to public
using (
  auth.uid() is not null
  and
  exists (
    select 1
    from public.group_members
    where group_members.group_id = groups.id
      and group_members.user_id = auth.uid()
  )
)
with check (
  auth.uid() is not null
  and
  exists (
    select 1
    from public.group_members
    where group_members.group_id = groups.id
      and group_members.user_id = auth.uid()
  )
);

drop policy if exists "users can read own group memberships" on public.group_members;
create policy "users can read own group memberships"
on public.group_members
for select
to public
using (auth.uid() = user_id);

drop policy if exists "users can create own group memberships" on public.group_members;
create policy "users can create own group memberships"
on public.group_members
for insert
to public
with check (auth.uid() = user_id);

drop function if exists public.create_group(text, integer);
drop function if exists public.join_group(uuid);

create or replace function public.create_group(name text, monthly_goal integer)
returns public.groups
language plpgsql
security definer
set search_path = public
as $$
declare
  created_group public.groups%rowtype;
begin
  if auth.uid() is null then
    raise exception 'authenticated user is required';
  end if;

  insert into public.groups (name, monthly_goal)
  values (create_group.name, create_group.monthly_goal)
  returning * into created_group;

  insert into public.group_members (group_id, user_id, role)
  values (created_group.id, auth.uid(), 'admin');

  return created_group;
end;
$$;

create or replace function public.join_group(group_id uuid)
returns public.groups
language plpgsql
security definer
set search_path = public
as $$
declare
  joined_group public.groups%rowtype;
begin
  if auth.uid() is null then
    raise exception 'authenticated user is required';
  end if;

  insert into public.group_members (group_id, user_id, role)
  values (join_group.group_id, auth.uid(), 'user');

  update public.groups
  set monthly_goal = monthly_goal + 3
  where groups.id = join_group.group_id
  returning * into joined_group;

  if joined_group.id is null then
    raise exception 'group not found';
  end if;

  return joined_group;
end;
$$;
