<script lang="ts">
	import Button from '#lib/components/ui/Button.svelte';
	import { ChevronLeft, ChevronRight } from '@lucide/svelte';

	interface Props {
		totalPages: number;
		currentPage: number;
		isLoading: boolean;
		loadPage: (page: number) => void;
	}

	let { totalPages, currentPage, isLoading, loadPage }: Props = $props();
</script>

{#if totalPages > 1}
	<div class="flex items-center gap-4 justify-center py-2">
		<Button
			variant="outline"
			size="sm"
			disabled={currentPage === 0 || isLoading}
			onclick={() => loadPage(currentPage - 1)}
		>
			<ChevronLeft class="size-4 mr-1" />
			Previous
		</Button>

		<span class="text-xs text-white/50">
			Page {currentPage + 1} of {totalPages}
		</span>

		<Button
			variant="outline"
			size="sm"
			disabled={currentPage >= totalPages - 1 || isLoading}
			onclick={() => loadPage(currentPage + 1)}
		>
			Next
			<ChevronRight class="size-4 ml-1" />
		</Button>
	</div>
{/if}
