<script lang="ts">
	import { page } from '$app/state';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { authStore } from '$lib/state/auth.svelte';
	import LoggedOut from './_components/LoggedOut.svelte';
	import LoggedIn from './_components/LoggedIn.svelte';
	import Title from '../Title.svelte';
	import favicon from '$lib/assets/favicon.webp';
	import TacticalPattern from '$lib/components/TacticalPattern.svelte';
	import Connecting from './_components/Connecting.svelte';
	import { ExternalLink } from '@lucide/svelte';

	const navItems = [
		{ name: 'Command Center', href: resolve('/') },
		{ name: 'Meet our AI', href: resolve('/meet-our-ai') },
		{ name: 'Profile', href: resolve('/profile') },
		{ name: 'Board Setups', href: resolve('/board-setups') },
		{ divider: true },
		{ name: 'Open Source', href: resolve('/opensource') },
		{
			name: 'Github',
			href: 'https://github.com/Thomas-More-Digital-Innovation/2526-DI-004-GoStrategy',
			external: true
		}
	];

	onMount(async () => {
		await authStore.check();
	});
</script>

<aside
	class="w-64 border-r border-white/5 bg-surface-elevated/30 backdrop-blur-xl flex flex-col fixed inset-y-0 overflow-hidden"
>
	<TacticalPattern variant="mixed" opacity={12} size={110} color="white" />

	<div class="px-4 py-8 flex justify-center z-10">
		<img src={favicon} alt="Logo" class="w-12 h-12" />
		<Title />
	</div>

	<!-- eslint-disable svelte/no-navigation-without-resolve -->
	<nav class="flex-1 space-y-1 px-4 relative z-10">
		{#each navItems as item, idx (item.name ?? idx)}
			{#if item.divider}
				<div class="border-t border-white/10 my-8"></div>
			{:else if item.external}
				<a
					href={item.href}
					target="_blank"
					rel="noopener noreferrer"
					class="flex items-center gap-1 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group text-white/50 hover:bg-white/5 hover:text-white"
				>
					<span>{item.name}</span>
					<ExternalLink class="w-3.5 h-3.5" />
				</a>
			{:else}
				<a
					href={item.href}
					class="flex items-center gap-1 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 group {item.href ===
					page.url.pathname
						? 'bg-brand-primary/10 text-brand-primary'
						: 'text-white/50 hover:bg-white/5 hover:text-white'}"
				>
					<span>{item.name}</span>
				</a>
			{/if}
		{/each}
	</nav>
	<!-- eslint-enable svelte/no-navigation-without-resolve -->

	<div
		class="px-4 py-2.5 flex items-center justify-between text-[11px] text-white/40 border-t border-white/5 bg-black/10 relative z-10"
	>
		<span>v{__VERSION__}</span>
		{#if __COMMIT_SHA__ !== 'local'}
			<a
				href={`https://github.com/${__REPO_NAME__}/commit/${__COMMIT_SHA__}`}
				target="_blank"
				rel="noopener noreferrer"
				class="font-mono text-white/40 hover:text-white transition-colors underline decoration-white/20 hover:decoration-white"
				title="View live commit on GitHub"
			>
				#{__COMMIT_SHA__}
			</a>
		{:else}
			<span class="font-mono text-white/30">#local</span>
		{/if}
	</div>

	<div class="group p-4 border-t border-brand-accent/20 bg-black/20 relative z-10">
		{#if authStore.loading}
			<Connecting />
		{:else if authStore.user}
			<LoggedIn />
		{:else}
			<LoggedOut />
		{/if}
	</div>
</aside>
