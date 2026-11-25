import { friends as friendsApi } from '#lib/api/client.js';
import { toastStore } from '#lib/state/toast.svelte.js';
import type { FriendEntry } from '#lib/types/generated/models.js';

export type FriendsApi = Pick<
	typeof friendsApi,
	'list' | 'getCode' | 'regenerateCode' | 'accept' | 'remove'
>;

export class FriendsState {
	friends = $state<FriendEntry[]>([]);
	incoming = $state<FriendEntry[]>([]);
	outgoing = $state<FriendEntry[]>([]);
	code = $state<string | null>(null);
	loading = $state(true);

	private busy = false;

	constructor(private api: FriendsApi = friendsApi) {}

	async load() {
		this.loading = true;
		try {
			const [list, code] = await Promise.all([this.api.list(), this.api.getCode()]);
			this.apply(list);
			this.code = code.code;
		} catch (e) {
			toastStore.handleApiMessage(e, 'Failed to load friends');
		} finally {
			this.loading = false;
		}
	}

	accept(id: number) {
		return this.mutate(() => this.api.accept(id), 'Failed to accept request');
	}

	// decline, cancel and unfriend are the same call: the backend deletes the row for any of them
	remove(id: number) {
		return this.mutate(() => this.api.remove(id), 'Failed to remove friend');
	}

	async regenerateCode() {
		try {
			this.code = (await this.api.regenerateCode()).code;
		} catch (e) {
			toastStore.handleApiMessage(e, 'Failed to regenerate friend code');
		}
	}

	// a second click while one call runs would hit a row that is already gone and toast a 404
	private async mutate(action: () => Promise<void>, fallback: string) {
		if (this.busy) return;
		this.busy = true;
		try {
			await action();
			this.apply(await this.api.list());
		} catch (e) {
			toastStore.handleApiMessage(e, fallback);
		} finally {
			this.busy = false;
		}
	}

	private apply(list: {
		friends: FriendEntry[];
		incoming: FriendEntry[];
		outgoing: FriendEntry[];
	}) {
		this.friends = list.friends;
		this.incoming = list.incoming;
		this.outgoing = list.outgoing;
	}
}
