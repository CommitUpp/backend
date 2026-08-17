ALTER TABLE public.chat_rooms
ALTER COLUMN movie_id DROP NOT NULL;

CREATE UNIQUE INDEX chat_rooms_group_id_general_unique
ON public.chat_rooms (group_id)
WHERE movie_id IS NULL;

CREATE OR REPLACE FUNCTION public.create_group(name text, monthly_goal integer)
RETURNS public.groups
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path TO 'public'
AS $$
DECLARE
  created_group public.groups%rowtype;
BEGIN
  IF auth.uid() IS NULL THEN
    RAISE EXCEPTION 'authenticated user is required';
  END IF;

  INSERT INTO public.groups (name, monthly_goal)
  VALUES (create_group.name, create_group.monthly_goal)
  RETURNING * INTO created_group;

  INSERT INTO public.group_members (group_id, user_id, role)
  VALUES (created_group.id, auth.uid(), 'admin');

  INSERT INTO public.chat_rooms (group_id, movie_id, movie_title, created_by)
  VALUES (created_group.id, NULL, '雑談', auth.uid());

  RETURN created_group;
END;
$$;
