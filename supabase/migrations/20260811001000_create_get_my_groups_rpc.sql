CREATE OR REPLACE FUNCTION "public"."get_groups_by_user"("target_user_id" "uuid")
RETURNS TABLE (
  "id" "uuid",
  "name" "text",
  "member_count" integer
)
LANGUAGE "sql" SECURITY DEFINER
SET "search_path" TO 'public'
AS $$
  SELECT
    "groups"."id",
    "groups"."name",
    count("active_members"."user_id")::integer AS "member_count"
  FROM "public"."group_members" AS "self_members"
  INNER JOIN "public"."groups"
    ON "groups"."id" = "self_members"."group_id"
  INNER JOIN "public"."group_members" AS "active_members"
    ON "active_members"."group_id" = "groups"."id"
    AND "active_members"."is_active" = true
  WHERE $1 = "auth"."uid"()
    AND "self_members"."user_id" = $1
    AND "self_members"."is_active" = true
  GROUP BY
    "groups"."id",
    "groups"."name",
    "self_members"."joined_at"
  ORDER BY "self_members"."joined_at" DESC;
$$;

ALTER FUNCTION "public"."get_groups_by_user"("target_user_id" "uuid") OWNER TO "postgres";

GRANT ALL ON FUNCTION "public"."get_groups_by_user"("target_user_id" "uuid") TO "anon";
GRANT ALL ON FUNCTION "public"."get_groups_by_user"("target_user_id" "uuid") TO "authenticated";
GRANT ALL ON FUNCTION "public"."get_groups_by_user"("target_user_id" "uuid") TO "service_role";
