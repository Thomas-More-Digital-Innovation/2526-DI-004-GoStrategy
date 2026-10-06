<script lang="ts">
	import {
		validatePasswordInput,
		getPasswordStrengthScore,
		isPasswordStrong
	} from '$lib/utils/authValidation';
	import { Check, X } from '@lucide/svelte';

	interface Props {
		password: string;
		isValid?: boolean;
		strength?: number;
	}

	let { password, isValid = $bindable(false), strength = $bindable(0) }: Props = $props();

	const checks = $derived(validatePasswordInput(password));
	const strengthScore = $derived(getPasswordStrengthScore(password));
	const isStrong = $derived(isPasswordStrong(password));

	$effect(() => {
		isValid = isStrong;
		strength = strengthScore;
	});

	const criteria = $derived([
		{ passed: checks.isMinLength, label: '8+ Characters' },
		{ passed: checks.hasUppercase, label: 'Uppercase Letter (A-Z)' },
		{ passed: checks.hasLowercase, label: 'Lowercase Letter (a-z)' },
		{ passed: checks.hasNumber, label: 'Numeric Digit (0-9)' },
		{
			passed: checks.isValidFormat,
			label: 'Valid Characters Only',
			fullWidth: true
		}
	]);
</script>

<div class="space-y-2 animate-fade-in py-2.5 px-3 bg-black/45 border border-white/10 rounded-2xl">
	<div
		class="text-[9px] uppercase tracking-widest text-white/40 font-bold border-b border-white/5 pb-1.5 mb-2"
	>
		Password Requirements
	</div>
	<div class="grid grid-cols-1 sm:grid-cols-2 gap-x-3 gap-y-2 text-xs font-semibold">
		{#each criteria as item}
			<div
				class="flex items-center gap-2 transition-colors duration-200 {item.passed
					? 'text-emerald-400'
					: 'text-white/40'} {item.fullWidth
					? 'col-span-1 sm:col-span-2 border-t border-white/5 pt-1.5 mt-0.5'
					: ''}"
			>
				{#if item.passed}
					<Check class="size-3.5 text-emerald-400 shrink-0 stroke-3" />
				{:else}
					<X class="size-3.5 text-red-500/70 shrink-0 stroke-3" />
				{/if}
				<span class="text-[10px] tracking-wide uppercase">{item.label}</span>
			</div>
		{/each}
	</div>
</div>
