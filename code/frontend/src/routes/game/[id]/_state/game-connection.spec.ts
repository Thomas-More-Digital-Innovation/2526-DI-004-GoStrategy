import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { games } from '#lib/api/client.js';
import { serverStore } from '#lib/state/server-status.svelte.js';
import { toastStore } from '#lib/state/toast.svelte.js';
import { GameConnectionManager } from './game-connection.svelte';
import { createFakeSocket } from './fake-socket';

vi.mock('#lib/api/client.js', () => ({ games: { abandon: vi.fn() } }));
vi.mock('#lib/state/server-status.svelte.js', () => ({
	serverStore: { check: vi.fn(), isOnline: true }
}));

describe('GameConnectionManager', () => {
	let socket: ReturnType<typeof createFakeSocket>;
	let conn: GameConnectionManager;

	beforeEach(() => {
		vi.useFakeTimers();
		vi.stubGlobal('window', { location: { href: '' } });
		toastStore.toasts = [];
		socket = createFakeSocket();
		conn = new GameConnectionManager(socket);
	});

	afterEach(() => {
		vi.unstubAllGlobals();
		vi.useRealTimers();
		vi.clearAllMocks();
	});

	it('connects to the given game and seat', async () => {
		await conn.connect('g', 0);
		expect(socket.connect).toHaveBeenCalledWith('g', 0);
		expect(conn.connected).toBe(true);
	});

	it.each([
		[true, 'Game session not found'],
		[false, 'Failed to connect to game server']
	])('on failure with server online=%s shows a toast and goes home', async (online, msg) => {
		socket.connect.mockRejectedValue(new Event('error'));
		vi.mocked(serverStore).isOnline = online;
		await conn.connect('g', 0);
		expect(conn.connected).toBe(false);
		expect(toastStore.toasts[0].message).toContain(msg);
		vi.advanceTimersByTime(3000);
		expect(window.location.href).toBe('/');
	});

	it('reconnects with exponential backoff', async () => {
		socket.connect.mockRejectedValueOnce(new Error('1')).mockResolvedValueOnce(undefined);
		const done = conn.attemptReconnect('g', 0);
		expect(conn.isReconnecting).toBe(true);
		await vi.advanceTimersByTimeAsync(2000);
		await done;
		expect(socket.connect).toHaveBeenCalledTimes(2);
		expect(conn.connected).toBe(true);
		expect(conn.isReconnecting).toBe(false);
		expect(conn.reconnectAttempts).toBe(0);
	});

	it('ignores a second reconnect while one is running', async () => {
		const first = conn.attemptReconnect('g', 0);
		await conn.attemptReconnect('g', 0);
		await first;
		expect(socket.connect).toHaveBeenCalledTimes(1);
	});

	it('gives up after max attempts', async () => {
		socket.connect.mockRejectedValue(new Error('down'));
		vi.mocked(serverStore).isOnline = false;
		const done = conn.attemptReconnect('g', 0);
		// why: 2+4+8+16+32s backoff; runAllTimers would also expire the toast
		await vi.advanceTimersByTimeAsync(62_000);
		await done;
		expect(socket.connect).toHaveBeenCalledTimes(conn.maxReconnectAttempts);
		expect(conn.isReconnecting).toBe(false);
		expect(toastStore.toasts[0].message).toContain('multiple attempts');
	});

	it('abandon disconnects, calls the api and goes home even on error', async () => {
		vi.mocked(games.abandon).mockRejectedValue(new Error('x'));
		vi.spyOn(console, 'error').mockImplementation(() => {});
		await conn.abandonAndQuit('g');
		expect(socket.disconnect).toHaveBeenCalled();
		expect(games.abandon).toHaveBeenCalledWith('g');
		expect(window.location.href).toBe('/');
	});
});
