import { beforeEach, describe, expect, it, vi } from 'vitest';
import { auth } from '#lib/api/client.js';
import { authStore } from './auth.svelte';

vi.mock('#lib/api/client.js', () => ({
	auth: {
		getMe: vi.fn(),
		login: vi.fn(),
		register: vi.fn(),
		logout: vi.fn(),
		changePassword: vi.fn()
	}
}));

const api = vi.mocked(auth);
const user = { id: 1, username: 'alice' } as never;

describe('authStore', () => {
	beforeEach(async () => {
		vi.clearAllMocks();
		api.logout.mockResolvedValue(undefined as never);
		await authStore.logout();
	});

	it('loads the current user', async () => {
		api.getMe.mockResolvedValue(user);
		await authStore.check();
		expect(authStore.user).toBe(user);
		expect(authStore.isLoggedIn).toBe(true);
		expect(authStore.loading).toBe(false);
	});

	it('clears the user when getMe fails', async () => {
		api.getMe.mockRejectedValue(new Error('401'));
		await authStore.check();
		expect(authStore.isLoggedIn).toBe(false);
		expect(authStore.loading).toBe(false);
	});

	it('dedupes concurrent checks', async () => {
		api.getMe.mockResolvedValue(user);
		await Promise.all([authStore.check(), authStore.check()]);
		expect(api.getMe).toHaveBeenCalledTimes(1);
	});

	it.each([
		['login', () => authStore.login('a', 'b')],
		['register', () => authStore.register('a', 'b')],
		['changePassword', () => authStore.changePassword('a', 'b', 'b')]
	] as const)('%s refreshes the user', async (method, run) => {
		api.getMe.mockResolvedValue(user);
		api[method].mockResolvedValue(undefined as never);
		await authStore.check();
		await run();
		expect(api[method]).toHaveBeenCalled();
		expect(api.getMe).toHaveBeenCalledTimes(2);
	});

	it('does not refresh when login fails', async () => {
		api.login.mockRejectedValue(new Error('bad'));
		await expect(authStore.login('a', 'b')).rejects.toThrow('bad');
		expect(api.getMe).not.toHaveBeenCalled();
	});

	it('logout clears the user', async () => {
		api.getMe.mockResolvedValue(user);
		await authStore.check();
		await authStore.logout();
		expect(authStore.user).toBeNull();
	});
});
