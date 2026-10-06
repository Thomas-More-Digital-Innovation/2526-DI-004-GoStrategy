<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { authStore } from '$lib/state/auth.svelte';
	import { toastStore } from '$lib/state/toast.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import TacticalBriefing from './_components/TacticalBriefing.svelte';
	import PasswordChecker from './_components/PasswordChecker.svelte';
	import { User, Lock } from '@lucide/svelte';

	let username = $state('');
	let password = $state('');
	let isLogin = $state(true);
	let loading = $state(false);
	let isPasswordValid = $state(false);
	let passwordStrength = $state(0);
	let isBlueCommander = $state(true);

	const isUsernameFormatValid = $derived(/^[a-zA-Z0-9_]*$/.test(username));
	const isUsernameLengthValid = $derived(username.length >= 3 && username.length <= 50);
	const isUsernameValid = $derived(isUsernameFormatValid && isUsernameLengthValid);

	const usernameError = $derived.by(() => {
		if (!username) return '';
		if (!isUsernameFormatValid) return 'Only letters, numbers, and underscores allowed';
		if (!isLogin && !isUsernameLengthValid) {
			return username.length < 3
				? 'Username must be at least 3 characters'
				: 'Username must be 50 characters or less';
		}
		return '';
	});

	const isFormValid = $derived(
		username.length > 0 &&
			isUsernameFormatValid &&
			(isLogin || isUsernameValid) &&
			password.length > 0 &&
			(isLogin || isPasswordValid)
	);

	async function handleSubmit() {
		if (!isFormValid) return;
		loading = true;

		try {
			if (isLogin) {
				await authStore.login(username, password);
			} else {
				await authStore.register(username, password);
			}
			goto(resolve('/'));
		} catch (e) {
			toastStore.handleApiMessage(e, 'Authentication failed');
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>GoStrategy — {isLogin ? 'Sign In' : 'Enlist Commander'}</title>
</svelte:head>

{#snippet PillButton(
	onclick: () => void,
	text: string,
	disabled: boolean,
	isBlueCommander: boolean
)}
	<button
		{onclick}
		{disabled}
		class="cursor-pointer flex-1 py-2.5 text-xs font-bold uppercase tracking-wider text-center z-10 transition-colors duration-200 {isBlueCommander
			? 'text-brand-primary hover:text-brand-primary/60'
			: 'text-brand-secondary hover:text-brand-secondary/60'}"
	>
		{text}
	</button>
{/snippet}

<div class="flex items-center justify-center min-h-[80vh] p-2 sm:p-4">
	<!-- main dual-pane container -->
	<div
		class="w-full max-w-4xl grid grid-cols-1 md:grid-cols-12 glass rounded-3xl overflow-hidden shadow-2xl relative border border-white/10 animate-fade-in-slide"
	>
		<!-- Left Pane: Tactical Briefing -->
		<TacticalBriefing bind:isBlueCommander />

		<!-- Right Pane: Authentication Form -->
		<div class="col-span-1 md:col-span-7 p-6 sm:p-10 flex flex-col justify-center relative z-10">
			<div
				class="text-center md:text-left space-y-1.5 mb-6 animate-fade-in-slide delay-100 opacity-0"
			>
				<h1 class="text-2xl font-black text-white uppercase tracking-widest">
					{isLogin ? 'Sign In' : 'Enlist Commander'}
				</h1>
				<p class="text-white/50 text-xs">
					{isLogin
						? 'Enter your credentials to command your squad'
						: 'Create your commander account to join the battle'}
				</p>
			</div>

			<!-- dynamic sliding pill switcher -->
			<div
				class="relative flex p-1 bg-black/40 border border-white/10 rounded-2xl mb-6 animate-fade-in-slide delay-200 opacity-0"
			>
				<div
					class="absolute top-1 bottom-1 left-1 rounded-xl {isBlueCommander
						? 'bg-brand-primary/10 border-brand-primary/30'
						: 'bg-brand-secondary/10 border-brand-secondary/30'} transition-all duration-300 ease-out shadow-inner"
					style:width="calc(50% - 4px)"
					style:transform="translateX({isLogin ? '0' : '100%'})"
				></div>

				{@render PillButton(
					() => {
						isLogin = true;
					},
					'Sign In',
					false,
					isBlueCommander
				)}
				{@render PillButton(
					() => {
						isLogin = false;
					},
					'Sign Up',
					false,
					isBlueCommander
				)}
			</div>

			<!-- credentials form -->
			<form
				onsubmit={(e) => {
					e.preventDefault();
					handleSubmit();
				}}
				class="space-y-4 animate-fade-in-slide delay-300 opacity-0"
			>
				<!-- username -->
				<Input
					label="Username"
					placeholder="Enter commander username"
					bind:value={username}
					disabled={loading}
					error={usernameError}
					focusColor={isBlueCommander ? 'primary' : 'secondary'}
				>
					{#snippet leadingIcon()}
						<User class="size-4" />
					{/snippet}
				</Input>

				<!-- password -->
				<Input
					type="password"
					label="Password"
					placeholder="Enter password"
					bind:value={password}
					disabled={loading}
					showPasswordToggle={true}
					focusColor={isBlueCommander ? 'primary' : 'secondary'}
				>
					{#snippet leadingIcon()}
						<Lock class="size-4" />
					{/snippet}
				</Input>

				<!-- password requirements (only in register mode) -->
				{#if !isLogin && password.length > 0}
					<PasswordChecker
						{password}
						bind:isValid={isPasswordValid}
						bind:strength={passwordStrength}
					/>
				{/if}

				<!-- submit authorization -->
				<div class="pt-2">
					<Button
						type="submit"
						variant={isBlueCommander ? 'primary' : 'secondary'}
						class="w-full relative overflow-hidden"
						disabled={loading || !isFormValid}
						disabledMessage={!isLogin && password.length > 0 && passwordStrength < 4
							? 'Complete security standards first'
							: 'Enter credentials'}
						{loading}
					>
						{isLogin ? 'Enter Battlefield' : 'Deploy Commander'}
					</Button>
				</div>
			</form>
		</div>
	</div>
</div>
