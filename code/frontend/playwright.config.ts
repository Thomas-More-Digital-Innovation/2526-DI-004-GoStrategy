import { defineConfig, devices } from '@playwright/test';

const baseURL = 'http://localhost:4173';

export default defineConfig({
	// why: in CI the production containers from docker-compose.e2e.yml already serve the app on :4173
	webServer: process.env.CI
		? undefined
		: {
				command: 'pnpm run build && pnpm run preview -- --port 4173',
				url: baseURL,
				reuseExistingServer: true
			},
	use: { baseURL, reducedMotion: 'reduce' },
	testDir: 'tests',
	testMatch: '**/*.e2e.{ts,js}',
	projects: [
		{
			name: 'chromium',
			use: { ...devices['Desktop Chrome'] }
		},
		{
			name: 'firefox',
			use: { ...devices['Desktop Firefox'] }
		},
		{
			name: 'webkit',
			timeout: 60_000,
			use: { ...devices['Desktop Safari'] }
		}
	]
});
