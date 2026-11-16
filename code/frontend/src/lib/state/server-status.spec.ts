import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { monitoring } from '#lib/api/client.js';
import { serverStore } from './server-status.svelte';

vi.mock('#lib/api/client.js', () => ({ monitoring: { health: vi.fn() } }));

const health = vi.mocked(monitoring.health);

describe('serverStore', () => {
	beforeEach(() => {
		vi.useFakeTimers();
		health.mockReset();
		vi.stubGlobal('navigator', { onLine: true });
	});

	afterEach(() => {
		serverStore.stopPolling();
		vi.unstubAllGlobals();
		vi.useRealTimers();
	});

	it('is online when health reports ok', async () => {
		health.mockResolvedValue({ status: 'ok' } as never);
		expect(await serverStore.check()).toBe(true);
		expect(serverStore.lastChecked).toBeInstanceOf(Date);
	});

	it('is offline when health fails or reports otherwise', async () => {
		health.mockRejectedValueOnce(new Error('down'));
		expect(await serverStore.check()).toBe(false);
		health.mockResolvedValueOnce({ status: 'degraded' } as never);
		expect(await serverStore.check()).toBe(false);
	});

	it('short-circuits when the browser is offline', async () => {
		vi.stubGlobal('navigator', { onLine: false });
		expect(await serverStore.check()).toBe(false);
		expect(health).not.toHaveBeenCalled();
	});

	it('polls on an interval and only starts once', async () => {
		health.mockResolvedValue({ status: 'ok' } as never);
		serverStore.startPolling(1000);
		serverStore.startPolling(1000);
		expect(health).toHaveBeenCalledTimes(1);
		await vi.advanceTimersByTimeAsync(2000);
		expect(health).toHaveBeenCalledTimes(3);
		serverStore.stopPolling();
		await vi.advanceTimersByTimeAsync(2000);
		expect(health).toHaveBeenCalledTimes(3);
	});
});
