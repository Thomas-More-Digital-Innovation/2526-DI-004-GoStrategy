<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { authStore } from '#lib/state/auth.svelte.js';
	import Button from '#lib/components/ui/Button.svelte';
	import Card from '#lib/components/ui/Card.svelte';
	import Input from '#lib/components/ui/Input.svelte';
	import AddFriendModal from './_components/AddFriendModal.svelte';
	import FriendRow from './_components/FriendRow.svelte';
	import ShareCode from './_components/ShareCode.svelte';
	import { FriendsState } from './_state/friends.svelte';

	const social = new FriendsState();
	let codeInput = $state('');
	let addCode = $state('');
	let addOpen = $state(false);

	onMount(async () => {
		await authStore.check();
		if (!authStore.user) {
			const back = encodeURIComponent(page.url.pathname + page.url.search);
			// eslint-disable-next-line svelte/no-navigation-without-resolve -- fixed internal path
			goto(`${resolve('login')}?next=${back}`);
			return;
		}
		const invited = page.url.searchParams.get('add');
		if (invited) openAdd(invited);
		await social.load();
	});

	function openAdd(code: string) {
		addCode = code;
		addOpen = true;
	}

	function lookUp() {
		const code = codeInput.trim();
		if (code) openAdd(code);
	}

	// why: drop ?add= so a reload does not reopen the invite
	function closeAdd() {
		addOpen = false;
		if (page.url.searchParams.has('add')) goto(resolve('friends'), { replaceState: true });
	}
</script>

<svelte:head>
	<title>GoStrategy — Friends</title>
</svelte:head>

{#if authStore.user}
	<AddFriendModal
		bind:isOpen={addOpen}
		code={addCode}
		onClose={closeAdd}
		onSent={() => {
			closeAdd();
			social.load();
		}}
	/>
	<div class="space-y-8">
		<ShareCode code={social.code} onRegenerate={() => social.regenerateCode()} />

		<Card class="space-y-3">
			<h2 class="text-lg font-bold uppercase tracking-wider text-white">Add a friend</h2>
			<form
				class="flex items-end gap-3"
				onsubmit={(e) => {
					e.preventDefault();
					lookUp();
				}}
			>
				<Input class="flex-1" placeholder="Enter a friend code" bind:value={codeInput} />
				<Button type="submit" disabled={!codeInput.trim()}>Look up</Button>
			</form>
		</Card>

		{#if social.incoming.length > 0}
			<Card class="space-y-3">
				<h2 class="text-lg font-bold uppercase tracking-wider text-white">Requests</h2>
				<ul class="space-y-2">
					{#each social.incoming as entry (entry.id)}
						<FriendRow {entry}>
							{#snippet actions()}
								<Button size="sm" onclick={() => social.accept(entry.id)}>Accept</Button>
								<Button size="sm" variant="ghost" onclick={() => social.remove(entry.id)}>
									Decline
								</Button>
							{/snippet}
						</FriendRow>
					{/each}
				</ul>
			</Card>
		{/if}

		<Card class="space-y-3">
			<h2 class="text-lg font-bold uppercase tracking-wider text-white">
				Friends ({social.friends.length})
			</h2>
			{#if social.loading}
				<p class="text-sm text-white/50">Loading…</p>
			{:else if social.friends.length === 0}
				<p class="text-sm text-white/50">No friends yet. Share your code to get started.</p>
			{:else}
				<ul class="space-y-2">
					{#each social.friends as entry (entry.id)}
						<FriendRow {entry}>
							{#snippet actions()}
								<Button size="sm" variant="ghost" onclick={() => social.remove(entry.id)}>
									Remove
								</Button>
							{/snippet}
						</FriendRow>
					{/each}
				</ul>
			{/if}
		</Card>

		{#if social.outgoing.length > 0}
			<Card class="space-y-3">
				<h2 class="text-lg font-bold uppercase tracking-wider text-white">Sent requests</h2>
				<ul class="space-y-2">
					{#each social.outgoing as entry (entry.id)}
						<FriendRow {entry}>
							{#snippet actions()}
								<Button size="sm" variant="ghost" onclick={() => social.remove(entry.id)}>
									Cancel
								</Button>
							{/snippet}
						</FriendRow>
					{/each}
				</ul>
			</Card>
		{/if}
	</div>
{/if}
