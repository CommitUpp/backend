SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

COMMENT ON SCHEMA "public" IS 'standard public schema';

CREATE EXTENSION IF NOT EXISTS "pg_stat_statements" WITH SCHEMA "extensions";

CREATE EXTENSION IF NOT EXISTS "pgcrypto" WITH SCHEMA "extensions";

CREATE EXTENSION IF NOT EXISTS "supabase_vault" WITH SCHEMA "vault";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA "extensions";

CREATE TYPE "public"."gender" AS ENUM (
    'male',
    'female',
    'other'
);

ALTER TYPE "public"."gender" OWNER TO "postgres";

COMMENT ON TYPE "public"."gender" IS '性別';

CREATE TYPE "public"."reason" AS ENUM (
    'survey',
    'ad',
    'goal',
    'watchlater'
);

ALTER TYPE "public"."reason" OWNER TO "postgres";

COMMENT ON TYPE "public"."reason" IS 'グループ上限解放の理由';

CREATE TYPE "public"."role" AS ENUM (
    'admin',
    'editor',
    'user'
);

ALTER TYPE "public"."role" OWNER TO "postgres";

CREATE TYPE "public"."watch_status" AS ENUM (
    'watching',
    'watched',
    'wanna_watch'
);

ALTER TYPE "public"."watch_status" OWNER TO "postgres";

COMMENT ON TYPE "public"."watch_status" IS '視聴フラグ';

SET default_tablespace = '';

SET default_table_access_method = "heap";

CREATE TABLE IF NOT EXISTS "public"."groups" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "name" "text" NOT NULL,
    "monthly_goal" smallint DEFAULT '0'::smallint NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL
);

ALTER TABLE "public"."groups" OWNER TO "postgres";

COMMENT ON TABLE "public"."groups" IS 'グループ管理テーブル';

COMMENT ON COLUMN "public"."groups"."monthly_goal" IS 'グループの月間視聴本数目標';

CREATE OR REPLACE FUNCTION "public"."create_group"("name" "text", "monthly_goal" integer) RETURNS "public"."groups"
    LANGUAGE "plpgsql" SECURITY DEFINER
    SET "search_path" TO 'public'
    AS $$
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

ALTER FUNCTION "public"."create_group"("name" "text", "monthly_goal" integer) OWNER TO "postgres";

CREATE OR REPLACE FUNCTION "public"."handle_new_user"() RETURNS "trigger"
    LANGUAGE "plpgsql" SECURITY DEFINER
    SET "search_path" TO 'public'
    AS $$
BEGIN
  -- インサート先を小文字の public.users に変更
  INSERT INTO public.users (id, name, avatar_url)
  VALUES (
    new.id,
    coalesce(new.raw_user_meta_data ->> 'full_name', 'Google User'),
    new.raw_user_meta_data ->> 'avatar_url'
  );
  RETURN new;
END;
$$;

ALTER FUNCTION "public"."handle_new_user"() OWNER TO "postgres";

CREATE OR REPLACE FUNCTION "public"."join_group"("group_id" "uuid") RETURNS "public"."groups"
    LANGUAGE "plpgsql" SECURITY DEFINER
    SET "search_path" TO 'public'
    AS $$
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

ALTER FUNCTION "public"."join_group"("group_id" "uuid") OWNER TO "postgres";

CREATE OR REPLACE FUNCTION "public"."update_modified_column"() RETURNS "trigger"
    LANGUAGE "plpgsql"
    AS $$
BEGIN
    NEW.updated_at = now(); -- ここで updated_at に現在の時刻を入れます
    RETURN NEW;
END;
$$;

ALTER FUNCTION "public"."update_modified_column"() OWNER TO "postgres";

CREATE TABLE IF NOT EXISTS "public"."chat_messages" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "room_id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "content" "text" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    CONSTRAINT "CHAT_MESSAGES_content_check" CHECK (("char_length"("content") <= 1000))
);

ALTER TABLE "public"."chat_messages" OWNER TO "postgres";

COMMENT ON TABLE "public"."chat_messages" IS 'チャットルーム内の投稿';

CREATE TABLE IF NOT EXISTS "public"."chat_rooms" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "group_id" "uuid" NOT NULL,
    "movie_id" "uuid" NOT NULL,
    "movie_title" "text" NOT NULL,
    "created_by" "uuid" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL
);

