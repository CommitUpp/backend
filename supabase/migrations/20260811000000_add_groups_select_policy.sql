CREATE POLICY "users can read joined groups"
ON "public"."groups"
FOR SELECT
TO "authenticated"
USING (
  EXISTS (
    SELECT 1
    FROM "public"."group_members"
    WHERE "group_members"."group_id" = "groups"."id"
      AND "group_members"."user_id" = "auth"."uid"()
      AND "group_members"."is_active" = true
  )
);
