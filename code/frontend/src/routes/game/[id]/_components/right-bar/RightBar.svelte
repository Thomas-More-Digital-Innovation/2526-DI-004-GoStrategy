<script lang="ts">
    import Card from "$lib/components/ui/Card.svelte";
    import GameControls from "./GameControls.svelte";
    import GameHistory from "./GameHistory.svelte";
</script>

<Card class="space-y-3 max-h-125 flex flex-col">
    <GameControls />
    <hr class="border-white/10" />
    <GameHistory />
</Card>
