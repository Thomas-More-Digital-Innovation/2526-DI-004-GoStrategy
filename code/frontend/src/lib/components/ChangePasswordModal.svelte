<script lang="ts">
	import { authStore } from '$lib/state/auth.svelte';
	import { toastStore } from '$lib/state/toast.svelte';
	import Button from './ui/Button.svelte';
	import Input from './ui/Input.svelte';
	import Modal from './ui/Modal.svelte';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
	}

	let { isOpen = $bindable(), onClose }: Props = $props();

	let oldPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let loading = $state(false);

	async function handleSubmit() {
		if (!oldPassword || !newPassword || !confirmPassword) {
			toastStore.error('Please fill in all fields');
			return;
		}

		if (newPassword !== confirmPassword) {
			toastStore.error('New passwords do not match');
			return;
		}

		if (newPassword.length < 8) {
			toastStore.error('New password must be at least 8 characters long');
			return;
		}

		loading = true;

		try {
			await authStore.changePassword(oldPassword, newPassword, confirmPassword);
			toastStore.success('Password changed successfully!');
			oldPassword = '';
			newPassword = '';
			confirmPassword = '';

			setTimeout(() => {
				if (isOpen) handleClose();
			}, 1500);
		} catch (e) {
			toastStore.handleApiMessage(e, 'Failed to change password');
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		isOpen = false;
		onClose?.();
	}
</script>

<Modal
	bind:isOpen
	title="Change Password"
	description="You will be logged out from other sessions after updating."
	onClose={handleClose}
	maxWidth="md"
>
	<form
		onsubmit={(e) => {
			e.preventDefault();
			handleSubmit();
		}}
		class="space-y-4"
	>
		<Input
			type="password"
			label="Current Password"
			placeholder="Enter current password"
			bind:value={oldPassword}
			disabled={loading}
			showPasswordToggle={true}
			sanitize="password"
		/>

		<div class="h-px w-full bg-white/5 my-2"></div>

		<Input
			type="password"
			label="New Password"
			placeholder="At least 8 characters"
			bind:value={newPassword}
			disabled={loading}
			showPasswordToggle={true}
			sanitize="password"
		/>

		<Input
			type="password"
			label="Confirm New Password"
			placeholder="Re-enter new password"
			bind:value={confirmPassword}
			disabled={loading}
			showPasswordToggle={true}
			sanitize="password"
		/>

		<div class="pt-4 flex flex-col-reverse sm:flex-row gap-3">
			<Button variant="ghost" class="w-full" onclick={handleClose} disabled={loading}>
				Cancel
			</Button>
			<Button type="submit" variant="primary" class="w-full" {loading}>Update Password</Button>
		</div>
	</form>
</Modal>
