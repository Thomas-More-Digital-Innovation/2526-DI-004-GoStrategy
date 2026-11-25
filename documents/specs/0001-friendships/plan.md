# 0001: Friendship system, plan

Three vertical PRs. Each one ships something a user can see and passes CI on its own. Follow `CLAUDE.md`: tests first, the 300-line file limit, swag annotations, run only affected tests.

## PR 1: login `?next=` redirect

Independent of the rest. It makes invite links survive login and registration.

1. Pure validator in `code/frontend/src/routes/login/_utils/next.ts`, with `next.spec.ts` covering the spec's open-redirect cases.
2. `code/frontend/src/routes/login/+page.svelte`: after login or register, `goto` the validated `next` instead of `/`.
3. Run `just web check` and the new spec file.

## PR 2: friendships end to end

The whole feature. Large, so keep one commit per layer (data → REST → UI) to make it reviewable commit by commit.

### Commit 1: data layer

1. `code/backend/internal/models/user.go`: add `FriendCode *string` with `json:"-"`, `size:8`, and a unique index `idx_users_friend_code`.
2. New `code/backend/internal/models/friendship.go`:
   - The `Friendship` model and the `FriendshipStatus` constants.
   - The pair index tag `uniqueIndex:idx_friendships_pair,expression:...`.
   - Response structs (`FriendCodeResponse`, `FriendEntry`, `FriendListResponse`).
3. `code/backend/internal/db/db.go`: add `&models.Friendship{}` to `AllModels`.
4. New `code/backend/internal/db/migrations/006_friendships_rls.sql`: the policies from the spec's RLS table, in the style of `004`.
5. New `code/backend/internal/db/friend_code.go`:
   - Generate, normalise and format, as pure functions.
   - `GetOrCreateFriendCode`, `RegenerateFriendCode`, `FindUserByFriendCode`.
6. New `code/backend/internal/db/friendship.go`:
   - `SendFriendRequest`, which handles crossed requests, caps, self-requests and the retry on a unique violation.
   - `AcceptFriendRequest`, `DeleteFriendship` and `ListFriendships`.
   - All of them run through `WithRLS`.
7. Tests:
   - `friend_code_test.go`: normalisation table, alphabet, format, collision retry.
   - `friendship_test.go`: every transition, crossed request, caps, soft-deleted user, concurrent send.
   - Postgres RLS cases in `rls_postgres_test.go`.

### Commit 2: REST

1. New `code/backend/internal/api/rest/handlers/friends.go`, with swag annotations. Split the friend-code handlers into `friend_code.go` if the file nears 300 lines.
2. `code/backend/internal/api/router.go`: register the spec's endpoint table, with `UserRateLimitMiddleware(actionLimiter)` on writes and `preview`.
3. Handler tests:
   - Status codes for each error in the endpoint table.
   - No `friend_code` in the `GET /users/:id` body.
   - Rate limit applied.
4. Run `just go docs` (swagger, asyncapi and tygo types) and commit the generated files.

### Commit 3: frontend

1. `code/frontend/src/lib/api/`: a `friends` client for the new endpoints, using the generated types.
2. New route `code/frontend/src/routes/friends/`:
   - `+page.svelte`: friends, incoming and outgoing lists; accept, decline, cancel and remove actions.
   - `_state/friends.svelte.ts`: a runes class that takes the API client as a constructor default.
   - `_components/`: the list item and `ShareCode` (code display, share button, regenerate with a confirmation).
   - `_utils/share.ts`: share-or-copy logic that takes a `navigator`-like object, plus `share.spec.ts`.
3. Confirmation modal and invite route:
   - `friends/_components/AddFriendModal.svelte`: reads the code, calls `preview`, shows the confirmation and sends on click.
   - `friends/add/+page.ts`: redirects `/friends/add?code=` to `/friends?add=`, which opens the modal.
   - Logged-out visitors go to `/login?next=...` (needs PR 1).
4. `code/frontend/src/lib/components/sidebar/SideBar.svelte`: add a `/friends` item to the nav.
5. `code/frontend/src/app.html`: generic `og:`/`twitter:` meta tags, with an image under `static/`.
6. Tests:
   - State spec with a fake client.
   - Playwright `code/frontend/tests/friends.e2e.ts` for the end-to-end flow in the spec's Verification table.

### Verify

`go test ./internal/db/ ./internal/api/...`, `just go lint`, `just web check`, the touched Vitest files, and `just e2e` for the new e2e.

## PR 3: profiles and friends leaderboard

1. `code/backend/internal/api/router.go`:
   - Fix `users.GET("/stats", ...)` to `users.GET("/:id/stats", ...)`. This is an existing bug and can be cherry-picked into its own `fix(api)` PR earlier.
   - Add `users.GET("/by-name/:username", ...)`.
2. `code/backend/internal/db/user.go`: add `GetUserByUsername`, case-insensitive.
3. Backend: `ListFriendships` joins `user_stats` for accepted entries, and `FriendEntry` gains `stats`.
4. Run `just go docs` to regenerate types.
5. New route `code/frontend/src/routes/profile/[username]/+page.svelte`:
   - Resolves through `by-name` and reuses the `profile/_components`.
   - Move a component to `src/lib/components` only if both routes now consume it, per `CLAUDE.md` §3.1.
6. Friends leaderboard component on the social tab: you plus friends, sorted as the spec describes. Friend names link to `/profile/[username]`.
7. Tests:
   - Handler tests: the fixed stats route returns the right user; `by-name` matches case-insensitively and has no `friend_code` in the body.
   - Sort order as a pure function with a spec.
   - Profile 404 for an unknown name.
8. `ROADMAP.md`: tick "Friendship System".
