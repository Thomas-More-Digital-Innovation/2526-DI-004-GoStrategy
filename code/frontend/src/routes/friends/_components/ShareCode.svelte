<script lang="ts">
	import Button from '#lib/components/ui/Button.svelte';
	import Card from '#lib/components/ui/Card.svelte';
	import Modal from '#lib/components/ui/Modal.svelte';
	import { toastStore } from '#lib/state/toast.svelte.js';
	import { Copy } from '@lucide/svelte';
	import { shareInvite } from '../_utils/share';

	interface Props {
		code: string | null;
		onRegenerate: () => Promise<void>;
	}

	let { code, onRegenerate }: Props = $props();
	let confirmOpen = $state(false);
	let regenerating = $state(false);

	async function share() {
		if (!code) return;
		try {
			const outcome = await shareInvite(navigator, { code, origin: location.origin });
			if (outcome === 'copied') toastStore.success('Invite copied to clipboard');
		} catch (e) {
			toastStore.handleApiMessage(e, 'Could not share your invite');
		}
	}

	async function copyCodeToClipboard() {
		if (!code) return;
		try {
			await navigator.clipboard.writeText(code);
			toastStore.success('Friend code copied');
		} catch (e) {
			toastStore.handleApiMessage(e, 'Could not copy your friend code');
		}
	}

	async function regenerate() {
		regenerating = true;
		await onRegenerate();
		regenerating = false;
		confirmOpen = false;
	}
</script>

<Card class="flex flex-wrap items-center justify-between gap-4">
	<div>
		<p class="text-xs uppercase tracking-wider text-white/50">Your friend code</p>
		<div class="flex">
			<p class="font-mono text-3xl font-bold tracking-widest text-white" data-testid="friend-code">
				{code ?? '····-····'}
			</p>
			<Button disabled={!code} onclick={copyCodeToClipboard} variant="ghost" size="sm">
				<Copy />
				<span class="sr-only">Copy friend code</span>
			</Button>
		</div>
	</div>
	<div class="flex gap-3">
		<Button disabled={!code} onclick={share}>Share invite</Button>
		<Button variant="outline" disabled={!code} onclick={() => (confirmOpen = true)}>
			Regenerate
		</Button>
	</div>
</Card>

<Modal
	bind:isOpen={confirmOpen}
	title="Regenerate friend code?"
	description="Your old code and every link you shared will stop working."
	maxWidth="md"
	onClose={() => (confirmOpen = false)}
>
	<p class="text-sm text-white/70">
		Existing friends and pending requests are not affected. Use this if someone is spamming you.
	</p>
	{#snippet actions()}
		<div class="flex justify-end gap-3">
			<Button variant="ghost" onclick={() => (confirmOpen = false)}>Cancel</Button>
			<Button loading={regenerating} onclick={regenerate}>Regenerate</Button>
		</div>
	{/snippet}
</Modal>