ALTER TABLE "public"."chat_rooms" OWNER TO "postgres";

COMMENT ON TABLE "public"."chat_rooms" IS 'チャットテーブル';

CREATE TABLE IF NOT EXISTS "public"."favorite_movies" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "user_id" "uuid" DEFAULT "gen_random_uuid"(),
    "movie_id" "uuid" DEFAULT "gen_random_uuid"(),
    "order_num" smallint,
    CONSTRAINT "chk_order_num_range" CHECK ((("order_num" >= 1) AND ("order_num" <= 3)))
);

ALTER TABLE "public"."favorite_movies" OWNER TO "postgres";

COMMENT ON TABLE "public"."favorite_movies" IS 'お気に入り映画（上限3個 / プロフィールに表示）';

CREATE TABLE IF NOT EXISTS "public"."group_goals" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "group_id" "uuid" DEFAULT "gen_random_uuid"(),
    "year_month" "text",
    "goal_count" integer,
    "achieved" boolean
);

ALTER TABLE "public"."group_goals" OWNER TO "postgres";

COMMENT ON TABLE "public"."group_goals" IS '月間目標の達成記録';

CREATE TABLE IF NOT EXISTS "public"."group_limit_logs" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "user_id" "uuid" DEFAULT "gen_random_uuid"(),
    "reason" "public"."reason"
);

ALTER TABLE "public"."group_limit_logs" OWNER TO "postgres";

COMMENT ON TABLE "public"."group_limit_logs" IS 'グループ上限増加の履歴';

CREATE TABLE IF NOT EXISTS "public"."group_members" (
    "group_id" "uuid" NOT NULL,
    "user_id" "uuid" NOT NULL,
    "role" "public"."role",
    "joined_at" timestamp without time zone DEFAULT "now"() NOT NULL,
    "is_active" boolean DEFAULT true NOT NULL
);

ALTER TABLE "public"."group_members" OWNER TO "postgres";

COMMENT ON TABLE "public"."group_members" IS 'グループとメンバーの中間テーブル';

COMMENT ON COLUMN "public"."group_members"."is_active" IS 'グループに所属しているかどうか（抜けたらfalseに変わる）';

CREATE TABLE IF NOT EXISTS "public"."movie_streamings" (
    "movie_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "service_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "crated_at" timestamp with time zone DEFAULT "now"() NOT NULL
);

ALTER TABLE "public"."movie_streamings" OWNER TO "postgres";

COMMENT ON TABLE "public"."movie_streamings" IS '映画とサブスクの中間テーブル';

CREATE TABLE IF NOT EXISTS "public"."movies" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "tmdb_id" "text" NOT NULL,
    "title" "text" NOT NULL,
    "poster_url" "text",
    "trailer_url" "text" NOT NULL,
    "overview" "text",
    "release_date" "date",
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "genres" "text"[] NOT NULL
);

ALTER TABLE "public"."movies" OWNER TO "postgres";

COMMENT ON TABLE "public"."movies" IS '映画管理テーブル';

COMMENT ON COLUMN "public"."movies"."genres" IS '映画のジャンル配列';

CREATE TABLE IF NOT EXISTS "public"."recommendation_groups" (
    "recommendation_id" "uuid" NOT NULL,
    "group_id" "uuid" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL
);

ALTER TABLE "public"."recommendation_groups" OWNER TO "postgres";

COMMENT ON TABLE "public"."recommendation_groups" IS 'おすすめの送信先グループ（複数選択可）';

CREATE TABLE IF NOT EXISTS "public"."recommendations" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "movie_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "from_user_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL
);

ALTER TABLE "public"."recommendations" OWNER TO "postgres";

COMMENT ON TABLE "public"."recommendations" IS 'おすすめテーブル';

CREATE TABLE IF NOT EXISTS "public"."streaming_services" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "name" "text" NOT NULL,
    "logo_url" "text",
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL
);

ALTER TABLE "public"."streaming_services" OWNER TO "postgres";

COMMENT ON TABLE "public"."streaming_services" IS 'サブスクテーブル';

CREATE TABLE IF NOT EXISTS "public"."survey_answers" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "answered_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "survey_id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "user_id" "uuid" DEFAULT "gen_random_uuid"(),
    "question_id" "uuid" DEFAULT "gen_random_uuid"(),
    "answer" "text"
);

