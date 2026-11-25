import { expect, test, type Page } from '@playwright/test';

// why: full-stack suite, needs the Go backend on VITE_API_BASE (default http://localhost:8080)
const PASSWORD = 'E2ePassword1';

async function signUp(page: Page, username: string) {
	await page.getByRole('button', { name: 'Sign Up' }).click();
	await page.getByPlaceholder('Enter commander username').fill(username);
	await page.getByPlaceholder('Enter password').fill(PASSWORD);
	await page.getByRole('button', { name: 'Deploy Commander' }).click();
}

test.describe('friends', () => {
	test('a logged-out invitee registers through the link, sends, and the inviter accepts', async ({
		browser,
		browserName
	}) => {
		const stamp = `${browserName.slice(0, 2)}${Date.now() % 1_000_000}`;
		const alice = `fa_${stamp}`;
		const bob = `fb_${stamp}`;

		const aliceCtx = await browser.newContext();
		const alicePage = await aliceCtx.newPage();
		await alicePage.goto('/login');
		await signUp(alicePage, alice);
		await expect(alicePage).toHaveURL('/');
		await alicePage.goto('/friends');
		// why: the code renders as a placeholder until the API answers, so wait for the real one
		const codeBox = alicePage.getByTestId('friend-code');
		await expect(codeBox).toHaveText(/^[0-9A-Z]{4}-[0-9A-Z]{4}$/);
		const code = (await codeBox.textContent())!.trim();

		await alicePage.goto(`/friends/add?code=${code}`);
		await expect(alicePage.getByText('That is your own friend code')).toBeVisible();
		await expect(alicePage.getByRole('button', { name: 'Send request' })).toHaveCount(0);
		await alicePage.goto('/friends');

		const bobCtx = await browser.newContext();
		const bobPage = await bobCtx.newPage();
		const sends: string[] = [];
		bobPage.on('request', (r) => {
			if (r.method() === 'POST' && r.url().endsWith('/friends/requests')) sends.push(r.url());
		});
		await bobPage.goto(`/friends/add?code=${code}`);
		await expect(bobPage).toHaveURL(/\/login\?next=/);
		await signUp(bobPage, bob);
		await expect(bobPage.getByRole('heading', { name: `Add ${alice}?` })).toBeVisible();
		// why: opening an invite link must never send, because prefetchers and previews open URLs
		expect(sends).toHaveLength(0);
		await bobPage.getByRole('button', { name: 'Send request' }).click();
		await expect(bobPage).toHaveURL(/\/friends$/);
		expect(sends).toHaveLength(1);
		await expect(bobPage.getByText('Sent requests')).toBeVisible();

		await alicePage.reload();
		await expect(alicePage.getByText(bob)).toBeVisible();
		await alicePage.getByRole('button', { name: 'Accept' }).click();
		await expect(alicePage.getByRole('heading', { name: 'Friends (1)' })).toBeVisible();

		await bobPage.reload();
		await expect(bobPage.getByRole('heading', { name: 'Friends (1)' })).toBeVisible();
		await expect(bobPage.getByText(alice)).toBeVisible();

		await aliceCtx.close();
		await bobCtx.close();
	});

	test('an unknown code shows a message instead of a confirmation', async ({
		page,
		browserName
	}) => {
		await page.goto('/login');
		await signUp(page, `fu_${browserName.slice(0, 2)}${Date.now() % 1_000_000}`);
		await expect(page).toHaveURL('/');
		await page.goto('/friends/add?code=ZZZZ-ZZZZ');
		await expect(page).toHaveURL(/\/friends\?add=/);
		await expect(page.getByRole('button', { name: 'Cancel' })).toBeVisible();
		await expect(page.getByText('Looking up friend code')).toHaveCount(0);
		await expect(page.getByRole('button', { name: 'Send request' })).toHaveCount(0);
	});
});
