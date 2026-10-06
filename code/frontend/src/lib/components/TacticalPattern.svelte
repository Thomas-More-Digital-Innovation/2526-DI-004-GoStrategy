<script lang="ts">
	import { onMount } from 'svelte';
	import flankUrl from '$lib/assets/patterns/tactical-flank.svg';
	import pincerUrl from '$lib/assets/patterns/tactical-pincer.svg';
	import infiltrateUrl from '$lib/assets/patterns/tactical-infiltrate.svg';
	import blitzUrl from '$lib/assets/patterns/tactical-blitz.svg';

	export type TacticalVariant = 'mixed' | 'random' | 'flank' | 'pincer' | 'infiltrate' | 'blitz';

	interface Props {
		variant?: TacticalVariant;
		opacity?: number; // 0 to 100
		size?: number; // pixel tile size per sector
		color?: string;
		class?: string;
	}

	let {
		variant = 'mixed',
		opacity = 12,
		size = 130,
		color = 'white',
		class: className = ''
	}: Props = $props();

	type PatternKey = 'flank' | 'pincer' | 'infiltrate' | 'blitz';

	const patterns: Record<PatternKey, string> = {
		flank: flankUrl,
		pincer: pincerUrl,
		infiltrate: infiltrateUrl,
		blitz: blitzUrl
	};

	const variantKeys: PatternKey[] = ['flank', 'pincer', 'infiltrate', 'blitz'];

	let randomPick = $state<PatternKey | null>(null);

	onMount(() => {
		if (variant === 'random') {
			randomPick = variantKeys[Math.floor(Math.random() * variantKeys.length)];
		}
	});

	const resolvedVariant = $derived(variant === 'random' ? (randomPick ?? 'mixed') : variant);

	const patternId = `tactical-pattern-${Math.random().toString(36).substring(2, 9)}`;
	const isMixed = $derived(resolvedVariant === 'mixed');
	const patternWidth = $derived(isMixed ? size * 2 : size);
	const patternHeight = $derived(isMixed ? size * 2 : size);
	const singlePatternUrl = $derived(resolvedVariant !== 'mixed' ? patterns[resolvedVariant] : '');
</script>

<div
	class="absolute inset-0 pointer-events-none select-none overflow-hidden z-0 {className}"
	style:opacity={opacity / 100}
>
	<svg width="100%" height="100%" xmlns="http://www.w3.org/2000/svg">
		<defs>
			<pattern
				id={patternId}
				width={patternWidth}
				height={patternHeight}
				patternUnits="userSpaceOnUse"
			>
				{#if isMixed}
					<image href={flankUrl} x="0" y="0" width={size} height={size} />
					<image href={pincerUrl} x={size} y="0" width={size} height={size} />
					<image href={infiltrateUrl} x="0" y={size} width={size} height={size} />
					<image href={blitzUrl} x={size} y={size} width={size} height={size} />
				{:else}
					<image href={singlePatternUrl} x="0" y="0" width={size} height={size} />
				{/if}
			</pattern>
		</defs>
		<rect width="100%" height="100%" fill="url(#{patternId})" />
	</svg>
</div>