ALTER TABLE "public"."survey_answers" OWNER TO "postgres";

CREATE TABLE IF NOT EXISTS "public"."survey_questions" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "survey_id" "uuid" DEFAULT "gen_random_uuid"(),
    "question" "text",
    "order_num" integer
);

ALTER TABLE "public"."survey_questions" OWNER TO "postgres";

COMMENT ON TABLE "public"."survey_questions" IS 'アンケートの質問';

CREATE TABLE IF NOT EXISTS "public"."surveys" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "title" "text",
    "is_active" boolean
);

ALTER TABLE "public"."surveys" OWNER TO "postgres";

COMMENT ON TABLE "public"."surveys" IS 'アンケートのマスタ';

CREATE TABLE IF NOT EXISTS "public"."users" (
    "id" "uuid" DEFAULT "gen_random_uuid"() NOT NULL,
    "name" "text" NOT NULL,
    "avatar_url" "text",
    "birth_date" "date",
    "group_limit" smallint DEFAULT '3'::smallint NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "gender" "public"."gender" DEFAULT 'other'::"public"."gender" NOT NULL
);

ALTER TABLE "public"."users" OWNER TO "postgres";

COMMENT ON TABLE "public"."users" IS 'ユーザー管理テーブル';

COMMENT ON COLUMN "public"."users"."group_limit" IS '所属可能なグループ上限数 デフォルト3件';

CREATE TABLE IF NOT EXISTS "public"."watch_statuses" (
    "user_id" "uuid" NOT NULL,
    "movie_id" "uuid" NOT NULL,
    "status" "public"."watch_status" NOT NULL,
    "created_at" timestamp with time zone DEFAULT "now"() NOT NULL,
    "updated_at" timestamp with time zone DEFAULT "now"()
);

ALTER TABLE "public"."watch_statuses" OWNER TO "postgres";

COMMENT ON TABLE "public"."watch_statuses" IS '視聴状態の管理テーブル';

ALTER TABLE ONLY "public"."chat_messages"
    ADD CONSTRAINT "CHAT_MESSAGES_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."chat_rooms"
    ADD CONSTRAINT "CHAT_ROOMS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."favorite_movies"
    ADD CONSTRAINT "FAVORITE_MOVIES_order_num_key" UNIQUE ("order_num");

ALTER TABLE ONLY "public"."favorite_movies"
    ADD CONSTRAINT "FAVORITE_MOVIES_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."groups"
    ADD CONSTRAINT "GROUPS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."group_goals"
    ADD CONSTRAINT "GROUP_GOALS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."group_limit_logs"
    ADD CONSTRAINT "GROUP_LIMIT_LOGS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."group_members"
    ADD CONSTRAINT "GROUP_MEMBERS_pkey" PRIMARY KEY ("group_id", "user_id");

ALTER TABLE ONLY "public"."movies"
    ADD CONSTRAINT "MOVIES_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."movies"
    ADD CONSTRAINT "MOVIES_tmdb_id_key" UNIQUE ("tmdb_id");

ALTER TABLE ONLY "public"."movie_streamings"
    ADD CONSTRAINT "MOVIE_STREAMINGS_pkey" PRIMARY KEY ("movie_id", "service_id");

ALTER TABLE ONLY "public"."recommendations"
    ADD CONSTRAINT "RECOMMENDATIONS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."recommendation_groups"
    ADD CONSTRAINT "RECOMMENDATION_GROUPS_pkey" PRIMARY KEY ("recommendation_id", "group_id");

ALTER TABLE ONLY "public"."streaming_services"
    ADD CONSTRAINT "STREAMING_SERVICES_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."surveys"
    ADD CONSTRAINT "SURVEYS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."survey_answers"
    ADD CONSTRAINT "SURVEY_ANSWERS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."survey_questions"
    ADD CONSTRAINT "SURVEY_QUESTIONS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."users"
    ADD CONSTRAINT "USERS_pkey" PRIMARY KEY ("id");

ALTER TABLE ONLY "public"."watch_statuses"
    ADD CONSTRAINT "WATCH_STATUSES_pkey" PRIMARY KEY ("user_id", "movie_id");

