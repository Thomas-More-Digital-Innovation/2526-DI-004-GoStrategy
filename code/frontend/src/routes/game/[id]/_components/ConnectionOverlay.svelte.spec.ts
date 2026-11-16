import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ConnectionOverlay from './ConnectionOverlay.svelte';

const base = { gameId: 'g-123', reconnectAttempts: 2, maxReconnectAttempts: 5 };

describe('ConnectionOverlay', () => {
	it('shows reconnect progress', async () => {
		const screen = await render(ConnectionOverlay, { ...base, isReconnecting: true });
		await expect.element(screen.getByRole('heading', { name: 'Connection Lost' })).toBeVisible();
		await expect.element(screen.getByText('Attempt 2 of 5')).toBeVisible();
	});

	it('shows game not found and returns to menu', async () => {
		const onReturnToMenu = vi.fn();
		const screen = await render(ConnectionOverlay, {
			...base,
			isReconnecting: false,
			onReturnToMenu
		});
		await expect.element(screen.getByRole('heading', { name: 'Game Not Found' })).toBeVisible();
		await screen.getByRole('button', { name: 'Return to Menu' }).click();
		expect(onReturnToMenu).toHaveBeenCalled();
	});
});
