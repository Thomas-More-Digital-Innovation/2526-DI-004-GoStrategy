# 0001: Friendship system

Status: approved

## Why

First item of the "Social & PvP" roadmap section. Players need a way to connect before direct challenges or matchmaking make sense. Players add each other with a personal friend code, shared as a code or as a link, the way most games do it.

## Scope / Non-goals

In scope:

- A personal friend code per user, created lazily, regenerable.
- Friend requests by code only: send, accept, decline, cancel, unfriend.
- A shareable invite link with a native share sheet and a clipboard fallback.
- A social tab listing friends and pending requests.
- A post-login `?next=` redirect so invite links survive login and registration.
- Public profile pages at `/profile/[username]` and a friends leaderboard.

Non-goals (do not build, do not refactor toward):

- Adding friends by username or user search. This was decided against: it turns the code into a formality and opens request spam.
- Blocking, a `declined` status, or ignore lists. Regenerating the code is the abuse control.
- Avatars or profile customization. The unvalidated `ProfilePicture` field stays untouched and unrendered. It belongs to a later avatar-presets spec.
- Online presence, challenges, notifications or push channels. Pending requests are discovered by loading the social tab.
- Global leaderboards, ratings or Elo. Current stats are AI-only and farmable.
- Personalised link previews. The frontend is a static SPA, so previews use one generic Open Graph set.
- Friends-only profile visibility. Profiles stay public.
- Changing user IDs. The integer primary key stays, and usernames identify users in URLs.

## Design

### Friend code

- Stored on `users.friend_code`: 8 characters from the Crockford base32 alphabet `0123456789ABCDEFGHJKMNPQRSTVWXYZ`, uppercase, no dash. Nullable, with a unique index. Multiple NULLs are allowed in both Postgres and SQLite.
- Generated with `crypto/rand`, never derived from the user ID. That gives about 40 bits, enough at this scale when combined with rate limiting.
- **Created lazily:** `POST /users/me/friend-code` returns the existing code or generates one. Users who never open the social tab never get a code. No backfill.
- **Regenerate:** `POST /users/me/friend-code/regenerate` replaces the code, so old codes and links stop working.
- **On collision** (unique violation), retry up to a named constant number of attempts (5), then fail with 500.
- **Input normalisation, backend only:** strip spaces and dashes, uppercase, map `O→0` and `I/L→1`. A normalised value that isn't 8 valid characters is a 400.
- **Output format:** the backend always returns codes formatted as `XXXX-XXXX`. The frontend never parses or formats codes, so there is no mirrored rule and no contract fixture.
- **Never public:** `User.FriendCode` gets `json:"-"`, and `UserDTO` does not include it. Only the owner's `/users/me/friend-code*` endpoints return it.

### Friendships table

One directed row per relationship:

```go
type Friendship struct {
    ID          int              `gorm:"primaryKey"`
    RequesterID int              `gorm:"not null;index"`
    AddresseeID int              `gorm:"not null;index"`
    Status      FriendshipStatus `gorm:"not null;size:20"` // "pending" | "accepted"
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

- **Pair uniqueness:** a unique expression index on `(least(requester_id, addressee_id), greatest(requester_id, addressee_id))`, declared in the GORM tag so SQLite tests exercise it. Username uniqueness uses the same technique: `idx_users_username_lower` in `internal/models/user.go`.
- **Hard delete:** no `gorm.DeletedAt`, unlike every other model. A soft-deleted row would keep the pair index occupied and block re-requesting after a decline. Say so in a comment on the model.
- **Status transitions:**
  - Creation: `pending`.
  - `pending → accepted`: addressee only.
  - Decline (addressee), cancel (requester) and unfriend (either party): all `DELETE` the row.
- **Crossed requests:** if B sends to A while A→B is pending, accept the existing row in the same transaction and return it. If A→B is already accepted, return 409.
- **Concurrent sends:** if an insert hits the pair unique index, re-run the crossed-request lookup once. A second violation is a 409.
- **Self-request:** 400.
- **Caps (named constants):**
  - `MaxFriends = 200`, counting accepted rows per user. Checked for both parties on send and on accept. Over the cap returns 409.
  - `MaxOutgoingPending = 50` per requester. Over the cap returns 429.
- **Soft-deleted users:** list queries join `users` and exclude `deleted_at IS NOT NULL`, so friends who deleted their account disappear. Sending to a soft-deleted user's code returns 404, because the user lookup is soft-delete-aware by default.

Rejected alternatives:

- Two mirrored rows per friendship: double writes that can drift.
- A canonical ordered pair plus `requested_by`: easy to get the direction wrong in every query.

### Row-level security

Follows the existing convention (see migrations `004`/`005`). A new `006_friendships_rls.sql` enables and forces RLS on `friendships`. The table and indexes come from `AutoMigrate`, which runs before the SQL files.

| Command | Policy (`me` = `app.current_user_id`) |
| --- | --- |
| SELECT, DELETE | `requester_id = me OR addressee_id = me` |
| INSERT | `WITH CHECK requester_id = me` |
| UPDATE | `USING addressee_id = me WITH CHECK addressee_id = me` |

- All friendship DB functions run inside `WithRLS` with `WithUserID` set.
- The crossed-request accept is an UPDATE by the addressee, so it fits the policy.
- `users` and `user_stats` stay publicly readable, so the code lookup and list joins need no extra policy.

### REST API

All under `RequireAuth`. Writes and `preview` share the existing `actionLimiter` (10/min per user). No new limiter.

| Method | Path | Body / query | Success | Errors |
| --- | --- | --- | --- | --- |
| POST | `/users/me/friend-code` | none | 200 `{code}` | 401 |
| POST | `/users/me/friend-code/regenerate` | none | 200 `{code}` | 401 |
| GET | `/friends` | none | 200 `{friends, incoming, outgoing}` | 401 |
| GET | `/friends/preview` | `?code=` | 200 `{user: UserDTO}` | 400, 404 |
| POST | `/friends/requests` | `{code}` | 201 pending / 200 auto-accepted | 400 (invalid code or self), 404, 409, 429 |
| POST | `/friends/requests/:id/accept` | none | 200 | 404 (not pending or not addressee), 409 (cap) |
| DELETE | `/friends/:id` | none | 204 | 404 |

- **Entries:** each entry in `friends`, `incoming` and `outgoing` is `{id, user: UserDTO, created_at}`. The `id` is the friendship ID, used for accept and delete.
- **Unrevealing 404s:** a 404 from accept or delete does not reveal whether the row exists for someone else. RLS hides it anyway.
- **Annotations:** every handler carries full swag annotations, and response structs live in `internal/models` so `tygo` generates their TypeScript types.

### Login `?next=`

- Login and register both honour `?next=<path>` after success.
- Only accept values that start with `/`, do not start with `//` or `/\`, and contain no scheme. Anything else falls back to `/`. This closes the open-redirect hole.
- Unauthenticated visitors to `/friends/add` are sent to `/login?next=/friends/add?code=...`, with the inner value URL-encoded.