ALTER TABLE ONLY "public"."favorite_movies"
    ADD CONSTRAINT "favorite_movies_user_movie_unique" UNIQUE ("user_id", "movie_id");

CREATE OR REPLACE TRIGGER "update_watch_statuses_updated_at" BEFORE UPDATE ON "public"."watch_statuses" FOR EACH ROW EXECUTE FUNCTION "public"."update_modified_column"();

ALTER TABLE ONLY "public"."chat_messages"
    ADD CONSTRAINT "CHAT_MESSAGES_room_id_fkey" FOREIGN KEY ("room_id") REFERENCES "public"."chat_rooms"("id");

ALTER TABLE ONLY "public"."chat_messages"
    ADD CONSTRAINT "CHAT_MESSAGES_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."chat_rooms"
    ADD CONSTRAINT "CHAT_ROOMS_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."chat_rooms"
    ADD CONSTRAINT "CHAT_ROOMS_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."groups"("id");

ALTER TABLE ONLY "public"."chat_rooms"
    ADD CONSTRAINT "CHAT_ROOMS_movie_id_fkey" FOREIGN KEY ("movie_id") REFERENCES "public"."movies"("id");

ALTER TABLE ONLY "public"."favorite_movies"
    ADD CONSTRAINT "FAVORITE_MOVIES_movie_id_fkey" FOREIGN KEY ("movie_id") REFERENCES "public"."movies"("id");

ALTER TABLE ONLY "public"."favorite_movies"
    ADD CONSTRAINT "FAVORITE_MOVIES_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."group_goals"
    ADD CONSTRAINT "GROUP_GOALS_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."groups"("id");

ALTER TABLE ONLY "public"."group_limit_logs"
    ADD CONSTRAINT "GROUP_LIMIT_LOGS_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."group_members"
    ADD CONSTRAINT "GROUP_MEMBERS_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."groups"("id");

ALTER TABLE ONLY "public"."group_members"
    ADD CONSTRAINT "GROUP_MEMBERS_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."movie_streamings"
    ADD CONSTRAINT "MOVIE_STREAMINGS_movie_id_fkey" FOREIGN KEY ("movie_id") REFERENCES "public"."movies"("id");

ALTER TABLE ONLY "public"."movie_streamings"
    ADD CONSTRAINT "MOVIE_STREAMINGS_service_id_fkey" FOREIGN KEY ("service_id") REFERENCES "public"."streaming_services"("id");

ALTER TABLE ONLY "public"."recommendations"
    ADD CONSTRAINT "RECOMMENDATIONS_from_user_id_fkey" FOREIGN KEY ("from_user_id") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."recommendations"
    ADD CONSTRAINT "RECOMMENDATIONS_movie_id_fkey" FOREIGN KEY ("movie_id") REFERENCES "public"."movies"("id");

ALTER TABLE ONLY "public"."recommendation_groups"
    ADD CONSTRAINT "RECOMMENDATION_GROUPS_group_id_fkey" FOREIGN KEY ("group_id") REFERENCES "public"."groups"("id");

ALTER TABLE ONLY "public"."recommendation_groups"
    ADD CONSTRAINT "RECOMMENDATION_GROUPS_recommendation_id_fkey" FOREIGN KEY ("recommendation_id") REFERENCES "public"."recommendations"("id");

ALTER TABLE ONLY "public"."survey_answers"
    ADD CONSTRAINT "SURVEY_ANSWERS_question_id_fkey" FOREIGN KEY ("question_id") REFERENCES "public"."survey_questions"("id");

ALTER TABLE ONLY "public"."survey_answers"
    ADD CONSTRAINT "SURVEY_ANSWERS_survey_id_fkey" FOREIGN KEY ("survey_id") REFERENCES "public"."surveys"("id");

ALTER TABLE ONLY "public"."survey_answers"
    ADD CONSTRAINT "SURVEY_ANSWERS_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");

ALTER TABLE ONLY "public"."survey_questions"
    ADD CONSTRAINT "SURVEY_QUESTIONS_survey_id_fkey" FOREIGN KEY ("survey_id") REFERENCES "public"."surveys"("id");

ALTER TABLE ONLY "public"."users"
    ADD CONSTRAINT "users_id_fkey" FOREIGN KEY ("id") REFERENCES "auth"."users"("id") ON DELETE CASCADE;

