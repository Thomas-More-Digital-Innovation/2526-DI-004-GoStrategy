<script lang="ts">
	import Card from '$lib/components/ui/Card.svelte';
	import Input from '$lib/components/ui/Input.svelte';

	interface Props {
		name: string;
		description: string;
		isDefault: boolean;
		disabled?: boolean;
	}

	let {
		name = $bindable(''),
		description = $bindable(''),
		isDefault = $bindable(false),
		disabled = false
	}: Props = $props();
</script>

<Card class="space-y-6">
	<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
		<Input
			label="Setup Name"
			placeholder="e.g., Aggressive Scout Rush"
			bind:value={name}
			{disabled}
			sanitize="generic"
		/>
		<Input
			label="Description (Optional)"
			placeholder="Briefly describe your strategy..."
			bind:value={description}
			{disabled}
			sanitize="generic"
		/>
	</div>

	<div class="flex items-center gap-3">
		<input
			id="isDefault"
			type="checkbox"
			bind:checked={isDefault}
			{disabled}
			class="w-5 h-5 rounded border-white/10 bg-white/5 text-brand-accent focus:ring-brand-accent/50 cursor-pointer"
		/>
		<label for="isDefault" class="text-sm font-medium text-white/70 cursor-pointer select-none">
			Set as default setup
		</label>
	</div>
</Card>
