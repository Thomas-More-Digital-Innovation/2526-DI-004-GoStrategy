<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { marked } from 'marked';
	import DocsNav from './_components/DocsNav.svelte';
	import DocsContent from './_components/DocsContent.svelte';

	interface NavItem {
		title: string;
		slug: string;
	}

	const fallbackNavItems: NavItem[] = [
		{ title: 'Overview', slug: 'Home' },
		{ title: 'Docker Deployment', slug: 'Docker-Deployment' },
		{ title: 'Configuration', slug: 'Configuration' },
		{ title: 'Architecture', slug: 'Architecture' }
	];

	let navItems = $state<NavItem[]>(fallbackNavItems);
	let activeSlug = $state('Home');
	let markdownContent = $state('');
	let renderedHtml = $state('');
	let loading = $state(true);
	let error = $state('');

	// parse _Sidebar.md bullet links into nav items
	function parseSidebarMarkdown(raw: string): NavItem[] {
		const items: NavItem[] = [];
		const regex = /^\*\s*\[(.*?)\]\((.*?)\)/gm;
		let match;
		while ((match = regex.exec(raw)) !== null) {
			items.push({
				title: match[1].trim(),
				slug: match[2].trim().replace(/\.md$/, '')
			});
		}
		return items.length > 0 ? items : fallbackNavItems;
	}

	async function loadSidebar() {
		try {
			const res = await fetch('/docs/_Sidebar.md');
			if (res.ok) {
				const text = await res.text();
				navItems = parseSidebarMarkdown(text);
			}
		} catch {
			// ponytail: fallback to hardcoded list if sidebar load fails
		}
	}

	async function loadDoc(slug: string) {
		loading = true;
		error = '';
		try {
			const res = await fetch(`/docs/${slug}.md`);
			if (!res.ok) throw new Error(`Document "${slug}" could not be loaded`);
			markdownContent = await res.text();
			const parsed = await marked.parse(markdownContent);
			renderedHtml = parsed;
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to load document';
			renderedHtml = '';
		} finally {
			loading = false;
		}
	}

	function handleSelect(slug: string) {
		const targetSlug = slug.replace(/\.md$/, '');
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`${resolve('/docs')}?page=${targetSlug}`, { keepFocus: true, noScroll: false });
	}

	function handleContentClick(event: MouseEvent) {
		const target = (event.target as HTMLElement).closest('a');
		if (!target) return;

		const href = target.getAttribute('href');
		if (
			!href ||
			href.startsWith('http://') ||
			href.startsWith('https://') ||
			href.startsWith('#')
		) {
			return;
		}

		// internal relative doc links (e.g., Docker-Deployment or ./Docker-Deployment.md)
		event.preventDefault();
		const cleanSlug = href.replace(/^\.\//, '').replace(/\.md$/, '');
		handleSelect(cleanSlug);
	}

	$effect(() => {
		const queryPage = page.url.searchParams.get('page') || 'Home';
		activeSlug = queryPage;
		loadDoc(queryPage);
	});

	$effect(() => {
		loadSidebar();
	});
</script>

<svelte:head>
	<title>GoStrategy — Documentation ({activeSlug})</title>
	<meta
		name="description"
		content="Official documentation for GoStrategy: deployment, configuration, and architecture."
	/>
</svelte:head>

<div class="space-y-8 pb-12">
	<div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
		<aside class="lg:col-span-1">
			<div class="sticky top-6">
				<DocsNav items={navItems} {activeSlug} onSelect={handleSelect} />
			</div>
		</aside>

		<div class="lg:col-span-3">
			{#if loading}
				<div class="p-12 text-center text-white/40 font-mono text-sm">Loading documentation...</div>
			{:else if error}
				<div
					class="p-8 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-sm space-y-2"
				>
					<p class="font-semibold">Failed to load documentation</p>
					<p class="text-xs text-red-400/80">{error}</p>
				</div>
			{:else}
				<!-- svelte-ignore a11y_click_events_have_key_events -->
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div onclick={handleContentClick}>
					<DocsContent html={renderedHtml} {activeSlug} items={navItems} onSelect={handleSelect} />
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	:global(.docs-prose h1) {
		font-size: 1.75rem;
		font-weight: 800;
		color: #ffffff;
		margin-bottom: 1rem;
		margin-top: 1.5rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
		padding-bottom: 0.5rem;
	}
	:global(.docs-prose h1:first-child) {
		margin-top: 0;
	}
	:global(.docs-prose h2) {
		font-size: 1.25rem;
		font-weight: 700;
		color: #ffffff;
		margin-top: 2rem;
		margin-bottom: 0.75rem;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		padding-bottom: 0.25rem;
	}
	:global(.docs-prose h3) {
		font-size: 1rem;
		font-weight: 600;
		color: #ffffff;
		margin-top: 1.5rem;
		margin-bottom: 0.5rem;
	}
	:global(.docs-prose p) {
		margin-bottom: 1rem;
		line-height: 1.7;
		color: rgba(255, 255, 255, 0.75);
	}
	:global(.docs-prose ul) {
		list-style-type: disc;
		padding-left: 1.25rem;
		margin-bottom: 1rem;
	}
	:global(.docs-prose ol) {
		list-style-type: decimal;
		padding-left: 1.25rem;
		margin-bottom: 1rem;
	}
	:global(.docs-prose li) {
		margin-bottom: 0.35rem;
		color: rgba(255, 255, 255, 0.75);
	}
	:global(.docs-prose a) {
		color: var(--color-brand-primary, #67e8f9);
		text-decoration: underline;
		text-underline-offset: 2px;
	}
	:global(.docs-prose a:hover) {
		opacity: 0.8;
	}
	:global(.docs-prose code) {
		background-color: rgba(0, 0, 0, 0.5);
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: 0.15rem 0.4rem;
		border-radius: 0.375rem;
		font-family: monospace;
		font-size: 0.85em;
		color: #f3f4f6;
	}
	:global(.docs-prose pre) {
		background-color: rgba(0, 0, 0, 0.6);
		border: 1px solid rgba(255, 255, 255, 0.1);
		padding: 1rem;
		border-radius: 0.5rem;
		overflow-x: auto;
		margin-bottom: 1.25rem;
	}
	:global(.docs-prose pre code) {
		background-color: transparent;
		border: none;
		padding: 0;
	}
	:global(.docs-prose hr) {
		border-color: rgba(255, 255, 255, 0.1);
		margin: 2rem 0;
	}
	:global(.docs-prose table) {
		width: 100%;
		border-collapse: collapse;
		margin-bottom: 1.5rem;
		font-size: 0.875rem;
	}
	:global(.docs-prose th) {
		text-align: left;
		padding: 0.6rem 0.85rem;
		background-color: rgba(255, 255, 255, 0.05);
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: #ffffff;
		font-weight: 600;
	}
	:global(.docs-prose td) {
		padding: 0.6rem 0.85rem;
		border: 1px solid rgba(255, 255, 255, 0.1);
		color: rgba(255, 255, 255, 0.75);
	}
</style>