ALTER TABLE ONLY "public"."watch_statuses"
    ADD CONSTRAINT "watch_statuses_movie_id_fkey" FOREIGN KEY ("movie_id") REFERENCES "public"."movies"("id");

ALTER TABLE ONLY "public"."watch_statuses"
    ADD CONSTRAINT "watch_statuses_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");

CREATE POLICY "Allow all authenticated users" ON "public"."watch_statuses" TO "authenticated" USING (("auth"."uid"() = "user_id")) WITH CHECK (("auth"."uid"() = "user_id"));

CREATE POLICY "Allow authenticated users to read movies" ON "public"."movies" FOR SELECT TO "authenticated" USING (true);

CREATE POLICY "authenticated users can create groups" ON "public"."groups" FOR INSERT WITH CHECK (("auth"."uid"() IS NOT NULL));

CREATE POLICY "authenticated users can update groups" ON "public"."groups" FOR UPDATE USING ((("auth"."uid"() IS NOT NULL) AND (EXISTS ( SELECT 1
   FROM "public"."group_members"
  WHERE (("group_members"."group_id" = "groups"."id") AND ("group_members"."user_id" = "auth"."uid"())))))) WITH CHECK ((("auth"."uid"() IS NOT NULL) AND (EXISTS ( SELECT 1
   FROM "public"."group_members"
  WHERE (("group_members"."group_id" = "groups"."id") AND ("group_members"."user_id" = "auth"."uid"()))))));

ALTER TABLE "public"."chat_messages" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."chat_rooms" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."favorite_movies" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."group_goals" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."group_limit_logs" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."group_members" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."groups" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."movie_streamings" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."movies" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."recommendation_groups" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."recommendations" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."streaming_services" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."survey_answers" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."survey_questions" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."surveys" ENABLE ROW LEVEL SECURITY;

ALTER TABLE "public"."users" ENABLE ROW LEVEL SECURITY;

CREATE POLICY "users can create own group memberships" ON "public"."group_members" FOR INSERT WITH CHECK (("auth"."uid"() = "user_id"));

CREATE POLICY "users can read own group memberships" ON "public"."group_members" FOR SELECT USING (("auth"."uid"() = "user_id"));

ALTER TABLE "public"."watch_statuses" ENABLE ROW LEVEL SECURITY;

ALTER PUBLICATION "supabase_realtime" OWNER TO "postgres";

ALTER PUBLICATION "supabase_realtime" ADD TABLE ONLY "public"."chat_messages";

GRANT USAGE ON SCHEMA "public" TO "postgres";
GRANT USAGE ON SCHEMA "public" TO "anon";
GRANT USAGE ON SCHEMA "public" TO "authenticated";
GRANT USAGE ON SCHEMA "public" TO "service_role";

GRANT ALL ON TABLE "public"."groups" TO "anon";
GRANT ALL ON TABLE "public"."groups" TO "authenticated";
GRANT ALL ON TABLE "public"."groups" TO "service_role";

GRANT ALL ON FUNCTION "public"."create_group"("name" "text", "monthly_goal" integer) TO "anon";
GRANT ALL ON FUNCTION "public"."create_group"("name" "text", "monthly_goal" integer) TO "authenticated";
GRANT ALL ON FUNCTION "public"."create_group"("name" "text", "monthly_goal" integer) TO "service_role";

GRANT ALL ON FUNCTION "public"."handle_new_user"() TO "anon";
GRANT ALL ON FUNCTION "public"."handle_new_user"() TO "authenticated";
GRANT ALL ON FUNCTION "public"."handle_new_user"() TO "service_role";

GRANT ALL ON FUNCTION "public"."join_group"("group_id" "uuid") TO "anon";
GRANT ALL ON FUNCTION "public"."join_group"("group_id" "uuid") TO "authenticated";
GRANT ALL ON FUNCTION "public"."join_group"("group_id" "uuid") TO "service_role";

GRANT ALL ON FUNCTION "public"."update_modified_column"() TO "anon";
GRANT ALL ON FUNCTION "public"."update_modified_column"() TO "authenticated";
GRANT ALL ON FUNCTION "public"."update_modified_column"() TO "service_role";