### Sharing

- **Invite URL:** `${origin}/friends/add?code=XXXX-XXXX`. A query param, because the static SPA serves `index.html` for every path.
- **Share button:**
  - If `navigator.canShare?.(data)` is true, call `navigator.share({ title: 'GoStrategy', text, url })`.
  - Otherwise, `navigator.clipboard.writeText(`${text} ${url}`)` plus a toast.
  - `AbortError` (the user closed the sheet) is ignored.
- **Text:** `I'm <username> on GoStrategy, add me as a friend! Friend code: XXXX-XXXX`. Don't repeat the URL inside `text`, or some share targets show it twice.
- **Confirmation step:** `/friends/add` shows the preview ("Add `<username>`?") and sends the request only on button press. It never sends on load, because links get prefetched and embedded.
- **Link previews:** static, generic Open Graph and Twitter meta tags in `app.html` (title, description, image from `static/`).

### Profiles and friends leaderboard

- `GET /users/by-name/:username` matches case-insensitively (`lower(username) = lower(?)`) and returns `UserDTO`. Public.
- **Route fix:** the router registers `GET /users/stats` while the handler and swag docs expect `/users/:id/stats`, so today the handler always gets ID 0. Re-register it as `/users/:id/stats`. Nothing in the frontend calls the broken path.
- **Profile route:** `/profile/[username]` resolves the user by name, then reuses the existing profile components with that user's stats and `/users/:id/games`. Visiting your own name renders the same view. `/profile` stays the "me" page.
- **Friends leaderboard:**
  - Lives on the social tab: you plus accepted friends.
  - Ordered by wins descending, then win rate, then username.
  - Shows W/L/D and total games.
  - `GET /friends` gains a `stats` field per accepted entry, read in one query joined with `user_stats` (no N+1).

## Risks

| Risk | Impact |
| --- | --- |
| RLS policies exist only in Postgres SQL. The default test suite is SQLite | A policy bug goes unnoticed in local runs |
| Expression indexes (pair and code) differ between dialects | Migration fails at server boot in prod |
| Race on simultaneous crossed requests or accepts | Duplicate rows, or a cap exceeded by one |
| Open redirect through `?next=` | Phishing through a trusted domain |
| Friend code brute force through `preview` and `requests` | Unwanted requests |
| The friend code leaks through a public DTO or generated type | The capability is exposed |
| A friendships list that is N+1 over users and stats | Slow social tab |
| `navigator.share` is missing (Firefox desktop) or rejects | The share button does nothing |

## Verification

| Risk | Caught by |
| --- | --- |
| RLS | New Postgres test next to `TestRowLevelSecurity`: a third user can't SELECT, UPDATE or DELETE a pair's row; the requester can't accept; the addressee can. CI sets `TEST_DB_DSN` (`backend-ci.yml`) |
| Dialects | The same index tags run under SQLite in `go test ./internal/db/`. Postgres coverage comes from the CI job above, which runs `AutoMigrate` |
| Race | DB test: concurrent A→B and B→A sends end in exactly one row. The unique index is the backstop, and a handler maps a unique violation on send to a retry of the crossed-request path, once |
| Open redirect | Unit test (`*.spec.ts`) for the `next` validator: `//evil.com`, `/\evil.com`, `https://evil.com` and `javascript:` all fall back to `/`; `/friends/add?code=X` passes |
| Brute force | Handler test: `actionLimiter` is applied to `preview` and `requests` (429 after the limit) |
| Code leak | Handler test: `GET /users/:id` and `GET /users/by-name/:username` bodies contain no `friend_code` |
| N+1 | Review of the query: one SELECT with joins per list |
| Share | Component or state test with an injected `navigator`-like object: share path, clipboard fallback, and `AbortError` swallowed |
| End to end | Playwright e2e: user A copies a code, user B (logged out) opens the link, registers, lands on the confirmation page, sends; A accepts; both see each other |

## Assumptions

- Cap values (`MaxFriends = 200`, `MaxOutgoingPending = 50`) and the retry count (5) are sensible defaults, not product decisions. Change the constants freely.
- The leaderboard sort order (wins, then win rate, then name) was not discussed.
- The share text wording is a placeholder the user may tweak.
- Login and register live on the same `/login` route (`routes/login/+page.svelte`), so one `next` handler covers both.
- The roadmap item "Social Sharing" refers to game invite links, not friend links. It stays open.
