import { beforeEach, describe, expect, it, vi } from 'vitest';
import { toastStore } from '#lib/state/toast.svelte.js';
import type { FriendEntry, FriendListResponse } from '#lib/types/generated/models.js';
import { FriendsState, type FriendsApi } from './friends.svelte';

function entry(id: number, username: string): FriendEntry {
	const user = { id: id * 10, username, created_at: '', updated_at: '' };
	return { id, user, created_at: '' };
}

const empty: FriendListResponse = { friends: [], incoming: [], outgoing: [] };

function fakeApi(list: FriendListResponse = empty) {
	const api = {
		list: vi.fn().mockResolvedValue(list),
		getCode: vi.fn().mockResolvedValue({ code: 'K7QX-M2PA' }),
		regenerateCode: vi.fn().mockResolvedValue({ code: 'ZZZZ-0000' }),
		accept: vi.fn().mockResolvedValue(undefined),
		remove: vi.fn().mockResolvedValue(undefined)
	};
	return { api, state: new FriendsState(api as FriendsApi) };
}

describe('FriendsState', () => {
	beforeEach(() => {
		toastStore.toasts = [];
	});

	it('loads the lists and the own code', async () => {
		const list = { friends: [entry(1, 'bob')], incoming: [entry(2, 'cy')], outgoing: [] };
		const { state } = fakeApi(list);

		await state.load();

		expect(state.friends.map((f) => f.user.username)).toEqual(['bob']);
		expect(state.incoming.map((f) => f.user.username)).toEqual(['cy']);
		expect(state.code).toBe('K7QX-M2PA');
		expect(state.loading).toBe(false);
	});

	it('stops loading and toasts when the load fails', async () => {
		const { api, state } = fakeApi();
		api.list.mockRejectedValue(new Error('boom'));

		await state.load();

		expect(state.loading).toBe(false);
		expect(toastStore.toasts.map((t) => t.message)).toEqual(['boom']);
	});

	it('refreshes the lists after accepting', async () => {
		const { api, state } = fakeApi();
		api.list.mockResolvedValueOnce({ ...empty, incoming: [entry(2, 'cy')] });
		await state.load();
		api.list.mockResolvedValueOnce({ ...empty, friends: [entry(2, 'cy')] });

		await state.accept(2);

		expect(api.accept).toHaveBeenCalledWith(2);
		expect(state.incoming).toEqual([]);
		expect(state.friends.map((f) => f.user.username)).toEqual(['cy']);
	});

	it('ignores a second click while the first call runs', async () => {
		const { api, state } = fakeApi({ ...empty, incoming: [entry(2, 'cy')] });
		await state.load();

		await Promise.all([state.accept(2), state.accept(2)]);

		expect(api.accept).toHaveBeenCalledTimes(1);
	});

	it('accepts again once the previous call has finished', async () => {
		const { api, state } = fakeApi();
		await state.accept(1);
		await state.accept(2);

		expect(api.accept).toHaveBeenCalledTimes(2);
	});

	it('keeps the lists and toasts when accepting fails', async () => {
		const { api, state } = fakeApi({ ...empty, incoming: [entry(2, 'cy')] });
		await state.load();
		api.accept.mockRejectedValue(new Error('Friend limit reached'));

		await state.accept(2);

		expect(state.incoming).toHaveLength(1);
		expect(api.list).toHaveBeenCalledTimes(1);
		expect(toastStore.toasts[0].message).toBe('Friend limit reached');
	});

	it('removes through the one delete call', async () => {
		const { api, state } = fakeApi({ ...empty, friends: [entry(1, 'bob')] });
		await state.load();
		api.list.mockResolvedValueOnce(empty);

		await state.remove(1);

		expect(api.remove).toHaveBeenCalledWith(1);
		expect(state.friends).toEqual([]);
	});

	it('replaces the code on regenerate and keeps it on failure', async () => {
		const { api, state } = fakeApi();
		await state.load();

		await state.regenerateCode();
		expect(state.code).toBe('ZZZZ-0000');

		api.regenerateCode.mockRejectedValue(new Error('slow down'));
		await state.regenerateCode();
		expect(state.code).toBe('ZZZZ-0000');
		expect(toastStore.toasts[0].message).toBe('slow down');
	});
});
