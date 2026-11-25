-- 006_friendships_rls.sql

-- The friendships table and its indexes come from GORM AutoMigrate, which runs before these files.
ALTER TABLE friendships ENABLE ROW LEVEL SECURITY;
ALTER TABLE friendships FORCE ROW LEVEL SECURITY;

-- Only the two people in a friendship can see it; a third user sees nothing.
DROP POLICY IF EXISTS friendships_select_policy ON friendships;
CREATE POLICY friendships_select_policy ON friendships
    FOR SELECT
    USING (
        requester_id = NULLIF(current_setting('app.current_user_id', true), '')::integer
        OR addressee_id = NULLIF(current_setting('app.current_user_id', true), '')::integer
    );

-- You can only send requests as yourself.
DROP POLICY IF EXISTS friendships_insert_policy ON friendships;
CREATE POLICY friendships_insert_policy ON friendships
    FOR INSERT
    WITH CHECK (requester_id = NULLIF(current_setting('app.current_user_id', true), '')::integer);

-- Only the addressee can answer a request.
DROP POLICY IF EXISTS friendships_update_policy ON friendships;
CREATE POLICY friendships_update_policy ON friendships
    FOR UPDATE
    USING (addressee_id = NULLIF(current_setting('app.current_user_id', true), '')::integer)
    WITH CHECK (addressee_id = NULLIF(current_setting('app.current_user_id', true), '')::integer);

-- Either side can decline, cancel or unfriend.
DROP POLICY IF EXISTS friendships_delete_policy ON friendships;
CREATE POLICY friendships_delete_policy ON friendships
    FOR DELETE
    USING (
        requester_id = NULLIF(current_setting('app.current_user_id', true), '')::integer
        OR addressee_id = NULLIF(current_setting('app.current_user_id', true), '')::integer
    );
