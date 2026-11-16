import { expect, test, type Page, type WebSocketRoute } from '@playwright/test';

// why: full-stack suite, needs the Go backend on VITE_API_BASE (default http://localhost:8080)
const PASSWORD = 'E2ePassword1';

async function register(page: Page, username: string) {
	await page.goto('/login');
	await page.getByRole('button', { name: 'Sign Up' }).click();
	await page.getByPlaceholder('Enter commander username').fill(username);
	await page.getByPlaceholder('Enter password').fill(PASSWORD);
	await page.getByRole('button', { name: 'Deploy Commander' }).click();
	await expect(page).toHaveURL('/');
	await expect(page.getByRole('link', { name: username })).toBeVisible();
}

async function startGameVsAi(page: Page) {
	await page.goto('/select-ai?mode=human_vs_ai');
	await page.getByRole('button', { name: /FAFO/ }).click();
	await expect(page).toHaveURL(/\/game\/\w+/);
	await page.getByRole('button', { name: 'Start Game' }).click();
	await expect(page.getByText('Move History')).toBeVisible();
}

async function makeAnyMove(page: Page) {
	// why: setups are random, so try front-row pieces until one has a legal move (FAFO style)
	for (let x = 0; x < 10; x++) {
		await page.locator(`#cell-${x}-6`).click();
		const target = page.locator('.valid-move').first();
		try {
			await target.waitFor({ timeout: 1000 });
		} catch {
			continue;
		}
		await target.click();
		return;
	}
	throw new Error('no movable piece in the front row');
}

test.describe('full stack', () => {
	test('register, play a move against the AI', async ({ page, browserName }) => {
		await register(page, `e2e_${browserName}_${Date.now()}`);
		await startGameVsAi(page);
		await makeAnyMove(page);
		await expect(page.getByText(/Move \d+ of \d+/)).toBeVisible();
	});

	test('rejects a wrong password', async ({ page }) => {
		await page.goto('/login');
		await page.getByPlaceholder('Enter commander username').fill('nobody_here');
		await page.getByPlaceholder('Enter password').fill('WrongPassword1');
		await page.getByRole('button', { name: 'Enter Battlefield' }).click();
		await expect(page.getByText('Invalid username or password')).toBeVisible();
		await expect(page).toHaveURL(/\/login/);
	});

	test('unknown game shows not found and returns home', async ({ page }) => {
		await page.goto('/game/01ARZ3NDEKTSV4RRFFQ69G5FAV?mode=human_vs_ai');
		await expect(page.getByRole('heading', { name: 'Game Not Found' })).toBeVisible();
		await expect(page).toHaveURL('/', { timeout: 10_000 });
	});

	test('reconnects after the websocket drops', async ({ page, browserName }) => {
		const sockets: WebSocketRoute[] = [];
		await page.routeWebSocket(/\/game\//, (ws) => {
			ws.connectToServer();
			sockets.push(ws);
		});

		await register(page, `e2e_rc_${browserName}_${Date.now()}`);
		await startGameVsAi(page);
		await sockets.at(-1)!.close();

		await expect(page.getByText('Reconnected to game session successfully.')).toBeVisible({
			timeout: 15_000
		});
		expect(sockets.length).toBeGreaterThan(1);
		await expect(page.getByText('Move History')).toBeVisible();
	});
});