GRANT ALL ON TABLE "public"."chat_messages" TO "anon";
GRANT ALL ON TABLE "public"."chat_messages" TO "authenticated";
GRANT ALL ON TABLE "public"."chat_messages" TO "service_role";

GRANT ALL ON TABLE "public"."chat_rooms" TO "anon";
GRANT ALL ON TABLE "public"."chat_rooms" TO "authenticated";
GRANT ALL ON TABLE "public"."chat_rooms" TO "service_role";

GRANT ALL ON TABLE "public"."favorite_movies" TO "anon";
GRANT ALL ON TABLE "public"."favorite_movies" TO "authenticated";
GRANT ALL ON TABLE "public"."favorite_movies" TO "service_role";

GRANT ALL ON TABLE "public"."group_goals" TO "anon";
GRANT ALL ON TABLE "public"."group_goals" TO "authenticated";
GRANT ALL ON TABLE "public"."group_goals" TO "service_role";

GRANT ALL ON TABLE "public"."group_limit_logs" TO "anon";
GRANT ALL ON TABLE "public"."group_limit_logs" TO "authenticated";
GRANT ALL ON TABLE "public"."group_limit_logs" TO "service_role";

GRANT ALL ON TABLE "public"."group_members" TO "anon";
GRANT ALL ON TABLE "public"."group_members" TO "authenticated";
GRANT ALL ON TABLE "public"."group_members" TO "service_role";

GRANT ALL ON TABLE "public"."movie_streamings" TO "anon";
GRANT ALL ON TABLE "public"."movie_streamings" TO "authenticated";
GRANT ALL ON TABLE "public"."movie_streamings" TO "service_role";

GRANT ALL ON TABLE "public"."movies" TO "anon";
GRANT ALL ON TABLE "public"."movies" TO "authenticated";
GRANT ALL ON TABLE "public"."movies" TO "service_role";

GRANT ALL ON TABLE "public"."recommendation_groups" TO "anon";
GRANT ALL ON TABLE "public"."recommendation_groups" TO "authenticated";
GRANT ALL ON TABLE "public"."recommendation_groups" TO "service_role";

GRANT ALL ON TABLE "public"."recommendations" TO "anon";
GRANT ALL ON TABLE "public"."recommendations" TO "authenticated";
GRANT ALL ON TABLE "public"."recommendations" TO "service_role";

GRANT ALL ON TABLE "public"."streaming_services" TO "anon";
GRANT ALL ON TABLE "public"."streaming_services" TO "authenticated";
GRANT ALL ON TABLE "public"."streaming_services" TO "service_role";

GRANT ALL ON TABLE "public"."survey_answers" TO "anon";
GRANT ALL ON TABLE "public"."survey_answers" TO "authenticated";
GRANT ALL ON TABLE "public"."survey_answers" TO "service_role";

GRANT ALL ON TABLE "public"."survey_questions" TO "anon";
GRANT ALL ON TABLE "public"."survey_questions" TO "authenticated";
GRANT ALL ON TABLE "public"."survey_questions" TO "service_role";

GRANT ALL ON TABLE "public"."surveys" TO "anon";
GRANT ALL ON TABLE "public"."surveys" TO "authenticated";
GRANT ALL ON TABLE "public"."surveys" TO "service_role";

GRANT ALL ON TABLE "public"."users" TO "anon";
GRANT ALL ON TABLE "public"."users" TO "authenticated";
GRANT ALL ON TABLE "public"."users" TO "service_role";

GRANT ALL ON TABLE "public"."watch_statuses" TO "anon";
GRANT ALL ON TABLE "public"."watch_statuses" TO "authenticated";
GRANT ALL ON TABLE "public"."watch_statuses" TO "service_role";

ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES TO "postgres";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES TO "anon";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES TO "authenticated";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON SEQUENCES TO "service_role";

ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS TO "postgres";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS TO "anon";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS TO "authenticated";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON FUNCTIONS TO "service_role";

ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES TO "postgres";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES TO "anon";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES TO "authenticated";
ALTER DEFAULT PRIVILEGES FOR ROLE "postgres" IN SCHEMA "public" GRANT ALL ON TABLES TO "service_role";

drop extension if exists "pg_net";

CREATE TRIGGER on_auth_user_created AFTER INSERT ON auth.users FOR EACH ROW EXECUTE FUNCTION public.handle_new_user();
