import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import Tabs from './Tabs.svelte';

const tabs = [
	{ id: 'a', label: 'Alpha' },
	{ id: 'b', label: 'Beta', count: 3 }
];

describe('Tabs', () => {
	it('renders every tab with its count', async () => {
		const screen = await render(Tabs, { tabs });
		await expect.element(screen.getByRole('button', { name: 'Alpha' })).toBeVisible();
		await expect.element(screen.getByRole('button', { name: 'Beta 3' })).toBeVisible();
	});

	it('calls onchange with the clicked tab id', async () => {
		const onchange = vi.fn();
		const screen = await render(Tabs, { tabs, onchange });
		await screen.getByRole('button', { name: /Beta/ }).click();
		expect(onchange).toHaveBeenCalledWith('b');
	});
});
