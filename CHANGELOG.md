# Changelog

All notable changes to this project will be documented in this file.

## [0.4.0] - 07-10-2026

### Added
- **AI Opponents**: Introduced Heuristic, Minimax, and Monte Carlo Tree Search (MCTS) engine implementations available for games and simulation benchmarking.
- **Board Setup Export**: Added export functionality for custom board setups to JSON and binary formats for analysis and external training.
- **Connection Status**: Added real-time WebSocket connection state indicator to the profile card in the sidebar.
- **Game Matchup Display**: Added matchup overview to the game info card displaying opponent and AI details.
- **Testing & Quality Tooling**: Integrated ESLint, Prettier, Vitest, and Playwright alongside Lefthook pre-commit hooks for comprehensive code quality and test automation.

### Changed
- **Frontend Architecture**: Overhauled frontend architecture into a feature-first structure, colocating route-specific components, state machines, and controllers within `routes/game/[id]`.
- **Game Identifiers**: Migrated custom game identifiers to standardized ULID generation across backend models and routes.
- **Iconography**: Replaced raw emoji icons across the user interface with Lucide SVG icons.
- **Visual Refinements**: Centered and refined the game board layout, incorporated the GoStrategy logo into the game view, and removed question mark placeholders on empty tiles.
- **CI/CD Workflows**: Corrected workflow permissions in CD and CodeQL pipelines and removed auto-commits from the backend CI pipeline.
- **Dependency Management**: Synchronized Bun and PNPM lockfiles, bumped pipeline tool versions, and updated core package dependencies.

### Fixed
- **Setup Mode Exit**: Fixed issue where navigating back to the menu during the board setup phase failed to terminate the active game session.
- **Board Alignment**: Resolved CSS layout issue causing the game board to render off-center on wide viewports.
- **Setup Validation**: Enforced backend piece count and board boundary validations on custom board setup submissions.

### Security
- **Vulnerability Audits**: Resolved automated audit security vulnerabilities across backend and frontend dependencies.

## [0.3.2] - 03-06-2026

### Added
- **Game Duration Tracking**: Recorded explicit `started_at` and `ended_at` timestamps for game sessions to track match durations.

### Changed
- **Codebase Architecture & DRY**: Consolidated packages into `internal/`, eliminated duplicate logic, and standardized unit tests using `stretchr/testify/assert`.
- **API Documentation**: Aligned OpenAPI/Swagger specifications with restructured API routes and schema contracts.
- **CI & Formatting**: Hardened backend CI workflows and integrated automated formatting checks.

### Fixed
- **AI Match Stats**: Fixed stat resolution to ensure AI vs AI matches do not inadvertently award player win/loss records.
- **Log-Identified Edge Cases**: Resolved runtime edge cases and error handling regressions uncovered during log analysis.

## [0.3.1] - 28-05-2026

### Added
- **Accessibility (a11y)**: Integrated comprehensive keyboard navigation and ARIA labeling on the game board.
- **Rate Limiting**: Added secondary user-level rate limits for sensitive actions (creation, password changes).
- **AI vs AI Auto-Pause**: Implemented automatic pausing of AI vs AI games when no observers are connected.
- **UX Feedback**: Added "Reconnecting..." status indicators for socket drops.
- **Reconnect Button**: Added a button on the home page to allow users to reconnect to active games waiting for cleanup.

### Changed
- **Authentication**: Redesigned and improved the login and registration screens for a premium feel.
- **Refactoring**: Refactored the backend code, specifically targeting the API package.
- **Security**: Hardened username input validation using stricter regex sanitization to prevent injection.
- **Reliability**: Configured exponential backoff for WebSocket reconnections with reactive UI states.
- **AI vs AI Lifecycle**: Added a 30-second wait interval for grace-period cleanup of AI vs AI games.
- **Testing**: Increased test coverage to >60% across API packages and WebSocket connections.

### Fixed
- **Game Limiting**: Fixed issue allowing multiple active games per user, restricting to one active game at a time.

## [0.3.0] - 07-05-2026

### Added
- **Data Integrity**: Implement `deleted_at` patterns (soft deletes) for users and game sessions.

### Changed
- **Database**: Remove schema.sql and create a migration system (gorm or similar).
- **Rebranding**: Rebranding to GoStrategy for legal reasons #StayOpenSource.

## [0.2.2] - 04-05-2026

### Added
- Message on screens smaller than 1200px telling users to switch to a larger screen.
- Add gosec & npm audit to github actions.

### Changed
- Changed the AI vs AI board setup interface to clearly indicate which AI you are setting up.
- Add highlighting around the board to indicate whose turn it is.

### Fixed
- Fixed broken CHANGELOG.md symlink.
- Color of setup board is opposite of what it should be during setup phase.

## [0.2.1] - 03-05-2026

### Added
- Ability to see the username or AI name in the game info.
- Ability to change the password of the currently logged in user.

### Changed
- Show full early access disclaimer on the home screen.
- Centralized toast messages for better UX.
- Remove alert at end of game.

### Fixed
- `Profile` & `Board Setup` pages showing login screen after reloading the page.
- Server status indicator in the sidebar not working.
- Name tag at bottom of sidebar not routing to profile page.
- CodeQL scanning paths optimized to reduce unnecessary CI compute.

## [0.2.0] - 02-05-2026

### Added
- **Database Migrations:** Implemented a robust SQL migration system using Go's `embed` filesystem.
- **Session Persistence:** Games and moves are now automatically saved to the database.
- **User Stats:** Automated tracking of wins, losses, and move counts for registered users.
- **Setup Phase Timeout:** Added a 5-minute automated setup phase with warnings to prevent stalled games.
- **JWT v5 Migration:** Upgraded security infrastructure to use the latest `golang-jwt/v5` standard.
- **IP Rate Limiting:** Added protection against brute-force attacks on authentication endpoints.

### Changed
- **Concurrency Overhaul:** Optimized `GameRunner` and `WSHub` with event-driven channels, eliminating busy-waiting and reducing mutex contention.
- **Standardized Logging:** Refactored logging to use a consistent `[Time] [Tag] [Location] [User] Message` format.
- **AI Turn Pacing:** Improved AI move synchronization to ensure smooth animations in the frontend.

### Fixed
- Illegal piece movement validation (e.g., Flags can no longer move).
- Game state broadcasting regressions during high-concurrency sessions.
- Authorization leaks where spectators could occasionally attempt to send moves.

## [0.1.1] - 13-04-2026

### Fixed
- Mutex contention issues identified in system audit.
- Empty CORS origin handling.
- API path parameter inconsistencies.

## [0.1.0] - 03-04-2026

### Added
- Initial functional release of GoStrategy with Human-vs-AI and AI-vs-AI modes.
- Basic WebSocket implementation for game updates.
- In-memory session management.
