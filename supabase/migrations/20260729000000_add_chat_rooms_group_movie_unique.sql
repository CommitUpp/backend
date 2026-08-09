ALTER TABLE public.chat_rooms
ADD CONSTRAINT chat_rooms_group_id_movie_id_unique
UNIQUE (group_id, movie_id);
