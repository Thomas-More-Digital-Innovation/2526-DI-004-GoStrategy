import { expect, test } from '@playwright/test';

test.describe('Smoke Tests', () => {
	test('landing page loads and renders Command Center header', async ({ page }) => {
		await page.goto('/');
		await expect(page).toHaveTitle(/GoStrategy/);
		const header = page.locator('h1');
		await expect(header).toBeVisible();
		await expect(header).toHaveText(/Command Center/i);
	});
});
