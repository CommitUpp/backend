-- APIテスト用データ
-- password: password123

INSERT INTO auth.users (
  instance_id,
  id,
  aud,
  role,
  email,
  encrypted_password,
  email_confirmed_at,
  confirmation_token,
  recovery_token,
  email_change_token_new,
  email_change,
  raw_app_meta_data,
  raw_user_meta_data,
  phone_change,
  phone_change_token,
  email_change_token_current,
  reauthentication_token,
  created_at,
  updated_at
)
VALUES
(
  '00000000-0000-0000-0000-000000000000',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'authenticated',
  'authenticated',
  'api-test-user-1@example.com',
  extensions.crypt('password123', extensions.gen_salt('bf')),
  '2026-07-03 02:17:25.74907+00',
  '',
  '',
  '',
  '',
  '{"provider": "email", "providers": ["email"]}'::jsonb,
  '{"sub": "762c6359-60a5-4596-ab54-199d282c1ef7", "email": "api-test-user-1@example.com", "email_verified": true, "phone_verified": false, "full_name": "API Test User 1"}'::jsonb,
  '',
  '',
  '',
  '',
  '2026-07-03 02:17:25.74907+00',
  '2026-07-03 02:17:25.74907+00'
),
(
  '00000000-0000-0000-0000-000000000000',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'authenticated',
  'authenticated',
  'api-test-user-2@example.com',
  extensions.crypt('password123', extensions.gen_salt('bf')),
  '2026-07-03 06:45:02+00',
  '',
  '',
  '',
  '',
  '{"provider": "email", "providers": ["email"]}'::jsonb,
  '{"sub": "b9c81a9b-33f2-485b-8fa3-45b2daf16fe3", "email": "api-test-user-2@example.com", "email_verified": true, "phone_verified": false, "full_name": "API Test User 2"}'::jsonb,
  '',
  '',
  '',
  '',
  '2026-07-03 06:45:02+00',
  '2026-07-03 06:45:02+00'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO auth.identities (
  provider_id,
  user_id,
  identity_data,
  provider,
  last_sign_in_at,
  created_at,
  updated_at
)
VALUES
(
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  '{"sub": "762c6359-60a5-4596-ab54-199d282c1ef7", "email": "api-test-user-1@example.com", "email_verified": true, "phone_verified": false, "full_name": "API Test User 1"}'::jsonb,
  'email',
  '2026-07-03 02:17:25.74907+00',
  '2026-07-03 02:17:25.74907+00',
  '2026-07-03 02:17:25.74907+00'
),
(
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  '{"sub": "b9c81a9b-33f2-485b-8fa3-45b2daf16fe3", "email": "api-test-user-2@example.com", "email_verified": true, "phone_verified": false, "full_name": "API Test User 2"}'::jsonb,
  'email',
  '2026-07-03 06:45:02+00',
  '2026-07-03 06:45:02+00',
  '2026-07-03 06:45:02+00'
)
ON CONFLICT (provider_id, provider) DO NOTHING;

INSERT INTO public.users (
  id,
  name,
  avatar_url,
  birth_date,
  group_limit,
  created_at
)
VALUES
(
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'API Test User 1',
  '/avatars/api-test-user-1.png',
  '1995-01-15',
  3,
  '2026-07-03 02:17:25.74907+00'
),
(
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'API Test User 2',
  '/avatars/api-test-user-2.png',
  '1996-02-20',
  4,
  '2026-07-03 06:45:02+00'
)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  avatar_url = EXCLUDED.avatar_url,
  birth_date = EXCLUDED.birth_date,
  group_limit = EXCLUDED.group_limit;

INSERT INTO public.movies (
  id,
  tmdb_id,
  title,
  poster_url,
  trailer_url,
  overview,
  release_date,
  created_at,
  genres
)
VALUES
(
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  'api-test-db31549b',
  'API Test Movie 1',
  '/api-test-movie-1-poster.jpg',
  '',
  'APIテスト用映画データ',
  null,
  '2026-07-03 02:17:25.74907+00',
  ARRAY['Test']
),
(
  '459198ae-5401-425f-8a8f-ead8ca264f0a',
  'api-test-459198ae',
  'API Test Movie 2',
  '/api-test-movie-2-poster.jpg',
  '',
  'APIテスト用映画データ',
  null,
  '2026-07-03 06:45:02+00',
  ARRAY['Test']
),
(
  '9f759b8e-90ec-4dc7-b791-7d5a5ff28f74',
  'api-test-9f759b8e',
  'API Test Movie 3',
  '/api-test-movie-3-poster.jpg',
  '',
  'APIテスト用映画データ',
  null,
  '2026-07-03 03:00:11.974636+00',
  ARRAY['Test']
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO public.favorite_movies (
  id,
  created_at,
  user_id,
  movie_id,
  order_num
)
VALUES
(
  '447455f1-df6c-40a9-b08d-4afb216e3ac2',
  '2026-07-03 02:17:25.74907+00',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  null
),
(
  '5039577f-9620-4db6-a869-9aae7e657c31',
  '2026-07-03 06:45:02+00',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  '459198ae-5401-425f-8a8f-ead8ca264f0a',
  null
),
(
  'de2452cb-1f32-4d15-b791-7d31b161f434',
  '2026-07-03 03:00:11.974636+00',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  '9f759b8e-90ec-4dc7-b791-7d5a5ff28f74',
  null
)
ON CONFLICT (user_id, movie_id) DO UPDATE
SET
  created_at = EXCLUDED.created_at,
  order_num = EXCLUDED.order_num;

INSERT INTO public.streaming_services (
  id,
  name,
  logo_url,
  created_at
)
VALUES
(
  '0b90430e-5987-47ae-b76d-85d15930f3a1',
  'Netflix',
  '/streaming/netflix.png',
  '2026-07-03 00:00:00+00'
),
(
  '2d79b9d0-b7e5-4073-bbd4-33eaa95d65e2',
  'Prime Video',
  '/streaming/prime-video.png',
  '2026-07-03 00:00:00+00'
),
(
  'a19f37b6-9728-43e4-a8cd-311f8dfc0d34',
  'Disney+',
  '/streaming/disney-plus.png',
  '2026-07-03 00:00:00+00'
)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  logo_url = EXCLUDED.logo_url;

INSERT INTO public.movie_streamings (
  movie_id,
  service_id,
  crated_at
)
VALUES
(
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  '0b90430e-5987-47ae-b76d-85d15930f3a1',
  '2026-07-03 02:20:00+00'
),
(
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  '2d79b9d0-b7e5-4073-bbd4-33eaa95d65e2',
  '2026-07-03 02:20:00+00'
),
(
  '459198ae-5401-425f-8a8f-ead8ca264f0a',
  'a19f37b6-9728-43e4-a8cd-311f8dfc0d34',
  '2026-07-03 06:50:00+00'
)
ON CONFLICT (movie_id, service_id) DO NOTHING;

INSERT INTO public.groups (
  id,
  name,
  monthly_goal,
  created_at
)
VALUES
(
  '18e15360-e4dd-4b18-8f0a-b84b63bd22f1',
  'API Test Group',
  10,
  '2026-07-03 07:00:00+00'
),
(
  '875d3193-f33f-4a20-b08a-ac3453404ef8',
  'API Test Archived Group',
  3,
  '2026-07-03 07:10:00+00'
)
ON CONFLICT (id) DO UPDATE
SET
  name = EXCLUDED.name,
  monthly_goal = EXCLUDED.monthly_goal;

INSERT INTO public.group_members (
  group_id,
  user_id,
  role,
  joined_at,
  is_active
)
VALUES
(
  '18e15360-e4dd-4b18-8f0a-b84b63bd22f1',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'admin',
  '2026-07-03 07:00:00+00',
  true
),
(
  '18e15360-e4dd-4b18-8f0a-b84b63bd22f1',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'user',
  '2026-07-03 07:05:00+00',
  true
),
(
  '875d3193-f33f-4a20-b08a-ac3453404ef8',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'admin',
  '2026-07-03 07:10:00+00',
  false
)
ON CONFLICT (group_id, user_id) DO UPDATE
SET
  role = EXCLUDED.role,
  joined_at = EXCLUDED.joined_at,
  is_active = EXCLUDED.is_active;

INSERT INTO public.watch_statuses (
  user_id,
  movie_id,
  status,
  created_at,
  updated_at
)
VALUES
(
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  'watching',
  '2026-07-03 08:00:00+00',
  '2026-07-03 08:00:00+00'
),
(
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  '459198ae-5401-425f-8a8f-ead8ca264f0a',
  'watched',
  '2026-07-03 08:10:00+00',
  '2026-07-03 08:10:00+00'
),
(
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  'watched',
  '2026-07-03 08:20:00+00',
  '2026-07-03 08:20:00+00'
),
(
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  '9f759b8e-90ec-4dc7-b791-7d5a5ff28f74',
  'wanna_watch',
  '2026-07-03 08:30:00+00',
  '2026-07-03 08:30:00+00'
)
ON CONFLICT (user_id, movie_id) DO UPDATE
SET
  status = EXCLUDED.status,
  created_at = EXCLUDED.created_at,
  updated_at = EXCLUDED.updated_at;

INSERT INTO public.group_goals (
  id,
  created_at,
  group_id,
  year_month,
  goal_count,
  achieved
)
VALUES
(
  '977fe41e-dd7f-4117-95c7-437d574d5f47',
  '2026-07-03 09:00:00+00',
  '18e15360-e4dd-4b18-8f0a-b84b63bd22f1',
  '2026-07',
  10,
  false
)
ON CONFLICT (id) DO UPDATE
SET
  group_id = EXCLUDED.group_id,
  year_month = EXCLUDED.year_month,
  goal_count = EXCLUDED.goal_count,
  achieved = EXCLUDED.achieved;

INSERT INTO public.group_limit_logs (
  id,
  created_at,
  user_id,
  reason
)
VALUES
(
  '1a5e4a65-999f-41f9-8b7d-68ff1b2127f6',
  '2026-07-03 09:10:00+00',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'survey'
)
ON CONFLICT (id) DO UPDATE
SET
  user_id = EXCLUDED.user_id,
  reason = EXCLUDED.reason;

INSERT INTO public.chat_rooms (
  id,
  group_id,
  movie_id,
  movie_title,
  created_by,
  created_at
)
VALUES
(
  '69904c3d-75f9-4dc3-a22c-f0fc4740a7c5',
  '18e15360-e4dd-4b18-8f0a-b84b63bd22f1',
  'db31549b-d1f9-4f9c-94c5-ce4cf3d0fa09',
  'API Test Movie 1',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  '2026-07-03 09:20:00+00'
)
ON CONFLICT (id) DO UPDATE
SET
  group_id = EXCLUDED.group_id,
  movie_id = EXCLUDED.movie_id,
  movie_title = EXCLUDED.movie_title,
  created_by = EXCLUDED.created_by;

INSERT INTO public.chat_messages (
  id,
  room_id,
  user_id,
  content,
  created_at
)
VALUES
(
  '5ca7cc1b-cebc-4288-a96a-d7551340dcaf',
  '69904c3d-75f9-4dc3-a22c-f0fc4740a7c5',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  'API test message from user 1',
  '2026-07-03 09:21:00+00'
),
(
  '0520fb65-ee43-478c-9a9c-b03570ef0b7b',
  '69904c3d-75f9-4dc3-a22c-f0fc4740a7c5',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  'API test message from user 2',
  '2026-07-03 09:22:00+00'
)
ON CONFLICT (id) DO UPDATE
SET
  room_id = EXCLUDED.room_id,
  user_id = EXCLUDED.user_id,
  content = EXCLUDED.content;

INSERT INTO public.recommendations (
  id,
  movie_id,
  from_user_id,
  created_at
)
VALUES
(
  'c44f2c91-451b-4b47-b852-fb996eb6a8ec',
  '9f759b8e-90ec-4dc7-b791-7d5a5ff28f74',
  '762c6359-60a5-4596-ab54-199d282c1ef7',
  '2026-07-03 09:30:00+00'
)
ON CONFLICT (id) DO UPDATE
SET
  movie_id = EXCLUDED.movie_id,
  from_user_id = EXCLUDED.from_user_id;

INSERT INTO public.recommendation_groups (
  recommendation_id,
  group_id,
  created_at
)
VALUES
(
  'c44f2c91-451b-4b47-b852-fb996eb6a8ec',
  '18e15360-e4dd-4b18-8f0a-b84b63bd22f1',
  '2026-07-03 09:31:00+00'
)
ON CONFLICT (recommendation_id, group_id) DO NOTHING;

INSERT INTO public.surveys (
  id,
  created_at,
  title,
  is_active
)
VALUES
(
  'd4c0247c-27c3-47d8-8efb-6c902abc69dc',
  '2026-07-03 09:40:00+00',
  'API Test Survey',
  true
)
ON CONFLICT (id) DO UPDATE
SET
  title = EXCLUDED.title,
  is_active = EXCLUDED.is_active;

INSERT INTO public.survey_questions (
  id,
  created_at,
  survey_id,
  question,
  order_num
)
VALUES
(
  '84dfed7c-2973-414f-9718-498a56745486',
  '2026-07-03 09:41:00+00',
  'd4c0247c-27c3-47d8-8efb-6c902abc69dc',
  'Which genre do you want to watch more?',
  1
)
ON CONFLICT (id) DO UPDATE
SET
  survey_id = EXCLUDED.survey_id,
  question = EXCLUDED.question,
  order_num = EXCLUDED.order_num;

INSERT INTO public.survey_answers (
  id,
  answered_at,
  survey_id,
  user_id,
  question_id,
  answer
)
VALUES
(
  '59c29b2a-df4e-41b5-b727-7d614c924876',
  '2026-07-03 09:42:00+00',
  'd4c0247c-27c3-47d8-8efb-6c902abc69dc',
  'b9c81a9b-33f2-485b-8fa3-45b2daf16fe3',
  '84dfed7c-2973-414f-9718-498a56745486',
  'Science Fiction'
)
ON CONFLICT (id) DO UPDATE
SET
  survey_id = EXCLUDED.survey_id,
  user_id = EXCLUDED.user_id,
  question_id = EXCLUDED.question_id,
  answer = EXCLUDED.answer;
